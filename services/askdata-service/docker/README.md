# 本服务内的 Docker 资源说明
#
# 本仓库是 monorepo：Elasticsearch / Kibana / Qdrant / Embedding 统一放在
#   deploy/local/docker-compose.yaml
#   deploy/local/docker-compose.infra.yaml
#   deploy/local/embedding/
#
# 本目录保留与独立「电商问数」工程一致的结构，便于放：
#   - mysql/        问数相关建表 SQL、种子数据
#   - elasticsearch/  索引 mapping、插件说明
#   - embedding/    若需服务私有模型挂载说明（默认仍用 deploy/local/embedding）
