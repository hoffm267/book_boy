export interface User {
  id: number
  username: string
  email: string
}

export interface Book {
  id: number
  isbn: string
  title: string
  total_pages: number
}

export interface Audiobook {
  id: number
  title: string
  total_length: string
}

export interface Progress {
  id: number
  user_id: number
  book_id: number | null
  audiobook_id: number | null
  book_page: number | null
  audiobook_time: string | null
}

export interface EnrichedProgress {
  Progress: Progress
  Book: Book | null
  Audiobook: Audiobook | null
}

export interface AuthResponse {
  token: string
  user: User
}

export interface BookFormData {
  isbn: string
  title: string
  total_pages: string | number
}

export interface AudiobookFormData {
  title: string
  total_length: string
}

export interface ProgressFormData {
  book_page: string | number
  audiobook_time: string
}
