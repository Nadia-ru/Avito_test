package operations

type OperationType int

const (
	OperationTypeAccrual     OperationType = 1
	OperationTypeOrder                     = 2
	OperationTypeTranslation               = 3
)

type OperationStatus int

const (
	OperationStatusAwaitingCompletion OperationStatus = 1
	OperationStatusSuccessfully                       = 2
	OperationStatusCanceled                           = 3
)
