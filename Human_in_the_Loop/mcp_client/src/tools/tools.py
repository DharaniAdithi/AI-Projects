from mcp_client.src.routes.routes import router
from mcp_client.src.repository.repository import ChatRepository

repo = ChatRepository()

@router.tool()
def place_order_tool(customer_id: str, item_name: str, quantity: int) -> str:
    """Places a coffee order for a customer."""
    customer = repo.get_customer(customer_id)
    if not customer:
        return "Customer not found."
    item = repo.get_item_by_name(item_name)
    if not item:
        return "Item not available."
    if item.stock_quantity < quantity:
        return f"Only {item.stock_quantity} items available in stock."

    total_price = item.price * quantity
    billing = repo.create_billing(total_price)
    order = repo.create_order(
        customer_id=customer.customer_id,
        billing_id=billing.billing_id
    )

    repo.add_order_item(
        order_id=order.order_id,
        item_id=item.item_id,
        quantity=quantity
    )
    repo.update_stock(
        item_id=item.item_id,
        quantity=quantity
    )

    return (
        f"Order placed successfully!\n"
        f"Order ID : {order.order_id}\n"
        f"Item     : {item.item_name}\n"
        f"Quantity : {quantity}\n"
        f"Total    : ₹{total_price}"
    )


@router.tool()
def check_order_status_tool(order_id: str) -> str:
    """Returns the current status of a specific order by order ID."""
    status = repo.get_order_status(order_id)
    if not status:
        return "Order not found."
    return f"Order {order_id} status: {status}"