"""MySQL 异步客户端占位（SQLAlchemy / asyncmy）。"""

from app.conf import get_config


def get_mysql_dsn() -> str:
    m = get_config().mysql
    return (
        f"mysql+asyncmy://{m.username}:{m.password}@{m.host}:{m.port}/{m.database}"
        f"?charset={m.get('charset', 'utf8mb4')}"
    )
