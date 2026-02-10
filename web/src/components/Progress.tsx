import { useState, useEffect } from 'react'
import { fetchEventSource } from '@microsoft/fetch-event-source'
import AddBookModal from './AddBookModal'
import AddAudiobookModal from './AddAudiobookModal'
import ProgressCard from './ProgressCard'
import ProgressEditModal from './ProgressEditModal'
import { useAuth } from '../AuthContext'
import type { Progress as ProgressType, EnrichedProgress, Book, ProgressFormData } from '../types'

function Progress() {
  const { token, apiUrl, api, onLogout } = useAuth()
  const [progressList, setProgressList] = useState<EnrichedProgress[]>([])
  const [editingProgress, setEditingProgress] = useState<ProgressType | null>(null)
  const [originalValues, setOriginalValues] = useState<{ book_page: string | number; audiobook_time: string }>({
    book_page: '',
    audiobook_time: ''
  })
  const [showBookModal, setShowBookModal] = useState(false)
  const [showAudiobookModal, setShowAudiobookModal] = useState(false)
  const [linkingProgressId, setLinkingProgressId] = useState<number | null>(null)

  useEffect(() => {
    fetchEnrichedProgress()

    const controller = new AbortController()

    fetchEventSource(`${apiUrl}/events?token=${token}`, {
      signal: controller.signal,
      onmessage(event) {
        if (event.event === 'book.metadata_fetched' && event.data) {
          const updatedBook: Book = JSON.parse(event.data)
          setProgressList(prev =>
            prev.map(prog =>
              prog.Book && prog.Book.id === updatedBook.id
                ? { ...prog, Book: updatedBook }
                : prog
            )
          )
        }
      },
      onerror() {
        onLogout()
      },
    })

    return () => {
      controller.abort()
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const fetchEnrichedProgress = async () => {
    try {
      const data = await api.get<EnrichedProgress[]>('/progress/enriched')
      const sortedData = Array.isArray(data) ? data.sort((a, b) => a.Progress.id - b.Progress.id) : []
      setProgressList(sortedData)
    } catch {
      setProgressList([])
    }
  }

  const handleDelete = async (id: number) => {
    if (!confirm('Delete this progress entry?')) return

    try {
      await api.delete(`/progress/${id}`)
      fetchEnrichedProgress()
    } catch {
      alert('Failed to delete')
    }
  }

  const handleEditClick = (enrichedProgress: EnrichedProgress) => {
    setEditingProgress(enrichedProgress.Progress)
    setOriginalValues({
      book_page: enrichedProgress.Progress.book_page || '',
      audiobook_time: enrichedProgress.Progress.audiobook_time || ''
    })
  }

  const handleSaveEdit = async (formData: ProgressFormData) => {
    if (!editingProgress) return

    try {
      if (formData.book_page && formData.book_page !== originalValues.book_page) {
        await api.patch(`/progress/${editingProgress.id}/page`, {
          page: parseInt(String(formData.book_page))
        })
      }

      if (formData.audiobook_time && formData.audiobook_time !== originalValues.audiobook_time) {
        await api.patch(`/progress/${editingProgress.id}/time`, {
          audiobook_time: formData.audiobook_time
        })
      }

      setEditingProgress(null)
      fetchEnrichedProgress()
    } catch {
      alert('Failed to update progress')
    }
  }

  const editingEnriched = editingProgress
    ? progressList.find(p => p.Progress.id === editingProgress.id)
    : undefined

  const linkingEnriched = linkingProgressId
    ? progressList.find(p => p.Progress.id === linkingProgressId)
    : undefined

  const handleModalSave = () => {
    setShowBookModal(false)
    setShowAudiobookModal(false)
    setLinkingProgressId(null)
    setEditingProgress(null)
    fetchEnrichedProgress()
  }

  return (
    <>
      <div className="card">
        <div className="card-header">
          <h2>Reading Progress</h2>
          <div className="button-group">
            <button onClick={() => setShowBookModal(true)} className="btn btn-primary">+ Book</button>
            <button onClick={() => setShowAudiobookModal(true)} className="btn btn-primary">+ Audiobook</button>
          </div>
        </div>

        {progressList.length === 0 ? (
          <p className="empty-state">
            No progress tracked yet. Click "Start Tracking" to begin!
          </p>
        ) : (
          <div className="grid">
            {progressList.map(enrichedProgress => (
              <ProgressCard
                key={enrichedProgress.Progress.id}
                enrichedProgress={enrichedProgress}
                onEdit={() => handleEditClick(enrichedProgress)}
                onDelete={() => handleDelete(enrichedProgress.Progress.id)}
              />
            ))}
          </div>
        )}
      </div>

      {editingProgress && (
        <ProgressEditModal
          progress={editingProgress}
          book={editingEnriched?.Book || null}
          audiobook={editingEnriched?.Audiobook || null}
          onClose={() => setEditingProgress(null)}
          onSave={handleSaveEdit}
          onAddBook={() => {
            setLinkingProgressId(editingProgress.id)
            setShowBookModal(true)
          }}
          onAddAudiobook={() => {
            setLinkingProgressId(editingProgress.id)
            setShowAudiobookModal(true)
          }}
        />
      )}

      {showBookModal && (
        <AddBookModal
          linkingToProgressId={linkingProgressId}
          prefillTitle={linkingProgressId ? linkingEnriched?.Audiobook?.title || '' : null}
          onClose={() => { setShowBookModal(false); setLinkingProgressId(null) }}
          onSave={handleModalSave}
        />
      )}

      {showAudiobookModal && (
        <AddAudiobookModal
          linkingToProgressId={linkingProgressId}
          prefillTitle={linkingProgressId ? linkingEnriched?.Book?.title || '' : null}
          onClose={() => { setShowAudiobookModal(false); setLinkingProgressId(null) }}
          onSave={handleModalSave}
        />
      )}
    </>
  )
}

export default Progress
