# Embedding 模型（TEI / bge-large-zh-v1.5）

`docker-compose` 中的 `embedding` 服务会挂载本目录下的模型：

```text
deploy/local/embedding/bge-large-zh-v1.5  →  /models/bge-large-zh-v1.5
```

## 下载模型

任选其一（需能访问 Hugging Face，或用镜像站）：

```bash
cd deploy/local/embedding
# 使用 huggingface-cli
pip install -U "huggingface_hub[cli]"
huggingface-cli download BAAI/bge-large-zh-v1.5 --local-dir bge-large-zh-v1.5
```

或使用 `git lfs`：

```bash
git lfs install
git clone https://huggingface.co/BAAI/bge-large-zh-v1.5 bge-large-zh-v1.5
```

## 启动

```bash
docker compose -f deploy/local/docker-compose.yaml --profile embedding up -d embedding embedding-proxy
# 或 infra：
docker compose -f deploy/local/docker-compose.infra.yaml --profile embedding up -d embedding embedding-proxy
```

探测：`curl http://127.0.0.1:8082/health`

> 注意：本机若跑着 VirtualBox/Homestead，**8081 常被 VBoxHeadless 占用**，因此默认映射为 **8082**。

说明：

- TEI 官方 CPU 镜像是 `linux/amd64`。在 Apple Silicon 上用 QEMU 跑时，**直接映射端口常会 `Connection reset`**；因此宿主机经 `embedding-proxy`（Envoy）访问。
- 首次 / 重启后需要暖机（可能数分钟）。日志出现 `Ready` 后再测。
- 可用 `docker logs -f mymall-embedding` 等到 `Ready`。

> 模型体积较大，请勿提交到 git（见同目录 `.gitignore`）。
