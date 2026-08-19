from typing import TypedDict
from app.entities.column_info import ColumnInfo
from app.entities.metric_info import MetricInfo
from app.entities.value_info import ValueInfo



class MetricInfoState(TypedDict):
    # 面向提示词的指标结构，不再保留内部 id
    name: str
    description: str
    relevant_columns: list[str]
    alias: list[str]


class ColumnInfoState(TypedDict):
    # 面向提示词的字段结构，只保留 SQL 生成需要理解的字段属性
    name: str
    type: str
    role: str
    examples: list
    description: str
    alias: list[str]


class TableInfoState(TypedDict):
    # 表信息会嵌套字段列表，方便后续按表组织 SQL 上下文
    name: str
    role: str
    description: str
    columns: list[ColumnInfoState]

class DateInfoState(TypedDict):
    date: str
    weekday: str
    quarter: str


class DBInfoState(TypedDict):
    dialect: str
    version: str

class DataAgentState(TypedDict):
    sql: str  # 生成的 SQL
    query: str  # 用户输入的查询
    keywords: list[str]  # 抽取的关键词

    retrieved_column_infos: list[ColumnInfo]  # 检索到的字段信息
    retrieved_metric_infos: list[MetricInfo]  # 检索到的指标信息
    retrieved_value_infos: list[ValueInfo]  # 检索到的取值信息

    table_infos: list[TableInfoState]  # 合并和补齐后的表结构上下文
    metric_infos: list[MetricInfoState]  # 合并后的指标上下文
    date_info: DateInfoState  # 当前日期 星期和季度信息
    db_info: DBInfoState  # 数据库方言和版本信息
    error: str | None  # 校验 SQL 时出现的错误信息
