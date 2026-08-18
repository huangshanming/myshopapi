"""
【案例】自定义流（custom）最简版：在节点内通过 get_stream_writer() 写入任意可序列化数据，stream 侧用 custom 接收。

对应教程章节：第 25 章 - LangGraph 高级特性 → 1、流式处理（Streaming）

知识点速览：
- 这是 `custom` 模式的最小案例，重点不是业务逻辑，而是先看懂“节点内部怎么主动写出一段流式消息”。
- `get_stream_writer()` 仅在图执行（stream/astream）过程中有效；调用 `graph.stream` 时，`stream_mode` 里必须包含 `custom`。
- 自定义块与状态更新是分开的：前者更适合 UI/日志/进度提示，后者仍然通过 State 和 Reducer 管理。
"""

from typing import TypedDict

from langgraph.config import get_stream_writer
from langgraph.graph import StateGraph, START, END


class State(TypedDict):
    query: str
    answer: str


def node(state: State):
    writer = get_stream_writer()
    writer({"custom_key": "欢迎来到线上Agent班级学习，O(∩_∩)O"})
    return {"answer": "some data"}


def main():
    graph = (
        StateGraph(State)
        .add_node(node)
        .add_edge(START, "node")
        .add_edge("node", END)
        .compile()
    )

    # 仅 custom：for chunk in graph.stream({"query": "example"}, stream_mode=["custom"]): print(chunk)
    # custom + updates：for mode, chunk in graph.stream(..., stream_mode=["updates", "custom"]): ...
    for chunk in graph.stream({"query": "example"}, stream_mode=["values", "custom"]):
        print(chunk)


if __name__ == "__main__":
    main()

"""
【输出示例】
('values', {'query': 'example'})
('custom', {'custom_key': '欢迎来到线上Agent班级学习，O(∩_∩)O'})
('values', {'query': 'example', 'answer': 'some data'})
"""
