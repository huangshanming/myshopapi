import asyncio

from langgraph.runtime import Runtime

from app.agent.context import DataAgentContext
from app.agent.state import DataAgentState


async def filter_metric(state: DataAgentState, runtime: Runtime[DataAgentContext]):
    """过滤与当前问题无关的指标，保留高相关指标。"""
    writer = runtime.stream_writer
    writer("过滤指标信息")
    await asyncio.sleep(0.5)
