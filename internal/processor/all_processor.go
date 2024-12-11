package processor

import (
	"net/http"

	"jianghai-hu/wallet-service/internal/common"
)

type Response struct {
	ErrorCode int         `json:"error_code"`
	ErrorMsg  string      `json:"error_msg"`
	Data      interface{} `json:"data"`
}

type APIProcessorConfig struct {
	Command   string
	Processor func(http.ResponseWriter, *http.Request)
	Method    string
}

func AllProcessorConfigs() []*APIProcessorConfig {
	return []*APIProcessorConfig{
		{
			Command:   common.COMMAND_DEPOSIT,
			Processor: DepositProcessor,
			Method:    "POST",
		},
		{
			Command:   common.COMMAND_WITHDRAW,
			Processor: WithdrawProcessor,
			Method:    "POST",
		},
		{
			Command:   common.COMMAND_TRANSFER,
			Processor: TransferProcessor,
			Method:    "POST",
		},
	}
}
