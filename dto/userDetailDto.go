package dto

type UserDetailDto struct {
	ID        int             `json:"id"`
	FirstName string          `json:"firstName"`
	LastName  string          `json:"lastName,omitempty"`
	Address   string          `json:"address,omitempty"`
	CreatedAt string          `json:"createdAt"`
	UserImage *UserImagesDto  `json:"userImage,omitempty"`
	Posts     DetailPostDto   `json:"posts,omitempty"`
	Followers DetailFollowDto `json:"followers,omitempty"`
	Following DetailFollowDto `json:"following,omitempty"`
}

type UserDetailRequestDto struct {
	ID   int    `json:"id" param:"id" query:"id"`
	Name string `json:"name" param:"name" query:"name"`
}

type UpdateUserRequestDto struct {
	UserId    int    `json:"userId" validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName"`
	Address   string `json:"Address"`
}
