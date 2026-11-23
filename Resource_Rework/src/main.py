from contextlib import asynccontextmanager
from typing import AsyncGenerator
import logging
import sys

from fastmcp import FastMCP

from src.settings import config, get_settings
from src.utils.logger import get_logger, _SingletonLogger
from src.repository.database import Database
from src.routes.mcp_route import router

import src.tools.mcp_tool
import src.resources.mcp_resources
import src.prompts.mcp_prompt 


logger = get_logger(__name__)

db_connection = None
app_config = None


@asynccontextmanager
async def lifespan(app: FastMCP) -> AsyncGenerator[None, None]:
    """
    Handles the startup and shutdown events of the FastMCP application.
    """
    settings = get_settings()
    _SingletonLogger.configure(level=settings.log_level)

    logger.info("Application startup complete.")
    yield
    logger.info("Application shutdown complete.")


def initialize_config() -> bool:
    """
    Initializes global application configuration and database connection.
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
    """
    try:
        from src.migrations.create_tables import Migration
        from src.migrations.seeders import Seeder

        logger.info("Running database migrations...")
        migration = Migration()
        migration.create_tables()

        logger.info("Running database seeder...")
        seeder = Seeder()
        seeder.seed_data()

        logger.info("Database setup complete")

    except Exception as e:
        logger.warning(f"Migration or seeding failed: {e}")


def start_mcp_server():
    """
    Configures and launches the FastMCP server.
    """
    mcp = FastMCP("cafebot-mcp-server", lifespan=lifespan)

    mcp.mount(router)

    logger.info("Starting CafeBot MCP Server")
    logger.info("MCP Endpoint: http://localhost:8000/mcp")

    mcp.run(transport="streamable-http")


if __name__ == "__main__":
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