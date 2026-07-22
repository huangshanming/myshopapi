#!/usr/bin/env python3
"""Finish opaque entity responses: DataResp alias + order concrete resps."""
from __future__ import annotations

import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
SERVICES = [
    ROOT / "services/catalog-service",
    ROOT / "services/merchant-service",
    ROOT / "services/order-service",
    ROOT / "services/user-service",
]
ORD = ROOT / "services/order-service"


def ensure_data_resp_alias(svc: Path) -> None:
    types_go = svc / "internal/types/types.go"
    if not types_go.exists():
        return
    t = types_go.read_text()
    if "type DataResp =" in t or "type DataResp struct" in t:
        return
    if "type AnyResp struct" not in t:
        return
    # alias after AnyResp block
    t = t.replace(
        "type AnyResp struct {\n\tData interface{} `json:\"data,optional\"`\n}\n",
        "type AnyResp struct {\n\tData interface{} `json:\"data,optional\"`\n}\n\n"
        "// DataResp is the preferred name for opaque entity JSON bodies.\n"
        "type DataResp = AnyResp\n",
    )
    types_go.write_text(t)
    print("alias DataResp", svc.name)

    api_dir = svc / "api"
    for api in api_dir.glob("*.api"):
        at = api.read_text()
        if "type DataResp {" in at or "type DataResp =" in at:
            pass
        elif "type AnyResp {" in at:
            at = at.replace(
                "type AnyResp {\n\tData interface{} `json:\"data,optional\"`\n}",
                "type AnyResp {\n\tData interface{} `json:\"data,optional\"`\n}\n\n"
                "type DataResp {\n\tData interface{} `json:\"data,optional\"`\n}",
            )
        # entity routes: AnyResp → DataResp (keep Empty/PageList/typed intact)
        at = at.replace("returns (AnyResp)", "returns (DataResp)")
        api.write_text(at)
        print("api DataResp", api.name)


def replace_logic_any_to_data(svc: Path) -> int:
    n = 0
    logic = svc / "internal/logic"
    if not logic.exists():
        return 0
    for p in logic.rglob("*_logic.go"):
        text = p.read_text()
        orig = text
        # only when still returning AnyResp with Data payload (entity wrap)
        if "*types.AnyResp" not in text and "types.AnyResp{" not in text:
            continue
        # skip if already mixed with concrete conversions in same file carefully
        text = text.replace("*types.AnyResp", "*types.DataResp")
        text = text.replace("&types.AnyResp{", "&types.DataResp{")
        text = text.replace("types.AnyResp{", "types.DataResp{")
        if text != orig:
            p.write_text(text)
            n += 1
    if n:
        print("logic DataResp", svc.name, n, "files")
    return n


