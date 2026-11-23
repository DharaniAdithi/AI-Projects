from ..models.models import AgentResponse
from src.utils.helpers import create_api_response,log_error
from src.repository.repository import CafeBotRepository

from dotenv import load_dotenv
from src.utils.error_codes import ErrorCode, HTTP_STATUS_MAP
from src_client.agent_client.agent import AgentManager



load_dotenv()



class ChatService:
    def __init__(self, repo : CafeBotRepository, agent : AgentManager) -> None:
        self.repo = repo
        self.agent = agent
    async def cafe_bot(self,request,request_id)-> AgentResponse:
        try:
            user_query=request.user_query
            result=await self.agent.ask(user_query)
            return create_api_response(
                code=200,
                status="success",
                message="Chatbot responded successfully",
                data={"CafeBot:...":result},
                error=None,
                request_id=request_id
            )
        except Exception as e:
             log_error(file_name="services.py",function_name="cafe_bot",exception=e)
             return create_api_response(
                code=HTTP_STATUS_MAP.get(ErrorCode.INTERNAL_SERVER_ERROR),
                status="error",
                message="Internal Server Error",
                data=None,
                error={"type": "ServiceException", "details": str(e)},
                request_id=request_id
            )


 
 