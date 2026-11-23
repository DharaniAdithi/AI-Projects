from fastapi import Request
from fastapi.exceptions import RequestValidationError
from fastapi.responses import JSONResponse
import uuid
from datetime import datetime
from src.utils.logger import get_logger

logger = get_logger(__name__)


async def validation_exception_handler(
    request: Request, exc: RequestValidationError
) -> JSONResponse:
    """
    Handle Pydantic/FastAPI request validation errors globally.

    Args:
        request (Request): Incoming FastAPI request.
        exc (RequestValidationError): Raised validation exception.

    Returns:
        JSONResponse: Standardized 422 error response.
    """
    error_details = []
    for error in exc.errors():
        field_path = " -> ".join(str(loc) for loc in error.get("loc", []))
        msg = error.get("msg", "")

        if "missing" in msg:
            message = f"Field '{field_path}' is required."
        elif "empty" in msg or "whitespace" in msg:
            message = f"Field '{field_path}' cannot be empty or whitespace only."
        else:
            message = f"{field_path}: {msg}"

        error_details.append({"code": "ZAP_VAL_022", "message": message})

    logger.warning("Validation error on request: %s", request.url)

    return JSONResponse(
        status_code=422,
        content={
            "status": "failure",
            "code": 422,
            "message": "Something went wrong",
            "request_id": str(uuid.uuid4()),
            "timestamp": datetime.utcnow().isoformat(),
            "errors": error_details,
        },
    )