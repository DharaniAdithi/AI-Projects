import boto3
import os
from dotenv import load_dotenv
from langchain_aws import ChatBedrock
from langchain.agents import create_agent
from langchain.agents.middleware import SummarizationMiddleware, HumanInTheLoopMiddleware
from langchain_mcp_adapters.client import MultiServerMCPClient
from langchain_mcp_adapters.tools import load_mcp_tools
from langchain_mcp_adapters.prompts import load_mcp_prompt
from langchain_mcp_adapters.resources import load_mcp_resources
from langchain_core.messages import HumanMessage
from langchain_core.runnables import RunnableConfig
from langgraph.checkpoint.postgres.aio import AsyncPostgresSaver
from agent.system_prompt import SYSTEM_PROMPT
from utils.helpers import custom_error_handler
from langchain.agents.structured_output import ToolStrategy
from models.models import OrderDetails, ChatResponseDTO
from typing import Union
from langgraph.types import Command
from utils.logger import logger

load_dotenv()


class AgentManager:

    def __init__(self):
        self.region = "us-east-1"
        self.model_id = os.getenv("MODEL_ID")
        self.provider = os.getenv("PROVIDER")
        self.mcp_server_url = "http://localhost:8000/mcp"
        self.db_uri = "postgresql://postgres:Mass%40342005@localhost:5432/cafe_chatbot_mcp"
        self.prompt_args = {
            "brew_buddy_prompt": {},
            "order_confirmation_prompt": {
                "item_name": "Cappuccino",
                "quantity": "1",
                "total_price": "180.0"
            },
        }

        logger.info("AgentManager initialized")

    def get_bedrock_client(self):
        try:
            logger.debug("Creating Bedrock client")
            return boto3.client("bedrock-runtime", region_name=self.region)
        except Exception as e:
            logger.exception("Failed to create Bedrock client")
            raise Exception(f"Error creating Bedrock client: {str(e)}")

    def get_llm(self, max_tokens=500, temperature=0):
        try:
            logger.debug("Initializing Bedrock LLM")
            return ChatBedrock(
                client=self.get_bedrock_client(),
                model_id=self.model_id,
                provider=self.provider,
                model_kwargs={"max_tokens": max_tokens, "temperature": temperature},
            )
        except Exception as e:
            logger.exception("Failed to initialize LLM")
            raise Exception(f"Error initializing LLM: {str(e)}")

    async def ask(self, question: str, thread_id: str = "default") -> str:

        #SQ 1.22 - SQ 1.37 : Initialize the empty string, create the MCP cient and open the streamable http connection and fetch the registered tools,prompts,resources
        collected_text = ""
        logger.info(f"ask() called thread_id={thread_id}")

        try:
            client = MultiServerMCPClient(
                {
                    "bean-brew": {
                        "url": self.mcp_server_url,
                        "transport": "streamable-http",
                    }
                }
            )

            async with client.session("bean-brew") as session:
                logger.info("MCP session started")

                tools = await load_mcp_tools(session)
                logger.info("MCP tools loaded")

                list_result = await session.list_prompts()
                available_prompts = list_result.prompts

                all_prompt_messages = []
                for prompt in available_prompts:
                    if prompt.name not in self.prompt_args:
                        continue
                    args = self.prompt_args.get(prompt.name, {})
                    messages = await load_mcp_prompt(session, prompt.name, arguments=args)
                    all_prompt_messages.extend(messages)

                blobs = await load_mcp_resources(session)
                logger.info("MCP resources loaded")

                resource_context = ""
                for blob in blobs:
                    uri = blob.metadata.get("uri", "unknown")
                    text = blob.as_bytes().decode("utf-8", errors="replace")
                    resource_context += f"\n[{uri}]\n{text}\n"

                #SQ 1.38 - SQ 1.47 : open the async postgreSQL connection then create the checkpoint tables i not exist, and then build the langchain agent with tools,checkpointer,human in the loop middleware and then invoke the agent asynchronously with message and thread id
                async with AsyncPostgresSaver.from_conn_string(self.db_uri) as checkpointer:
                    logger.info("PostgreSQL checkpointer connected")
                    await checkpointer.setup()

                    agent = create_agent(
                        model=self.get_llm(),
                        tools=tools,
                        system_prompt=SYSTEM_PROMPT,
                        checkpointer=checkpointer,
                        middleware=[
                            SummarizationMiddleware(
                                model=self.get_llm(max_tokens=150, temperature=1),
                                trigger=("tokens", 150),
                                keep=("messages", 2),
                                summary_prompt="""
                                Summarize the conversation clearly and concisely in plain text.
                                IMPORTANT: Always preserve:
                                - Exact coffee items discussed with prices in Indian Rupees
                                - Exact order IDs if any order was placed
                                - Shop hours, location and phone if discussed
                                - Customer preferences
                                """
                            ),
                            HumanInTheLoopMiddleware( interrupt_on={"place_order_tool": True})
                        ],
                        # response_format=ToolStrategy(
                        #     schema=Union[OrderDetails, ChatResponseDTO],handle_errors=custom_error_handler
                        # ),
                    )

                    config = RunnableConfig(configurable={"thread_id": thread_id})

                    checkpoint = await checkpointer.aget(
                        {"configurable": {"thread_id": thread_id}}
                    )

                    # if checkpoint is None:
                    #     context = all_prompt_messages
                    #     if resource_context:
                    #         context = context + [HumanMessage(content=f"Available resources:\n{resource_context}")]
                    #     input_messages = context + [HumanMessage(content=question)]
                    # else:
                    #     input_messages = [HumanMessage(content=question)]
                        
                    response = await agent.ainvoke(
                        {"messages": question},
                        config
                    )

                    #SQ 1.48 - 1.58 : ask the customer to confirm the order the if yes means invoke the agent store the result in some variable and finally return that variable

                    print("\nRAW RESULT:", response)


                    if "__interrupt__" in response:
                        logger.warning("Human approval required")
                        print("\nHUMAN APPROVAL REQUIRED!")
                        print(response)

                        for intr in response["__interrupt__"]:
                            print("Interrupt:", intr.value)

                        while True:
                            user_input = input("\n Please confirm the order? (yes/no): ").strip().lower()
                            if user_input in ["yes", "no"]:
                                break
                            print("Please type 'yes' or 'no'")

                        decision = "approve" if user_input == "yes" else "reject"
                        final = await agent.ainvoke(Command(resume={"decisions": [{"type": decision}]}),config=config,)
                        print("\nFINAL RESULT:", final)

                        final_response = final["messages"][-1].content

                        if isinstance(final_response, list):
                            for block in final_response:
                                if isinstance(block, dict) and "text" in block:
                                    collected_text += block["text"]
                        elif isinstance(final_response, str):
                            collected_text = final_response

                    else:
                        final_response = response["messages"][-1].content

                        if isinstance(final_response, list):
                            for block in final_response:
                                if isinstance(block, dict) and "text" in block:
                                    collected_text += block["text"]
                        elif isinstance(final_response, str):
                            collected_text = final_response
            

                print(f"\n chatbot: {collected_text}")
                logger.info("Agent response generated successfully")
                return collected_text
        #SQ 1.59 : raise exception - agent processing failed
        except Exception as e:
            logger.exception("Agent processing failed")
            raise Exception(f"Agent processing failed: {str(e)}")

        
