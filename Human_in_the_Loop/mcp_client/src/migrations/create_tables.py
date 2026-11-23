from sqlalchemy.orm import Session
from mcp_client.src.repository.database import Database
from mcp_client.src.repository.schemas.schema import (
    Customer,
    Item,
    Order,
    OrderItem,
    Billing,
    Error,
    Base
)


def create_tables() -> None:
    print("create tables")
    db = Database()
    Base.metadata.create_all(bind=db.engine)
    print("table created successfully")
    _seed_data()


def _seed_data() -> None:
    session: Session | None = None
    try:
        db = Database()
        session = db.SessionLocal()

        # prevent duplicate seed
        if session.query(Customer).first():
            return

        # Customers
        customers = [
            Customer(id=1, customer_name="Durai", phone_number="9876543210", address="Chennai", email="durai@gmail.com", created_by="system", updated_by="admin"),
            Customer(id=2, customer_name="Prakash", phone_number="9123456780",address="Coimbatore", email="prakash@gmail.com",created_by="system", updated_by="admin"),
            Customer(id=3, customer_name="Siva", phone_number="9988776655",address="Madurai", email="siva@gmail.com",created_by="system", updated_by="admin"),
            Customer(id=4, customer_name="Akshay", phone_number="9871234560",address="Salem", email="akshay@gmail.com",created_by="system", updated_by="admin"),
            Customer(id=5, customer_name="Viswa", phone_number="9765432100",address="Trichy", email="viswa@gmail.com",created_by="system", updated_by="admin"),
        ]

        session.add_all(customers)
        session.flush()

        # Items
        items = [
            Item(id=1, item_name="coffee", stock_quantity=100, price=120, created_by="system", updated_by="durai"),
            Item(id=2, item_name="Cappuccino", stock_quantity=80, price=180, created_by="system", updated_by="durai"),
            Item(id=3, item_name="tea", stock_quantity=70, price=200, created_by="system", updated_by="durai"),
            Item(id=4, item_name="milk", stock_quantity=60, price=220, created_by="system", updated_by="durai"),
            Item(id=5, item_name="cold coffee", stock_quantity=50, price=240, created_by="system", updated_by="durai"),
        ]

        session.add_all(items)
        session.flush()

        # Billing
        billings = [
            Billing(id=1, date="2026-03-12", amount=300, created_by="system", updated_by="durai"),
            Billing(id=2, date="2026-03-12", amount=200, created_by="system", updated_by="durai"),
        ]

        session.add_all(billings)
        session.flush()

        # Orders
        orders = [
            Order(
                id=1, customer_id=customers[0].customer_id, billing_id=billings[0].billing_id, status="PREPARING", created_by="system", updated_by="durai"
            ),
            Order( id=2, customer_id=customers[1].customer_id, billing_id=billings[1].billing_id, status="READY", created_by="system", updated_by="durai"
            ),
        ]

        session.add_all(orders)
        session.flush()

        # Order Items
        order_items = [
            OrderItem(id=1, order_id=orders[0].order_id, item_id=items[0].item_id, quantity=1, created_by="system", updated_by="system"),
            OrderItem(id=2, order_id=orders[0].order_id, item_id=items[1].item_id, quantity=2, created_by="system", updated_by="system"),
            OrderItem(id=3, order_id=orders[1].order_id, item_id=items[2].item_id, quantity=1, created_by="system", updated_by="system"),
            OrderItem(id=4, order_id=orders[1].order_id, item_id=items[3].item_id, quantity=1, created_by="system", updated_by="system"),
        ]

        session.add_all(order_items)
        session.commit()

    except Exception as e:
        if session:
            session.rollback()
        raise
    finally:
        if session:
            session.close()