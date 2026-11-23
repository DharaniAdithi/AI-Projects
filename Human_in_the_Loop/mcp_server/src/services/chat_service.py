import datetime
import uuid
from agent.mcp_client import AgentManager
from utils.helpers import log_error
from langchain.agents.middleware._redaction import (
    apply_strategy,
    PIIDetectionError,
    detect_email,
    detect_credit_card
)
from utils.logger import logger

class ChatService:
    def __init__(self):
        self.agent = AgentManager()

#SQ 1.8 - 1.18: Instantiate ChatService, validate non‑empty message, generate request ID and timestamp, enrich message with customer UUID, and prepare message for guardrail checks.
    async def chat_with_agent(self, message: str, thread_id: str) -> dict:
        print("service")
        logger.info("chat_with_agent called")

        if not message or message.strip() == "":
            logger.error("Message cannot be empty")
            log_error(file_name="chat_service.py", function_name="chat_with_agent", exception="Message cannot be empty")
            raise ValueError("Message cannot be empty")

        try:
            request_id = str(uuid.uuid4())
            timestamp = datetime.datetime.now().isoformat()
            
            cust_msg = f"{message}\n My customer_id is: {thread_id}"
            logger.debug(f"Customer message received: {cust_msg}")
            print(f"message: {cust_msg}")

            try:
                matches = (
                    detect_email(cust_msg)
                    + detect_credit_card(cust_msg)
                )
                if not matches:
                    logger.info("No PII detected in message")
                    print("No PII detected")
                else:
                    for m in matches:
                        print(f"  [{m['type']:10}] {m['value']}")
                        logger.warning(f"PII detected - Type: {m['type']}")

                    strategy = "mask"
                    cust_msg = apply_strategy(cust_msg, matches, strategy)
                    print(f"Strategy '{strategy}' applied")
                    print(f"After: {cust_msg[:100]}")
                    logger.info(f"PII masking strategy '{strategy}' applied")

            except PIIDetectionError as e:
                logger.error(f"PII detection failed: {str(e)}")
                return {
                    "data": ["you message contains sensitive information so can't processing"],
                    "request_id": request_id,
                    "timestamp": timestamp
                }
            logger.info("Sending message to agent")

            #SQ 1.21 : Send sanitized message and thread ID to AgentManager to AgentManager
            final_response = await self.agent.ask(cust_msg, thread_id)
            logger.info("Response received from agent")

            #SQ 1.60 : build the response with data, request_id and timestamp
            return {
                "data": [final_response],
                "request_id": request_id,
                "timestamp": timestamp
            }
        
        #SQ 1.62 - SQ 1.65 : log the value error and unexpected exception 
        except ValueError as ve:
            logger.error(f"Validation error: {str(ve)}")
            log_error(file_name="chat_service.py", function_name="chat_with_agent", exception=ve)
            raise ValueError(f"Invalid input: {str(ve)}")

        except Exception as e:
            logger.exception("Unhandled exception in chat_with_agent")
            log_error(file_name="chat_service.py", function_name="chat_with_agent", exception=e)
            raise Exception(f"Chat processing failed: {str(e)}")