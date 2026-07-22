// Package zrpcx helpers for go-zero zrpc dial / server conf.
// Prefer Etcd when Hosts is non-empty; otherwise dial Endpoints directly.
package zrpcx

import (
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
)

const (
	KeyUser     = "user.rpc"
	KeyCatalog  = "catalog.rpc"
	KeyMerchant = "merchant.rpc"
)

// ParseHosts splits comma/space separated etcd hosts.
func ParseHosts(v string) []string {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}
	parts := strings.FieldsFunc(v, func(r rune) bool {
		return r == ',' || r == ' ' || r == ';'
	})
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// Dial builds a zrpc client: Etcd discovery when etcdHosts+key set, else direct Endpoints.
func Dial(endpoint string, etcdHosts []string, key string) (zrpc.Client, error) {
	c := zrpc.RpcClientConf{NonBlock: true}
	if len(etcdHosts) > 0 && strings.TrimSpace(key) != "" {
		c.Etcd = discov.EtcdConf{Hosts: etcdHosts, Key: key}
	} else if strings.TrimSpace(endpoint) != "" {
		c.Endpoints = []string{endpoint}
	} else {
		return nil, fmt.Errorf("zrpcx: need endpoint or etcd hosts+key")
	}
	cli, err := zrpc.NewClient(c)
	if err != nil {
		return nil, err
	}
	return cli, nil
}

// ServerConf builds RpcServerConf; registers to etcd when etcdHosts+key set.
func ServerConf(listenOn string, etcdHosts []string, key string) zrpc.RpcServerConf {
	c := zrpc.RpcServerConf{ListenOn: listenOn}
	if len(etcdHosts) > 0 && strings.TrimSpace(key) != "" {
		c.Etcd = discov.EtcdConf{Hosts: etcdHosts, Key: key}
	}
	c.Mode = service.DevMode
	c.Log.Mode = "console"
	c.Log.Encoding = "plain"
	return c
}
