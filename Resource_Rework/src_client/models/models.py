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
