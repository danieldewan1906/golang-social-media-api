package dto

type LikeDto struct {
	UserId    int    `json:"userId"`
	CreatedAt string `json:"createdAt"`
}

type LikeRequestDto struct {
	UserId int `json:"userId" validate:"required"`
	PostId int `json:"postId" validate:"required"`
}
