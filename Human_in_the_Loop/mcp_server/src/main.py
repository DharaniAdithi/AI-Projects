from fastapi import FastAPI
from contextlib import asynccontextmanager
from dotenv import load_dotenv
import uvicorn
from routes.routes import router as router
from utils.logger import logger

load_dotenv()

# # Lifespan Event (Startup)
# @asynccontextmanager
# async def lifespan(app: FastAPI):

#     create_tables()
#     print("Database tables created successfully")
#     yield

#SQ 1.0 - SQ 1.3 : Initialize FastAPI application, configure metadata, and register API routes under /api/v1.

# FastAPI App
app = FastAPI(
    title="Brew Buddy AI Coffee Assistant",
    version="1.0.0"
)


# Register Routes
app.include_router(router)
logger.info("Routes registered successfully")



if __name__ == "__main__":
    """
    Main entry point for running the FastAPI application locally.
    """
    logger.info("Starting cafebot...")
    uvicorn.run(
        "main:app",
        host="0.0.0.0",
        port=8080,
        reload=True,
        log_level="info"
    )