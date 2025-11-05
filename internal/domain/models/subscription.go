package models


type Subscription struct {
	Id int64 `json:"id"`
	ServiceName string  `json:"service_name"`
	Price       int `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
}
