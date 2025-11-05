package requests


type CreateSubscriptionRequest struct {
	ServiceName string  `json:"service_name" binding:"required"`
	Price       int `json:"price" binding:"required"`
	UserID      string  `json:"user_id" binding:"required"`
	StartDate   string  `json:"start_date" binding:"required"`
}


type UpdateSubscriptionRequest struct {
	ServiceName string  `json:"service_name"`
	Price       int `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
}




type SumSubscriptionRequest struct {
	ServiceName string  `json:"service_name" binding:"required"`
	UserID      string  `json:"user_id" binding:"required"`
	StartDate   string  `json:"start_date" binding:"required"`
	EndDate     string  `json:"end_date" binding:"required"`
}