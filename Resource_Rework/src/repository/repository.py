from src.repository.database import Database
from src.repository.schemas.schema import MenuItems, Orders, OrderItems, Customers, ErrorLogs
from datetime import datetime
from typing import Dict, Any, List
from src.utils.logger import get_logger
from src.constants import TOOL_GET_MENU, MSG_ORDER_NOT_FOUND
from src.utils.helpers import log_error

logger = get_logger(__name__)


class CafeBotRepository:
    """
    Repository for all cafe bot database operations.

    Responsibilities:
        - Menu retrieval.
        - Order creation, status check, update.
    """

    def __init__(self, database : Database) -> None:
        """Initialize with shared Database singleton."""
        self.db = database

    def _get_session(self):
        """
        Return a new database session.

        Returns:
            Session: Active SQLAlchemy session.
        """
        return self.db.SessionLocal()

    def get_menu_items(self) -> List[Dict[str, float]]:
        """
        Fetch all available menu items.

        Returns:
            List[Dict[str, Any]]: Menu items with name and price.

        Raises:
            Exception: On database query failure.
        """
        session = self._get_session()
        try:
            items = session.query(MenuItems).all()
            return [{"item_name": i.item_name, "price": i.price} for i in items]
        except Exception as e:
            log_error(file_name="repository.py", function_name=TOOL_GET_MENU, exception=e)
            raise
        finally:
            session.close()

    def create_order(
        self, customer_name: str, phone_no: str, item_name: str, quantity: int
    ) -> Dict[str, Any]:
        """
        Create a new customer order.

        Args:
            customer_name (str): Customer's name.
            phone_no (str): Customer's phone number.
            item_name (str): Menu item name.
            quantity (int): Number of items.

        Returns:
            Dict[str, Any]: Created order details.

        Raises:
            Exception: On database failure.
        """
        session = self._get_session()
        try:
            customer = Customers(name=customer_name, phone_no=phone_no)
            session.add(customer)
            session.flush()

            menu = session.query(MenuItems).filter(MenuItems.item_name == item_name).first()

            order = Orders(customer_id=customer.id)
            session.add(order)
            session.flush()

            order_item = OrderItems(order_id=order.id, item_id=menu.id, quantity=quantity)
            session.add(order_item)
            session.commit()

            logger.info("Order created: order_id=%s item=%s", order.id, item_name)
            return {
                "order_id": order.id,
                "item": item_name,
                "quantity": quantity,
                "status": "PENDING",
            }
        except Exception as e:
            log_error(file_name="repository.py", function_name="create_order", exception=e)
            raise
        finally:
            session.close()

    def get_order_status(self, order_id: int) -> Dict[str, str]:
        """
        Retrieve the status of an order.

        Args:
            order_id (int): The order's primary key.

        Returns:
            Dict[str, Any]: Order ID and status, or not-found message.

        Raises:
            Exception: On database failure.
        """
        session = self._get_session()
        try:
            order = session.query(Orders).filter(Orders.id == order_id).first()
            if not order:
                return {"message": MSG_ORDER_NOT_FOUND}
            return {"order_id": order_id, "status": order.status}
        except Exception as e:
            log_error(file_name="repository.py", function_name="get_order_status", exception=e)
            raise
        finally:
            session.close()

    def update_order(self, order_id: int, quantity: int) -> Dict[str, str]:
        """
        Update the quantity of an existing order.

        Args:
            order_id (int): The order's primary key.
            quantity (int): New quantity.

        Returns:
            Dict[str, Any]: Updated order summary with total price.

        Raises:
            Exception: On database failure.
        """
        session = self._get_session()
        try:
            order = session.query(Orders).filter(Orders.id == order_id).first()
            if not order:
                return {"message": MSG_ORDER_NOT_FOUND}

            order_item = session.query(OrderItems).filter(OrderItems.order_id == order_id).first()
            menu = session.query(MenuItems).filter(MenuItems.id == order_item.item_id).first()

            cost = menu.price * quantity
            order_item.quantity = quantity
            order_item.updated_at = datetime.utcnow()
            session.commit()

            logger.info("Order updated: order_id=%s quantity=%s total=%s", order_id, quantity, cost)
            return {
                "order_id": order_id,
                "item_name": menu.item_name,
                "quantity": quantity,
                "total_price": cost,
            }
        except Exception as e:
            log_error(file_name="repository.py", function_name="update_order", exception=e)
            raise
        finally:
            session.close()


class ErrorLogRepository:
    """
    Repository for persisting error logs.

    Responsibilities:
        - Store error records in the database.
    """

    def __init__(self, database : Database) -> None:
        """Initialize with shared Database singleton."""
        self.db = database

    def _get_session(self):
        """
        Return a new database session.

        Returns:
            Session: Active SQLAlchemy session.
        """
        return self.db.SessionLocal()

    def store_error(self, file_name: str, function_name: str, error: str) -> None:
        """
        Persist an error log entry.

        Args:
            file_name (str): File where the error occurred.
            function_name (str): Function where the error occurred.
            error (str): Full error message and stack trace.
        """
        session = self._get_session()
        try:
            log = ErrorLogs(file_name=file_name, function_name=function_name, error=error)
            session.add(log)
            session.commit()
        except Exception as db_error:
            session.rollback()
            logger.exception("Failed to store error log: %s", db_error)
        finally:
            session.close()