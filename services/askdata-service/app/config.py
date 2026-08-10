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

    askdata_http_port: int = Field(
        default=8889, validation_alias="MYMALL_ASKDATA_HTTP_PORT"
    )
    askdata_mode: str = Field(default="dev", validation_alias="MYMALL_ASKDATA_MODE")

    user_http: str = Field(
        default="http://127.0.0.1:8881", validation_alias="MYMALL_USER_HTTP"
    )
    catalog_http: str = Field(
        default="http://127.0.0.1:8882", validation_alias="MYMALL_CATALOG_HTTP"
    )
    order_http: str = Field(
        default="http://127.0.0.1:8883", validation_alias="MYMALL_ORDER_HTTP"
    )
    merchant_http: str = Field(
        default="http://127.0.0.1:8884", validation_alias="MYMALL_MERCHANT_HTTP"
    )

    # LLM（后续接 LangGraph / Ollama / OpenAI 时使用）
    openai_api_key: str = Field(default="", validation_alias="OPENAI_API_KEY")
    openai_base_url: str = Field(
        default="https://api.openai.com/v1", validation_alias="OPENAI_BASE_URL"
    )
    askdata_model: str = Field(
        default="qwen2.5:7b", validation_alias="MYMALL_ASKDATA_MODEL"
    )
    ollama_base_url: str = Field(
        default="http://127.0.0.1:11434", validation_alias="MYMALL_OLLAMA_BASE_URL"
    )

    @property
    def port(self) -> int:
        return self.askdata_http_port

    @property
    def mode(self) -> str:
        return self.askdata_mode


@lru_cache
def get_settings() -> Settings:
    return Settings()
