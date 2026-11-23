import uuid
from sqlalchemy import Column, String, Integer, Boolean, Date, ForeignKey, TIMESTAMP, Text, DATETIME,Identity
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy.orm import relationship, declarative_base
from sqlalchemy.sql import func
from datetime import date as date
 
Base = declarative_base()

class Customer(Base):
    __tablename__ = "customers"
 
    customer_id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    id = Column(Integer, index = True, nullable=False)
    customer_name = Column(String(100), nullable=False)
    phone_number = Column(String(15), unique=True)
    address = Column(Text, nullable=False)
    email = Column(String(255), unique=True)
    is_active = Column(Boolean, default=True)
    created_by = Column(String(100), nullable=False)
    created_at = Column(TIMESTAMP, server_default=func.now())
    updated_by = Column(String(100), nullable=False)
    updated_at = Column(TIMESTAMP, server_default=func.now(), onupdate=func.now())
 
    orders = relationship("Order", back_populates="customer")
 
 
class Billing(Base):
    __tablename__ = "billings"
 
    billing_id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    id = Column(Integer, index = True, nullable=False)
    date = Column(Date, nullable=False, default=date.today)
    amount = Column(Integer, nullable=False)
    is_active = Column(Boolean, default=True)
    created_by = Column(String(100), nullable=False)
    created_at = Column(TIMESTAMP, server_default=func.now())
    updated_by = Column(String(100), nullable=False)
    updated_at = Column(TIMESTAMP, server_default=func.now(), onupdate=func.now())
 
    orders = relationship("Order", back_populates="billing")
 
 
class Order(Base):
    __tablename__ = "orders"
 
    order_id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    id = Column(Integer, index = True, nullable=False)
    status = Column(String(25), nullable=False)
    customer_id = Column(UUID(as_uuid=True), ForeignKey("customers.customer_id"), nullable=False)
    billing_id = Column(UUID(as_uuid=True), ForeignKey("billings.billing_id"))
    is_active = Column(Boolean, default=True)
    created_by = Column(String(100), nullable=False)
    created_at = Column(TIMESTAMP, server_default=func.now())
    updated_by = Column(String(100), nullable=False)
    updated_at = Column(TIMESTAMP, server_default=func.now(), onupdate=func.now())
    customer = relationship("Customer", back_populates="orders")
    billing = relationship("Billing", back_populates="orders")
    order_items = relationship("OrderItem", back_populates="order")
 
 
class Item(Base):
    __tablename__ = "items"
 
    item_id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    id = Column(Integer, index = True, nullable=False)
    item_name = Column(String(100), nullable=False)
    stock_quantity = Column(Integer, nullable=False)
    price = Column(Integer, nullable=False)
    is_active = Column(Boolean, default=True)
    created_by = Column(String(100), nullable=False)
    created_at = Column(TIMESTAMP, server_default=func.now())
    updated_by = Column(String(100), nullable=False)
    updated_at = Column(TIMESTAMP, server_default=func.now(), onupdate=func.now())
 
    order_items = relationship("OrderItem", back_populates="item")
 
 
class OrderItem(Base):
    __tablename__ = "order_items"
 
    order_items_id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    id = Column(Integer, index = True, nullable=False)
    quantity = Column(Integer, nullable=False)
    order_id = Column(UUID(as_uuid=True), ForeignKey("orders.order_id"))
    item_id = Column(UUID(as_uuid=True), ForeignKey("items.item_id"))
    is_active = Column(Boolean, default=True)
    created_by = Column(String(100), nullable=False)
    created_at = Column(TIMESTAMP, server_default=func.now())
    updated_by = Column(String(100), nullable=False)
    updated_at = Column(TIMESTAMP, server_default=func.now(), onupdate=func.now())
    order = relationship("Order", back_populates="order_items")
    item = relationship("Item", back_populates="order_items")
 
 
class Error(Base):
    __tablename__ = "errors"
 
    error_id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    id = Column(Integer, index = True, nullable=False)
    file_name = Column(String(100), nullable=False)
    function_name = Column(String(100), nullable=False)
    error_name = Column(String, nullable=False)
    is_active = Column(Boolean, default=True)
    created_by = Column(String(100), nullable=False)
    created_at = Column(TIMESTAMP, server_default=func.now())
    updated_by = Column(String(100), nullable=False)
    updated_at = Column(TIMESTAMP, server_default=func.now(), onupdate=func.now())
 