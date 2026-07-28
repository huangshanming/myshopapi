from fastapi import FastAPI

from app.api.routes import health, recommend, track


def create_app() -> FastAPI:
    app = FastAPI(
        title="mymall recommend-service",
        version="0.1.0",
        description="Collaborative recommend (ItemCF) — skeleton",
    )
    app.include_router(health.router)
    app.include_router(recommend.router, prefix="/api/v1/recommend")
    app.include_router(track.router, prefix="/api/v1/recommend")
    return app


app = create_app()
