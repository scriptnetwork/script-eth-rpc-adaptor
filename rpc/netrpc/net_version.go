package netrpc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/scripttoken/script-eth-rpc-adaptor/common"
	hexutil "github.com/scripttoken/script/common/hexutil"
	"github.com/scripttoken/script/ledger/types"
	trpc "github.com/scripttoken/script/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

type chainIDResultWrapper struct {
	chainID string
}

// ------------------------------- net_version -----------------------------------

func (e *NetRPCService) Version(ctx context.Context) (result string, err error) {
	logger.Infof("net_version called")

	client := rpcc.NewRPCClient(common.GetScriptRPCEndpoint())
	rpcRes, rpcErr := client.Call("script.GetStatus", trpc.GetStatusArgs{})
	var blockHeight uint64
	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.GetStatusResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		re := chainIDResultWrapper{
			chainID: trpcResult.ChainID,
		}
		blockHeight = uint64(trpcResult.LatestFinalizedBlockHeight)
		return re, nil
	}

	resultIntf, err := common.HandleScriptRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return "", err
	}
	scriptChainIDResult, ok := resultIntf.(chainIDResultWrapper)
	if !ok {
		return "", fmt.Errorf("failed to convert chainIDResultWrapper")
	}

	scriptChainID := scriptChainIDResult.chainID
	ethChainID := types.MapChainID(scriptChainID, blockHeight).Uint64()
	result = hexutil.EncodeUint64(ethChainID)

	return result, nil
}
