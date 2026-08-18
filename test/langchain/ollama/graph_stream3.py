"""
【案例】messages 流模式：从图中调用 LLM 的节点逐 token（或片段）推送输出，便于打字机效果。

对应教程章节：第 25 章 - LangGraph 高级特性 → 1、流式处理（Streaming）

知识点速览：
- stream_mode="messages" 时，每次迭代一般为 (message_chunk, metadata)：chunk 为模型输出片段，metadata 标明节点等上下文。
- 这个案例最适合用来建立“LangGraph Streaming 不只流状态，也能流模型输出”这层认知。
- 流式消费侧通常关心 `chunk.content` 和 `metadata`；前者是输出片段，后者帮助你知道这些片段来自哪个节点。
- 需配置环境变量（如 aliQwen-api）与网络；模型、base_url 按你本地教程为准。
"""

import os
from typing import TypedDict

from langchain.chat_models import init_chat_model
from langgraph.graph import StateGraph, START

from dotenv import load_dotenv

load_dotenv(encoding="utf-8")


class State(TypedDict):
    query: str
    answer: str


def node(state: State):
    print("开始调用 node 节点")

    model = init_chat_model(
        model="qwen2.5:7b",
        model_provider="ollama",
        base_url="http://127.0.0.1:11434",
    )

    llm_result = model.invoke([("user", state["query"])])
    print("llm invoke 结束", end="\n\n")

    return {"answer": llm_result}


def main():
    graph = (
        StateGraph(state_schema=State).add_node(node).add_edge(START, "node").compile()
    )

    inputs = {"query": "帮我生成一个200字的小学生作文，主题为我的一天"}

    # messages：从图内触发的大模型调用处流式输出；(chunk, metadata) 见官方文档
    for chunk, _metadata in graph.stream(inputs, stream_mode="messages"):
        # print(f"type of chunk:{type(chunk)}")  # 调试时可打开
        print(chunk.content, end="")
        # print(chunk, end="")


if __name__ == "__main__":
    main()

"""
【输出示例】
(.venv) didilili@DidililiMacBook-Pro streaming % python3 StreamLLMTokens.py
开始调用 node 节点
我的一天  

清晨，阳光悄悄爬上窗台，我伸个懒腰起床了！吃完妈妈做的香喷喷的煎蛋和牛奶，背上书包去上学。课堂上，我认真听讲，积极举手回答问题；课间和好朋友跳皮筋、讲故事，笑声像铃铛一样清脆。中午吃食堂的番茄炒蛋盖饭，暖暖的真好吃！放学后，我先完成作业，再陪小猫“团团”玩一会儿毛线球。晚饭后，我和爸爸一起读绘本，妈妈教我折了一只纸鹤，翅膀还微微翘着呢！临睡前，我刷牙洗脸，把小书包整理好，明天还要早起升旗呢！这一天像一颗甜甜的水果糖——有学习的酸、玩耍的甜、家人的暖，还有成长的光。我爱这充实又快乐的一天！（198字）llm invoke 结束
"""
