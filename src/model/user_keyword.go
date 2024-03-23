package model

type UserKeyword struct {
	UserID    string `json:"user_id" gorm:"primaryKey"`
	KeywordID string `json:"keyword_id" gorm:"primaryKey"`
	FileName  string `json:"file_name"`
}
