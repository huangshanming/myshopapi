import asyncio

from langgraph.runtime import Runtime

from app.agent.context import DataAgentContext
from app.agent.state import DataAgentState
from app.prompt.prompt_loader import load_prompt
from langchain_core.prompts import PromptTemplate
from langchain_core.output_parsers import StrOutputParser
from app.agent.llm import llm
from app.agent.sql_utils import extract_sql
from app.core.log import logger
import yaml


async def generate_sql(state: DataAgentState, runtime: Runtime[DataAgentContext]):
    """根据召回与过滤后的元数据，调用 LLM 生成 SQL。"""
    writer = runtime.stream_writer
    writer("生成SQL")
    
    # 这些上下文都由前置节点准备完成，模型只在给定表、字段、指标口径范围内生成 SQL
    table_infos = state["table_infos"]
    metric_infos = state["metric_infos"]
    date_info = state["date_info"]
    db_info = state["db_info"]
    query = state["query"]

    prompt = PromptTemplate(
        template=load_prompt("generate_sql"),
        input_variables=["table_infos", "metric_infos", "date_info", "db_info", "query"],
    )

    # 生成 SQL 只需要一段纯文本，所以这里使用 StrOutputParser
    output_parser = StrOutputParser()
    chain = prompt | llm | output_parser
    body =  {
        # YAML 更适合放进提示词：保留嵌套结构、顺序和中文说明，方便模型理解表字段关系
        "table_infos": yaml.dump(table_infos, allow_unicode=True, sort_keys=False),
        "metric_infos": yaml.dump(metric_infos, allow_unicode=True, sort_keys=False),
        "date_info": yaml.dump(date_info, allow_unicode=True, sort_keys=False),
        "db_info": yaml.dump(db_info, allow_unicode=True, sort_keys=False),
        "query": query,
    }
    result = await chain.ainvoke(
        body
    )

    sql = extract_sql(result)
    logger.info(f"生成的SQL Body：{body}")
    logger.info(f"生成的SQL：{sql}")
    return {"sql": sql}