from langchain_core.messages import SystemMessage, HumanMessage, AIMessage, ToolMessage

messages = [
    SystemMessage(content="你是一位乐于助人的智能小助手"),
    HumanMessage(content="你好，请你介绍一下你自己"),
    AIMessage(content="我是一名人工智能助手，请问您有什么想问的吗？"),
    ToolMessage(
        content='{"population": 21540000, "area": "16410平方公里"}',
        tool_call_id="call_abc123",
    ),
]
print(messages)
