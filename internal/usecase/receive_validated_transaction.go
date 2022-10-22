package usecase

import (
	"encoding/json"
	"nextclan/transaction-gateway/transaction-receive-service/internal/entity"
	"nextclan/transaction-gateway/transaction-receive-service/pkg/logger"
)

//Use cases include:
/*
Given verified transaction When submit with validatedTransaction command Then the transaction should complete without errors.
Given verified transaction When receive transaction with verifiedTransaction command Then the transaction should publish into rabbitMQ without error.
*/

//Receive Verified Txn
//Publish to RabbitMQ

type ReceiveValidatedTransactionUseCase struct {
	log logger.Interface
}

func NewReceiveValidatedTransaction(l logger.Interface) *ReceiveValidatedTransactionUseCase {
	return &ReceiveValidatedTransactionUseCase{log: l}
}

//Receive new verified txn from clients
func (txn *ReceiveValidatedTransactionUseCase) Execute(validatedTransaction entity.ValidatedTransaction) {
	data, _ := json.Marshal(validatedTransaction)
	err := MessagingClient.PublishOnQueue(data, "txt.gw", "topic", "validated.transaction")
	if err != nil {
		txn.log.Error("cannot publis to queue for %s ,reason: %v", validatedTransaction.TransactionId, err)
	}
}
