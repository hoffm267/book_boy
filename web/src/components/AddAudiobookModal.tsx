import { useState, FormEvent, ChangeEvent } from 'react'
import { useAuth } from '../AuthContext'
import type { AudiobookFormData } from '../types'

interface AddAudiobookModalProps {
  onClose: () => void
  onSave: () => void
  linkingToProgressId: number | null
  prefillTitle: string | null
}

function AddAudiobookModal({ onClose, onSave, linkingToProgressId, prefillTitle }: AddAudiobookModalProps) {
  const { api } = useAuth()
  const [formData, setFormData] = useState<AudiobookFormData>({
    title: prefillTitle || '',
    total_length: '',
  })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      const path = linkingToProgressId ? `/audiobooks?pgId=${linkingToProgressId}` : '/audiobooks'
      await api.post(path, formData)
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
          <h2>Add Audiobook</h2>
          <button onClick={onClose} className="close-btn">&times;</button>
        </div>

        {error && <div className="error">{error}</div>}

        <form onSubmit={handleSubmit}>
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

export default AddAudiobookModal
