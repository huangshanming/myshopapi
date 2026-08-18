from langchain.chat_models import (
    init_chat_model,
)

import os
from dotenv import load_dotenv

load_dotenv(encoding="utf-8")  # encoding 指定 utf-8，避免 .env 中中文注释乱码

model = init_chat_model(
    model="qwen-plus",
    model_provider="openai",
    api_key=os.getenv("QWEN_API_KEY"),
    base_url="https://dashscope.aliyuncs.com/compatible-mode/v1",  # 阿里百炼 OpenAI 兼容接口地址
)

for chunk in model.stream("你是谁？"):
    print(chunk.content, end="", flush=True)