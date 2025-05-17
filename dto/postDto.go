package dto

type DetailPostDto struct {
	Total int       `json:"total,omitempty"`
	Data  []PostDto `json:"data,omitempty"`
}

type PostDto struct {
	ID        int64        `json:"id"`
	UserId    int          `json:"user_id"`
	Content   string       `json:"content"`
	ImageUrl  string       `json:"image_url"`
	IsActive  bool         `json:"is_active"`
	CreatedAt string       `json:"created_at"`
	UpdatedAt string       `json:"updated_at"`
	Likes     []LikeDto    `json:"likes"`
	Comments  []CommentDto `json:"comments"`
}

type PostRequestDto struct {
	UserId   int    `json:"userId" form:"userId"`
	Content  string `json:"content" form:"content"`
	Filename string `json:"filename" form:"filename"`
}

type PostDeleteRequestDto struct {
	ID           int  `json:"id" form:"id"`
	UserId       int  `json:"userId" form:"userId"`
	IsDeletePost bool `json:"isDeletePost" form:"isDeletePost"`
}
