package health

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
)

type Checker func(ctx context.Context) error

type Registry struct {
	mu       sync.RWMutex
	checkers map[string]Checker
}

func NewRegistry() *Registry {
	return &Registry{checkers: make(map[string]Checker)}
}

func (r *Registry) Register(name string, checker Checker) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.checkers[name] = checker
}

func (r *Registry) CheckAll(ctx context.Context) (ready bool, failed map[string]string) {
	failed = make(map[string]string)
	r.mu.RLock()
	defer r.mu.RUnlock()
	for name, checker := range r.checkers {
		if err := checker(ctx); err != nil {
			failed[name] = err.Error()
		}
	}
	return len(failed) == 0, failed
}

func (r *Registry) ReadyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ok, failed := r.CheckAll(req.Context())
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if !ok {
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"status": "not_ready", "checks": failed})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
	}
}
