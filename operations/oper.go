package operations

import (
	"avito_test_fix/connections"
	"context"
)

func AddOperation(UserId, Amount int, status OperationStatus, operationType OperationType) (int, error) {
	var operationId int
	err := connections.PostgresConn.QueryRow(context.Background(),
		"insert into operations (user_id, amount, status, start_date, operation_type) values ($1, $2, $3, now(), $4) returning id",
		UserId, Amount, status, operationType).Scan(&operationId)
	return operationId, err
}

func UpdateOperation(operationId int, status OperationStatus) {
	connections.PostgresConn.QueryRow(context.Background(),
		"update operations set status = $2 where id = $1",
		operationId, status)
}
