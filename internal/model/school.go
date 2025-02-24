package model

type School struct {
	ID   int    `json:"id"`
	Name string `json:"name" binding:"required,min=3"`
}
