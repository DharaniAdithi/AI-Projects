from src.routes.mcp_route import router
from src.repository.repository import CafeBotRepository
from src.utils.helpers import log_error
import json


@router.resource("cafe://menu-items")
def menu_items_resource():
    """
    Resource: All available coffee menu items.
    """

    try:
        repo = CafeBotRepository()

        items = repo.get_menu_items()

        return json.dumps(items)

    except Exception as e:
        log_error(
            file_name="mcp_resources.py",
            function_name="menu_items_resource",
            exception=e
        )

        return json.dumps([])
@router.resource("cafe://shop-hours")
def shop_hours():
    """
    Resource: Coffee shop opening hours.
    """

    hours = {
        "Monday-Friday": "8 AM - 8 PM",
        "Saturday": "9 AM - 10 PM",
        "Sunday": "9 AM - 6 PM"
    }

    return json.dumps(hours)
@router.resource("cafe://shop-info")
def shop_info():
    """
    Resource: Basic information about Bean & Brew Coffee Shop.
    """

    info = {
        "name": "Bean & Brew",
        "specialty": "Artisan coffee and handcrafted beverages",
        "location": "Navalur, Chennai",
        "famous_drink": "Cappuccino"
    }

    return json.dumps(info)