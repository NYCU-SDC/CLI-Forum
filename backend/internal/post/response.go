package post

type Response struct {
	ID       string `json:"id"`
	AuthorID string `json:"author_id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	CreateAt string `json:"create_at"`
}
