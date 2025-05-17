package dto

type UserImagesDto struct {
	ID        int64  `json:"id,omitempty"`
	UserId    int    `json:"userId,omitempty"`
	ImageUrl  string `json:"imageUrl,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
	Extension string `json:"extension,omitempty"`
	BaseUrl   string `json:"baseUrl,omitempty"`
}

type UserImageRequestDto struct {
	Filename  string
	Extension string
}
