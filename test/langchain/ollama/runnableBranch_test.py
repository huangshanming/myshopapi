import os

from langchain_core.runnables import RunnableBranch
from langchain.chat_models import init_chat_model
from langchain_core.prompts import ChatPromptTemplate
from dotenv import load_dotenv
from langchain_core.output_parsers import StrOutputParser

load_dotenv(encoding="utf-8")
model = init_chat_model(
    model="qwen:4b",
    model_provider="ollama",
    base_url="http://127.0.0.1:11434",
)

def determine_language(inputs):
    query = inputs["query"]
    if "日语" in query:
        return "japanese"
    elif "英语" in query:
        return "english"
    elif "韩语" in query:
        return "korean"
    else:
        return "chinese"

# 英语分支：提示词模板 + 占位符 query
english_prompt = ChatPromptTemplate.from_messages(
    [("system", "你是一个英语翻译专家，你叫小英"), ("human", "{query}")]
)

japanese_prompt = ChatPromptTemplate.from_messages(
    [("system", "你是一个日语翻译专家，你叫小日"), ("human", "{query}")]
)

korean_prompt = ChatPromptTemplate.from_messages(
    [("system", "你是一个韩语翻译专家，你叫小韩"), ("human", "{query}")]
)
parser = StrOutputParser()

chain = RunnableBranch(
   (lambda x: determine_language(x) == "japanese", japanese_prompt | model | parser),
   (lambda x: determine_language(x) == "english", english_prompt | model | parser),
   (lambda x: determine_language(x) == "korean", korean_prompt | model | parser),
)

test_queries = [
    {"query": '请你用韩语翻译这句话:"见到你很高兴"'},
    {"query": '请你用日语翻译这句话:"见到你很高兴"'},
    {"query": '请你用英语翻译这句话:"见到你很高兴"'},
]

for query_input in test_queries:
    lang = determine_language(query_input)
    print(f"检测到语言类型: {lang}")

    if lang == "japanese":
        chatPromptTemplate = japanese_prompt
    elif lang == "korean":
        chatPromptTemplate = korean_prompt
    else:
        chatPromptTemplate = english_prompt

    formatted_messages = chatPromptTemplate.format_messages(**query_input)

    for msg in formatted_messages:
        print(msg.content)

    result = chain.invoke(query_input)
    print(result)