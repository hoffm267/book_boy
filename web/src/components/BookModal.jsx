import { useState } from 'react'

function BookModal({ type, item, token, apiUrl, onClose, onSave, queryParams }) {
  const [formData, setFormData] = useState(
    item || (type === 'book'
      ? { isbn: '', title: '', total_pages: '' }
      : { title: '', total_length: '' })
  )
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      const endpoint = type === 'book' ? 'books' : 'audiobooks'
      const method = item ? 'PUT' : 'POST'
      let url = item ? `${apiUrl}/${endpoint}/${item.id}` : `${apiUrl}/${endpoint}`

      if (queryParams && !item) {
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
        throw new Error(data.error || 'Failed to save')
      }

      const result = await response.json()
      onSave(result)
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value })
  }

  return (
    <div className="modal" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h2>{item ? `Edit ${type === 'book' ? 'Book' : 'Audiobook'}` : `Add ${type === 'book' ? 'Book by ISBN' : 'Audiobook'}`}</h2>
          <button onClick={onClose} className="close-btn">&times;</button>
        </div>

        {error && <div className="error">{error}</div>}

        <form onSubmit={handleSubmit}>
          {type === 'book' && (
            <div className="form-group">
              <label>ISBN *</label>
              <input
                type="text"
                name="isbn"
                value={formData.isbn}
                onChange={handleChange}
                required
              />
            </div>
          )}

          <div className="form-group">
            <label>Title {type === 'book' && '(optional - filled by worker)'}</label>
            <input
              type="text"
              name="title"
              value={formData.title}
              onChange={handleChange}
              required={type !== 'book'}
              placeholder={type === 'book' ? 'Leave empty to auto-fill' : ''}
            />
          </div>

          {type === 'book' ? (
            <div className="form-group">
              <label>Total Pages (optional - filled by worker)</label>
              <input
                type="number"
                name="total_pages"
                value={formData.total_pages}
                onChange={handleChange}
                min="1"
                placeholder="Leave empty to auto-fill"
              />
            </div>
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
