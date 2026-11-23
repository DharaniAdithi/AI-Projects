from sqlalchemy import Column, Integer, text, DateTime, String, ForeignKey, Boolean, Enum
from sqlalchemy.orm import relationship, DeclarativeBase  # ← fix deprecated import
from enum import Enum as PyEnum


class Base(DeclarativeBase):
    """Base class for all ORM models."""
    pass


class OrderStatus(PyEnum):
    """Valid order status values."""
    processing = "processing"
    ready = "ready"


class Customers(Base):
    """ORM model for customers table."""

    __tablename__ = "customers"

    id = Column(Integer, primary_key=True, index=True)
    uuid = Column(String, server_default=text("gen_random_uuid()"), unique=True, index=True)
    name = Column(String)
    phone_no = Column(String)
    created_at = Column(DateTime, server_default=text("CURRENT_TIMESTAMP"))
    created_by = Column(String, default="SYSTEM")
    updated_at = Column(DateTime, server_default=text("CURRENT_TIMESTAMP"))
    updated_by = Column(String, default="SYSTEM")
    is_active = Column(Boolean, default=True)

    orders = relationship("Orders", back_populates="customer")


class Orders(Base):
    """ORM model for orders table."""

    __tablename__ = "orders"

    id = Column(Integer, primary_key=True, index=True)
    uuid = Column(String, server_default=text("gen_random_uuid()"), unique=True, index=True)
    customer_id = Column(Integer, ForeignKey("customers.id"))
    status = Column(Enum(OrderStatus), default=OrderStatus.processing)
    created_at = Column(DateTime, server_default=text("CURRENT_TIMESTAMP"))
    created_by = Column(String, default="SYSTEM")
    updated_at = Column(DateTime, server_default=text("CURRENT_TIMESTAMP"))
    updated_by = Column(String, default="SYSTEM")
    is_active = Column(Boolean, default=True)

    customer = relationship("Customers", back_populates="orders")
    order_items = relationship("OrderItems", back_populates="order")


class MenuItems(Base):
    """ORM model for menu_items table."""

    __tablename__ = "menu_items"

    id = Column(Integer, primary_key=True, index=True)
    uuid = Column(String, server_default=text("gen_random_uuid()"), unique=True, index=True)
    item_name = Column(String)
    price = Column(Integer)
    created_at = Column(DateTime, server_default=text("CURRENT_TIMESTAMP"))
    created_by = Column(String, default="SYSTEM")
    updated_at = Column(DateTime, server_default=text("CURRENT_TIMESTAMP"))
    updated_by = Column(String, default="SYSTEM")
    is_active = Column(Boolean, default=True)

    order_items = relationship("OrderItems", back_populates="menu_item")


class OrderItems(Base):
    """ORM model for order_items table."""

    __tablename__ = "order_items"

    id = Column(Integer, primary_key=True, index=True)
    uuid = Column(String, server_default=text("gen_random_uuid()"), unique=True, index=True)
    order_id = Column(Integer, ForeignKey("orders.id"))
    item_id = Column(Integer, ForeignKey("menu_items.id"))
    quantity = Column(Integer)
    created_at = Column(DateTime, server_default=text("CURRENT_TIMESTAMP"))
    created_by = Column(String, default="SYSTEM")
    updated_at = Column(DateTime, server_default=text("CURRENT_TIMESTAMP"))
    updated_by = Column(String, default="SYSTEM")
    is_active = Column(Boolean, default=True)

    order = relationship("Orders", back_populates="order_items")
    menu_item = relationship("MenuItems", back_populates="order_items")


class ErrorLogs(Base):
    """ORM model for error_logs table."""

    __tablename__ = "error_logs"

    id = Column(Integer, primary_key=True, index=True)
    uuid = Column(String, server_default=text("gen_random_uuid()"), unique=True, index=True)
    file_name = Column(String)
    function_name = Column(String)
    error = Column(String)
    created_at = Column(DateTime, server_default=text("CURRENT_TIMESTAMP"))
    created_by = Column(String, default="SYSTEM")
    updated_at = Column(DateTime, server_default=text("CURRENT_TIMESTAMP"))
    updated_by = Column(String, default="SYSTEM")
    is_active = Column(Boolean, default=True)
