# -------------------------------------------------
# Tool Names
# -------------------------------------------------
TOOL_GET_MENU = "get_menu_items"
TOOL_ORDER_CREATION = "order_creation"
TOOL_CHECK_ORDER_STATUS = "check_order_status"
TOOL_ORDER_UPDATION = "order_updation"

# -------------------------------------------------
# General Messages
# -------------------------------------------------
MSG_AGENT_ERROR = "Sorry! Something went wrong while processing your request. Please try again."
MSG_ORDER_NOT_FOUND = "Order not found"
MSG_MENU_FETCH_ERROR = "Unable to fetch menu items right now."
MSG_ORDER_CREATE_ERROR = "Order could not be created."
MSG_ORDER_STATUS_ERROR = "Unable to check order status."
MSG_ORDER_UPDATE_ERROR = "Unable to update order."

# -------------------------------------------------
# Config Keys
# -------------------------------------------------
DEFAULT_HOST = "0.0.0.0"
DEFAULT_PORT = 8001
DEFAULT_REGION = "us-east-1"
DEFAULT_MAX_TOKENS = 500
DEFAULT_TEMPERATURE = 0
DEFAULT_LOG_LEVEL = "info"
MODEL_ID = "arn:aws:bedrock:us-east-1:451433485314:application-inference-profile/eyhgqon7j4gk"
PROVIDER = "amazon"
MCP_SERVER_URL = "http://localhost:8000/mcp"
API_PREFIX = "/cafebot/api/v1"