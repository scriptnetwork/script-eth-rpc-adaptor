package ethrpc

import (
	"context"
	"math"
	"math/big"

	"github.com/scripttoken/script-eth-rpc-adaptor/common"

	hexutil "github.com/scripttoken/script/common/hexutil"
	trpc "github.com/scripttoken/script/rpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_getBlockTransactionCountByNumber -----------------------------------
func (e *EthRPCService) GetBlockTransactionCountByNumber(ctx context.Context, numberStr string) (result hexutil.Uint64, err error) {
	logger.Infof("eth_getBlockTransactionCountByNumber called")
	height := common.GetHeightByTag(numberStr)
	if height == math.MaxUint64 {
		height, err = common.GetCurrentHeight()
		if err != nil {
			return result, err
		}
	}

	chainIDStr, err := e.ChainId(ctx)
	if err != nil {
		logger.Errorf("Failed to get chainID\n")
		return result, nil
	}
	chainID := new(big.Int)
	chainID.SetString(chainIDStr, 16)

	client := rpcc.NewRPCClient(common.GetScriptRPCEndpoint())
	rpcRes, rpcErr := client.Call("script.GetBlockByHeight", trpc.GetBlockByHeightArgs{
		Height: height})
	block, err := GetBlockFromTRPCResult(chainID, rpcRes, rpcErr, false)
	return hexutil.Uint64(len(block.Transactions)), err
}
