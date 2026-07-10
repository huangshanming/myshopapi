package health

import (
	"context"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
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

func (r *Registry) ReadyHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		failed := make(map[string]string)
		r.mu.RLock()
		defer r.mu.RUnlock()
		for name, checker := range r.checkers {
			if err := checker(ctx); err != nil {
				failed[name] = err.Error()
			}
		}
		if len(failed) > 0 {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready", "checks": failed})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	}
}
