import { useState, useEffect } from 'react'
import BookModal from './BookModal'

function Books({ token, apiUrl }) {
  const [books, setBooks] = useState([])
  const [audiobooks, setAudiobooks] = useState([])
  const [showModal, setShowModal] = useState(false)
  const [modalType, setModalType] = useState('book')
  const [editingItem, setEditingItem] = useState(null)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    fetchBooks()
    fetchAudiobooks()
  }, [])

  const fetchBooks = async () => {
    try {
      const response = await fetch(`${apiUrl}/books`, {
        headers: { 'Authorization': `Bearer ${token}` }
      })
      const result = await response.json()
      const data = result.data || result
      setBooks(Array.isArray(data) ? data : [])
    } catch (err) {
      console.error('Failed to fetch books:', err)
      setBooks([])
    }
  }

  const fetchAudiobooks = async () => {
    try {
      const response = await fetch(`${apiUrl}/audiobooks`, {
        headers: { 'Authorization': `Bearer ${token}` }
      })
      const result = await response.json()
      const data = result.data || result
      setAudiobooks(Array.isArray(data) ? data : [])
    } catch (err) {
      console.error('Failed to fetch audiobooks:', err)
      setAudiobooks([])
    }
  }

  const handleAddClick = (type) => {
    setModalType(type)
    setEditingItem(null)
    setShowModal(true)
  }

  const handleEditClick = (item, type) => {
    setModalType(type)
    setEditingItem(item)
    setShowModal(true)
  }

  const handleDelete = async (id, type) => {
    if (!confirm(`Delete this ${type}?`)) return

    try {
      const endpoint = type === 'book' ? 'books' : 'audiobooks'
      await fetch(`${apiUrl}/${endpoint}/${id}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${token}` }
      })
      type === 'book' ? fetchBooks() : fetchAudiobooks()
    } catch (err) {
      alert('Failed to delete')
    }
  }

  const handleModalSave = () => {
    setShowModal(false)
    fetchBooks()
    fetchAudiobooks()
  }

  return (
    <>
      <div className="card">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
          <h2>Books</h2>
          <button onClick={() => handleAddClick('book')} className="btn btn-primary">
            + Add Book
          </button>
        </div>
        {books.length === 0 ? (
          <p style={{ color: '#7f8c8d', padding: '20px' }}>No books yet. Click "+ Add Book" to get started!</p>
        ) : (
          <div className="grid">
            {books.map(book => (
              <div key={book.id} className="book-card">
              <h3>{book.title}</h3>
              <p>ISBN: {book.isbn}</p>
              <p>Total Pages: {book.total_pages}</p>
                <div className="book-actions">
                  <button onClick={() => handleEditClick(book, 'book')} className="btn btn-primary">
                    Edit
                  </button>
                  <button onClick={() => handleDelete(book.id, 'book')} className="btn btn-danger">
                    Delete
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      <div className="card">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
          <h2>Audiobooks</h2>
          <button onClick={() => handleAddClick('audiobook')} className="btn btn-primary">
            + Add Audiobook
          </button>
        </div>
        <div className="grid">
          {audiobooks.map(audiobook => (
            <div key={audiobook.id} className="book-card">
              <h3>{audiobook.title}</h3>
              <p>Duration: {audiobook.total_length}</p>
              <div className="book-actions">
                <button onClick={() => handleEditClick(audiobook, 'audiobook')} className="btn btn-primary">
                  Edit
                </button>
                <button onClick={() => handleDelete(audiobook.id, 'audiobook')} className="btn btn-danger">
                  Delete
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>

      {showModal && (
        <BookModal
          type={modalType}
          item={editingItem}
          token={token}
          apiUrl={apiUrl}
          onClose={() => setShowModal(false)}
          onSave={handleModalSave}
        />
      )}
    </>
  )
}

export default Books
