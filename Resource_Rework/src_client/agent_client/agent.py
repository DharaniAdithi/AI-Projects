import re
import boto3
from typing import Any
from langchain_aws import ChatBedrock
from langgraph.prebuilt import create_react_agent
from langchain_mcp_adapters.client import MultiServerMCPClient
from langchain_mcp_adapters.prompts import load_mcp_prompt
from langchain_mcp_adapters.resources import load_mcp_resources
from langchain_core.messages import HumanMessage
from langchain_mcp_adapters.tools import load_mcp_tools
from src_client.agent_client.system_prompt import SYSTEM_PROMPT
from src.utils.helpers import log_error
from src.utils.logger import get_logger
from src.constants import (
    DEFAULT_REGION as REGION, 
    MODEL_ID, PROVIDER, 
    MCP_SERVER_URL,
    DEFAULT_MAX_TOKENS
)

logger = get_logger(__name__)

class AgentManager:

    def get_bedrock_client(self):
        try:
            return boto3.client(
                service_name="bedrock-runtime",
                region_name=REGION,
            )
        except Exception as e:
            log_error("agent.py", "get_bedrock_client", e)
            raise


    def get_llm(self, max_tokens=DEFAULT_MAX_TOKENS, temperature=0):
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


    def clean_response(self, text: Any) -> str:
        """
        Ensures output is a string and removes <thinking> tags.
        """
        if not isinstance(text, str):
            if isinstance(text, list):
                text = " ".join([str(t.get("text", t)) if isinstance(t, dict) else str(t) for t in text])
            else:
                text = str(text)

        return re.sub(r"<thinking>.*?</thinking>", "", text, flags=re.DOTALL).strip()


    async def ask(self, question: str):
        try:
            logger.info(f"\n🧑 Customer: {question}")

            client = MultiServerMCPClient(
                {
                    "cafebot": {
                        "transport": "streamable_http",
                        "url": MCP_SERVER_URL,
                    }
                }
            )

            async with client.session("cafebot") as session:
            
                tools = await load_mcp_tools(session)
                logger.info(f"\nLoaded {len(tools)} MCP tools")
                for t in tools:
                    logger.info(" • %s", t.name)

                tool_names = [t.name for t in tools]
                if "get_menu_items" not in tool_names:
                    logger.warning("⚠️ get_menu_items tool is MISSING — orders will fail!")
                
                prompt_messages = []
                prompt_messages += await load_mcp_prompt(session, "brew_buddy")
                prompt_messages += await load_mcp_prompt(session, "order_validation")
                prompt_messages += await load_mcp_prompt(session, "cafe_guardrails")

                
                blobs = await load_mcp_resources(
                    session,
                    uris=[
                        "cafe://menu-items",
                        "cafe://shop-hours",
                        "cafe://shop-info",
                    ],
                )

                resource_context = ""
                for blob in blobs:
                    try:
                        text = blob.as_bytes().decode("utf-8", errors="replace")
                        uri = blob.metadata.get("uri", "resource")
                        resource_context += f"\n[{uri}]\n{text}\n"
                    except Exception:
                        continue

                if resource_context:
                    prompt_messages.append(
                        HumanMessage(content=f"Available cafe resources:\n{resource_context}")
                    )

                llm = self.get_llm()
                agent = create_react_agent(
                    model=llm,
                    tools=tools,
                    prompt=SYSTEM_PROMPT
                )

                result = await agent.ainvoke(
                    {"messages": prompt_messages + [HumanMessage(content=question)]}
                )

                raw_content = result["messages"][-1].content
                final_text = self.clean_response(raw_content)

                logger.info(f"🤖 CafeBot: {final_text}")
                return final_text

        except Exception as e:
            import traceback
            traceback.print_exc()
            log_error("agent.py", "ask", e)
            return {
                "status": "error",
                "message": "I'm having trouble connecting to the coffee machine. Please try again!"
            }
