package {{.pkgName}}

import (
	{{.imports}}
)

type {{.logic}} struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

{{if .hasDoc}}{{.doc}}{{end}}
func New{{.logic}}(svcCtx *svc.ServiceContext) *{{.logic}} {
	return &{{.logic}}{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *{{.logic}}) {{.function}}(ctx context.Context{{if .request}}, {{.request}}{{end}}) {{.responseType}} {
	// todo: add your logic here and delete this line

	{{.returnString}}
}
