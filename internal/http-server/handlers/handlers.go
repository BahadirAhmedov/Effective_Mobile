package handlers

import (
	"errors"
	"strconv"

	"log/slog"
	"net/http"

	"github.com/BahadirAhmedov/data-aggregation/internal/domain/models"
	"github.com/BahadirAhmedov/data-aggregation/internal/lib/logger/sl"
	"github.com/BahadirAhmedov/data-aggregation/internal/storage"
	"github.com/BahadirAhmedov/data-aggregation/internal/transport/http/requests"
	"github.com/BahadirAhmedov/data-aggregation/internal/transport/http/responses"
	_ "github.com/BahadirAhmedov/data-aggregation/cmd/data-aggregation/docs"
	"github.com/BahadirAhmedov/data-aggregation/internal/lib/api/httputil"


	"github.com/gin-gonic/gin"
)
type Subscription struct{
	SubscriptionProvider Subscriptioner
}


type Subscriptioner interface{
	Create(requests.CreateSubscriptionRequest) (int64, error)
	Read(Id int64) (models.Subscription, error)
	Update(req requests.UpdateSubscriptionRequest, Id int64) (int64, error)
	Delete(Id int64) (int64, error)
	List() ([]models.Subscription, error)
	Sum(req requests.SumSubscriptionRequest) (int64, error)

}


func New(
	subscriptionCreator Subscriptioner,
) *Subscription {
	return &Subscription{
		SubscriptionProvider: subscriptionCreator,
	}
}

// CreateSubscription godoc
// @Summary      Create subscription
// @Description  create subscription
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        input body requests.CreateSubscriptionRequest true "Subscription Info"
// @Success      201  {object}  responses.CreateSubscriptionResponse
// @Failure      400  {object}  httputil.Response
// @Failure      500  {object}  httputil.Response
// @Router       /subscriptions [post]
func (s *Subscription) CreateSubscription(log *slog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
	const op = "http-server.handlers.CreateSubscription"
		
	var request requests.CreateSubscriptionRequest

	log.With(
		slog.String("op", op),
		slog.String("user_id", request.UserID),
	)

	err := ctx.BindJSON(&request)
	if err != nil {
		log.Error("failed to decode request body", sl.Err(err))

		ctx.JSON(http.StatusBadRequest, httputil.Error("failed to decode request body"))
		
		return 
	}

	log.Info("request body decoded", slog.Any("request", request))

	
	id, err := s.SubscriptionProvider.Create(request)
	if errors.Is(err, storage.ErrSubscriptionExists) {
			log.Error("subscription already exists", sl.Err(err))

			ctx.JSON(http.StatusBadRequest, httputil.Error("subscription already exists"))

			return 
	}
	
	if errors.Is(err, storage.ErrInvalidStartDateFormat) {
			log.Error("invalid start_date format", sl.Err(err))

			ctx.JSON(http.StatusBadRequest, httputil.Error("invalid start_date format"))

			return 
	}


	if err != nil {
		log.Error("failed to save subscription", sl.Err(err))

		ctx.JSON(http.StatusBadRequest, httputil.Error("failed to save subscription"))

		return 
	}
	

	resp := responses.CreateSubscriptionResponse{
		Id:          id,
		ServiceName: request.ServiceName,
		Price:       request.Price,
		UserID:      request.UserID,
		StartDate:   request.StartDate,
	}
	ctx.JSON(http.StatusCreated, resp)
	
	}	
}

// ReadSubscription godoc
// @Summary      Show subscription
// @Description  get subscription by ID
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Subscription ID"
// @Success      200  {object}  models.Subscription
// @Failure      400  {object}  httputil.Response
// @Failure      500  {object}  httputil.Response
// @Router       /subscriptions/{id} [get]
func (s *Subscription) ReadSubscription(log *slog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
	const op = "http-server.handlers.ReadSubscription"
		
	log.With(
		slog.String("op", op),
	)
	
	subscriptionId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {

		ctx.JSON(http.StatusBadRequest, httputil.Error("Could not parse subscription id"))

		return
	}
	
	subscription, err := s.SubscriptionProvider.Read(subscriptionId) 
	if errors.Is(err, storage.ErrSubscriptionNotFound) {
		log.Error("subscription not found", sl.Err(err))

		ctx.JSON(http.StatusBadRequest, httputil.Error("subscription not found"))

		return
	}
	
	if  err != nil {
		log.Error("internal server error", sl.Err(err))

		ctx.JSON(http.StatusInternalServerError, httputil.Error("internal server error"))

		return 
	}

	ctx.JSON(http.StatusOK, subscription)
	}	
}

// ListSubscription godoc
// @Summary      Show subscriptions
// @Description  show list of subscriptions
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Success      200  {object}  []models.Subscription
// @Failure      500  {object}  httputil.Response
// @Router       /subscriptions [get]
func (s *Subscription) ListSubscription(log *slog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
	const op = "http-server.handlers.ListSubscription"
		
	log.With(
		slog.String("op", op),
	)

	subscriptions, err := s.SubscriptionProvider.List() 
	if  err!= nil {
		log.Error("internal server error", sl.Err(err))	
		ctx.JSON(http.StatusInternalServerError, httputil.Error("internal server error"))
		return 
	}

	ctx.JSON(http.StatusOK, subscriptions)
	}	
}

