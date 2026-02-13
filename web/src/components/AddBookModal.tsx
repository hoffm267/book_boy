import { useState, FormEvent, ChangeEvent } from 'react'
import { useAuth } from '../AuthContext'
import type { BookFormData } from '../types'

interface AddBookModalProps {
  onClose: () => void
  onSave: () => void
  linkingToProgressId: number | null
  prefillTitle: string | null
}

function AddBookModal({ onClose, onSave, linkingToProgressId, prefillTitle }: AddBookModalProps) {
  const { api } = useAuth()
  const [useISBN, setUseISBN] = useState(!prefillTitle && !linkingToProgressId)
  const [formData, setFormData] = useState<BookFormData>({
    isbn: '',
    title: prefillTitle || '',
    total_pages: '',
  })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      const path = linkingToProgressId ? `/books?pgId=${linkingToProgressId}` : '/books'
      await api.post(path, {
        ...formData,
        isbn: formData.isbn.replace(/[^0-9]/g, ''),
        total_pages: parseInt(String(formData.total_pages)),
      })
      onSave()
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message)
      }
    } finally {
      setLoading(false)
    }
  }

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value })
  }

  return (
    <div className="modal" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h2>Add {useISBN ? 'Book by ISBN' : 'Book Manually'}</h2>
          <button onClick={onClose} className="close-btn">&times;</button>
        </div>

        {error && <div className="error">{error}</div>}

        <form onSubmit={handleSubmit}>
          {!linkingToProgressId && (
            <div className="form-group">
              <label className="checkbox-label">
                <input
                  type="checkbox"
                  checked={useISBN}
                  onChange={(e) => setUseISBN(e.target.checked)}
                />
                <span>Auto-fill by ISBN (fetches metadata automatically)</span>
              </label>
            </div>
          )}

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

          {(!useISBN || linkingToProgressId) && (
            <>
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
            </>
          )}

          <div className="modal-actions">
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

export default AddBookModal
