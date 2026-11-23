from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
import uvicorn
import logging
import asyncio
import selectors
from .settings.settings import settings
from .routes.routes import router
from .repository.database.database_connection import Database
from .migration.migration import Migration


app=None
logger = logging.getLogger(__name__)

#SQ 1.1 TO 1.9 description:Initializes application settings, database connection, and logging configuration.
def initialize_config():
    global app_config, db_connection
    try:
        app_config = settings
        db_connection = Database()
        log_level = getattr(logging, app_config.log_level.upper(), logging.INFO)
        logging.getLogger().setLevel(log_level)
        logger.info("Configuration initialized successfully")
        return True
    except Exception as e:
        logger.error(f"Failed to initialize configuration: {e}")
        return False
    
#SQ 1.24 to SQ 1.31 description:Creates and configures the FastAPI application with CORS middleware and API routes
def create_app():
    app = FastAPI(
        title=" Multi-Agent Code Pipeline API",
        description="API for an agentic workflow that automates code generation, review,and contextual research using multiple AI agents.",
        version="1.0.0",
    )

    app.add_middleware(
        CORSMiddleware,
        allow_origins=["*"],
        allow_credentials=False,
        allow_methods=["GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"],
        allow_headers=[
            "Accept",
            "Accept-Language",
            "Content-Language",
            "Content-Type",
            "Authorization",
            "accessToken",
            "deviceIdentifier",
            "X-RequestID",
            "X-Requested-With",
            "Origin",
            "Cache-Control",
            "Pragma"
        ],
        expose_headers=["*"],
        max_age=3600,
    )
    app.include_router(router)

    return app
#SQ 1.1 to SQ 1.35 Initializes configuration, validates database connectivity, runs migrations, and starts the application server. 
def run_local_server():
    logger.info("Starting in local development mode...")

    if not initialize_config():
        logger.error("Failed to initialize configuration")
        exit(1)

    if not db_connection.test_connection():
        logger.error("Database connection test failed")
        exit(1)
    try:
        migration = Migration()
        migration.create_tables()
        logger.info("Database tables created successfully")
    except Exception as e:
        logger.warning(f"Migration failed: {e}")
    port = int(app_config.port)
    host = app_config.host

    logger.info(f"Starting server on {host}:{port}")
    logger.info(f"Database: {app_config.db_host}:{app_config.db_port}/{app_config.db_name}")

    app = create_app()
    
    #SQ 1.33 to SQ 1.35 description:Starts the Uvicorn server asynchronously using a selector-based event loop.
    async def serve():
        config = uvicorn.Config(
            app,
            host=host,
            port=port,
            reload=False,
            log_level=app_config.log_level.lower(),
            access_log=True,
        )
        server = uvicorn.Server(config)
        await server.serve()

    asyncio.run(
        serve(),
        loop_factory=lambda: asyncio.SelectorEventLoop(selectors.SelectSelector())
    )

if __name__ == "__main__":
    try:
        run_local_server()
    except KeyboardInterrupt:
        logger.info("Server stopped by user")
    except Exception as e:
        logger.error(f"Server startup failed: {e}")
        exit(1)