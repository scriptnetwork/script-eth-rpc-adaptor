package ethrpc

import (
	"context"
	"encoding/json"

	"github.com/scripttoken/script-eth-rpc-adaptor/common"

	trpc "github.com/scripttoken/script/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_sendRawTransaction -----------------------------------

func (e *EthRPCService) SendRawTransaction(ctx context.Context, txBytes string) (result string, err error) {
	logger.Infof("eth_sendRawTransaction called")

	client := rpcc.NewRPCClient(common.GetScriptRPCEndpoint())
	rpcRes, rpcErr := client.Call("script.BroadcastRawEthTransactionAsync", trpc.BroadcastRawTransactionAsyncArgs{TxBytes: txBytes})

	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.BroadcastRawTransactionAsyncResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		return trpcResult.TxHash, nil
	}

	resultIntf, err := common.HandleScriptRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		logger.Errorf("eth_sendRawTransaction, err: %v", err)
		return "", err
	}
	result = resultIntf.(string)

	logger.Infof("eth_sendRawTransaction, result: %v\n", result)

	return result, nil
}
