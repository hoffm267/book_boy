import { useState, FormEvent } from 'react'
import type { Progress, Book, Audiobook, ProgressFormData } from '../types'

interface ProgressEditModalProps {
  progress: Progress
  book: Book | null
  audiobook: Audiobook | null
  onClose: () => void
  onSave: (formData: ProgressFormData) => void
  onAddBook: () => void
  onAddAudiobook: () => void
}

function ProgressEditModal({ progress, book, audiobook, onClose, onSave, onAddBook, onAddAudiobook }: ProgressEditModalProps) {
  const [formData, setFormData] = useState<ProgressFormData>({
    book_page: progress.book_page || '',
    audiobook_time: progress.audiobook_time || ''
  })

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    onSave(formData)
  }

  return (
    <div className="modal" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h2>Update Progress</h2>
          <button onClick={onClose} className="close-btn">&times;</button>
        </div>

        <form onSubmit={handleSubmit}>
          {progress.book_id && book && (
            <div className="form-group">
              <label>Page ({book.total_pages} total)</label>
              <input
                type="number"
                value={formData.book_page}
                onChange={(e) => setFormData({...formData, book_page: e.target.value})}
                min="1"
                max={book.total_pages}
              />
            </div>
          )}

          {progress.audiobook_id && audiobook && (
            <div className="form-group">
              <label>Time (HH:MM:SS)</label>
              <input
                type="text"
                value={formData.audiobook_time}
                onChange={(e) => setFormData({...formData, audiobook_time: e.target.value})}
                pattern="[0-9]{2}:[0-9]{2}:[0-9]{2}"
                placeholder="00:00:00"
              />
            </div>
          )}

          {progress.book_id && !progress.audiobook_id && (
            <div className="form-group">
              <button
                type="button"
                onClick={onAddAudiobook}
                className="btn"
                style={{ width: '100%' }}
              >
                + Add Audiobook
              </button>
            </div>
          )}

          {progress.audiobook_id && !progress.book_id && (
            <div className="form-group">
              <button
                type="button"
                onClick={onAddBook}
                className="btn"
                style={{ width: '100%' }}
              >
                + Add Book
              </button>
            </div>
          )}

          <div style={{ display: 'flex', gap: '10px', justifyContent: 'flex-end' }}>
            <button type="button" onClick={onClose} className="btn btn-secondary">
              Cancel
            </button>
            <button type="submit" className="btn btn-primary">
              Save
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

export default ProgressEditModal
