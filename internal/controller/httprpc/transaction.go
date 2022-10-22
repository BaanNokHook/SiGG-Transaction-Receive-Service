package httprpc

import (
	"nextclan/transaction-gateway/transaction-receive-service/internal/entity"
	usecase "nextclan/transaction-gateway/transaction-receive-service/internal/usecase"
	rpc "nextclan/transaction-gateway/transaction-receive-service/pkg/httprpc"
	"nextclan/transaction-gateway/transaction-receive-service/pkg/logger"
	"strconv"

	"github.com/gin-gonic/gin"
)

type rawTransaction struct {
	hexstring  string
	maxFeeRate int
}

type TransactionRpc struct {
	rrt usecase.ReceiveRawTransaction
	rvt usecase.ReceiveValidatedTransaction
	gau usecase.GetAddressUTXO
	l   logger.Interface
}

//These method will match with the json body from the client POST request.
func (t *TransactionRpc) Sendrawtransaction(hexdata string, maxGasFee int) string {
	result, err := t.rrt.Execute(entity.RawTransaction{TransactionData: hexdata, MaxGasFee: strconv.Itoa(maxGasFee)})
	if err != nil {
		return ""
	}
	return result
}

func (t *TransactionRpc) Sendvalidatedrawtransaction(transactionId string, signature string) string {
	t.rvt.Execute(entity.ValidatedTransaction{TransactionId: transactionId, Signature: signature})
	return transactionId
}

func (t *TransactionRpc) Getaddressutxos(params interface{}) interface{} {

	para := params.(map[string]interface{})
	rpcParam := map[string]interface{}{
		"addresses": para["addresses"],
		"chainInfo": para["chainInfo"],
	}
	response, _ := t.gau.Execute(rpcParam)
	return response.Result
}

func newTransactionRoutes(h *gin.Engine, rrt usecase.ReceiveRawTransaction, rvt usecase.ReceiveValidatedTransaction, gau usecase.GetAddressUTXO, l logger.Interface) {
	transactionRpc := TransactionRpc{rrt, rvt, gau, l}
	h.POST("/", func(c *gin.Context) { rpc.ProcessJsonRPC(c, &transactionRpc) })
}
