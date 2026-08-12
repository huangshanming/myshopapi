from __future__ import annotations

from loguru import logger
from omegaconf import DictConfig

from app.conf import ROOT_DIR, get_config


def setup_logging(cfg: DictConfig | None = None) -> None:
    """初始化 loguru：控制台 + logs/ 文件。"""
    cfg = cfg or get_config()
    log_cfg = cfg.get("logging", {})
    level = str(log_cfg.get("level", "INFO"))
    log_dir = ROOT_DIR / str(log_cfg.get("dir", "logs"))
    log_dir.mkdir(parents=True, exist_ok=True)

    logger.remove()
    logger.add(
        lambda msg: print(msg, end=""),
        level=level,
        colorize=True,
        format="<green>{time:HH:mm:ss}</green> | <level>{level}</level> | {message}",
    )
    logger.add(
        log_dir / "askdata.log",
        level=level,
        rotation=str(log_cfg.get("rotation", "10 MB")),
        retention=str(log_cfg.get("retention", "7 days")),
        encoding="utf-8",
    )
