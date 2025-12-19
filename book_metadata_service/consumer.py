import os
import json
import asyncio
from typing import Optional
import httpx
from aio_pika import connect_robust, ExchangeType
from events import BookCreatedEvent, BookMetadataFetchedEvent, EventPublisher

RABBITMQ_URL = os.getenv('RABBITMQ_URL', 'amqp://guest:guest@rabbitmq:5672/')
EXCHANGE = 'book_events'


async def fetch_metadata_from_isbn(isbn: str) -> dict:
    """Fetch book metadata from OpenLibrary API"""
    try:
        async with httpx.AsyncClient(follow_redirects=True) as client:
            response = await client.get(
                f"https://openlibrary.org/isbn/{isbn}.json",
                timeout=10.0
            )

            response.raise_for_status()
            data = response.json()

            publisher_data = None
            if "publishers" in data and len(data["publishers"]) > 0:
                publisher_data = data["publishers"][0]

            metadata = {
                'title': data.get("title", ""),
                'total_pages': data.get("number_of_pages", 0),
                'publisher': publisher_data,
            }

            return metadata

    except httpx.HTTPStatusError as e:
        print(f"HTTP error: {e.response.status_code}")
        return {}
    except httpx.RequestError as e:
        print(f"Request error: {e}")
        return {}



async def process_message(message, publisher):
    """Process a book.created event"""
    async with message.process():
        try:
            data = json.loads(message.body)
            event = BookCreatedEvent.model_validate(data)

            print(f"Processing book {event.book_id} with ISBN {event.isbn}")

            metadata = await fetch_metadata_from_isbn(event.isbn)

            if metadata and metadata.get('title'):
                result = BookMetadataFetchedEvent(
                    book_id=event.book_id,
                    isbn=event.isbn,
                    title=metadata.get('title', ''),
                    total_pages=metadata.get('total_pages', 0),
                    author=metadata.get('author'),
                    publisher=metadata.get('publisher'),
                    success=True
                )
            else:
                result = BookMetadataFetchedEvent(
                    book_id=event.book_id,
                    isbn=event.isbn,
                    title='',
                    total_pages=0,
                    success=False,
                    error='Failed to fetch metadata from OpenLibrary'
                )

            await publisher.publish('book.metadata_fetched', result)

        except Exception as e:
            print(f"Error processing message: {e}")
            raise


async def main():
    print("Metadata service consumer starting...")

    max_retries = 10
    retry_delay = 2

    for attempt in range(max_retries):
        try:
            connection = await connect_robust(RABBITMQ_URL)
            break
        except Exception as e:
            if attempt < max_retries - 1:
                print(f"Failed to connect to RabbitMQ (attempt {attempt + 1}/{max_retries}): {e}")
                await asyncio.sleep(retry_delay)
                retry_delay = min(retry_delay * 2, 30)
            else:
                raise

    publisher = EventPublisher(connection, EXCHANGE)
    await publisher.setup()

    channel = await connection.channel()
    await channel.set_qos(prefetch_count=1)

    exchange = await channel.declare_exchange(
        EXCHANGE,
        ExchangeType.TOPIC,
        durable=True
    )

    queue = await channel.declare_queue('metadata_service_jobs', durable=True)
    await queue.bind(exchange, routing_key='book.created')

    print("Metadata service consumer started, waiting for book.created events...")

    async with queue.iterator() as queue_iter:
        async for message in queue_iter:
            await process_message(message, publisher)


if __name__ == '__main__':
    asyncio.run(main())
