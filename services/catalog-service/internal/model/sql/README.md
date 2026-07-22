# 本服务相关 DDL
#
# - 权威迁移脚本：`scripts/*.sql`
# - goctl model 输入：`scripts/ddl/catalog-tables.sql`
# - 生成：`./scripts/gen-model.sh catalog`
# - 输出：`internal/modelgen`（扁平生成；各域 `product/content/shopops/notify` 的 model 仍可手写或逐步切 alias）
