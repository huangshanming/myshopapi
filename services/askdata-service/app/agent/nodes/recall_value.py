from langchain_core.output_parsers import JsonOutputParser
from langchain_core.prompts import PromptTemplate
from langgraph.runtime import Runtime

from app.agent.context import DataAgentContext
from app.agent.llm import llm
from app.agent.state import DataAgentState
from app.entities.value_info import ValueInfo
from app.prompt.prompt_loader import load_prompt
from app.core.log import logger


async def recall_value(state: DataAgentState, runtime: Runtime[DataAgentContext]):
    """基于关键词从 Elasticsearch 召回字段真实取值。"""
    writer = runtime.stream_writer
    writer("召回字段取值")
    query = state["query"]
    keywords = state["keywords"]
    # 字段取值走 Elasticsearch，不需要 embedding_client
    value_es_repository = runtime.context["value_es_repository"]
     # 取值扩展关注“真实值”，例如华北地区可扩展出华北
    prompt = PromptTemplate(
        template=load_prompt("extend_keywords_for_value_recall"),
        input_variables=["query"],
    )
    output_parser = JsonOutputParser()
    chain = prompt | llm | output_parser

    result = await chain.ainvoke({"query": query})
    keywords = set(keywords + result)
    value_infos_map: dict[str, ValueInfo] = {}
    for keyword in keywords:
        # 直接用关键词查 Elasticsearch 字段值索引
        current_value_infos: list[ValueInfo] = await value_es_repository.search(keyword)

        for current_value_info in current_value_infos:
            if current_value_info.id not in value_infos_map:
                value_infos_map[current_value_info.id] = current_value_info

    retrieved_value_infos: list[ValueInfo] = list(value_infos_map.values())
    logger.info(f"检索到字段取值：{list(value_infos_map.keys())}")
    return {"retrieved_value_infos": retrieved_value_infos}