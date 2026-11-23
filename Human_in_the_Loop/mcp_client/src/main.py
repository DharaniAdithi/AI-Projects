from fastmcp import FastMCP
from contextlib import asynccontextmanager
from mcp_client.src.routes.routes import router
from mcp_client.src.migrations.create_tables import create_tables

@asynccontextmanager
async def lifespan(context):
    create_tables()
    print("Database tables created successfully")
    yield

mcp = FastMCP("chatbot-mcp-server", lifespan=lifespan)
mcp.mount(router)

if __name__ == "__main__":
    mcp.run(transport="streamable-http", host="0.0.0.0", port=8000)