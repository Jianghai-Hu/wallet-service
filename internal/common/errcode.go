package common

// TODO: move this part to proto file
// errcode range reserved for this service[15500000,15600000)

//nolint:revive,stylecheck // fix in future
const (
	Constant_ERROR_UNKNOW             = 15500000
	Constant_ERROR_SERVICE_INTERNAL   = 15500001
	Constant_ERROR_INSUFFIENT_BALANCE = 15500002
	Constant_ERROR_PARAM              = 15500003
)
