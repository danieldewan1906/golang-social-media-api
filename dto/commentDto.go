package dto

type CommentRequestDto struct {
	UserId      int    `json:"userId"`
	PostId      int    `json:"postId"`
	TextComment string `json:"textComment"`
}

type CommentDto struct {
	ID          int    `json:"id"`
	UserId      int    `json:"userId"`
	PostId      int    `json:"postId"`
	TextComment string `json:"textComment"`
	CreatedAt   string `json:"createdAt"`
}
