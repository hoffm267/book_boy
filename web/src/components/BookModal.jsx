import { useState, useEffect } from 'react'

function BookModal({ type, item, token, apiUrl, onClose, onSave, queryParams }) {
    const [useISBN, setUseISBN] = useState(true)
    const [formData, setFormData] = useState(
        item || (type === 'book'
            ? { isbn: '', title: '', total_pages: '' }
            : { title: '', total_length: '' })
    )
    const [error, setError] = useState('')
    const [loading, setLoading] = useState(false)
    const [existingBook, setExistingBook] = useState(null)
    const [checkingISBN, setCheckingISBN] = useState(false)
    const [checkISBNTimeout, setCheckISBNTimeout] = useState(null)

    useEffect(() => {
        return () => {
            if (checkISBNTimeout) {
                clearTimeout(checkISBNTimeout)
            }
        }
    }, [checkISBNTimeout])

    const handleSubmit = async (e) => {
        e.preventDefault()
        setError('')
        setLoading(true)

        try {
            const endpoint = type === 'book' ? 'books' : 'audiobooks'
            const method = item ? 'PUT' : 'POST'
            let url = item ? `${apiUrl}/${endpoint}/${item.id}` : `${apiUrl}/${endpoint}`

            if (queryParams && !item && type === 'book' && useISBN) {
                url += `?${queryParams}`
            }

            const body = type === 'book'
                ? { ...formData, total_pages: parseInt(formData.total_pages) }
                : formData

            const response = await fetch(url, {
                method,
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify(body)
            })

            if (!response.ok) {
                const data = await response.json()
                if (response.status === 409) {
                    throw new Error('You are already tracking this book. Check your reading list below.')
                }
                throw new Error(data.error || 'Failed to save')
            }

            const result = await response.json()
            const savedItem = result.data || result

            onSave(savedItem)
        } catch (err) {
            setError(err.message)
        } finally {
            setLoading(false)
        }
    }

    const handleChange = (e) => {
        setFormData({ ...formData, [e.target.name]: e.target.value })
    }

    const checkISBN = async (isbn) => {
        if (!isbn || isbn.length < 10) {
            setExistingBook(null)
            return
        }

        setCheckingISBN(true)
        try {
            const response = await fetch(`${apiUrl}/books/filter?isbn=${encodeURIComponent(isbn)}`, {
                headers: { 'Authorization': `Bearer ${token}` }
            })

            if (response.ok) {
                const books = await response.json()
                if (books && books.length > 0) {
                    setExistingBook(books[0])
                } else {
                    setExistingBook(null)
                }
            }
        } catch (err) {
            console.error('ISBN check failed:', err)
        } finally {
            setCheckingISBN(false)
        }
    }

    const handleISBNChange = (e) => {
        const isbn = e.target.value
        handleChange(e)

        if (checkISBNTimeout) {
            clearTimeout(checkISBNTimeout)
        }

        const timeout = setTimeout(() => checkISBN(isbn), 500)
        setCheckISBNTimeout(timeout)
    }

    return (
        <div className="modal" onClick={onClose}>
        <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
        <h2>
        {item
            ? `Edit ${type === 'book' ? 'Book' : 'Audiobook'}`
            : `Add ${type === 'book' ? (useISBN ? 'Book by ISBN' : 'Book Manually') : 'Audiobook'}`
        }
        </h2>
        <button onClick={onClose} className="close-btn">&times;</button>
        </div>

        {error && <div className="error">{error}</div>}

        <form onSubmit={handleSubmit}>
        {type === 'book' && !item && (
            <div className="form-group" style={{ marginBottom: '20px' }}>
            <label style={{ display: 'flex', alignItems: 'center', gap: '10px', cursor: 'pointer' }}>
            <input
            type="checkbox"
            checked={useISBN}
            onChange={(e) => setUseISBN(e.target.checked)}
            style={{ width: 'auto', cursor: 'pointer' }}
            />
            <span>Auto-fill by ISBN (fetches metadata automatically)</span>
            </label>
            </div>
        )}

        {type === 'book' && (
            <div className="form-group">
            <label>ISBN *</label>
            <input
            type="text"
            name="isbn"
            value={formData.isbn}
            onChange={handleISBNChange}
            required
            />
            {checkingISBN && <p style={{ fontSize: '0.9em', color: '#7f8c8d', marginTop: '5px' }}>Checking ISBN...</p>}
            {existingBook && !item && (
                <div style={{
                    marginTop: '10px',
                        padding: '10px',
                        background: '#ecf0f1',
                        borderRadius: '5px',
                        fontSize: '0.9em'
                }}>
                <strong>Book already exists:</strong> {existingBook.title || 'Untitled'} ({existingBook.total_pages || 0} pages)
                <br />
                <span style={{ color: '#7f8c8d' }}>Creating a new progress entry for this book</span>
                </div>
            )}
            </div>
        )}

        {(!useISBN || type !== 'book') && (
            <div className="form-group">
            <label>Title *</label>
            <input
            type="text"
            name="title"
            value={formData.title}
            onChange={handleChange}
            required
            />
            </div>
        )}

        {type === 'book' ? (
            !useISBN && (
                <div className="form-group">
                <label>Total Pages *</label>
                <input
                type="number"
                name="total_pages"
                value={formData.total_pages}
                onChange={handleChange}
                min="1"
                required
                />
                </div>
            )
        ) : (
            <div className="form-group">
            <label>Duration (e.g., "5 hours 30 minutes")</label>
            <input
            type="text"
            name="total_length"
            value={formData.total_length}
            onChange={handleChange}
            required
            placeholder="5 hours 30 minutes"
            />
            </div>
        )}

        <div style={{ display: 'flex', gap: '10px', justifyContent: 'flex-end' }}>
        <button type="button" onClick={onClose} className="btn btn-secondary">
        Cancel
        </button>
        <button type="submit" className="btn btn-primary" disabled={loading}>
        {loading ? 'Saving...' : 'Save'}
        </button>
        </div>
        </form>
        </div>
        </div>
    )
}

export default BookModal
