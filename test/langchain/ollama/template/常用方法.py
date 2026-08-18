from langchain_core.prompts import PromptTemplate, ChatPromptTemplate, MessagesPlaceholder
# 对话提示词
from langchain_core.messages import HumanMessage, AIMessage

# template = PromptTemplate.from_template(
#     "你是一个专业的{role}工程师，请回答我的问题，我的问题是：{question}"
# )

# prompt_str = template.format(role="python开发", question="二分法怎么写")
# print(prompt_str)


# chat_prompt = ChatPromptTemplate.from_messages([
#     ("system", "你是一个专业的{role}工程师，请回答我的问题"),
#     ("user", "{question}"),
# ])

# prompt_str = chat_prompt.format_messages(role="python开发", question="二分法怎么写")
# print(prompt_str)


prompt = ChatPromptTemplate.from_messages([
    ("system", "你是一个资深的Python应用开发工程师，请认真回答我提出的Python相关的问题"),
    MessagesPlaceholder("memory"),
    ("human", "{question}")
])

prompt_value = prompt.invoke({
    "memory": [
        HumanMessage(content="我的名字叫亮仔，是一名程序员"),
        AIMessage(content="好的，亮仔你好")
    ],
    "question": "请问我的名字叫什么？"
})

print(prompt_value.to_string())