// UpdateSubscription godoc
// @Summary      Update subscription
// @Description  update subscription by id
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Subscription ID"
// @Param        input body requests.UpdateSubscriptionRequest true "subscription info"
// @Success      200  {object}  responses.CreateSubscriptionResponse
// @Failure      400  {object}  httputil.Response
// @Failure      500  {object}  httputil.Response
// @Router       /subscriptions/{id} [put]
func (s *Subscription) UpdateSubscription(log *slog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
	const op = "http-server.handlers.UpdateSubscription"

	var request requests.UpdateSubscriptionRequest

	log.With(
		slog.String("op", op),
	)

	subscriptionId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		
		ctx.JSON(http.StatusBadRequest, httputil.Error("Could not parse subscription id"))

		return
	}

	err = ctx.BindJSON(&request)
	if err != nil {
		log.Error("failed to decode request body", sl.Err(err))

		ctx.JSON(http.StatusBadRequest, httputil.Error("failed to decode request body"))

		return 
	}

	log.Info("request body decoded", slog.Any("request", request))
	
	id, err := s.SubscriptionProvider.Update(request, subscriptionId)
	if errors.Is(err, storage.ErrSubscriptionExists) {

		log.Error("subscription already exists", sl.Err(err))

		ctx.JSON(http.StatusInternalServerError, httputil.Error("subscription already exists"))
	
		return 
	}

	if errors.Is(err, storage.ErrInvalidStartDateFormat) {
		
		log.Error("invalid start_date format", sl.Err(err))

		ctx.JSON(http.StatusBadRequest, httputil.Error("invalid start_date format"))

		return 
	}


	if errors.Is(err, storage.ErrSubscriptionNotFound) {
	
		log.Error("subscription not found", sl.Err(err))
	
		ctx.JSON(http.StatusBadRequest, httputil.Error("subscription not found"))
		
		return 
	}


	if err != nil {
	
		log.Error("failed to update subscription", sl.Err(err))
		
		ctx.JSON(http.StatusInternalServerError, httputil.Error("failed to update subscription"))

		return 
	}

	resp := responses.CreateSubscriptionResponse{
		Id:          id,
		ServiceName: request.ServiceName,
		Price:       request.Price,
		UserID:      request.UserID,
		StartDate:   request.StartDate,
	}

	ctx.JSON(http.StatusOK, resp)

	}	
}

// DeleteSubscription godoc
// @Summary      Delete subscription
// @Description  delete subscription by id
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Subscription ID"
// @Success      200  {object}  responses.DeleteSubscriptionResponse
// @Failure      400  {object}  httputil.Response
// @Failure      500  {object}  httputil.Response
// @Router       /subscriptions/{id} [delete]
func (s *Subscription) DeleteSubscription(log *slog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
	const op = "http-server.handlers.DeleteSubscription"
		
	log.With(
		slog.String("op", op),
	)

	subscriptionId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {

		log.Error("Could not parse subscription id", sl.Err(err))

		ctx.JSON(http.StatusBadRequest, httputil.Error("Could not parse subscription id"))

		return
	}


	id, err := s.SubscriptionProvider.Delete(subscriptionId)
	
	if errors.Is(err, storage.ErrSubscriptionNotFound) {

		log.Error("subscription not found", sl.Err(err))
		
		ctx.JSON(http.StatusInternalServerError, httputil.Error("subscription not found"))

		return
	}

	if err != nil  {
		log.Error("failed to delete subscription", sl.Err(err))

		ctx.JSON(http.StatusInternalServerError, httputil.Error("failed to delete subscription"))

		return

	}
	resp := responses.DeleteSubscriptionResponse{
		Message: "subscription deleted successfully",
		Id: id,
	}

	ctx.JSON(http.StatusOK, resp)
	}	
}



// SumSubscription godoc
// @Summary     Sum subscriptions
// @Description  sum price of all subscriptions for the selected period filtered by user_id and service_name
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        input body requests.SumSubscriptionRequest true "Subscription Info"
// @Success      200  {object}  responses.SumSubscriptionResponse
// @Failure      400  {object}  httputil.Response
// @Failure      500  {object}  httputil.Response
// @Router       /subscriptions/sum [post]
func (s *Subscription) SumSubscriptions(log *slog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
	const op = "http-server.handlers.DeleteSubscription"

	var request requests.SumSubscriptionRequest

	log.With(
		slog.String("op", op),
	)

	err := ctx.BindJSON(&request)
	if err != nil {
		log.Error("failed to decode request body", sl.Err(err))

		ctx.JSON(http.StatusInternalServerError, httputil.Error("failed to decode request body"))

		return 
	}

	totalSum, err := s.SubscriptionProvider.Sum(request)
	if errors.Is(err, storage.ErrUnableToCalculateSum) {
		log.Error("unable to calculate sum", sl.Err(err))

		ctx.JSON(http.StatusInternalServerError, httputil.Error("unable to calculate sum"))

		return
	}

	if errors.Is(err, storage.ErrInvalidStartDateFormat) {
		log.Error("invalid start_date format", sl.Err(err))

		ctx.JSON(http.StatusBadRequest, httputil.Error("invalid start_date format"))

		return 
	}
	
	if errors.Is(err, storage.ErrInvalidEndDateFormat) {
		log.Error("invalid end_date format", sl.Err(err))

		ctx.JSON(http.StatusBadRequest, httputil.Error("invalid end_date format"))

		return 
	}

	resp := responses.SumSubscriptionResponse{
		TotalSum: totalSum,
	}
	ctx.JSON(http.StatusOK, resp)

	}	
}

