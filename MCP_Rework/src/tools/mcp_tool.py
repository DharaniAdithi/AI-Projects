from langchain_core.tools import tool
from typing import Dict, Any
from src.repository.repository import CafeBotRepository
from src.utils.helpers import log_error
from src.utils.logger import get_logger
from src.constants import (
    TOOL_GET_MENU, 
    TOOL_CHECK_ORDER_STATUS, 
    TOOL_ORDER_CREATION, 
    TOOL_ORDER_UPDATION, 
    MSG_MENU_FETCH_ERROR,
    MSG_ORDER_CREATE_ERROR,
    MSG_ORDER_STATUS_ERROR,
    MSG_ORDER_UPDATE_ERROR
)
from src.routes.mcp_route import router

logger = get_logger(__name__)

@router.tool()
def get_menu_items() -> Dict[str, Any]:
    """
    Fetch all available menu items from the database.

    Returns:
        Dict[str, Any]: Menu items list or error response.
    """
    try:
        repo = CafeBotRepository()
        items = repo.get_menu_items()
        logger.info("Fetched %d menu items.", len(items))
        return {"items": items}
    except Exception as e:
        logger.exception("get_menu_items failed: %s", e)
        log_error(file_name="tools.py", function_name=TOOL_GET_MENU, exception=e)
        return {"status": "error", "message": MSG_MENU_FETCH_ERROR}


@router.tool()
def order_creation(
    customer_name: str, phone_no: str, item_name: str, quantity: int
) -> Dict[str, Any]:
    """
    Create a new order for a customer.

    Args:
        customer_name (str): Customer's full name.
        phone_no (str): Customer's phone number.
        item_name (str): Menu item to order.
        quantity (int): Number of items.

    Returns:
        Dict[str, Any]: Created order details or error response.
    """
    try:
        repo = CafeBotRepository()
        result = repo.create_order(customer_name, phone_no, item_name, quantity)
        logger.info("Order created for customer: %s", customer_name)
        return result
    except Exception as e:
        logger.exception("order_creation failed: %s", e)
        log_error(file_name="tools.py", function_name=TOOL_ORDER_CREATION, exception=e)
        return {"status": "error", "message": MSG_ORDER_CREATE_ERROR}


@router.tool()
def check_order_status(order_id: int) -> Dict[str, Any]:
    """
    Check the current status of an order.

    Args:
        order_id (int): The unique order identifier.

    Returns:
        Dict[str, Any]: Order status or error response.
    """
    try:
        repo = CafeBotRepository()
        result = repo.get_order_status(order_id)
        logger.info("Checked status for order_id: %s", order_id)
        return result
    except Exception as e:
        logger.exception("check_order_status failed: %s", e)
        log_error(file_name="tools.py", function_name=TOOL_CHECK_ORDER_STATUS, exception=e)
        return {"status": "error", "message": MSG_ORDER_STATUS_ERROR}


@router.tool()
def order_updation(order_id: int, quantity: int) -> Dict[str, Any]:
    """
    Update the quantity of an existing order.

    Args:
        order_id (int): The unique order identifier.
        quantity (int): New quantity.

    Returns:
        Dict[str, Any]: Updated order summary or error response.
    """
    try:
        repo = CafeBotRepository()
        result = repo.update_order(order_id, quantity)
        logger.info("Order updated: order_id=%s quantity=%s", order_id, quantity)
        return result
    except Exception as e:
        logger.exception("order_updation failed: %s", e)
        log_error(file_name="tools.py", function_name=TOOL_ORDER_UPDATION, exception=e)
        return {"status": "error", "message": MSG_ORDER_UPDATE_ERROR}