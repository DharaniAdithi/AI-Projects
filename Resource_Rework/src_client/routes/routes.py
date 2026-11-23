from fastapi import APIRouter,Response
import uuid

from src.repository.database import Database
from ..models.models import AgentResponse, RequestPayload
from src_client.services.service import ChatService
from src.utils.error_codes import ApplicationError
from src.utils.error_codes import ErrorCode, HTTP_STATUS_MAP
from src.repository.repository import CafeBotRepository
from src_client.agent_client.agent import AgentManager
from src.utils.helpers import create_api_response,log_error

router = APIRouter(tags=["CafeChatBot"])


def handle_application_error(e: ApplicationError, request_id: str) -> AgentResponse:
   
    return create_api_response(
        code=e.http_status,
        status="error",
        message=e.message,
        error=e.to_dict().get("error"),
        request_id=request_id
    )


@router.post("/chat", response_model=AgentResponse)
async def cafe_bot(request: RequestPayload,response:Response):
    request_id = str(uuid.uuid4())
    database = Database()
    service = ChatService(CafeBotRepository(database), AgentManager())
    try:
        result=await service.cafe_bot(request,request_id)
        response.status_code=result.status_code
        return result
    except Exception as e:
        log_error(file_name="routes.py",function_name="cafe_bot",exception=e)
        error_resp=create_api_response(
            code=HTTP_STATUS_MAP.get(ErrorCode.INTERNAL_SERVER_ERROR),
            status="error",
            message="Internal Server Error",
            error={"type": "UnhandledException", "details": str(e)},
            request_id=request_id
        )
        response.status_code=error_resp.status_code
        return error_resp
