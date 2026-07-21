package {{.pkgName}}

import (
	{{.imports}}
)

type {{.logic}} struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

{{if .hasDoc}}{{.doc}}{{end}}
func New{{.logic}}(ctx context.Context, svcCtx *svc.ServiceContext) *{{.logic}} {
	return &{{.logic}}{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *{{.logic}}) {{.function}}(ctx context.Context{{if .request}}, {{.request}}{{end}}) {{.responseType}} {
	// todo: add your logic here and delete this line

	{{.returnString}}
}
