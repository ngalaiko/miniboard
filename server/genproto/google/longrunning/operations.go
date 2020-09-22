package longrunning

import "google.golang.org/genproto/googleapis/longrunning"

type (
	Operation              = longrunning.Operation
	ListOperationsRequest  = longrunning.ListOperationsRequest
	ListOperationsResponse = longrunning.ListOperationsResponse
	GetOperationRequest    = longrunning.GetOperationRequest
	DeleteOperationRequest = longrunning.DeleteOperationRequest
	CancelOperationRequest = longrunning.CancelOperationRequest
	WaitOperationRequest   = longrunning.WaitOperationRequest
	Operation_Error        = longrunning.Operation_Error
	Operation_Response     = longrunning.Operation_Response
	OperationsClient       = longrunning.OperationsClient
	OperationsServer       = longrunning.OperationsServer
)

var NewOperationsClient = longrunning.NewOperationsClient
