from functools import lru_cache
from urllib.parse import quote_plus

from pydantic import Field
from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    model_config = SettingsConfigDict(
        env_file=".env",
        env_file_encoding="utf-8",
        extra="ignore",
        populate_by_name=True,
    )

    reco_http_port: int = Field(default=8888, validation_alias="MYMALL_RECO_HTTP_PORT")
    reco_mode: str = Field(default="dev", validation_alias="MYMALL_RECO_MODE")

    catalog_http: str = Field(default="http://127.0.0.1:8882", validation_alias="MYMALL_CATALOG_HTTP")

    # MySQL — local defaults align with lottery/catalog; override via MYMALL_MYSQL_*
    mysql_host: str = Field(default="127.0.0.1", validation_alias="MYMALL_MYSQL_HOST")
    mysql_port: int = Field(default=3306, validation_alias="MYMALL_MYSQL_PORT")
    mysql_username: str = Field(default="homestead", validation_alias="MYMALL_MYSQL_USERNAME")
    mysql_password: str = Field(default="secret", validation_alias="MYMALL_MYSQL_PASSWORD")
    mysql_dbname: str = Field(default="mymall", validation_alias="MYMALL_MYSQL_DBNAME")
    mysql_charset: str = Field(default="utf8mb4", validation_alias="MYMALL_MYSQL_CHARSET")

    redis_host: str = Field(default="127.0.0.1", validation_alias="MYMALL_REDIS_HOST")
    redis_port: int = Field(default=6379, validation_alias="MYMALL_REDIS_PORT")
    redis_password: str = Field(default="", validation_alias="MYMALL_REDIS_PASSWORD")
    redis_db: int = Field(default=0, validation_alias="MYMALL_REDIS_DB")

    rabbitmq_url: str = Field(
        default="amqp://mymall:mymall@127.0.0.1:5672/",
        validation_alias="MYMALL_RABBITMQ_URL",
    )

    # Milvus — local standalone :19530; override via MYMALL_MILVUS_*
    milvus_host: str = Field(default="127.0.0.1", validation_alias="MYMALL_MILVUS_HOST")
    milvus_port: int = Field(default=19530, validation_alias="MYMALL_MILVUS_PORT")
    milvus_user: str = Field(default="", validation_alias="MYMALL_MILVUS_USER")
    milvus_password: str = Field(default="", validation_alias="MYMALL_MILVUS_PASSWORD")
    milvus_db_name: str = Field(default="default", validation_alias="MYMALL_MILVUS_DB")
    # default collection for item embeddings (MF / dual-tower)
    milvus_item_collection: str = Field(
        default="reco_item_emb",
        validation_alias="MYMALL_MILVUS_ITEM_COLLECTION",
    )
    milvus_embedding_dim: int = Field(default=64, validation_alias="MYMALL_MILVUS_EMBEDDING_DIM")

    @property
    def port(self) -> int:
        return self.reco_http_port

    @property
    def mode(self) -> str:
        return self.reco_mode

    @property
    def redis_url(self) -> str:
        auth = f":{self.redis_password}@" if self.redis_password else ""
        return f"redis://{auth}{self.redis_host}:{self.redis_port}/{self.redis_db}"

    @property
    def milvus_uri(self) -> str:
        return f"http://{self.milvus_host}:{self.milvus_port}"

    @property
    def mysql_dsn(self) -> str:
        """SQLAlchemy / generic URL style (password URL-encoded)."""
        user = quote_plus(self.mysql_username)
        password = quote_plus(self.mysql_password)
        return (
            f"mysql+pymysql://{user}:{password}"
            f"@{self.mysql_host}:{self.mysql_port}/{self.mysql_dbname}"
            f"?charset={self.mysql_charset}"
        )

    def mysql_connect_kwargs(self) -> dict:
        """kwargs for pymysql.connect(**settings.mysql_connect_kwargs())."""
        return {
            "host": self.mysql_host,
            "port": self.mysql_port,
            "user": self.mysql_username,
            "password": self.mysql_password,
            "database": self.mysql_dbname,
            "charset": self.mysql_charset,
        }


@lru_cache
def get_settings() -> Settings:
    return Settings()
