import { useState, useEffect } from 'react'
import { fetchWithAuth, UnauthorizedError } from '../utils/api'

function BookModal({ type, item, token, apiUrl, onClose, onSave, queryParams, userId, linkingToProgressId, prefillTitle, onLogout }) {
    const [useISBN, setUseISBN] = useState(!prefillTitle && !linkingToProgressId)
    const [formData, setFormData] = useState(
        item || (type === 'book'
            ? { isbn: '', title: prefillTitle || '', total_pages: '' }
            : { title: prefillTitle || '', total_length: '' })
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

            if (linkingToProgressId) {
                url += `?pgId=${linkingToProgressId}`
            } else if (queryParams && !item && type === 'book' && useISBN) {
                url += `?${queryParams}`
            }

            const body = type === 'book'
                ? { ...formData, isbn: formData.isbn.replace(/[^0-9]/g, ''), total_pages: parseInt(formData.total_pages) }
                : formData

            const response = await fetchWithAuth(
                url,
                {
                    method,
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${token}`
                    },
                    body: JSON.stringify(body)
                },
                onLogout
            )

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
            if (err.name !== 'UnauthorizedError') {
                setError(err.message)
            }
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
            const response = await fetchWithAuth(
                `${apiUrl}/books/filter?isbn=${encodeURIComponent(isbn)}`,
                { headers: { 'Authorization': `Bearer ${token}` } },
                onLogout
            )

            if (response.ok) {
                const books = await response.json()
                if (books && books.length > 0) {
                    setExistingBook(books[0])
                } else {
                    setExistingBook(null)
                }
            }
        } catch (err) {
            if (err.name !== 'UnauthorizedError') {
                console.error('ISBN check failed:', err)
            }
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
                    {type === 'book' && !item && !linkingToProgressId && (
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
                                onChange={linkingToProgressId ? handleChange : handleISBNChange}
                                required
                            />
                        </div>
                    )}

                    {((!useISBN && type === 'book') || type === 'audiobook' || linkingToProgressId) && (
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
                        (!useISBN || linkingToProgressId) && (
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
                            <label>Duration (e.g., "5:30")</label>
                            <input
                                type="text"
                                name="total_length"
                                value={formData.total_length}
                                onChange={handleChange}
                                required
                                placeholder="00:00:00"
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
