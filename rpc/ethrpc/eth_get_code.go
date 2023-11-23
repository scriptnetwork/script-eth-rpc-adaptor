package ethrpc

import (
	"context"
	"encoding/json"
	"math"
	"strings"

	"github.com/scripttoken/script-eth-rpc-adaptor/common"

	trpc "github.com/scripttoken/script/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_getCode -----------------------------------

func (e *EthRPCService) GetCode(ctx context.Context, address string, tag string) (result string, err error) {
	logger.Infof("eth_getCode called")

	height := common.GetHeightByTag(tag)
	if height == math.MaxUint64 {
		height = 0 // 0 is interpreted as the last height by the script.GetAccount method
	}

	client := rpcc.NewRPCClient(common.GetScriptRPCEndpoint())

	// maxRetry := 3
	maxRetry := 1
	for i := 0; i < maxRetry; i++ { // It might take some time for a tx to be finalized, retry a few times

		rpcRes, rpcErr := client.Call("script.GetCode", trpc.GetCodeArgs{Address: address, Height: height})

		parse := func(jsonBytes []byte) (interface{}, error) {
			trpcResult := trpc.GetCodeResult{}
			json.Unmarshal(jsonBytes, &trpcResult)
			return trpcResult.Code, nil
		}

		resultIntf, err := common.HandleScriptRPCResponse(rpcRes, rpcErr, parse)
		if err != nil {
			return result, err
		}

		result = resultIntf.(string)
		if result == "" { // might need to wait for the tx to be finalized
			// time.Sleep(blockInterval) // one block duration
		}
	}

	if result == "" {
		result = "0x"
	}

	if !strings.HasPrefix(result, "0x") {
		result = "0x" + result
	}

	return result, nil
}
