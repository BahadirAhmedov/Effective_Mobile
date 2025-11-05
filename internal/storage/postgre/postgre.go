package postgre

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/BahadirAhmedov/data-aggregation/internal/domain/models"
	"github.com/BahadirAhmedov/data-aggregation/internal/storage"
	"github.com/BahadirAhmedov/data-aggregation/internal/transport/http/requests"
	"github.com/lib/pq"
)


type Storage struct{
	db *sql.DB
}

func New(host string, port int, user string, password string, dbname string)(*Storage){

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)
	
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
  		panic(err)
	}  	
	
	err = db.Ping()
  	if err != nil {
    	panic(err)
  	}

  	fmt.Println("Successfully connected!")

	return &Storage{db: db}
}


func (s *Storage) Create(req requests.CreateSubscriptionRequest)(int64, error){
	const op = "storage.postgre.Create"

	layout := "01-2006"
	parsed, err := time.Parse(layout, req.StartDate)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrInvalidStartDateFormat)
	}

	var id int64

	err = s.db.QueryRow("INSERT INTO subscriptions(serviceName, price, userId, startDate) VALUES($1, $2, $3, $4)  RETURNING id", req.ServiceName, req.Price, req.UserID, parsed).Scan(&id)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == storage.UniqueViolation{
			return 0, fmt.Errorf("%s: %w", op, storage.ErrSubscriptionExists)			
		}		
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}


func (s *Storage) Read(id int64) (models.Subscription, error){
	const op = "storage.postgre.Read"

	row := s.db.QueryRow("SELECT id, serviceName, price, userId, TO_CHAR(startDate, 'MM-YYYY') AS startDate FROM subscriptions WHERE id = $1", id)


	var subscription models.Subscription

	err := row.Scan(&subscription.Id ,&subscription.ServiceName, &subscription.Price, &subscription.UserID, &subscription.StartDate)
	

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println(err)
			return models.Subscription{}, fmt.Errorf("%s: %w", op, storage.ErrSubscriptionNotFound)
		}
		return models.Subscription{}, fmt.Errorf("%s: %w", op, err)
	}

	return subscription, nil
}




func (s *Storage) List() ([]models.Subscription, error){
	const op = "storage.postgre.List"

	rows, err  := s.db.Query("SELECT id, serviceName, price, userId, TO_CHAR(startDate, 'MM-YYYY') AS startDate FROM subscriptions")
	if err != nil {
		return []models.Subscription{}, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var subscriptions []models.Subscription

	for rows.Next() {
		var subscription models.Subscription
		err := rows.Scan(&subscription.Id ,&subscription.ServiceName, &subscription.Price, &subscription.UserID, &subscription.StartDate)
		if err != nil {
			return []models.Subscription{}, fmt.Errorf("%s: %w", op, err)
		}
		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, nil
}


func (s *Storage) Update(req requests.UpdateSubscriptionRequest, Id int64) (int64, error){
	const op = "storage.postgre.Update"

	layout := "01-2006"
	parsed, err := time.Parse(layout, req.StartDate)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrInvalidStartDateFormat)
	}

	var id int64

	err = s.db.QueryRow("UPDATE subscriptions SET serviceName = $1, price = $2, userId = $3, startDate = $4 WHERE id = $5 RETURNING id", req.ServiceName, req.Price, req.UserID, parsed, Id).Scan(&id)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == storage.UniqueViolation{
			return 0, fmt.Errorf("%s: %w", op, storage.ErrSubscriptionExists)			
		}

		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrSubscriptionNotFound)			
		}
	
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}



func (s *Storage) Delete(Id int64) (int64, error){
	const op = "storage.postgre.Delete"

	var id int64 

	err := s.db.QueryRow("DELETE FROM subscriptions WHERE id = $1 RETURNING id", Id).Scan(&id)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrSubscriptionNotFound)			
		}		
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}


func (s *Storage) Sum(req requests.SumSubscriptionRequest) (int64, error){
	const op = "storage.postgre.Sum"

	layout := "01-2006"
	
	startDateParsed, err := time.Parse(layout, req.StartDate)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrInvalidStartDateFormat)
	}

	endDateParsed, err := time.Parse(layout, req.EndDate)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrInvalidEndDateFormat)
	}

	var totalSum int64
	
	err = s.db.QueryRow(
	   `SELECT SUM(price) 
		FROM subscriptions 
		WHERE userId = $1 
			AND serviceName = $2 
			AND startDate BETWEEN $3 AND $4`, req.UserID, req.ServiceName, startDateParsed, endDateParsed).Scan(&totalSum)
	
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrUnableToCalculateSum)
	}

	return totalSum, nil
}