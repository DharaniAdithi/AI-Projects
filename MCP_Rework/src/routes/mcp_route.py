from fastmcp import FastMCP
from contextlib import asynccontextmanager
from typing import AsyncGenerator

@asynccontextmanager
async def lifespan(app):
    from src.utils.logger import get_logger, _SingletonLogger
    from src.settings import get_settings
    settings = get_settings()
    _SingletonLogger.configure(level=settings.log_level)
    yield

router = FastMCP("cafebot-router", lifespan=lifespan)

import src.tools.mcp_tool  # keep — forces @router.tool() decorators to run