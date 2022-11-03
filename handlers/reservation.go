package handlers

import (
	"avito_test_fix/connections"
	"avito_test_fix/operations"
	"avito_test_fix/utilities"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type ReservationRequest struct {
	UserId    int `json:"user_id"`
	ServiceId int `json:"service_id"`
	OrderId   int `json:"order_id"`

	Amount int `json:"amount"`
}

type ReservationResponse struct {
	Status     string `json:"status"`
	NewBalance int    `json:"new_balance"`
	Details    string `json:"details,omitempty"`
}

func Reservation(c *fiber.Ctx) error {
	var request ReservationRequest
	var response ReservationResponse
	err := c.BodyParser(&request)
	if err != nil {
		return c.JSON(utilities.RequestErrorInputData())
	}

	var userBalance int
	err = connections.PostgresConn.QueryRow(context.Background(),
		"select b.balance from balance b where user_id = $1;", request.UserId).Scan(&userBalance)
	if err != nil {
		return c.JSON(utilities.RequestErrorQueryPg())
	}

	if userBalance-request.Amount < 0 {
		return c.JSON(map[string]string{
			"status": "err",
			"detail": "Insufficient user balance",
		})
	}

	operationId, err := operations.AddOperation(request.UserId, request.Amount,
		operations.OperationStatusAwaitingCompletion, operations.OperationTypeOrder)
	if err != nil {
		return c.JSON(utilities.RequestErrorInputData())
	}

	_, err = connections.PostgresConn.Exec(context.Background(),
		"insert into deals (operation_id, service_id, order_id) values ($1, $2, $3)",
		operationId, request.ServiceId, request.OrderId)
	if err != nil {
		fmt.Printf("\nerr insert into deals (operation_id, service_id, order_id) values ($1, $2, $3)\nERR: %v", err)
		return c.JSON(utilities.RequestErrorInputData())
	}

	// Списание средств с кошелька
	err = connections.PostgresConn.QueryRow(context.Background(),
		"update balance b set balance = (b.balance - $1) where user_id = $2 returning balance;",
		request.Amount, request.UserId).Scan(&response.NewBalance)
	if err != nil {
		return c.JSON(utilities.RequestErrorQueryPg())
	}

	response.Status = "success"
	return c.JSON(response)
}
