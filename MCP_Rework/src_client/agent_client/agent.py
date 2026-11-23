import re
import boto3
from langchain_aws import ChatBedrock
from langgraph.prebuilt import create_react_agent
from langchain_mcp_adapters.client import MultiServerMCPClient

from src_client.agent_client.system_prompt import SYSTEM_PROMPT
from src.utils.helpers import log_error
from src.constants import DEFAULT_REGION, MODEL_ID, PROVIDER, MCP_CONFIG
from src.utils.logger import get_logger
from langchain_mcp_adapters.tools import load_mcp_tools  # add this import

logger = get_logger(__name__)

class AgentManager:
    def get_bedrock_client(self):
        try:
            return boto3.client(
                service_name="bedrock-runtime",
                region_name=DEFAULT_REGION,
            )
        except Exception as e:
            log_error("agent.py", "get_bedrock_client", e)
            raise



    def get_llm(self, max_tokens=500, temperature=0.7):
        try:
            return ChatBedrock(
                client=self.get_bedrock_client(),
                model_id=MODEL_ID,
                provider=PROVIDER,
                model_kwargs={
                    "max_tokens": max_tokens,
                    "temperature": temperature,
                },
            )
        except Exception as e:
            log_error("agent.py", "get_llm", e)
            raise


    def clean_response(self, text) -> str:
        if not isinstance(text, str):
            if isinstance(text, list):
                text = " ".join(
                    [str(t.get("text", t)) if isinstance(t, dict) else str(t) for t in text]
                )
            else:
                text = str(text)
        return re.sub(r"<thinking>.*?</thinking>", "", text, flags=re.DOTALL).strip()



    async def ask(self, question: str):
        try:
            logger.info(f"\n🧑 Customer: {question}")

            client = MultiServerMCPClient(MCP_CONFIG)

            async with client.session("cafebot") as session:
                raw = await session.list_tools()
                logger.info(f"Raw tools from server: {raw}")

                tools = await load_mcp_tools(session)
                logger.info(f"Loaded {len(tools)} MCP tools")

                agent = create_react_agent(
                    model=self.get_llm(),
                    tools=tools,
                    prompt=SYSTEM_PROMPT,
                )

                result = await agent.ainvoke(
                    {"messages": [("user", question)]}
                )

            raw_content = result["messages"][-1].content
            final_text = self.clean_response(raw_content)

            logger.info(f"🤖 CafeBot: {final_text}")
            return final_text

        except Exception as e:
            import traceback
            traceback.print_exc()
            log_error("agent.py", "ask", e)
            return "I'm having trouble connecting to the coffee machine. Please try again!"