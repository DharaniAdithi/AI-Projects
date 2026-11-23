import logging
from langchain_core.messages import HumanMessage
from langchain.agents import create_agent
from .system_prompts import System_prompt
from ..utils.helper import get_llm,log_error

logger = logging.getLogger(__name__)

class CodeAgent:
    
    # SQ: 1.47 to SQ: 1.51, description: Initializes the CodeAgent with LLM, system prompt
    def __init__(self):
        print("code init")
        self.llm = get_llm(max_tokens=500, temperature=0)
        self.prompt = System_prompt()

    # SQ: 1.103 to SQ: 1.109, description: Builds the agent executes the query, and returns the final response.
    async def _run_agent(self, question: str,context:str,review_comments:str) -> str:
        try:
            agent = create_agent(
                model=self.llm,
                tools=[],
                system_prompt=self.prompt.code_system_prompt,
            )

            response = await agent.ainvoke(
                {"messages": [HumanMessage(content=f"user query:{question},context and review comments can be empty context:{context} review comments:{review_comments}")]}
            )

            final_content = response["messages"][-1].content
            logger.info(f"Final messages code agent --------: {final_content}")

            return final_content

        except Exception as e:
            log_error("code_agent.py", "_run_agent", e)
            raise e

    # SQ: 1.102 to SQ: 1.111, description:Invokes the agent and returns the generated code response.
    async def ask(self, question: str,context:str,review_comments:str) -> str:
        try:
            logger.info(f"Customer: {question}")

            final = await self._run_agent(question,context,review_comments)

            logger.info(f"CodeBot: {final}")
            return final

        except Exception as e:
            log_error("code_agent.py", "ask", e)
            logger.error(e)
            raise e