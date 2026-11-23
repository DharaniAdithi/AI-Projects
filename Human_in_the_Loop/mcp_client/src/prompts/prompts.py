
from mcp_client.src.routes.routes import router

@router.prompt()
def brew_buddy_prompt() -> str:
    """
    Returns the Brew Buddy system prompt for the AI agent.
    Defines the personality and rules for the chatbot.
    """
    return """You are Brew Buddy, a cheerful assistant for Bean & Brew coffee shop.
    Rules:
    - Use tools to place orders and check order status.
    - Use resources to fetch menu and shop information.
    - If the question is NOT related to the coffee shop, politely say you can only help with coffee shop queries.
    """


@router.prompt()
def order_confirmation_prompt(item_name: str, quantity: int, total_price: float) -> str:
    """
    Returns an order confirmation message for the customer
    before placing the order.
    """
    return (
        f"Please confirm your order:\n"
        f"Item     : {item_name}\n"
        f"Quantity : {quantity}\n"
        f"Total    : ₹{total_price}\n\n"
        f"Reply yes to confirm or no to cancel."
    )