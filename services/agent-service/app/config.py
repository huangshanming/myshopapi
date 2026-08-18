from functools import lru_cache

from pydantic import Field
from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    model_config = SettingsConfigDict(
        env_file=".env",
        env_file_encoding="utf-8",
        extra="ignore",
        populate_by_name=True,
    )

    agent_http_port: int = Field(default=8886, validation_alias="MYMALL_AGENT_HTTP_PORT")
    agent_mode: str = Field(default="dev", validation_alias="MYMALL_AGENT_MODE")

    user_http: str = Field(default="http://127.0.0.1:8881", validation_alias="MYMALL_USER_HTTP")
    catalog_http: str = Field(default="http://127.0.0.1:8882", validation_alias="MYMALL_CATALOG_HTTP")
    order_http: str = Field(default="http://127.0.0.1:8883", validation_alias="MYMALL_ORDER_HTTP")
    merchant_http: str = Field(default="http://127.0.0.1:8884", validation_alias="MYMALL_MERCHANT_HTTP")

    openai_api_key: str = Field(default="", validation_alias="OPENAI_API_KEY")
    openai_base_url: str = Field(default="https://api.openai.com/v1", validation_alias="OPENAI_BASE_URL")
    agent_model: str = Field(default="gpt-4o-mini", validation_alias="MYMALL_AGENT_MODEL")

    @property
    def port(self) -> int:
        return self.agent_http_port

    @property
    def mode(self) -> str:
        return self.agent_mode


@lru_cache
def get_settings() -> Settings:
    return Settings()
