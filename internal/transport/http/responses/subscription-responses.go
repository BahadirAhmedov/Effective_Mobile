package responses

type CreateSubscriptionResponse struct {
	Id int64 `json:"id"`
	ServiceName string  `json:"service_name"`
	Price       int `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
}


type UpdateSubscriptionResponse struct {
	Id int64 `json:"id"`
	ServiceName string  `json:"service_name"`
	Price       int `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
}


type DeleteSubscriptionResponse struct {
	Message string `json:"message"`
	Id int64  `json:"id"`
}


type SumSubscriptionResponse struct {
	TotalSum int64 `json:"total_sum"`
}
