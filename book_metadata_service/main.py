from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
import httpx
from typing import Optional
from pydantic import BaseModel

app = FastAPI(title="Book Metadata Service")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


class BookMetadata(BaseModel):
    isbn: str
    title: Optional[str] = None
    author: Optional[str] = None
    total_pages: Optional[int] = None
    cover_url: Optional[str] = None
    publisher: Optional[str] = None
    publish_date: Optional[str] = None


@app.get("/")
async def root():
    return {"service": "Book Metadata Service", "status": "healthy"}


@app.get("/health")
async def health():
    return {"status": "healthy"}


@app.get("/books/isbn/{isbn}", response_model=BookMetadata)
async def get_book_by_isbn(isbn: str):
    """
    Get book metadata by ISBN from Open Library API
    """
    try:
        async with httpx.AsyncClient(follow_redirects=True) as client:
            response = await client.get(
                f"https://openlibrary.org/isbn/{isbn}.json",
                timeout=10.0
            )

            if response.status_code == 404:
                raise HTTPException(status_code=404, detail="Book not found")

            response.raise_for_status()
            data = response.json()

            publisher_data = None
            if "publishers" in data and len(data["publishers"]) > 0:
                publisher_data = data["publishers"][0]

            book = BookMetadata(
                isbn=isbn,
                title=data.get("title"),
                total_pages=data.get("number_of_pages"),
                publisher=publisher_data,
                publish_date=data.get("publish_date"),
            )

            if "authors" in data and len(data["authors"]) > 0:
                author_key = data["authors"][0]["key"]
                author_response = await client.get(
                    f"https://openlibrary.org{author_key}.json",
                    timeout=10.0
                )
                if author_response.status_code == 200:
                    author_data = author_response.json()
                    book.author = author_data.get("name")

            if "covers" in data and len(data["covers"]) > 0:
                cover_id = data["covers"][0]
                book.cover_url = f"https://covers.openlibrary.org/b/id/{cover_id}-L.jpg"

            return book

    except httpx.HTTPError as e:
        raise HTTPException(status_code=500, detail=f"Error fetching book data: {str(e)}")


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
