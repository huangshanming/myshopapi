"""FastAPI 依赖注入（后续挂 DB session / 鉴权等）。"""

from app.conf import get_config


def get_app_config():
    return get_config()
