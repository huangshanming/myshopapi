from typing import TypedDict

class DataAgentState(TypedDict):
    query: str # 用户输入的查询
    error: str # 校验SQL时出现的错误信息

