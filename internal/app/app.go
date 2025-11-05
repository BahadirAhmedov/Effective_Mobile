package app

import (
	"log/slog"

	"github.com/BahadirAhmedov/data-aggregation/internal/http-server/handlers"
	"github.com/BahadirAhmedov/data-aggregation/internal/storage/postgre"
)


func New(
	log *slog.Logger,
	host string, 
	port int, 
	user string, 
	password string, 
	dbname string,
) (*handlers.Subscription) {
	storage := postgre.New(host, port, user, password, dbname)
	
	SubscriptionHandlers := handlers.New(storage)
	return SubscriptionHandlers
}