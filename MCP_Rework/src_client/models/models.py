from pydantic import BaseModel
from typing import Optional, Dict


class AgentResponse(BaseModel):
    """
    Standardized API response model.

    Attributes:
        status_code (int): HTTP status code.
        status (str): 'success' or 'error'.
        message (str): Human-readable message.
        data (Optional[Dict]): Response payload.
        error (Optional[Dict]): Error details.
        request_id (str): Unique request identifier.
        timestamp (str): ISO 8601 timestamp.
    """

    status_code: int
    status: str
    message: str
    data: Optional[Dict] = None
    error: Optional[Dict] = None
    request_id: str
    timestamp: str


class RequestPayload(BaseModel):
    """
    Incoming chat request payload.

    Attributes:
        user_query (str): The user's message.
    """

    user_query: str


# class OrderDetails(BaseModel):
#     """
#     Structured output model enforced by ToolStrategy.

#     Attributes:
#         item_name (str): Ordered menu item name.
#         quantity (int): Number of items.
#         status (str): Current order status.
#         order_id (Optional[int]): Order ID after creation.
#         total_price (Optional[int]): Total cost in INR.
#         message (Optional[str]): Friendly response to customer.
#     """

#     item_name: str = ""
#     quantity: int = 0
#     status: str = ""
#     order_id: Optional[int] = None
#     total_price: Optional[int] = None
#     message: Optional[str] = None
