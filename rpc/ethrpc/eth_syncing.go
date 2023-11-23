package ethrpc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/scripttoken/script-eth-rpc-adaptor/common"

	"github.com/scripttoken/script/common/hexutil"
	trpc "github.com/scripttoken/script/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

type syncingResultWrapper struct {
	*common.EthSyncingResult
	syncing bool
}

// ------------------------------- eth_syncing -----------------------------------
func (e *EthRPCService) Syncing(ctx context.Context) (result interface{}, err error) {
	logger.Infof("eth_syncing called")
	client := rpcc.NewRPCClient(common.GetScriptRPCEndpoint())
	rpcRes, rpcErr := client.Call("script.GetStatus", trpc.GetStatusArgs{})
	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.GetStatusResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		re := syncingResultWrapper{&common.EthSyncingResult{}, false}
		re.syncing = trpcResult.Syncing
		if trpcResult.Syncing {
			re.StartingBlock = 1
			re.CurrentBlock = hexutil.Uint64(trpcResult.CurrentHeight)
			re.HighestBlock = hexutil.Uint64(trpcResult.LatestFinalizedBlockHeight)
			re.PulledStates = re.CurrentBlock
			re.KnownStates = re.CurrentBlock
		}
		return re, nil
	}

	resultIntf, err := common.HandleScriptRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return "", err
	}
	scriptSyncingResult, ok := resultIntf.(syncingResultWrapper)
	if !ok {
		return nil, fmt.Errorf("failed to convert syncingResultWrapper")
	}
	if !scriptSyncingResult.syncing {
		result = false
	} else {
		result = scriptSyncingResult.EthSyncingResult
	}

	return result, nil
}
