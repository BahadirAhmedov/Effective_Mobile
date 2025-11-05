package storage

import(
	"errors"
)

const (
	UniqueViolation = "23505"
)

var (
	ErrSubscriptionExists = errors.New("subscription exists")
	ErrSubscriptionNotFound = errors.New("subscription not found")
	ErrInvalidStartDateFormat = errors.New("invalid start_date format")
	ErrInvalidEndDateFormat = errors.New("invalid end_date format")
	ErrUnableToCalculateSum = errors.New("unable to calculate the total cost of all subscriptions for a selected period")
)