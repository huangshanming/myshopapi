import asyncio

from langgraph.runtime import Runtime

from app.agent.context import DataAgentContext
from app.agent.state import DataAgentState


async def filter_table(state: DataAgentState, runtime: Runtime[DataAgentContext]):
    """过滤无关表，保留生成 SQL 所需的核心表集合。"""
    writer = runtime.stream_writer
    writer("过滤表信息")
    await asyncio.sleep(0.5)
