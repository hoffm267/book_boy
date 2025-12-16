package infra

type MetadataFetchJob struct {
	BookID int    `json:"book_id"`
	ISBN   string `json:"isbn"`
}
