from pydantic import BaseModel, Field
from typing import Optional


class ChatRequestDTO(BaseModel):
    query: str
    customer_id: str

class OrderDetails(BaseModel):
    """structure for the order details
    use this model when user ask about i want to place order
    item: name the item , 
    quantity: count of the item, 
    price: total price of the item"""
    item: str
    quantity: int
    price:float

class ChatResponseDTO(BaseModel):
    """"
    structure for the response
    use this model when user ask about menus, status of the item, normal things other than place order
    """
    status:str
    response:str