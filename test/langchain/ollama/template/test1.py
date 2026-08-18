import asyncio

import os
from dotenv import load_dotenv
from langchain.chat_models import init_chat_model
from langchain_core.messages import HumanMessage, SystemMessage

load_dotenv()

model = init_chat_model(
    model="qwen:4b",
    model_provider="ollama",
    base_url="http://127.0.0.1:11434",
    temperature=0.7
)

def demo_message_objects():
    """推荐：显式 Message 对象，角色与字段最清晰。"""
    message = [
        SystemMessage(content="你是一个专业的数学助手，回答要简短。"),
        HumanMessage(content="你好，你是谁？")
    ]
    resp = model.invoke(message)
    print(resp.content)


def demo_tuple_list():
    """元组列表：与 ChatPromptTemplate.from_messages 的写法一致。"""
    messages = [
        ("system", "你是一个专业的数学助手，回答要简短。"),
        ("human", "你好，你是谁？")
    ]
    resp = model.invoke(messages)
    print(resp.content)

def demo_dict_list():
    """字典列表：与 ChatPromptTemplate.from_messages 的写法一致。"""
    messages = [
        {"role": "system", "content": "你是一个专业的数学助手，回答要简短。"},
        {"role": "human", "content": "你好，你是谁？"}
    ]
    resp = model.invoke(messages)
    print(resp.content)

if __name__ == "__main__":
    # demo_message_objects()
    # demo_tuple_list()
    demo_dict_list()