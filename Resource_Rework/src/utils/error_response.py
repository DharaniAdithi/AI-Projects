from typing import Optional, Dict, Any
from src.models.models import AgentResponse
from datetime import datetime


def build_error_response(
    code: int,
    message: str,
    error: Optional[Dict[str, Any]],
    request_id: str,
) -> AgentResponse:
    """
    Build a standardized error AgentResponse.

    Args:
        code (int): HTTP status code.
        message (str): Human-readable message.
        error (Optional[Dict[str, Any]]): Error detail dict.
        request_id (str): Unique request identifier.

    Returns:
        AgentResponse: Populated error response model.
    """
    return AgentResponse(
        status_code=code,
        status="error",
        message=message,
        data=None,
        error=error,
        request_id=request_id,
        timestamp=datetime.utcnow().isoformat(),
    )