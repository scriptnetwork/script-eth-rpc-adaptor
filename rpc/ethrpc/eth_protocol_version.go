package ethrpc

import (
	"context"
	"encoding/json"

	"github.com/scripttoken/script-eth-rpc-adaptor/common"

	trpc "github.com/scripttoken/script/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_protocolVersion -----------------------------------

func (e *EthRPCService) ProtocolVersion(ctx context.Context) (result string, err error) {
	logger.Infof("eth_protocolVersion called")

	client := rpcc.NewRPCClient(common.GetScriptRPCEndpoint())
	rpcRes, rpcErr := client.Call("script.GetVersion", trpc.GetVersionArgs{})

	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.GetVersionResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		return trpcResult.Version, nil
	}

	resultIntf, err := common.HandleScriptRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return "", err
	}
	result = resultIntf.(string)

	return result, nil
}
