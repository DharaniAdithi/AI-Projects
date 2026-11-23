from contextlib import asynccontextmanager
from typing import AsyncGenerator
from fastapi import FastAPI, Request
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse
import uvicorn
import uuid
from datetime import datetime
from fastapi.exceptions import RequestValidationError
from src.settings import get_settings
from src.utils.logger import get_logger, _SingletonLogger
from src_client.routes.routes import router
from src.constants import API_PREFIX, DEFAULT_HOST, DEFAULT_PORT
import logging

logger = get_logger(__name__)



@asynccontextmanager
async def lifespan(app: FastAPI) -> AsyncGenerator[None, None]:
    """
    FastAPI async lifespan — logger initialization only.

    Args:
        app (FastAPI): The FastAPI application instance.

    Yields:
        None
    """
    settings = get_settings()
    _SingletonLogger.configure(level=settings.log_level)
    logger.info("Application startup complete.")
    yield
    logger.info("Application shutdown complete.")


logger = logging.getLogger(__name__)


def create_app():
   

    app = FastAPI(
        title="CafeBot Client API",
        description="FastAPI client that connects to MCP Server",
        version="1.0.0"
    )

    app.add_middleware(
        CORSMiddleware,
        allow_origins=["*"],
        allow_methods=["*"],
        allow_headers=["*"],
    )
    app.include_router(prefix= API_PREFIX , router=router)

    return app


app = create_app()

@app.exception_handler(RequestValidationError)
async def validation_exception_handler(request: Request, exc: RequestValidationError):
    request_id = str(uuid.uuid4())

    error_details = [
        {
            "field": ".".join(map(str, err["loc"])),
            "message": err["msg"]
        }
        for err in exc.errors()
    ]

    return JSONResponse(
        status_code=422,
        content={
            "status_code": 422,
            "status": "error",
            "message": "Validation Error",
            "error": {
                "type": "ValidationError",
                "details": error_details
            },
            "request_id": request_id,
            "timestamp": datetime.utcnow().isoformat()
        }
    )

@app.exception_handler(Exception)
async def generic_exception_handler(request: Request, exc: Exception) -> JSONResponse:
    """
    Global catch-all for unhandled exceptions within the MCP server.

    Args:
        request (Request): The incoming HTTP request.
        exc (Exception): The error instance raised.

    Returns:
        JSONResponse: A standardized 500 Internal Server Error response.
    """
    logger.error(f"Unhandled exception on {request.url}: {exc}")
    return JSONResponse(
        status_code=500,
        content={
            "status_code": 500,
            "status": "error",
            "message": "Internal Server Error",
            "error": str(exc),
            "request_id": str(uuid.uuid4()),
            "timestamp": datetime.utcnow().isoformat()
        },
    )


if __name__ == "__main__":
    try: 
        logger.info("\n Starting CafeBot Client API")
        logger.info("API Docs: http://localhost:8001/docs")
        logger.info("Chat Endpoint: http://localhost:8001/cafebot/api/v1/chat\n")

        uvicorn.run(
            "src_client.main:app",
            host=DEFAULT_HOST,
            port=DEFAULT_PORT,
            reload=True
        )
    except Exception as exec:
        logger.error(exec.message)