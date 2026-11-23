from enum import Enum
from typing import Optional, Dict, Any


class ErrorCode(Enum):
    """Enumeration of application-level error codes."""

    DATABASE_ERROR = "ZAP_SQL_005"
    VALIDATION_ERROR = "ZAP_VAL_022"
    NOT_FOUND = "ZAP_GEN_030"
    UNAUTHORIZED = "ZAP_AUTH_018"
    FORBIDDEN = "ZAP_AUTH_026"
    INTERNAL_SERVER_ERROR = "ZAP_SRV_006"


HTTP_STATUS_MAP: Dict[ErrorCode, int] = {
    ErrorCode.DATABASE_ERROR: 500,
    ErrorCode.VALIDATION_ERROR: 400,
    ErrorCode.NOT_FOUND: 404,
    ErrorCode.UNAUTHORIZED: 401,
    ErrorCode.FORBIDDEN: 403,
    ErrorCode.INTERNAL_SERVER_ERROR: 500,
}


class ApplicationError(Exception):
    """
    Structured application exception with error code and HTTP status.

    Attributes:
        code (str): Application error code string.
        http_status (int): Corresponding HTTP status code.
        message (str): Human-readable error message.
        details (Dict[str, Any]): Additional error context.
    """

    def __init__(
        self,
        error_code: ErrorCode,
        message: str,
        details: Optional[Dict[str, Any]] = None,
    ) -> None:
        """
        Initialize ApplicationError.

        Args:
            error_code (ErrorCode): The error code enum value.
            message (str): Human-readable error message.
            details (Optional[Dict[str, Any]]): Extra context for the error.
        """
        self.code: str = error_code.value
        self.http_status: int = HTTP_STATUS_MAP.get(error_code, 500)
        self.message: str = message
        self.details: Dict[str, Any] = details or {}
        super().__init__(self.message)

    def to_dict(self) -> Dict[str, Any]:
        """
        Serialize error to dictionary.

        Returns:
            Dict[str, Any]: Error representation.
        """
        return {
            "error": {
                "code": self.code,
                "message": self.message,
                "details": self.details,
            }
        }
