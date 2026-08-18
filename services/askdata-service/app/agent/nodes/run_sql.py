import asyncio

from langgraph.runtime import Runtime

from app.agent.context import DataAgentContext
from app.agent.state import DataAgentState


async def run_sql(state: DataAgentState, runtime: Runtime[DataAgentContext]):
    """执行最终 SQL 并回收查询结果。"""
    writer = runtime.stream_writer
    writer("执行 SQL")
    await asyncio.sleep(0.5)
