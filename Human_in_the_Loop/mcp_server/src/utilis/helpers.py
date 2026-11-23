from langchain.agents.structured_output import StructuredOutputValidationError
from langchain.agents.structured_output import MultipleStructuredOutputsError
import traceback
import logging

logger = logging.getLogger(__name__)


def log_error(file_name: str, function_name: str, exception: Exception):
    try:
        error_message = str(exception)
        stack_trace = traceback.format_exc()
        full_error = f"{error_message} | TRACE: {stack_trace}"
        logger.error(f"[{file_name}] [{function_name}] {full_error}")
        print(f"ERROR [{file_name}] [{function_name}] {error_message}")
    except Exception as e:
        logger.exception("Failed to log error: %s", e)


def custom_error_handler(error: Exception) -> str:
    if isinstance(error, StructuredOutputValidationError):
        return "There was an issue with the format. Try again."
    elif isinstance(error, MultipleStructuredOutputsError):
        return "Multiple structured outputs were returned. Pick the most relevant one."
    else:
        return f"Error: {str(error) }"