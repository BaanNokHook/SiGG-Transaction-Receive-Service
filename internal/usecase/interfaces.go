package usecase

import (
	"nextclan/transaction-gateway/transaction-receive-service/internal/entity"
	messaging "nextclan/transaction-gateway/transaction-receive-service/pkg/rabbitmq"

	"github.com/KeisukeYamashita/go-jsonrpc"
)

type (
	ReceiveRawTransaction interface {
		Execute(entity.RawTransaction) (string, error)
	}

	ReceiveValidatedTransaction interface {
		Execute(entity.ValidatedTransaction)
	}

	GetAddressUTXO interface {
		Execute(interface{}) (*jsonrpc.RPCResponse, error)
	}
)

var MessagingClient messaging.IMessagingClient
