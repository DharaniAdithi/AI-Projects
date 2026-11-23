import traceback
from datetime import datetime
from typing import Optional, Dict, Any
from src.models.models import AgentResponse
from src.repository.database import Database
from src.utils.logger import get_logger

logger = get_logger(__name__)


def create_api_response(
    code: int,
    status: str,
    message: str,
    data: Optional[Dict] = None,
    error: Optional[Dict] = None,
    request_id: Optional[str] = None,
) -> AgentResponse:
    """
    Build a standardized AgentResponse.

    Args:
        code (int): HTTP status code.
        status (str): 'success' or 'error'.
        message (str): Human-readable message.
        data (Optional[Dict]): Response payload.
        error (Optional[Dict]): Error detail.
        request_id (Optional[str]): Unique request ID.

    Returns:
        AgentResponse: Populated response model.
    """
    return AgentResponse(
        status_code=code,
        status=status,
        message=message,
        data=data,
        error=error,
        request_id=request_id,
        timestamp=datetime.utcnow().isoformat(),
    )


def log_error(file_name: str, function_name: str, exception: Exception) -> None:
    """
    Persist error log to the database.

    Args:
        file_name (str): Source file where the error occurred.
        function_name (str): Function where the error occurred.
        exception (Exception): The raised exception.
    """
    from src.repository.repository import ErrorLogRepository
    try:
        repo = ErrorLogRepository(Database())
        error_message = str(exception)
        stack_trace = traceback.format_exc()
        full_error = f"{error_message}|TRACE: {stack_trace}"
        repo.store_error(file_name=file_name, function_name=function_name, error=full_error)
    except Exception as repo_exc:
        logger.exception("Failed to store error via ErrorLogRepository: %s", repo_exc)