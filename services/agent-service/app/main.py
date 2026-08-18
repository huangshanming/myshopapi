from fastapi import FastAPI

from app.api.routes import agents, health


def create_app() -> FastAPI:
    app = FastAPI(
        title="mymall agent-service",
        version="0.1.0",
        description="Python AI agents (shopping guide, etc.)",
    )
    app.include_router(health.router)
    app.include_router(agents.router, prefix="/api/v1/agents")
    return app


app = create_app()
