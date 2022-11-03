package handlers

import (
	"avito_test_fix/connections"
	"avito_test_fix/utilities"
	"context"
	"github.com/gofiber/fiber/v2"
)

type GetBalanceRequest struct {
	UserId int `json:"user_id"`
}

type GetBalanceResponse struct {
	Status  string `json:"status"`
	Balance int    `json:"balance"`
}

func GetBalance(c *fiber.Ctx) error {
	var request GetBalanceRequest
	var response GetBalanceResponse
	err := c.BodyParser(&request)
	if err != nil {
		return c.JSON(utilities.RequestErrorInputData())
	}

	err = connections.PostgresConn.QueryRow(context.Background(),
		"select b.balance from balance b where user_id = $1;", request.UserId).Scan(&response.Balance)
	
	if err != nil {
		return c.JSON(utilities.RequestErrorQueryPg())
	}

	response.Status = "OK"
	return c.JSON(response)
}
