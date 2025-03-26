package post

type CreateRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
