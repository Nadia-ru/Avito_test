package main

import (
	"avito_test_fix/connections"
	h "avito_test_fix/handlers"
	"avito_test_fix/utilities"
	"github.com/gofiber/fiber/v2"
)

func main() {
	utilities.CheckEnvFile()
	connections.InitPostgresConnection()

	app := fiber.New()

	app.Post("/accrual", h.Accrual)
	app.Post("/get_balance", h.GetBalance)
	app.Post("/reservation", h.Reservation)
	app.Post("/un_reservation", h.UnReservation)

	err := app.Listen("0.0.0.0:8000")
	utilities.PanicIfErr(err)
}
