import asyncio

from langgraph.runtime import Runtime

from app.agent.context import DataAgentContext
from app.agent.state import DataAgentState


async def generate_sql(state: DataAgentState, runtime: Runtime[DataAgentContext]):
    """根据召回与过滤后的元数据，调用 LLM 生成 SQL。"""
    writer = runtime.stream_writer
    writer("生成 SQL")
    await asyncio.sleep(0.5)
