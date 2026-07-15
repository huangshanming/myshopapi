package httpserver

import (
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/rest"
)

// NewRest 创建 go-zero rest Server（不打印无用 banner）
func NewRest(port int, mode string) *rest.Server {
	c := rest.RestConf{}
	c.Host = "0.0.0.0"
	c.Port = port
	c.Timeout = 30000
	if mode == "release" || mode == "prod" || mode == "production" {
		c.Mode = "pro"
	} else {
		c.Mode = "dev"
	}
	c.Log.Mode = "console"
	c.Log.Encoding = "plain"
	return rest.MustNewServer(c, rest.WithCors())
}

func Healthz(service string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = fmt.Fprintf(w, `{"status":"ok","service":"%s"}`, service)
	}
}
