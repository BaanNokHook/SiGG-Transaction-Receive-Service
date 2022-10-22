package usecase

import (
	"encoding/json"
	"nextclan/transaction-gateway/transaction-receive-service/internal/entity"
	"nextclan/transaction-gateway/transaction-receive-service/pkg/logger"
	"nextclan/transaction-gateway/transaction-receive-service/pkg/redis"
)

//Use cases include:
/*
Given raw transaction When create with createTransaction command Then the transaction should complete without errors.
Given raw transaction When receive transaction with createTransaction command Then the transaction should publish into rabbitMQ without error.
*/

//Receive Raw Txn
//Publish to RabbitMQ

type ReceiveRawTransactionUsecase struct {
	log        logger.Interface
	redisCache *redis.RedisCache
}

func NewReceiveRawTransaction(l logger.Interface, redis *redis.RedisCache) *ReceiveRawTransactionUsecase {
	return &ReceiveRawTransactionUsecase{log: l, redisCache: redis}
}

//Receive new raw txn from clients
func (txn *ReceiveRawTransactionUsecase) Execute(rawTransaction entity.RawTransaction) (string, error) {
	transacionId, err := rawTransaction.Hash()
	if err != nil {
		txn.log.Error("invalid transaction for %s ,reason: %v", rawTransaction.TransactionData, err)
		return "", err
	}

	//TODO save transaction to redis
	data, err := json.Marshal(rawTransaction)
	err = txn.redisCache.Set(transacionId, data)
	if err != nil {
		txn.log.Error("cannot save to cache for %s ,reason: %v", rawTransaction.TransactionData, err)
		return "", err
	}
	//publish to MQ
	data, err = json.Marshal(entity.ValidateTransaction{TransactionId: transacionId})
	err = MessagingClient.PublishOnQueue(data, "txt.gw", "topic", "raw.transaction")
	if err != nil {
		txn.log.Error("cannot publish to queue for %s ,reason: %v", rawTransaction.TransactionData, err)
		return "", err
	}
	return transacionId, nil
}
