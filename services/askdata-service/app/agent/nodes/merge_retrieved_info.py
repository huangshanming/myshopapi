import asyncio

from langgraph.runtime import Runtime

from app.agent.context import DataAgentContext
from app.agent.state import DataAgentState


async def merge_retrieved_info(
    state: DataAgentState, runtime: Runtime[DataAgentContext]
):
    """合并字段、指标、取值等召回结果，并补齐主外键等关联信息。"""
    writer = runtime.stream_writer
    writer("合并召回信息")
    await asyncio.sleep(0.5)