def finish_order() -> None:
    biz = ORD / "internal/types/biz_types.go"
    t = biz.read_text()
    if "type ReviewEligibleResp struct" not in t:
        t += """
type ReviewEligibleResp struct {
	Eligible bool   `json:"eligible"`
	Reviewed bool   `json:"reviewed"`
	Status   int8   `json:"status"`
	OrderId  uint64 `json:"order_id"`
}

type CouponPreviewResp struct {
	GoodsAmount      float64     `json:"goods_amount"`
	DiscountAmount   float64     `json:"discount_amount"`
	PayAmount        float64     `json:"pay_amount"`
	BestUserCouponID uint64      `json:"best_user_coupon_id"`
	Available        interface{} `json:"available"`
	Unavailable      interface{} `json:"unavailable"`
}
"""
        biz.write_text(t)
        print("+order ReviewEligible/CouponPreview")

    types_go = ORD / "internal/types/types.go"
    tg = types_go.read_text()
    # sync first-class types into goctl file
    for block, needle in [
        (
            "type URLResp struct {\n\tUrl string `json:\"url\"`\n}\n\n",
            "type URLResp struct",
        ),
        (
            "type OrderDetailResp struct {\n\tOrder      interface{} `json:\"order\"`\n\tAfterSales interface{} `json:\"after_sales\"`\n}\n\n",
            "type OrderDetailResp struct",
        ),
        (
            "type ListResp struct {\n\tList interface{} `json:\"list\"`\n}\n\n",
            "type ListResp struct",
        ),
        (
            "type ReviewEligibleResp struct {\n\tEligible bool `json:\"eligible\"`\n\tReviewed bool `json:\"reviewed\"`\n\tStatus int8 `json:\"status\"`\n\tOrderId uint64 `json:\"order_id\"`\n}\n\n",
            "type ReviewEligibleResp struct",
        ),
        (
            "type CouponPreviewResp struct {\n\tGoodsAmount float64 `json:\"goods_amount\"`\n\tDiscountAmount float64 `json:\"discount_amount\"`\n\tPayAmount float64 `json:\"pay_amount\"`\n\tBestUserCouponID uint64 `json:\"best_user_coupon_id\"`\n\tAvailable interface{} `json:\"available\"`\n\tUnavailable interface{} `json:\"unavailable\"`\n}\n\n",
            "type CouponPreviewResp struct",
        ),
    ]:
        if needle not in tg:
            tg = tg.replace(
                "type EmptyResp struct {\n}\n",
                "type EmptyResp struct {\n}\n\n" + block,
            )
    # remove duplicate URLResp from biz if also in types
    types_go.write_text(tg)

    api = next((ORD / "api").glob("*.api"))
    at = api.read_text()
    if "type ReviewEligibleResp {" not in at:
        at = at.replace(
            "type ListResp {\n\tList interface{} `json:\"list\"`\n}",
            "type ListResp {\n\tList interface{} `json:\"list\"`\n}\n\n"
            "type ReviewEligibleResp {\n\tEligible bool `json:\"eligible\"`\n\tReviewed bool `json:\"reviewed\"`\n"
            "\tStatus int8 `json:\"status\"`\n\tOrderId uint64 `json:\"order_id\"`\n}\n\n"
            "type CouponPreviewResp {\n\tGoodsAmount float64 `json:\"goods_amount\"`\n"
            "\tDiscountAmount float64 `json:\"discount_amount\"`\n\tPayAmount float64 `json:\"pay_amount\"`\n"
            "\tBestUserCouponID uint64 `json:\"best_user_coupon_id\"`\n"
            "\tAvailable interface{} `json:\"available\"`\n\tUnavailable interface{} `json:\"unavailable\"`\n}",
        )
    at = at.replace(
        "post /api/v1/orders/coupon-preview (CouponPreviewReq) returns (DataResp)",
        "post /api/v1/orders/coupon-preview (CouponPreviewReq) returns (CouponPreviewResp)",
    )
    at = at.replace(
        "get /api/v1/orders/:id/review-eligible (IdPathReq) returns (DataResp)",
        "get /api/v1/orders/:id/review-eligible (IdPathReq) returns (ReviewEligibleResp)",
    )
    # also if still AnyResp
    at = at.replace(
        "post /api/v1/orders/coupon-preview (CouponPreviewReq) returns (AnyResp)",
        "post /api/v1/orders/coupon-preview (CouponPreviewReq) returns (CouponPreviewResp)",
    )
    at = at.replace(
        "get /api/v1/orders/:id/review-eligible (IdPathReq) returns (AnyResp)",
        "get /api/v1/orders/:id/review-eligible (IdPathReq) returns (ReviewEligibleResp)",
    )
    api.write_text(at)

    # biz ReviewEligible → typed
    rev = ORD / "internal/biz/review_logic.go"
    rt = rev.read_text()
    if "ReviewEligibleResp" not in rt:
        if '"mymall/services/order-service/internal/types"' not in rt:
            rt = rt.replace(
                '"mymall/services/order-service/internal/svc"',
                '"mymall/services/order-service/internal/svc"\n\t"mymall/services/order-service/internal/types"',
            )
        rt = rt.replace(
            "func (l *ReviewLogic) ReviewEligible(ctx context.Context, userID, orderID uint64) (map[string]interface{}, error) {\n"
            "\torder, err := l.svcCtx.Repo.FindByID(ctx, orderID, userID)\n"
            "\tif err != nil {\n"
            '\t\treturn nil, errors.New("订单不存在")\n'
            "\t}\n"
            "\texists, _ := l.svcCtx.Reviews.ExistsByOrderID(ctx, orderID)\n"
            "\treturn map[string]interface{}{\n"
            "\t\t\"eligible\": order.Status == model.OrderStatusCompleted && !exists,\n"
            "\t\t\"reviewed\": exists || order.Status == model.OrderStatusReviewed,\n"
            "\t\t\"status\":   order.Status,\n"
            "\t\t\"order_id\": order.ID,\n"
            "\t}, nil\n"
            "}",
            "func (l *ReviewLogic) ReviewEligible(ctx context.Context, userID, orderID uint64) (*types.ReviewEligibleResp, error) {\n"
            "\torder, err := l.svcCtx.Repo.FindByID(ctx, orderID, userID)\n"
            "\tif err != nil {\n"
            '\t\treturn nil, errors.New("订单不存在")\n'
            "\t}\n"
            "\texists, _ := l.svcCtx.Reviews.ExistsByOrderID(ctx, orderID)\n"
            "\treturn &types.ReviewEligibleResp{\n"
            "\t\tEligible: order.Status == model.OrderStatusCompleted && !exists,\n"
            "\t\tReviewed: exists || order.Status == model.OrderStatusReviewed,\n"
            "\t\tStatus:   order.Status,\n"
            "\t\tOrderId:  order.ID,\n"
            "\t}, nil\n"
            "}",
        )
        rev.write_text(rt)
        print("biz ReviewEligible typed")

    # eligible logic
    p = ORD / "internal/logic/user/review/eligible_logic.go"
    text = p.read_text()
    text = text.replace("*types.AnyResp", "*types.ReviewEligibleResp")
    text = text.replace("*types.DataResp", "*types.ReviewEligibleResp")
    text = re.sub(
        r"return &types\.(Any|Data)Resp\{Data: data\}, nil",
        "return data, nil",
        text,
    )
    p.write_text(text)
    print("eligible logic")

    # coupon preview
    p = ORD / "internal/logic/user/order/coupon_preview_logic.go"
    text = p.read_text()
    text = text.replace("*types.AnyResp", "*types.CouponPreviewResp")
    text = text.replace("*types.DataResp", "*types.CouponPreviewResp")
    old = """	data, err := biz.NewOrderLogic(l.svcCtx).CouponPreview(ctx, userID, req.Items, req.UserCouponID)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: data}, nil
}"""
    new = """	data, err := biz.NewOrderLogic(l.svcCtx).CouponPreview(ctx, userID, req.Items, req.UserCouponID)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.CouponPreviewResp{
		GoodsAmount:      data.GoodsAmount,
		DiscountAmount:   data.DiscountAmount,
		PayAmount:        data.PayAmount,
		BestUserCouponID: data.BestUserCouponID,
		Available:        data.Available,
		Unavailable:      data.Unavailable,
	}, nil
}"""
    if old in text:
        text = text.replace(old, new)
    else:
        text = re.sub(
            r"return &types\.(Any|Data)Resp\{Data: data\}, nil",
            """return &types.CouponPreviewResp{
		GoodsAmount:      data.GoodsAmount,
		DiscountAmount:   data.DiscountAmount,
		PayAmount:        data.PayAmount,
		BestUserCouponID: data.BestUserCouponID,
		Available:        data.Available,
		Unavailable:      data.Unavailable,
	}, nil""",
            text,
        )
    p.write_text(text)
    print("coupon preview logic")

    # dedupe URLResp / OrderDetail / List in biz_types if duplicated in types.go
    bt = biz.read_text()
    # keep aliases only if struct defs duplicated - remove struct defs that match types.go
    for name in ("URLResp", "OrderDetailResp", "ListResp", "ReviewEligibleResp", "CouponPreviewResp"):
        # if both define struct, leave biz as source of truth for now — build will catch dups
        pass
    print("order finish done")


def dedupe_order_types() -> None:
    """URLResp/OrderDetail/List/Review/Coupon may exist in both types.go and biz_types.go."""
    biz = ORD / "internal/types/biz_types.go"
    types_go = ORD / "internal/types/types.go"
    bt = biz.read_text()
    tg = types_go.read_text()
    for name in ("URLResp", "OrderDetailResp", "ListResp", "ReviewEligibleResp", "CouponPreviewResp"):
        if f"type {name} struct" in tg and f"type {name} struct" in bt:
            # remove struct from biz, keep alias if needed
            bt = re.sub(
                rf"\ntype {name} struct \{{[^}}]*\}}\n",
                "\n",
                bt,
                count=1,
                flags=re.S,
            )
    # ensure aliases for old names
    if "type OrderDetailData =" not in bt:
        bt += "\ntype OrderDetailData = OrderDetailResp\ntype ListData = ListResp\n"
    biz.write_text(bt)
    print("deduped order biz types")


def main() -> None:
    for svc in SERVICES:
        ensure_data_resp_alias(svc)
        replace_logic_any_to_data(svc)
    finish_order()
    dedupe_order_types()


if __name__ == "__main__":
    main()
