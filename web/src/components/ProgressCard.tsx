import type { EnrichedProgress } from '../types'

interface ProgressCardProps {
  enrichedProgress: EnrichedProgress
  onEdit: () => void
  onDelete: () => void
}

function calculateProgress(enrichedProgress: EnrichedProgress): number {
  const progress = enrichedProgress.Progress
  const book = enrichedProgress.Book
  const audiobook = enrichedProgress.Audiobook

  if (progress.book_id && progress.book_page && book && book.total_pages) {
    return Math.round((progress.book_page / book.total_pages) * 100)
  }

  if (progress.audiobook_id && progress.audiobook_time && audiobook && audiobook.total_length) {
    const totalSeconds = (timeString: string): number =>
      timeString.split(':').reduce(
        (sum, part, i) => sum + Number(part) * [3600, 60, 1][i],
        0
      )
    return Math.round(
      totalSeconds(progress.audiobook_time) /
      totalSeconds(audiobook.total_length) * 100
    )
  }

  return 0
}

function ProgressCard({ enrichedProgress, onEdit, onDelete }: ProgressCardProps) {
  const percent = calculateProgress(enrichedProgress)
  const { Progress: progress, Book: book, Audiobook: audiobook } = enrichedProgress
  const title = book?.title || audiobook?.title || 'Unknown'

  return (
    <div className="book-card">
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
      <p className="progress-percent">{percent}% complete</p>
      <div className="book-actions">
        <button onClick={onEdit} className="btn btn-primary">
          Update
        </button>
        <button onClick={onDelete} className="btn btn-danger">
          Delete
        </button>
      </div>
    </div>
  )
}

export default ProgressCard
