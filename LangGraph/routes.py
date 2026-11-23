from fastapi import APIRouter, Response
from typing import Any
import uuid
from ..models.models import  PipelineResponse, PipelineRequest
from ..services.service import CodePipelineService
from ..utils.helper import error_response,log_error
from ..agent.code_agent import CodeAgent
from ..agent.review_agent import ReviewAgent
from ..agent.research_agent import ResearchAgent
router = APIRouter(prefix="/api/v1/agentic", tags=["multi-agent pipeline"])

# SQ: 1.38 to SQ: 1.181, description: Handles incoming API requests, triggers the pipeline service, and returns the response.
@router.post("/pipeline/code", response_model=PipelineResponse, status_code=200,response_model_exclude_none=True)
async def code_pipeline(
        payload: PipelineRequest,         
        response: Response,           
    ) -> Any:
        request_id = None
        try:
            request_id=str(uuid.uuid4())
            print("route")
            service = CodePipelineService(CodeAgent,ReviewAgent,ResearchAgent,request_id)
            result = await service.code_pipeline(payload.user_query,payload.user_id)
            response.status_code=result.status_code
            return result
        except Exception as e:
            print("________________error_________________")
            response.status_code = 500
            log_error(file_name="route.py",function_name="code_pipeline",exception=e)
            return error_response(
                status_code=500,
                message=str(e),
                error=str(e),
                request_id=request_id,
            )