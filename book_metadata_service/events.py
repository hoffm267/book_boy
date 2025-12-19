import json
from aio_pika import connect_robust, Message, DeliveryMode, ExchangeType
from pydantic import BaseModel
from typing import Optional


class BookCreatedEvent(BaseModel):
    book_id: int
    isbn: str
    created_at: str


class BookMetadataFetchedEvent(BaseModel):
    book_id: int
    isbn: str
    title: str
    total_pages: int
    author: Optional[str] = None
    publisher: Optional[str] = None
    success: bool = True
    error: Optional[str] = None


class EventPublisher:
    def __init__(self, connection, exchange_name: str):
        self.connection = connection
        self.exchange_name = exchange_name
        self.channel = None
        self.exchange = None

    async def setup(self):
        self.channel = await self.connection.channel()
        self.exchange = await self.channel.declare_exchange(
            self.exchange_name,
            ExchangeType.TOPIC,
            durable=True
        )

    async def publish(self, routing_key: str, event: BaseModel):
        body = json.dumps(event.model_dump()).encode()
        message = Message(
            body,
            delivery_mode=DeliveryMode.PERSISTENT,
            content_type='application/json'
        )

        #TODO fix this
        if self.exchange is None:
            return
        
        await self.exchange.publish(message, routing_key=routing_key)
        print(f"Published event: {routing_key}")

    async def close(self):
        if self.channel:
            await self.channel.close()
