package usecase

import (
	"nextclan/transaction-gateway/transaction-receive-service/pkg/loaffinity"
	"nextclan/transaction-gateway/transaction-receive-service/pkg/logger"

	"github.com/KeisukeYamashita/go-jsonrpc"
)

type GetAddressUTXOUseCase struct {
	log        logger.Interface
	loaffinity loaffinity.ILoaffinityClient
}

func NewGetAddressUTXO(l logger.Interface, loaf *loaffinity.LoaffinityClient) *GetAddressUTXOUseCase {
	return &GetAddressUTXOUseCase{log: l, loaffinity: loaf}
}

func (gau *GetAddressUTXOUseCase) Execute(params interface{}) (*jsonrpc.RPCResponse, error) {
	return gau.loaffinity.GetAddressUTXO(params)
}
