import { useState, useEffect } from 'react'
import BookModal from './BookModal'
import ProgressEditModal from './ProgressEditModal'
import { fetchWithAuth } from '../utils/api'
import type { Progress as ProgressType, EnrichedProgress, Book, ProgressFormData } from '../types'

interface ProgressProps {
  token: string
  apiUrl: string
  userId: number
  onLogout: () => void
}

function Progress({ token, apiUrl, userId, onLogout }: ProgressProps) {
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

    console.log('Setting up SSE connection...')
    const eventSource = new EventSource(`${apiUrl}/events?token=${token}`)

    eventSource.onopen = () => {
      console.log('SSE connection opened')
    }

    eventSource.addEventListener('book.metadata_fetched', (e: MessageEvent) => {
      console.log('Received book.metadata_fetched event:', e.data)
      const updatedBook: Book = JSON.parse(e.data)
      console.log('Updating book:', updatedBook)
      setProgressList(prev => {
        const updated = prev.map(prog => {
          if (prog.Book && prog.Book.id === updatedBook.id) {
            return { ...prog, Book: updatedBook }
          }
          return prog
        })
        console.log('Updated progress list:', updated)
        return updated
      })
    })

    eventSource.onerror = (err) => {
      console.error('SSE error:', err)
      console.error('EventSource readyState:', eventSource.readyState)
    }

    return () => {
      console.log('Closing SSE connection')
      eventSource.close()
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const fetchEnrichedProgress = async () => {
    try {
      const response = await fetchWithAuth(
        `${apiUrl}/progress/enriched`,
        { headers: { 'Authorization': `Bearer ${token}` } },
        onLogout
      )
      const data: EnrichedProgress[] = await response.json()
      const sortedData = Array.isArray(data) ? data.sort((a, b) => a.Progress.id - b.Progress.id) : []
      setProgressList(sortedData)
    } catch (err) {
      if (err instanceof Error && err.name !== 'UnauthorizedError') {
        console.error('Failed to fetch enriched progress:', err)
        setProgressList([])
      }
    }
  }

  const handleDelete = async (id: number) => {
    if (!confirm('Delete this progress entry?')) return

    try {
      await fetchWithAuth(
        `${apiUrl}/progress/${id}`,
        {
          method: 'DELETE',
          headers: { 'Authorization': `Bearer ${token}` }
        },
        onLogout
      )
      fetchEnrichedProgress()
    } catch (err) {
      if (err instanceof Error && err.name !== 'UnauthorizedError') {
        alert('Failed to delete')
      }
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
        await fetchWithAuth(
          `${apiUrl}/progress/${editingProgress.id}/page`,
          {
            method: 'PATCH',
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ page: parseInt(String(formData.book_page)) })
          },
          onLogout
        )
      }

      if (formData.audiobook_time && formData.audiobook_time !== originalValues.audiobook_time) {
        await fetchWithAuth(
          `${apiUrl}/progress/${editingProgress.id}/time`,
          {
            method: 'PATCH',
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ audiobook_time: formData.audiobook_time })
          },
          onLogout
        )
      }

      setEditingProgress(null)
      fetchEnrichedProgress()
    } catch (err) {
      if (err instanceof Error && err.name !== 'UnauthorizedError') {
        alert('Failed to update progress')
      }
    }
  }

  const calculateProgress = (enrichedProgress: EnrichedProgress): number => {
    const progress = enrichedProgress.Progress
    const book = enrichedProgress.Book
    const audiobook = enrichedProgress.Audiobook

    if (progress.book_id && progress.book_page && book && book.total_pages) {
      return Math.round((progress.book_page / book.total_pages) * 100)
    }
    else if (progress.audiobook_id && progress.audiobook_time && audiobook && audiobook.total_length) {
      const totalSeconds = (timeString: string): number =>
        timeString.split(':').reduce(
          (sum, part, i) => sum + Number(part) * [3600, 60, 1][i],
          0
        )
      const percentage = (
        totalSeconds(progress.audiobook_time) /
        totalSeconds(audiobook.total_length) * 100
      )
      return Math.round(percentage)
    }
    return 0
  }

  return (
    <>
      <div className="card">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
          <h2>Reading Progress</h2>
          <div style={{ display: 'flex', gap: '10px' }}>
            <button onClick={() => setShowBookModal(true)} className="btn btn-primary">+ Book</button>
            <button onClick={() => setShowAudiobookModal(true)} className="btn btn-primary">+ Audiobook</button>
          </div>
        </div>

        {progressList.length === 0 ? (
          <p style={{ color: '#7f8c8d', textAlign: 'center', padding: '40px' }}>
            No progress tracked yet. Click "Start Tracking" to begin!
          </p>
        ) : (
          <div className="grid">
            {progressList.map(enrichedProgress => {
              const percent = calculateProgress(enrichedProgress)
              const { Progress: progress, Book: book, Audiobook: audiobook } = enrichedProgress
              const title = book?.title || audiobook?.title || 'Unknown'

              return (
                <div key={progress.id} className="book-card">
                  <h3>{title}</h3>
                  {progress.book_id && book && (
                    <p>Page: {progress.book_page} / {book.total_pages}</p>
                  )}
                  {progress.audiobook_id && audiobook && (
                    <p>Time: {progress.audiobook_time} / {audiobook.total_length}</p>
                  )}
                  <div className="progress-bar">
                    <div className="progress-fill" style={{ width: `${percent}%` }}></div>
                  </div>
                  <p style={{ textAlign: 'center', fontSize: '12px' }}>{percent}% complete</p>
                  <div className="book-actions">
                    <button onClick={() => handleEditClick(enrichedProgress)} className="btn btn-primary">
                      Update
                    </button>
                    <button onClick={() => handleDelete(progress.id)} className="btn btn-danger">
                      Delete
                    </button>
                  </div>
                </div>
              )
            })}
          </div>
        )}
      </div>

      {editingProgress && (() => {
        const enriched = progressList.find(p => p.Progress.id === editingProgress.id)
        const book = enriched?.Book || null
        const audiobook = enriched?.Audiobook || null

        return (
          <ProgressEditModal
            progress={editingProgress}
            book={book}
            audiobook={audiobook}
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
        )
      })()}

      {showBookModal && (() => {
        const enriched = progressList.find(p => p.Progress.id === linkingProgressId)
        const prefillTitle = enriched?.Audiobook?.title || ''

        return (
          <BookModal
            type="book"
            token={token}
            apiUrl={apiUrl}
            userId={userId}
            linkingToProgressId={linkingProgressId}
            prefillTitle={linkingProgressId ? prefillTitle : null}
            onClose={() => { setShowBookModal(false); setLinkingProgressId(null) }}
            onSave={() => {
              setShowBookModal(false)
              setLinkingProgressId(null)
              setEditingProgress(null)
              fetchEnrichedProgress()
            }}
            onLogout={onLogout}
          />
        )
      })()}

      {showAudiobookModal && (() => {
        const enriched = progressList.find(p => p.Progress.id === linkingProgressId)
        const prefillTitle = enriched?.Book?.title || ''

        return (
          <BookModal
            type="audiobook"
            token={token}
            apiUrl={apiUrl}
            userId={userId}
            linkingToProgressId={linkingProgressId}
            prefillTitle={linkingProgressId ? prefillTitle : null}
            onClose={() => { setShowAudiobookModal(false); setLinkingProgressId(null) }}
            onSave={() => {
              setShowAudiobookModal(false)
              setLinkingProgressId(null)
              setEditingProgress(null)
              fetchEnrichedProgress()
            }}
            onLogout={onLogout}
          />
        )
      })()}
    </>
  )
}

export default Progress
