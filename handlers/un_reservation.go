package handlers

import (
	"avito_test_fix/connections"
	"avito_test_fix/operations"
	"avito_test_fix/utilities"
	"context"
	"github.com/gofiber/fiber/v2"
)

type UnReservationRequest struct {
	UserId    int `json:"user_id"`
	ServiceId int `json:"service_id"`
	OrderId   int `json:"order_id"`

	Amount int `json:"amount"`
}

type UnReservationResponse struct {
	Status  string `json:"status"`
	Details string `json:"details,omitempty"`
}

func UnReservation(c *fiber.Ctx) error {
	var request UnReservationRequest
	var response UnReservationResponse
	err := c.BodyParser(&request)
	if err != nil {
		return c.JSON(utilities.RequestErrorInputData())
	}

	var operationId, operationAmount, dealId int
	err = connections.PostgresConn.QueryRow(context.Background(),
		`select  operation_id, oper.amount, d.id from deals d
				inner join operations oper on d.operation_id = oper.id
				where oper.user_id = $1 and d.order_id = $2 and d.service_id = $3`,
		request.UserId, request.OrderId, request.ServiceId).Scan(&operationId, &operationAmount, &dealId)
	if err != nil {
		return c.JSON(utilities.RequestErrorQueryPg())
	}

	if request.Amount != 0 {
		if request.Amount == operationAmount {

			connections.PostgresConn.Exec(context.Background(),
				"update deals set end_date = now(), final_amount = $1 where id = $2;", operationAmount, dealId)

			operations.UpdateOperation(operationId, operations.OperationStatusSuccessfully)
			response.Status = "OK"
			return c.JSON(response)
		} else {
			// Если запрошено больше средств, чем было зарезервировано
			if operationAmount-request.Amount < 0 {
				return c.JSON(map[string]string{
					"status": "err",
					"detail": "not enough funds in reserve",
				})
			}

			connections.PostgresConn.Exec(context.Background(),
				"update deals set end_date = now(), final_amount = $1 where id = $2;", request.Amount, dealId)

			connections.PostgresConn.QueryRow(context.Background(),
				"update balance b set balance = (b.balance + $1) where user_id = $2;",
				operationAmount-request.Amount, request.UserId)
			operations.UpdateOperation(operationId, operations.OperationStatusSuccessfully)

			response.Status = "OK"
			return c.JSON(response)
		}
	} else { // Если не передана сумма по завершению, считается что она не изменилась.
		connections.PostgresConn.Exec(context.Background(),
			"update deals set end_date = now(), final_amount = $1 where id = $2;", operationAmount, dealId)

		operations.UpdateOperation(operationId, operations.OperationStatusSuccessfully)
		response.Status = "OK"
		return c.JSON(response)
	}

}
