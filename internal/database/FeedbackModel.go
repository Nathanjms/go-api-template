package database

import "github.com/jmoiron/sqlx"

type Feedback struct {
	ID          int64  `json:"id"`
	Name        string `json:"name" form:"name"`
	UserID      int64  `json:"userId"`
	Type        string `json:"type" form:"type"`
	Description string `json:"description" form:"description"`
	CreatedAt   string `json:"createdAt"`
}

type FeedbackModel struct {
	*sqlx.DB
}

func (model *FeedbackModel) Save(name string, userID int64, feedbackType string, description string) error {
	_, err := model.DB.Exec("INSERT INTO feedback (name, user_id, type, description) VALUES (?, ?, ?, ?)", name, userID, feedbackType, description)

	return err
}
