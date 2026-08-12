"""配置加载入口：YAML（conf/）+ 环境变量占位。"""

from app.conf.loader import CONF_DIR, ROOT_DIR, get_config, get_settings

__all__ = ["CONF_DIR", "ROOT_DIR", "get_config", "get_settings"]
