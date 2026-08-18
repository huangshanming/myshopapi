# 输出解析器
from langchain_core.output_parsers import JsonOutputParser
from langchain_core.prompts import PromptTemplate, ChatPromptTemplate
from langchain.chat_models import init_chat_model
from dotenv import load_dotenv
from pydantic import BaseModel, Field
from loguru import logger
import asyncio

load_dotenv()
model = init_chat_model(
    model="qwen:4b",
    model_provider="ollama",
    base_url="http://127.0.0.1:11434",
    temperature=0.7
)
# 只用输出解析器


# parser = JsonOutputParser()

# result = model.invoke("请用 JSON 返回：name、age")
# data = parser.invoke(result)
# print(data)

# 结构化输出
from typing import TypedDict

# class Person(TypedDict):
#     name: str
#     age: int

# structured_model = model.with_structured_output(Person)

# result = structured_model.invoke("请返回一个世界上最帅的人")
# print(result)

# chat_prompt = ChatPromptTemplate.from_messages([
#     ("system", "你是一个{role}，请简单回答我的问题，结果返回json格式,q字段表示问题，a字段表示答案。"),
#     ("human", "请回答：{question}")
# ])

# prompt = chat_prompt.format_messages(role="losser", question="请回答：你为啥如此失败？")

# result = model.invoke(prompt)
# print(result.content)


# class Person(BaseModel):
#     """定义一条「新闻」的结构：时间、人物、事件。用于约束模型输出的 JSON 形状。"""
#     time: str = Field(description="时间")
#     person: str = Field(description="人物")
#     event: str = Field(description="事件")

# parser = JsonOutputParser(pydantic_object=Person)

# format_instructions = parser.get_format_instructions()

# chat_prompt = ChatPromptTemplate.from_messages([
#     ("system", "你是一个AI助手，你只能输出结构化JSON数据。"),
#     ("human", "请生成一个关于{topic}的新闻。{format_instructions}"),
# ])
# prompt = chat_prompt.format_messages(
#     topic="小米su7跑车", format_instructions=format_instructions
# )
# print(model.invoke(prompt).content)

# Pydantic + Annotated
from typing import Annotated
from pydantic import BaseModel, Field, ValidationError

# Age = Annotated[int, Field(ge=0, le=150,description="年龄, 范围0-150")]
# Name = Annotated[str, Field(description="名字")]

# class Person(BaseModel):
#     age: int
#     name: str
#     age2: Age
# try:
#     person = Person(age=100, name="张三", age2=100)
# except ValidationError as e:
#     print("数据校验失败：")
#     print(e)



schema = {
    "type": "object",
    "properties": {
        "name": {"type": "string"},
        "skill": {"type": "string"},
        "work_year": {"type": "integer", "minimum": 0}
    },
    "required": ["name", "skill", "work_year"]
}
# 直接把schema传给结构化输出
chain_schema = model.with_structured_output(schema)
res5 = chain_schema.invoke("生成一名云原生开发人员信息")
print("【5. JSON Schema 跨语言结构输出】")
print(type(res5), res5)
