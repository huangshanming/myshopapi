# Redis key helpers for recommend (fill in during ItemCF wiring).

I2I_PREFIX = "reco:i2i:"
U2I_PREFIX = "reco:u2i:"
HOT_GLOBAL = "reco:hot:global"
HOT_CATE_PREFIX = "reco:hot:cate:"
BOUGHT_FILTER_PREFIX = "reco:filter:bought:"
CACHE_PREFIX = "reco:cache:req:"
META_VERSION = "reco:meta:version"


def i2i_key(item_id: int) -> str:
    return f"{I2I_PREFIX}{item_id}"


def u2i_key(user_id: int) -> str:
    return f"{U2I_PREFIX}{user_id}"


def hot_cate_key(cate_id: int) -> str:
    return f"{HOT_CATE_PREFIX}{cate_id}"


def bought_filter_key(user_id: int) -> str:
    return f"{BOUGHT_FILTER_PREFIX}{user_id}"
