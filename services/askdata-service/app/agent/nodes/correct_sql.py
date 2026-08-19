import asyncio

from langgraph.runtime import Runtime

from app.agent.context import DataAgentContext
from app.agent.state import DataAgentState
from app.agent.sql_utils import extract_sql
from langchain_core.prompts import PromptTemplate
from langchain_core.output_parsers import StrOutputParser
import yaml
from app.core.log import logger
from app.prompt.prompt_loader import load_prompt
from app.agent.llm import llm



async def correct_sql(state: DataAgentState, runtime: Runtime[DataAgentContext]):
    writer = runtime.stream_writer
    writer("校正SQL")

    # 校正 SQL 仍然需要完整上下文，避免模型只根据报错修语法却改丢业务语义
    table_infos = state["table_infos"]
    metric_infos = state["metric_infos"]
    date_info = state["date_info"]
    db_info = state["db_info"]
    query = state["query"]

    # sql 是待修正的候选 SQL，error 是数据库 explain 返回的具体错误信息
    sql = state["sql"]
    error = state["error"]

    prompt = PromptTemplate(
        template=load_prompt("correct_sql"),
        input_variables=[
            "table_infos",
            "metric_infos",
            "date_info",
            "db_info",
            "query",
            "sql",
            "error",
        ],
    )

    output_parser = StrOutputParser()
    chain = prompt | llm | output_parser

    result = await chain.ainvoke(
        {
            # 与生成节点保持一致，用 YAML 向模型提供稳定、可读的结构化上下文
            "table_infos": yaml.dump(table_infos, allow_unicode=True, sort_keys=False),
            "metric_infos": yaml.dump(metric_infos, allow_unicode=True, sort_keys=False),
            "date_info": yaml.dump(date_info, allow_unicode=True, sort_keys=False),
            "db_info": yaml.dump(db_info, allow_unicode=True, sort_keys=False),
            "query": query,
            "sql": sql,
            "error": error,
        }
    )

    logger.info(f"校正后的SQL：{result}")
    return {"sql": result}

