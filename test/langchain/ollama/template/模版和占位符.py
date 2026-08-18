from langchain_core.prompts import PromptTemplate
from langchain.chat_models import init_chat_model
from dotenv import load_dotenv

load_dotenv()

model = init_chat_model(
    model="qwen:4b",
    model_provider="ollama",
    base_url="http://127.0.0.1:11434",
    temperature=0.7
)

template = "用不超过50字介绍：{topic} 是什么？"

prompt = template.format(topic="Python")
print(prompt)

# 多角色消息列表
from langchain_core.messages import (SystemMessage, HumanMessage, AIMessage)

messages = [
    SystemMessage(content="你是一个专业的数学助手，回答要简短。"),
    HumanMessage(content="你好，你是谁？")
]

resp = model.invoke(messages)
print(resp)

## 元组列表与字典列表
messages_as_tuples = [
    ("system", "你是一个专业的数学助手，回答要简短。"),
    ("user", "你好，你是谁？")
]

messages_as_dicts = [
    {"role": "system", "content": "你是一个专业的数学助手，回答要简短。"},
    {"role": "user", "content": "你好，你是谁？"}
]

resp = model.invoke(messages_as_tuples)
print(resp)

resp = model.invoke(messages_as_dicts)
print(resp)