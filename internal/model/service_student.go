package model

import "time"

type ServiceStudent struct {
	ID          int       `json:"id,omitempty"`
	StudentID   int       `json:"student_id" binding:"required"`
	UserID      int       `json:"user_id" binding:"required"`
	DateService time.Time `json:"date_service" validate:"ISO8601date"`
	Reason      string    `json:"reason" binding:"required"`
	File        string    `json:"file"`
}

type ServiceStudentUpdate struct {
	ID        int    `json:"id,omitempty"`
	StudentID int    `json:"student_id" `
	UserID    int    `json:"user_id" `
	Reason    string `json:"reason" binding:"required"`
	File      string `json:"file"`
}

type ServiceStudentResponse struct {
	ID          int                    `json:"id"`
	Student     StudentServiceResponse `json:"student"`
	User        UserServiceResponse    `json:"user"`
	DateService time.Time              `json:"date_service"`
	Reason      string                 `json:"reason"`
	File        string                 `json:"file"`
}
