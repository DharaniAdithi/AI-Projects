import logging
from pathlib import Path
from langchain.agents import create_agent
from langchain_core.messages import HumanMessage
from langchain_mcp_adapters.client import MultiServerMCPClient
from langchain_mcp_adapters.tools import load_mcp_tools
from .system_prompts import System_prompt
from ..utils.helper import get_llm,log_error
from ..settings.settings import settings
from..constant.constant import PYLINT_WORKSPACE 
logger = logging.getLogger(__name__)
 
 
class ReviewAgent:
    
    # SQ: 1.54 to SQ: 1.56, description: Initializes the ResearchAgent with LLM, review prompt and pylint workspace directory
    def __init__(self):
        logger.info("review init")
        self.llm = get_llm(max_tokens=4000, temperature=0)
        self.prompt = System_prompt()
        self.pylint_workspace=PYLINT_WORKSPACE
    
    # SQ: 1.123 to SQ: 1.129, description: Creates and runs the review agent with tools, returning the final reviewed output.
    async def _run_agent(self, file_path_str: str, tools) -> str:
        try:
            agent = create_agent(
                model=self.llm,
                tools=tools,
                system_prompt=self.prompt.review_system_prompt,
    )
 
            response = await agent.ainvoke(
                {
                    "messages": [
                        HumanMessage(content=(
                                f"Analyse this Python file: {file_path_str}\n"
                                "Use the available pylint tools to review it.\n"
                                "DO NOT rewrite or modify the code.\n"
                                "Return review_status as 'pass' or 'fail' based on score "
                                "if review score is above 6.0 consider as pass"
                                "and review_comment summarizing all issues found ignore  newline issues"
                                
                            )) ]   })
 
            last_message = response["messages"][-1].content
            logger.info(f"Final review response: {last_message}")
            return last_message
 
        except Exception as e:
            logger.error(e)
            log_error("review_agent.py", "_run_agent", e)
            raise e
 
    #SQ: 1.113 to SQ: 1.131, description: Connects to pylint MCP server, loads tools, executes the review flow, and returns the response. 
    async def ask(self, code: str) -> str:
 
        try:
            temp_file = Path(self.pylint_workspace)
            
            lines = code.splitlines()
            if len(lines) > 2:
                trimmed_code = "\n".join(lines[1:-1])
            else:
                trimmed_code = lines
            temp_file.write_text(trimmed_code, encoding="utf-8")
            logger.info(f"Written review file: {temp_file}")
            file_path_str = temp_file.as_posix()
            client = MultiServerMCPClient(
                {
                    "pylint-server": {
                        "transport": "streamable_http",
                        "url":settings.pylint_mcp_server_url,
                    }
                }
            )
 
            async with client.session("pylint-server") as session:
                
                tools = await load_mcp_tools(session)
                logger.info(f"Loaded {len(tools)} tools: {[t.name for t in tools]}")
                response = await self._run_agent(file_path_str, tools)
                return str(response)
 
        except Exception as e:
            log_error("review_agent.py", "ask", e)
            logger.error(e)
            raise e