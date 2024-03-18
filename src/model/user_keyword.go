package model

type UserKeyword struct {
	UserID   string `json:"user_id"`
	Keyword  string `json:"keyword"`
	FileName string `json:"file_name"`
}
