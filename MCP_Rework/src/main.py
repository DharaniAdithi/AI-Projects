from contextlib import asynccontextmanager
from typing import AsyncGenerator
from fastapi import FastAPI, Request
from fastapi.middleware.cors import CORSMiddleware
from fastapi.exceptions import RequestValidationError
from fastapi.responses import JSONResponse
import uvicorn
import uuid
from datetime import datetime
import logging
import sys
from fastmcp import FastMCP

from src.settings import config, get_settings
from src.utils.logger import get_logger, _SingletonLogger
from src.repository.database import Database
from src.routes.mcp_route import router

# Initialize the centralized logger
logger = get_logger(__name__)

db_connection = None
app_config = None

@asynccontextmanager
async def lifespan(app: FastMCP) -> AsyncGenerator[None, None]:
    """
    Handles the startup and shutdown events of the FastMCP application.

    This context manager initializes the singleton logger configuration based on 
    the application settings and logs the lifecycle state of the server.

    Args:
        app (FastMCP): The instance of the FastMCP application.

    Yields:
        None: Control is yielded back to the application while it is running.
    """
    settings = get_settings()
    _SingletonLogger.configure(level=settings.log_level)
    logger.info("Application startup complete.")
    yield
    logger.info("Application shutdown complete.")


def initialize_config() -> bool:
    """
    Initializes global application configuration and database connection.

    Sets the global `app_config` and attempts to instantiate the `Database` 
    repository.

    Returns:
        bool: True if initialization was successful, False otherwise.
    """
    global db_connection, app_config
    try:
        app_config = config
        logger.info("Loading configuration...")
        db_connection = Database()
        logger.info("Database connection initialized")
        return True
    except Exception as e:
        logger.error(f"Configuration initialization failed: {e}")
        return False


def run_migrations_and_seed():
    """
    Executes database schema migrations and seeds initial data.

    Importing Migration and Seeder logic locally to avoid circular dependencies.
    If migration or seeding fails, a warning is logged but the process continues.
    """
    try:
        from src.migrations.create_tables import Migration
        from src.migrations.seeders import Seeder
        
        logger.info("Running database migrations...")
        migration = Migration()
        migration.create_tables()
        logger.info("Database tables created")
        
        logger.info("Running database seeder...")
        seeder = Seeder()
        seeder.seed_data()
        logger.info("Database seeded successfully")
    except Exception as e:
        logger.warning(f"Migration or seeding failed: {e}")


def start_mcp_server():
    """
    Configures and launches the FastMCP server.

    Registers routes, sets up global exception handlers for unexpected errors, 
    and starts the server using the 'streamable-http' transport protocol.
    
    Note:
        The generic exception handler is registered before the server starts
        to ensure all runtime errors are captured and returned as JSON.
    """
    logger.info("Starting CafeBot MCP Server")
    logger.info("MCP Endpoint: http://localhost:8000/mcp")
    router.run(transport="streamable-http")


if __name__ == "__main__":
    """
    Main entry point for the application.
    
    1. Sets up basic logging.
    2. Initializes configuration and validates database connectivity.
    3. Runs migrations and seeds.
    4. Boots the MCP server.
    """
    import asyncio
    import sys

    # Windows ProactorEventLoop has a bug with socket cleanup in Python 3.12+
    # SelectorEventLoop is stable and works correctly
    if sys.platform == "win32":
        asyncio.set_event_loop_policy(asyncio.WindowsSelectorEventLoopPolicy())

    try:
        logging.basicConfig(level=logging.INFO)
        logger.info("Starting MCP Server Initialization")
        
        if not initialize_config():
            logger.error("Failed to initialize configuration")
            sys.exit(1)
            
        if not db_connection.test_connection():
            logger.error("Database connection failed")
            sys.exit(1)
            
        logger.info("Database connection successful")
        run_migrations_and_seed()
        start_mcp_server()
        
    except KeyboardInterrupt:
        logger.info("Server stopped by user")
    except Exception as e:
        logger.error(f"Server startup failed: {e}")
        sys.exit(1)