package dto

type DetailFollowDto struct {
	Total int         `json:"total,omitempty"`
	Data  []FollowDto `json:"data,omitempty"`
}

type FollowDto struct {
	FollowerId  int    `json:"followerId,omitempty"`
	FollowingId int    `json:"followingId,omitempty"`
	CreatedAt   string `json:"createdAt,omitempty"`
}

type FollowRequestDto struct {
	UserId      int `json:"userId" validate:"required"`
	FollowingId int `json:"followingId" validate:"required"`
}
