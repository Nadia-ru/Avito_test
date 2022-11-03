package handlers

import (
	"avito_test_fix/connections"
	"avito_test_fix/operations"
	"avito_test_fix/utilities"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type AccrualRequest struct {
	UserId int `json:"user_id"`
	Amount int `json:"amount"`
}

type AccrualResponse struct {
	Status     string `json:"status"`
	NewBalance int    `json:"new_balance,omitempty"`
}

func Accrual(c *fiber.Ctx) error {
	var request AccrualRequest
	var response AccrualResponse
	err := c.BodyParser(&request)
	if err != nil {
		return c.JSON(utilities.RequestErrorInputData())
	}

	var userExist bool
	err = connections.PostgresConn.QueryRow(context.Background(),
		"select exists(select true from balance where user_id = $1);", request.UserId).Scan(&userExist)
	if err != nil {
		fmt.Printf("err elect exists(select true from balance where user_id = $1)")
		return c.JSON(utilities.RequestErrorQueryPg())
	}

	if userExist {
		// Пополнение баланса
		err = connections.PostgresConn.QueryRow(context.Background(),
			"update balance b set balance = (b.balance + $1) where user_id = $2 returning balance;",
			request.Amount, request.UserId).Scan(&response.NewBalance)
		if err != nil {
			return c.JSON(utilities.RequestErrorQueryPg())
		}

		operations.AddOperation(request.UserId, request.Amount, operations.OperationStatusSuccessfully, operations.OperationTypeAccrual)

		response.Status = "OK"
		return c.JSON(response)

	} else {
		// Создание новой записи с балансом
		_, err = connections.PostgresConn.Exec(context.Background(),
			"insert into balance (user_id, balance) values ($1, $2)", request.UserId, request.Amount)
		if err != nil {
			return c.JSON(utilities.RequestErrorQueryPg())
		}

		operations.AddOperation(request.UserId, request.Amount, operations.OperationStatusSuccessfully, operations.OperationTypeAccrual)

		response.Status = "OK"
		response.NewBalance = request.Amount
		return c.JSON(response)
	}
}
