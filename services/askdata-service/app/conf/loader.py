from __future__ import annotations

from functools import lru_cache
from pathlib import Path

from dotenv import load_dotenv
from omegaconf import DictConfig, OmegaConf

# askdata-service/ 根目录
ROOT_DIR = Path(__file__).resolve().parents[2]
CONF_DIR = ROOT_DIR / "conf"

load_dotenv(ROOT_DIR / ".env")


def _load_yaml_dir(conf_dir: Path) -> DictConfig:
    files = sorted(conf_dir.glob("*.yaml")) + sorted(conf_dir.glob("*.yml"))
    if not files:
        return OmegaConf.create({})
    cfgs = [OmegaConf.load(p) for p in files]
    return OmegaConf.merge(*cfgs)


@lru_cache
def get_config() -> DictConfig:
    """合并 conf/*.yaml，并解析 ${oc.env:...}。"""
    cfg = _load_yaml_dir(CONF_DIR)
    OmegaConf.resolve(cfg)
    return cfg


def get_settings() -> DictConfig:
    """兼容旧调用名：返回完整配置树。"""
    return get_config()
