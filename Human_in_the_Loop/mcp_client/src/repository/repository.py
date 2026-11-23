from sqlalchemy.orm import Session
from sqlalchemy import func
from datetime import date

from mcp_client.src.repository.database import Database
from mcp_client.src.repository.schemas.schema import (
    Customer,
    Item,
    Order,
    OrderItem,
    Billing,
    Error
)


class ChatRepository:

    def __init__(self):
        self.db = Database()

    def _get_session(self) -> Session:
        try:
            return self.db.SessionLocal()
        except Exception as e:
            raise "Database connection failed"

    def get_customer(self, customer_id: str):
        session = None
        try:
            session = self._get_session()
            return session.query(Customer).filter(
                Customer.customer_id == customer_id
            ).first()
        except Exception as e:
            raise
        finally:
            if session:
                session.close()

    def get_all_customers(self):
        session = None
        try:
            session = self._get_session()
            return session.query(Customer).all()
        finally:
            if session:
                session.close()

    def get_item_by_name(self, item_name: str):
        session = None
        try:
            session = self._get_session()
            return session.query(Item).filter(
                Item.item_name.ilike(item_name)
            ).first()
        finally:
            if session:
                session.close()

    def get_all_items(self):
        session = None
        try:
            session = self._get_session()
            return session.query(Item).all()
        finally:
            if session:
                session.close()

    def check_stock(self, item_id: str):
        session = None
        try:
            session = self._get_session()
            item = session.query(Item).filter(
                Item.item_id == item_id
            ).first()
            return item.stock_quantity if item else None
        finally:
            if session:
                session.close()

    def _get_next_id(self, session: Session, model) -> int:
        max_id = session.query(func.max(model.id)).scalar()
        return (max_id or 0) + 1

    def create_billing(self, amount: float):
        session = None
        try:
            session = self._get_session()
            billing = Billing(
                id=self._get_next_id(session, Billing),
                date=date.today(),
                amount=amount,
                created_by="durai",
                updated_by="durai"
            )
            session.add(billing)
            session.commit()
            session.refresh(billing)
            return billing
        except Exception as e:
            session.rollback()
            raise
        finally:
            if session:
                session.close()

    def create_order(self, customer_id: str, billing_id: str, status: str = "PLACED"):
        session = None
        try:
            session = self._get_session()
            order = Order(
                id=self._get_next_id(session, Order),
                customer_id=customer_id,
                billing_id=billing_id,
                status=status,
                created_by="system",
                updated_by="system"
            )
            session.add(order)
            session.commit()
            session.refresh(order)
            return order
        except Exception as e:
            session.rollback()
            raise
        finally:
            if session:
                session.close()

    def add_order_item(self, order_id: str, item_id: str, quantity: int):
        session = None
        try:
            session = self._get_session()
            order_item = OrderItem(
                id=self._get_next_id(session, OrderItem),
                order_id=order_id,
                item_id=item_id,
                quantity=quantity,
                created_by="system",
                updated_by="system"
            )
            session.add(order_item)
            session.commit()
            session.refresh(order_item)
            return order_item
        except Exception as e:
            session.rollback()
            raise
        finally:
            if session:
                session.close()

    def update_stock(self, item_id: str, quantity: int):
        session = None
        try:
            session = self._get_session()
            item = session.query(Item).filter(
                Item.item_id == item_id
            ).first()
            if item:
                item.stock_quantity -= quantity
                session.commit()
                session.refresh(item)
            return item
        except Exception as e:
            session.rollback()
            raise
        finally:
            if session:
                session.close()

    def get_order_status(self, order_id: str):
        session = None
        try:
            session = self._get_session()
            order = session.query(Order).filter(
                Order.order_id == order_id
            ).first()
            return order.status if order else None
        finally:
            if session:
                session.close()


class ErrorLogRepository:
    def __init__(self):
        self.db = Database()

    def _get_session(self) -> Session:
        try:
            return self.db.SessionLocal()
        except Exception as e:
            raise "Database connection failed"

    def store_error(self,file_name:str,function_name:str,error:str):
        session=None
        try:
            session=self._get_session()
            log=Error(file_name=file_name,function_name=function_name,error=error)
            session.add(log)
            session.commit()
        except Exception as db_error:
            session.rollback()
            print("Error while saving log",db_error)
        finally:
            if session:
                session.close()  
