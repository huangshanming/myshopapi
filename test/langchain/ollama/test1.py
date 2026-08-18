from langchain.chat_models import init_chat_model
from langchain_core.messages import (HumanMessage, SystemMessage)

llm = init_chat_model(
    model="qwen:4b",
    model_provider="ollama",
    base_url="http://127.0.0.1:11434",
    temperature=0.7
)
message = [SystemMessage(content="你是一个AI助手，你叫黄小鸣")]
message.append(HumanMessage(content="请介绍一下你自己"))
query = llm.invoke(message)
print(query.content)