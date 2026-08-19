import asyncio

from langgraph.runtime import Runtime

from app.agent.context import DataAgentContext
from app.agent.state import DataAgentState
from app.agent.sql_utils import extract_sql
from app.repositories.mysql.dw.dw_mysql_repository import DWMySQLRepository
from app.core.log import logger


async def validate_sql(state: DataAgentState, runtime: Runtime[DataAgentContext]):
    """校验 SQL 语法与表字段合法性（如 EXPLAIN）。"""
    writer = runtime.stream_writer
    writer("校验 SQL")
    # 读取 generate_sql 写入状态的 SQL，并去掉模型可能残留的 markdown 围栏。
    sql = extract_sql(state["sql"])

     # SQL 可用性必须交给真实数仓判断，这里从运行时 context 中取 DW Repository。
    dw_mysql_repository: DWMySQLRepository = runtime.context["dw_mysql_repository"]
    try:
        # validate 内部使用 explain <sql>，只关心数据库能否成功解析这条 SQL。
        await dw_mysql_repository.validate(sql)
        logger.info("SQL语法正确")
        return {"error": None}
    except Exception as e:
        # 不直接抛异常，而是把错误信息写入 state，交给条件边判断。
        logger.info(f"SQL语法错误：{str(e)}")
        return {"error": str(e)}