from sqlalchemy import inspect
from src.repository.database import Database
from src.repository.schemas.schema import Customers, Orders, MenuItems, OrderItems, ErrorLogs
from src.utils.logger import get_logger

logger = get_logger(__name__)

TABLE_ORDER_CREATION: list[str] = [
    Customers.__tablename__,
    Orders.__tablename__,
    MenuItems.__tablename__,
    OrderItems.__tablename__,
    ErrorLogs.__tablename__,
]

MODEL_CLASSES: dict[str, type] = {
    Customers.__tablename__: Customers,
    Orders.__tablename__: Orders,
    MenuItems.__tablename__: MenuItems,
    OrderItems.__tablename__: OrderItems,
    ErrorLogs.__tablename__: ErrorLogs,
}


class Migration:
    """
    Handles database table creation at startup.

    Responsibilities:
        - Check existing tables.
        - Create missing tables in correct order.
    """

    def __init__(self) -> None:
        """Initialize Migration with shared Database singleton."""
        self.db = Database()
        self.engine = self.db.engine
        self.inspector = inspect(self.engine)

    def create_tables(self) -> None:
        """
        Create all required tables if they do not already exist.
        """
        for table_name in TABLE_ORDER_CREATION:
            if not self.inspector.has_table(table_name):
                MODEL_CLASSES[table_name].__table__.create(bind=self.engine)
                logger.info("Created table: %s", table_name)
            else:
                logger.info("Table already exists, skipping: %s", table_name)