from langchain_openai import ChatOpenAI

import os
from dotenv import load_dotenv

load_dotenv(encoding="utf-8")  # encoding 指定 utf-8，避免 .env 中中文注释乱码

llm = ChatOpenAI(model="deepseek-v3.2", api_key=os.getenv("QWEN_API_KEY"),
    base_url="https://dashscope.aliyuncs.com/compatible-mode/v1",  # 阿里百炼 OpenAI 兼容接口地址
)

response = llm.invoke("你是谁？")

print(response)