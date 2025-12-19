import { useState, useEffect } from 'react'
import ProgressModal from './ProgressModal'

function Progress({ token, apiUrl, userId }) {
    const [progressList, setProgressList] = useState([])
    const [books, setBooks] = useState([])
    const [audiobooks, setAudiobooks] = useState([])
    const [showModal, setShowModal] = useState(false)
    const [editingProgress, setEditingProgress] = useState(null)

    useEffect(() => {
        fetchProgress()
        fetchBooks()
        fetchAudiobooks()
        console.log('Setting up SSE connection...')
        const eventSource = new EventSource(`${apiUrl}/events?token=${token}`)

        eventSource.onopen = () => {
            console.log('SSE connection opened')
        }

        eventSource.addEventListener('book.metadata_fetched', (e) => {
            console.log('Received book.metadata_fetched event:', e.data)
            const updatedBook = JSON.parse(e.data)
            console.log('Updating book:', updatedBook)
            setBooks(prev => {
                const updated = prev.map(book =>
                    book.id === updatedBook.id ? updatedBook : book
                )
                console.log('Updated books list:', updated)
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
    }, [])

    const fetchProgress = async () => {
        try {
            const response = await fetch(`${apiUrl}/progress/filter?user_id=${userId}`, {
                headers: { 'Authorization': `Bearer ${token}` }
            })
            const result = await response.json()
            const data = result.data || result
            setProgressList(Array.isArray(data) ? data : [])
        } catch (err) {
            console.error('Failed to fetch progress:', err)
            setProgressList([])
        }
    }

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

    const handleDelete = async (id) => {
        if (!confirm('Delete this progress entry?')) return

        try {
            await fetch(`${apiUrl}/progress/${id}`, {
                method: 'DELETE',
                headers: { 'Authorization': `Bearer ${token}` }
            })
            fetchProgress()
        } catch (err) {
            alert('Failed to delete')
        }
    }

    const handleEditClick = (progress) => {
        setEditingProgress(progress)
        setShowModal(true)
    }

    const handleModalSave = () => {
        setShowModal(false)
        fetchProgress()
        fetchBooks()
        fetchAudiobooks()
    }

    const getBook = (bookId) => {
        if (!Array.isArray(books)) return null
        return books.find(b => b.id === bookId)
    }

    const getAudiobook = (audiobookId) => {
        if (!Array.isArray(audiobooks)) return null
        return audiobooks.find(ab => ab.id === audiobookId)
    }

    const calculateProgress = (progress) => {
        if (progress.book_id && progress.book_page && Array.isArray(books)) {
            const book = books.find(b => b.id === progress.book_id)
            if (book && book.total_pages) {
                return Math.round((progress.book_page / book.total_pages) * 100)
            }
        }
        else if (progress.audiobook_id && progress.audiobook_time && Array.isArray(audiobooks)) {
            const audiobook = audiobooks.find(ab => ab.id === progress.audiobook_id)
            if (audiobook && audiobook.total_length) {
                const totalSeconds = timeString =>
                    timeString.split(':').reduce(
                        (sum, part, i) => sum + part * [3600, 60, 1][i],
                        0
                    );
                const percentage = (
                    totalSeconds(progress.audiobook_time) /
                    totalSeconds(audiobook.total_length) * 100
                ).toFixed(2);
                return Math.round(percentage)
            }
        }
        return 0
    }

    return (
        <>
        <div className="card">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
        <h2>Reading Progress</h2>
        <button onClick={() => { setEditingProgress(null); setShowModal(true); }} className="btn btn-primary">
        + Start Tracking
        </button>
        </div>

        {progressList.length === 0 ? (
            <p style={{ color: '#7f8c8d', textAlign: 'center', padding: '40px' }}>
            No progress tracked yet. Click "Start Tracking" to begin!
            </p>
        ) : (
            <div className="grid">
            {progressList.map(progress => {
                const percent = calculateProgress(progress)
                const book = progress.book_id ? getBook(progress.book_id) : null
                const audiobook = progress.audiobook_id ? getAudiobook(progress.audiobook_id) : null
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
                    {progress.book_id && (
                        <>
                        <div className="progress-bar">
                        <div className="progress-fill" style={{ width: `${percent}%` }}></div>
                        </div>
                        <p style={{ textAlign: 'center', fontSize: '12px' }}>{percent}% complete</p>
                        </>
                    )}
                    {progress.audiobook_id && !progress.book_id && (
                        <>
                        <div className="progress-bar">
                        <div className="progress-fill" style={{ width: `${percent}%` }}></div>
                        </div>
                        <p style={{ textAlign: 'center', fontSize: '12px' }}>{percent}% complete</p>
                        </>
                    )}
                    <div className="book-actions">
                    <button onClick={() => handleEditClick(progress)} className="btn btn-primary">
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

        {showModal && (
            <ProgressModal
            progress={editingProgress}
            books={books}
            audiobooks={audiobooks}
            userId={userId}
            token={token}
            apiUrl={apiUrl}
            onClose={() => setShowModal(false)}
            onSave={handleModalSave}
            />
        )}
        </>
    )
}

export default Progress
