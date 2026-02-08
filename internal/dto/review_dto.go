package dto

type ReviewCreateRequest struct {
	Text   string `json:"text" binding:"required,min=3"`
	Rating int    `json:"rating" binding:"required,min=1,max=5"`
}
type ReviewUpdateRequest struct {
	Text   *string `json:"text"`
	Rating *int    `json:"rating"`
}

type ReviewListItem struct {
	ID       uint
	AuthorID uint
	TripID   uint
	Rating   int
	Text     string
}
