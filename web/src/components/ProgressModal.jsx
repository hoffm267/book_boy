import { useState, useEffect } from 'react'
import BookModal from './BookModal'

function ProgressModal({ progress, books, audiobooks, userId, token, apiUrl, onClose, onSave }) {
    const [formData, setFormData] = useState(
        progress || { user_id: userId != null ? userId : 0, book_id: '', audiobook_id: '', book_page: '', audiobook_time: '' }
    )
    const ChangeType = { BOOK: 0, AUDIOBOOK: 1 }
    const [toChange, setToChange] = useState(ChangeType.BOOK)
    const [error, setError] = useState('')
    const [loading, setLoading] = useState(false)
    const [newBook, setNewBook] = useState({ title: '', isbn: '', total_pages: '' })
    const [newAudiobook, setNewAudiobook] = useState({ title: '', total_length: '' })
    const [showNewBookForm, setShowNewBookForm] = useState(false)
    const [showNewAudiobookForm, setShowNewAudiobookForm] = useState(false)
    const [timeInput, setTimeInput] = useState({ hours: '', minutes: '', seconds: '' })
    const [totalLengthInput, setTotalLengthInput] = useState({ hours: '', minutes: '', seconds: '' })
    const [trackingMode, setTrackingMode] = useState(null)
    const [createdProgressId, setCreatedProgressId] = useState(null)

    useEffect(() => {
        if (progress && progress.audiobook_time) {
            const parts = progress.audiobook_time.split(':')
            if (parts.length === 3) {
                setTimeInput({ hours: parts[0], minutes: parts[1], seconds: parts[2] })
            } else if (parts.length === 2) {
                setTimeInput({ hours: parts[0], minutes: parts[1], seconds: '00' })
            }
        }
    }, [progress])

    const getSelectedBookTitle = () => {
        if (formData.book_id && Array.isArray(books)) {
            const book = books.find(b => b.id === parseInt(formData.book_id))
            return book ? book.title : ''
        }
        return ''
    }

    const getSelectedAudiobookTitle = () => {
        if (formData.audiobook_id && Array.isArray(audiobooks)) {
            const audiobook = audiobooks.find(ab => ab.id === parseInt(formData.audiobook_id))
            return audiobook ? audiobook.title : ''
        }
        return ''
    }

    const handleModeSelection = (mode) => {
        setTrackingMode(mode)
        if (mode === 'book' || mode === 'both') {
            const audiobookTitle = getSelectedAudiobookTitle()
            setNewBook({ title: audiobookTitle, isbn: '', total_pages: '' })
            setShowNewBookForm(true)
        }
        if (mode === 'audiobook') {
            const bookTitle = getSelectedBookTitle()
            setNewAudiobook({ title: bookTitle, total_length: '' })
            setShowNewAudiobookForm(true)
        }
    }

    const handleAddBookClick = () => {
        console.log("Come one, please")
        const audiobookTitle = getSelectedAudiobookTitle()
        setNewBook({ title: audiobookTitle, isbn: '', total_pages: '' })
        setShowNewBookForm(true)
    }

    const handleAddAudiobookClick = () => {
        const bookTitle = getSelectedBookTitle()
        setNewAudiobook({ title: bookTitle, total_length: '' })
        setShowNewAudiobookForm(true)
    }

    const handleCreateBook = async () => {
        setError('')
        setLoading(true)

        try {
            if (progress) {
                const bookData = {
                    title: newBook.title,
                    isbn: newBook.isbn,
                    total_pages: parseInt(newBook.total_pages)
                }

                const response = await fetch(`${apiUrl}/books?pgId=${progress.id}`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${token}`
                    },
                    body: JSON.stringify(bookData)
                })

                if (!response.ok) {
                    const data = await response.json()
                    throw new Error(data.error || 'Failed to create book')
                }

                setShowNewBookForm(false)
                setNewBook({ title: '', isbn: '', total_pages: '' })
                onClose()
                onSave()
            } else {
                const trackingData = {
                    format: 'book',
                    title: newBook.title,
                    isbn: newBook.isbn,
                    total_pages: parseInt(newBook.total_pages),
                    current_page: 1
                }

                const response = await fetch(`${apiUrl}/tracking/start`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${token}`
                    },
                    body: JSON.stringify(trackingData)
                })

                if (!response.ok) {
                    const data = await response.json()
                    throw new Error(data.error || 'Failed to start tracking book')
                }

                const result = await response.json()
                setShowNewBookForm(false)
                setNewBook({ title: '', isbn: '', total_pages: '' })

                if (trackingMode === 'both') {
                    setCreatedProgressId(result.id)
                    setNewAudiobook({ title: newBook.title, total_length: '' })
                    setShowNewAudiobookForm(true)
                } else {
                    onClose()
                    onSave()
                }
            }
        } catch (err) {
            setError(err.message)
        } finally {
            setLoading(false)
        }
    }

    const handleCreateAudiobook = async () => {
        setError('')
        setLoading(true)

        try {
            const progressId = progress?.id || createdProgressId

            if (progressId) {
                const response = await fetch(`${apiUrl}/audiobooks?pgId=${progressId}`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${token}`
                    },
                    body: JSON.stringify(newAudiobook)
                })

                if (!response.ok) {
                    const data = await response.json()
                    throw new Error(data.error || 'Failed to create audiobook')
                }

                setShowNewAudiobookForm(false)
                setNewAudiobook({ title: '', total_length: '' })
                onClose()
                onSave()
            } else {
                const trackingData = {
                    format: 'audiobook',
                    title: newAudiobook.title,
                    total_length: newAudiobook.total_length,
                    current_time: '00:00:00'
                }

                const response = await fetch(`${apiUrl}/tracking/start`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${token}`
                    },
                    body: JSON.stringify(trackingData)
                })

                if (!response.ok) {
                    const data = await response.json()
                    throw new Error(data.error || 'Failed to start tracking audiobook')
                }

                setShowNewAudiobookForm(false)
                setNewAudiobook({ title: '', total_length: '' })
                onClose()
                onSave()
            }
        } catch (err) {
            setError(err.message)
        } finally {
            setLoading(false)
        }
    }

    const handleTimeChange = (field, value) => {
        const numValue = value.replace(/\D/g, '')
        let finalValue = numValue

        if (field === 'hours') {
            finalValue = numValue.slice(0, 2)
        } else if (field === 'minutes' || field === 'seconds') {
            finalValue = numValue.slice(0, 2)
            if (parseInt(finalValue) > 59) {
                finalValue = '59'
            }
        }

        const newTimeInput = { ...timeInput, [field]: finalValue }
        setTimeInput(newTimeInput)

        const hours = newTimeInput.hours || '0'
        const minutes = newTimeInput.minutes || '0'
        const seconds = newTimeInput.seconds || '0'
        const timeString = `${hours.padStart(2, '0')}:${minutes.padStart(2, '0')}:${seconds.padStart(2, '0')}`

        setFormData({ ...formData, audiobook_time: timeString })
        setToChange(ChangeType.AUDIOBOOK)
    }

    const handleTotalLengthChange = (field, value) => {
        const numValue = value.replace(/\D/g, '')
        let finalValue = numValue

        if (field === 'hours') {
            finalValue = numValue.slice(0, 2)
        } else if (field === 'minutes' || field === 'seconds') {
            finalValue = numValue.slice(0, 2)
            if (parseInt(finalValue) > 59) {
                finalValue = '59'
            }
        }

        const newTotalLengthInput = { ...totalLengthInput, [field]: finalValue }
        setTotalLengthInput(newTotalLengthInput)

        const hours = newTotalLengthInput.hours || '0'
        const minutes = newTotalLengthInput.minutes || '0'
        const seconds = newTotalLengthInput.seconds || '0'
        const timeString = `${hours.padStart(2, '0')}:${minutes.padStart(2, '0')}:${seconds.padStart(2, '0')}`

        setNewAudiobook({ ...newAudiobook, total_length: timeString })
    }

    const handleSubmit = async (e) => {
        e.preventDefault()
        setError('')
        setLoading(true)

        try {
            const method = progress ? 'PUT' : 'POST'
            const url = progress ? `${apiUrl}/progress/${progress.id}` : `${apiUrl}/progress`

            if (toChange == ChangeType.BOOK) {
                formData.audiobook_id = null
                formData.audiobook_time = null
            }
            else {
                formData.book_id = null
                formData.book_page = null
            }
            const body = {
                user_id: parseInt(formData.user_id),
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

    const handleChange = (e) => {
        if (e.target.name == "book_page") {
            setToChange(ChangeType.BOOK)
        }
        else if (e.target.name == "audiobook_time") {
            setToChange(ChangeType.AUDIOBOOK)
        }
        setFormData({ ...formData, [e.target.name]: e.target.value })
    }

    return (
        <div className="modal" onClick={onClose}>
        <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
        <h2>{progress ? 'Update' : 'Start'} Progress</h2>
        <button onClick={onClose} className="close-btn">&times;</button>
        </div>

        {error && <div className="error">{error}</div>}

        {showNewBookForm && (
            <BookModal
            type="book"
            item={null}
            token={token}
            apiUrl={apiUrl}
            queryParams={progress ? `pgId=${progress.id}` : 'skipProgress=true'}
            onClose={() => {
                setShowNewBookForm(false);
                if (!progress) setTrackingMode(null);
            }}
            onSave={async (createdBook) => {
                setShowNewBookForm(false);

                if (!progress && createdBook?.id) {
                    try {
                        const progressData = {
                            user_id: userId,
                            book_id: createdBook.id,
                            book_page: 1,
                            audiobook_id: null,
                            audiobook_time: null
                        };

                        const response = await fetch(`${apiUrl}/progress`, {
                            method: 'POST',
                            headers: {
                                'Content-Type': 'application/json',
                                'Authorization': `Bearer ${token}`
                            },
                            body: JSON.stringify(progressData)
                        });

                        if (response.ok) {
                            const result = await response.json();
                            if (trackingMode === 'both') {
                                setCreatedProgressId(result.id || result.data?.id);
                                setNewAudiobook({ title: createdBook.title || '', total_length: '' });
                                setShowNewAudiobookForm(true);
                            } else {
                                onClose();
                                onSave();
                            }
                        }
                    } catch (err) {
                        console.error('Failed to create progress:', err);
                    }
                } else {
                    onClose();
                    onSave();
                }
            }}
            />
        )}

        <form onSubmit={handleSubmit}>

        {showNewAudiobookForm && (
            <div className="form-group">
            <label>Add New Audiobook</label>
            <input
            type="text"
            placeholder="Title"
            value={newAudiobook.title}
            onChange={(e) => setNewAudiobook({ ...newAudiobook, title: e.target.value })}
            style={{ marginBottom: '10px' }}
            required
            />
            <label>Total Length</label>
            <div style={{ display: 'flex', gap: '8px', alignItems: 'center', marginBottom: '10px' }}>
            <div style={{ flex: 1 }}>
            <input
            type="text"
            placeholder="HH"
            value={totalLengthInput.hours}
            onChange={(e) => handleTotalLengthChange('hours', e.target.value)}
            style={{ width: '100%', textAlign: 'center' }}
            maxLength="2"
            />
            <div style={{ fontSize: '10px', textAlign: 'center', marginTop: '2px' }}>Hours</div>
            </div>
            <span style={{ fontSize: '20px', fontWeight: 'bold' }}>:</span>
            <div style={{ flex: 1 }}>
            <input
            type="text"
            placeholder="MM"
            value={totalLengthInput.minutes}
            onChange={(e) => handleTotalLengthChange('minutes', e.target.value)}
            style={{ width: '100%', textAlign: 'center' }}
            maxLength="2"
            />
            <div style={{ fontSize: '10px', textAlign: 'center', marginTop: '2px' }}>Minutes</div>
            </div>
            <span style={{ fontSize: '20px', fontWeight: 'bold' }}>:</span>
            <div style={{ flex: 1 }}>
            <input
            type="text"
            placeholder="SS"
            value={totalLengthInput.seconds}
            onChange={(e) => handleTotalLengthChange('seconds', e.target.value)}
            style={{ width: '100%', textAlign: 'center' }}
            maxLength="2"
            />
            <div style={{ fontSize: '10px', textAlign: 'center', marginTop: '2px' }}>Seconds</div>
            </div>
            </div>
            <div style={{ display: 'flex', gap: '5px' }}>
            <button type="button" onClick={handleCreateAudiobook} className="btn btn-primary" disabled={loading}>
            Create Audiobook
            </button>
            <button type="button" onClick={() => { setShowNewAudiobookForm(false); if (!progress && trackingMode !== 'both') setTrackingMode(null); }} className="btn btn-secondary">
            {!progress && trackingMode !== 'both' ? 'Back' : 'Cancel'}
            </button>
            </div>
            </div>
        )}

        {!showNewBookForm && !showNewAudiobookForm && progress && (
            <>
            {formData.book_id && (
                <div className="form-group">
                <label>Current Page</label>
                <input
                type="number"
                name="book_page"
                value={formData.book_page}
                onChange={handleChange}
                min="1"
                />
                </div>
            )}

            {formData.audiobook_id && (
                <div className="form-group">
                <label>Current Time</label>
                <div style={{ display: 'flex', gap: '8px', alignItems: 'center' }}>
                <div style={{ flex: 1 }}>
                <input
                type="text"
                placeholder="HH"
                value={timeInput.hours}
                onChange={(e) => handleTimeChange('hours', e.target.value)}
                style={{ width: '100%', textAlign: 'center' }}
                maxLength="2"
                />
                <div style={{ fontSize: '10px', textAlign: 'center', marginTop: '2px' }}>Hours</div>
                </div>
                <span style={{ fontSize: '20px', fontWeight: 'bold' }}>:</span>
                <div style={{ flex: 1 }}>
                <input
                type="text"
                placeholder="MM"
                value={timeInput.minutes}
                onChange={(e) => handleTimeChange('minutes', e.target.value)}
                style={{ width: '100%', textAlign: 'center' }}
                maxLength="2"
                />
                <div style={{ fontSize: '10px', textAlign: 'center', marginTop: '2px' }}>Minutes</div>
                </div>
                <span style={{ fontSize: '20px', fontWeight: 'bold' }}>:</span>
                <div style={{ flex: 1 }}>
                <input
                type="text"
                placeholder="SS"
                value={timeInput.seconds}
                onChange={(e) => handleTimeChange('seconds', e.target.value)}
                style={{ width: '100%', textAlign: 'center' }}
                maxLength="2"
                />
                <div style={{ fontSize: '10px', textAlign: 'center', marginTop: '2px' }}>Seconds</div>
                </div>
                </div>
                </div>
            )}
            </>
        )}

        {!showNewBookForm && !showNewAudiobookForm && progress && (
            <>
            {!formData.book_id && !formData.audiobook_id && (
                <div style={{ display: 'flex', gap: '10px', marginBottom: '15px' }}>
                <button type="button" onClick={handleAddBookClick} className="btn btn-primary">
                + Track as Book
                </button>
                <button type="button" onClick={handleAddAudiobookClick} className="btn btn-primary">
                + Track as Audiobook
                </button>
                </div>
            )}

            {formData.book_id && !formData.audiobook_id && (
                <button type="button" onClick={handleAddAudiobookClick} className="btn btn-secondary" style={{ marginBottom: '15px' }}>
                + Also Track as Audiobook
                </button>
            )}

            {formData.audiobook_id && !formData.book_id && (
                <button type="button" onClick={handleAddBookClick} className="btn btn-secondary" style={{ marginBottom: '15px' }}>
                + Also Track as Book
                </button>
            )}
            </>
        )}

        {!showNewBookForm && !showNewAudiobookForm && !progress && !trackingMode && (
            <div>
            <p style={{ marginBottom: '15px', textAlign: 'center', color: '#555' }}>
            What would you like to track?
            </p>
            <div style={{ display: 'flex', flexDirection: 'column', gap: '10px' }}>
            <button
            type="button"
            onClick={() => handleModeSelection('book')}
            className="btn btn-primary"
            style={{ padding: '15px', fontSize: '16px' }}
            >
            📖 Book Only
            </button>
            <button
            type="button"
            onClick={() => handleModeSelection('audiobook')}
            className="btn btn-primary"
            style={{ padding: '15px', fontSize: '16px' }}
            >
            🎧 Audiobook Only
            </button>
            <button
            type="button"
            onClick={() => handleModeSelection('both')}
            className="btn btn-primary"
            style={{ padding: '15px', fontSize: '16px' }}
            >
            📖 + 🎧 Both (Same Title)
            </button>
            </div>
            </div>
        )}

        {!showNewBookForm && !showNewAudiobookForm && (formData.book_id || formData.audiobook_id) && (
            <div style={{ display: 'flex', gap: '10px', justifyContent: 'flex-end' }}>
            <button type="button" onClick={onClose} className="btn btn-secondary">
            Cancel
            </button>
            <button type="submit" className="btn btn-primary" disabled={loading}>
            {loading ? 'Saving...' : 'Save'}
            </button>
            </div>
        )}
        </form>
        </div>
        </div>
    )
}

export default ProgressModal
