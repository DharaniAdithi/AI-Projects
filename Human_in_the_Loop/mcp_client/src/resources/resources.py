from mcp_client.src.routes.routes import router
from mcp_client.src.repository.repository import ChatRepository

repo = ChatRepository()

@router.resource("resource://menu")
def get_menu_tool() -> str:
    """Returns the list of available coffee items."""
    items = repo.get_all_items()
    if not items:
        return "No coffee items available."
    menu = []
    for item in items:
        menu.append(
            f"{item.item_name} - ₹{item.price} (Stock: {item.stock_quantity})"
        )
    return "Available Coffees:\n" + "\n".join(menu)

@router.resource("resource://shop-info")
def get_shop_info() -> str:
    """
    Returns Bean & Brew shop details like name,
    hours, location and contact number.
    """
    return (
        "Shop Name  : Bean & Brew\n"
        "Hours      : Mon-Sat 8AM-10PM | Sun 9AM-8PM\n"
        "Location   : 12, MG Road, Chennai\n"
        "Phone      : +91 98765 43210\n"
        "Email      : hello123@gmail.com.com"
    )