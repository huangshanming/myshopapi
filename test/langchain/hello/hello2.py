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

response = model.invoke("你是谁？").content

print(response)
print("*" * 50)

model2 = init_chat_model(
    model="deepseek-v3.2",
    model_provider="openai",
    api_key=os.getenv("QWEN_API_KEY"),
    base_url="https://dashscope.aliyuncs.com/compatible-mode/v1",  # 阿里百炼 OpenAI 兼容接口地址
)

response2 = model2.invoke("你是谁？").content

print(response2)