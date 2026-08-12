from fastapi import FastAPI

from app.api.routes import ask, health
from app.core.lifespan import lifespan


def create_app() -> FastAPI:
    app = FastAPI(
        title="mymall askdata-service（电商问数）",
        version="0.1.0",
        description="自然语言查数 Agent：订单/GMV/库存等指标问答",
        lifespan=lifespan,
    )
    app.include_router(health.router)
    app.include_router(ask.router, prefix="/api/v1/askdata")
    return app


app = create_app()
