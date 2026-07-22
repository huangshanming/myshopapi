# 本服务相关 DDL
#
# - 权威迁移脚本：`scripts/*.sql`
# - goctl model 输入：`scripts/ddl/<svc>-tables.sql`
# - 生成：`./scripts/gen-model.sh user|catalog|order|merchant|all`
# - 输出：`internal/modelgen`（goctl sqlx entity + 默认 CRUD；业务 repository 仍手写）
#
# user-service 的 `internal/model` 以 type alias 引用 modelgen；
# 其余服务若 ALTER 列尚未完全并入 ddl，运行时实体可仍在 `internal/model`，modelgen 作为 goctl 产物保留。
