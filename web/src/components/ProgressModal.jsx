import { useState, useEffect } from 'react'

function ProgressModal({ progress, books, audiobooks, userId, token, apiUrl, onClose, onSave }) {
    const [formData, setFormData] = useState({
        book_id: progress?.book_id || '',
        audiobook_id: progress?.audiobook_id || '',
        book_page: progress?.book_page || '',
        audiobook_time: progress?.audiobook_time || ''
    })
    const [error, setError] = useState('')
    const [loading, setLoading] = useState(false)

    const handleSubmit = async (e) => {
        e.preventDefault()
        setError('')
        setLoading(true)

        try {
            const method = progress ? 'PUT' : 'POST'
            const url = progress ? `${apiUrl}/progress/${progress.id}` : `${apiUrl}/progress`

            const body = {
                user_id: userId,
                book_id: formData.book_id ? parseInt(formData.book_id) : null,
                audiobook_id: formData.audiobook_id ? parseInt(formData.audiobook_id) : null,
                book_page: formData.book_page ? parseInt(formData.book_page) : null,
                audiobook_time: formData.audiobook_time || null
            }

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
                throw new Error(data.error || 'Failed to save progress')
            }

            onSave()
        } catch (err) {
            setError(err.message)
        } finally {
            setLoading(false)
        }
    }

    return (
        <div className="modal" onClick={onClose}>
            <div className="modal-content" onClick={(e) => e.stopPropagation()}>
                <div className="modal-header">
                    <h2>{progress ? 'Update Progress' : 'Start Tracking'}</h2>
                    <button onClick={onClose} className="close-btn">&times;</button>
                </div>

                {error && <div className="error">{error}</div>}

                <form onSubmit={handleSubmit}>
                    <div className="form-group">
                        <label>Book (optional)</label>
                        <select
                            value={formData.book_id}
                            onChange={(e) => setFormData({ ...formData, book_id: e.target.value })}
                        >
                            <option value="">None</option>
                            {books.map(book => (
                                <option key={book.id} value={book.id}>
                                    {book.title} ({book.total_pages} pages)
                                </option>
                            ))}
                        </select>
                    </div>

                    {formData.book_id && (
                        <div className="form-group">
                            <label>Current Page</label>
                            <input
                                type="number"
                                value={formData.book_page}
                                onChange={(e) => setFormData({ ...formData, book_page: e.target.value })}
                                min="1"
                                required
                            />
                        </div>
                    )}

                    <div className="form-group">
                        <label>Audiobook (optional)</label>
                        <select
                            value={formData.audiobook_id}
                            onChange={(e) => setFormData({ ...formData, audiobook_id: e.target.value })}
                        >
                            <option value="">None</option>
                            {audiobooks.map(ab => (
                                <option key={ab.id} value={ab.id}>
                                    {ab.title} ({ab.total_length})
                                </option>
                            ))}
                        </select>
                    </div>

                    {formData.audiobook_id && (
                        <div className="form-group">
                            <label>Current Time (HH:MM:SS)</label>
                            <input
                                type="text"
                                value={formData.audiobook_time}
                                onChange={(e) => setFormData({ ...formData, audiobook_time: e.target.value })}
                                placeholder="00:00:00"
                                pattern="[0-9]{2}:[0-9]{2}:[0-9]{2}"
                                required
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

export default ProgressModal
