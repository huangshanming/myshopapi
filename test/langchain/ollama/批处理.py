import re
from langchain_core.prompts import PromptTemplate
from langchain.chat_models import init_chat_model
from dotenv import load_dotenv
import asyncio

load_dotenv()

model = init_chat_model(
    model="qwen:4b",
    model_provider="ollama",
    base_url="http://127.0.0.1:11434",
    temperature=0.7
)

from langchain_core.messages import (SystemMessage, HumanMessage, AIMessage)

# questions = [
#     "什么是Redis？简洁100字以内",
#     "Python的生成器是做什么的"
# ]

# response = model.batch(questions)
# for q, r in zip(questions, response):
#     print(f"问题：{q}\n回答：{r.content}\n")

# 异步批量
questions = [
    "什么是redis?简洁回答，字数控制在100以内",
    "Python的生成器是做什么的？简洁回答，字数控制在100以内",
    "解释一下Docker和Kubernetes的关系?简洁回答，字数控制在100以内",
]

async def async_batch_call():
    response = await model.abatch(questions)
    print(f"响应类型：{type(response)}")

    for q, r in zip(questions, response):
        print(f"问题：{q}\n回答：{r.content}\n")

if __name__ == "__main__":
    asyncio.run(async_batch_call())