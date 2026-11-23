from sqlalchemy.orm import Session
from src.repository.database import Database
from src.repository.schemas.schema import MenuItems
from src.utils.logger import get_logger

logger = get_logger(__name__)

INITIAL_MENU_ITEMS: list[dict] = [
    {"item_name": "Espresso", "price": 120},
    {"item_name": "Cappuccino", "price": 150},
    {"item_name": "Latte", "price": 160},
    {"item_name": "Cold Coffee", "price": 180},
    {"item_name": "Sandwich", "price": 100},
    {"item_name": "Brownie", "price": 90},
]


class Seeder:
    """
    Handles initial data seeding.

    Responsibilities:
        - Seed menu items if table is empty.
    """

    def __init__(self) -> None:
        """Initialize Seeder with shared Database singleton."""
        self.db = Database()

    def seed_data(self) -> None:
        """
        Seed initial menu items into the database.
        Skips seeding if data already exists.
        """
        session = Session(self.db.engine)
        try:
            if session.query(MenuItems).first():
                logger.info("Menu items already seeded, skipping.")
                return

            menu_items = [
                MenuItems(item_name=item.get("item_name"), price=item.get("price"))
                for item in INITIAL_MENU_ITEMS
            ]
            session.add_all(menu_items)
            session.commit()
            logger.info("Seeded %d menu items.", len(menu_items))

        except Exception as e:
            session.rollback()
            logger.exception("Seeding failed: %s", e)
            raise

        finally:
            session.close()