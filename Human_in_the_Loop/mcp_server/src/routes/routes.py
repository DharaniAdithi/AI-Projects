from fastapi import APIRouter, Depends
from fastapi.responses import JSONResponse
from models.models import ChatRequestDTO
from services.chat_service import ChatService
from utils.helpers import log_error
import uuid
import datetime
from utils.logger import logger

router = APIRouter(prefix="/api/v1", tags=["Chat"])

# def get_chat_service() -> ChatService:
#     try:
#         return ChatService()
#     except Exception as e:
#         raise Exception(f"Failed to initialize ChatService: {str(e)}")

#SQ 1.4 - SQ 1.7 :Receive POST /api/v1/chat request from client, validate request payload using Pydantic DTOs, and prepare request for service layer.
@router.post("/chat")
async def chat(
    request: ChatRequestDTO,
    service: ChatService = Depends(ChatService)
    ):
    try:
        logger.info("Received chat request", {"customer_id": request.customer_id})
        result = await service.chat_with_agent(message=request.query, thread_id=request.customer_id)

        #SQ 1.66 - SQ 1.68 : build the json response and forward the final response to client
        response = {
            "code": 200,
            "status": "success",
            "message": "Chatbot responded successfully",
            "data": [result["data"][0]],
            "error": None,
            "request_id": result.get("request_id"),
            "timestamp": result.get("timestamp")
        }

        logger.info("Chat response sent successfully",{"request_id": result.get("request_id")})
        return JSONResponse(status_code=200, content=response)
    
    #SQ 1.69 - SQ 1.72: log the unhandled exception, build the error json response and forward to the client
    except Exception as e:
        logger.error("Error occurred while processing chat request", e)
        log_error(file_name="routes.py", function_name="chat", exception=e)
        error_response = {
            "code": 500,
            "status": "error",
            "message": "Chatbot failed to respond",
            "data": [],
            "error": str(e),
            "request_id": str(uuid.uuid4()),
            "timestamp": datetime.datetime.now().isoformat()
        }
        return JSONResponse(status_code=500, content=error_response)