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
# 普通调用
# message = [
#     SystemMessage(content="你是一个法律助手，只回答法律问题，超出范围回答：非法律问题无可奉告"),
#     HumanMessage(content="简单介绍下广告法，一句话 50 字以内")
# ]
# resp = model.invoke(message)
# print(type(resp))
# print(resp.content)



# 异步调用
# async def main():
#     """异步主函数：必须用 async def，内部用 await 调用 ainvoke。"""
#     resp = await model.ainvoke("解释一下LangChain是什么，简介回答")
#     print(f"响应类型：{type(resp)}")
#     print(f"响应内容：{resp.content_blocks}")

# if __name__ == "__main__":
#     asyncio.run(main())

# 流式调用
# message = [
#     SystemMessage(content="你是一个法律助手，只回答法律问题，超出范围回答：非法律问题无可奉告"),
#     HumanMessage(content="简单介绍下广告法，一句话 50 字以内")
# ]
# for chunk in model.stream(message): 
#     print(chunk.content, end="", flush=True)
# print()

message = [
    SystemMessage(content="你叫小问，是一个乐于助人的小朋友"),
    HumanMessage(content="是你谁")
]

# 异步流式调用
async def async_stream_call():
    response = model.astream(message)
    print(f"响应类型：{type(response)}")

    async for chunk in response:
        print(chunk.content, end="", flush=True)

if __name__ == "__main__":
    asyncio.run(async_stream_call())