from __future__ import annotations

from pydantic import BaseModel, Field


class AskDataRequest(BaseModel):
    question: str = Field(
        ..., min_length=1, description="自然语言问数，如：昨天订单 GMV 是多少？"
    )
    merchant_id: int | None = Field(default=None, description="可选商家维度")
    days: int = Field(default=7, ge=1, le=90, description="默认统计窗口（天）")


class AskDataResponse(BaseModel):
    answer: str
    steps: list[str] = Field(default_factory=list)
    metrics: dict = Field(default_factory=dict)
    sql_preview: str | None = None
    source: str = "rule_stub"


class AskDataAgent:
    """电商问数 Agent。

    现状：关键词规则 stub，返回可解释步骤与占位指标。
    后续：LangGraph + tools（查 order / catalog / 数仓）、NL2SQL、权限与审计。
    """

    async def query(self, req: AskDataRequest) -> dict:
        q = req.question.strip()
        steps = [
            f"理解问题：{q}",
            "识别指标与时间范围",
            "规划数据查询（占位，未真实打库）",
            "汇总并生成可读回答",
        ]

        metrics: dict = {"window_days": req.days}
        if req.merchant_id is not None:
            metrics["merchant_id"] = req.merchant_id

        # 极简意图识别，便于联调前端 / 网关
        if any(k in q for k in ("gmv", "GMV", "成交额", "销售额", "营收")):
            metrics["metric"] = "gmv"
            metrics["value"] = None
            answer = (
                f"已识别 GMV 类问题（近 {req.days} 天）。"
                "当前为规则占位，尚未连接订单库；后续将生成 SQL / 调用 order-service。"
            )
            sql_preview = (
                "SELECT DATE(created_at) AS d, SUM(pay_amount) AS gmv "
                f"FROM orders WHERE created_at >= CURDATE() - INTERVAL {req.days} DAY "
                "GROUP BY d"
            )
        elif any(k in q for k in ("订单", "单量", "order")):
            metrics["metric"] = "order_count"
            metrics["value"] = None
            answer = (
                f"已识别订单量类问题（近 {req.days} 天）。"
                "当前为规则占位，尚未连接订单库。"
            )
            sql_preview = (
                "SELECT COUNT(*) AS order_count FROM orders "
                f"WHERE created_at >= CURDATE() - INTERVAL {req.days} DAY"
            )
        elif any(k in q for k in ("库存", "缺货", "stock")):
            metrics["metric"] = "inventory"
            metrics["value"] = None
            answer = "已识别库存类问题。当前为规则占位，后续可对接 catalog / inventory。"
            sql_preview = None
        else:
            metrics["metric"] = "unknown"
            answer = (
                f"暂未识别「{q}」对应的标准指标。"
                "可尝试询问：GMV、订单量、库存等；后续将接入完整 NL2SQL Agent。"
            )
            sql_preview = None

        resp = AskDataResponse(
            answer=answer,
            steps=steps,
            metrics=metrics,
            sql_preview=sql_preview,
            source="rule_stub",
        )
        return resp.model_dump()
