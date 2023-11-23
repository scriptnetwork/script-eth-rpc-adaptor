package ethrpc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/scripttoken/script-eth-rpc-adaptor/common"
	tcommon "github.com/scripttoken/script/common"
	"github.com/scripttoken/script/common/hexutil"
	"github.com/scripttoken/script/ledger/types"
	trpc "github.com/scripttoken/script/rpc"

	rpcc "github.com/ybbus/jsonrpc"
)

// ------------------------------- eth_getTransactionByBlockHashAndIndex -----------------------------------
func (e *EthRPCService) GetTransactionByBlockHashAndIndex(ctx context.Context, hashStr string, txIndexStr string) (result common.EthGetTransactionResult, err error) {
	logger.Infof("GetTransactionByBlockHashAndIndex called")
	txIndex := common.GetHeightByTag(txIndexStr)
	client := rpcc.NewRPCClient(common.GetScriptRPCEndpoint())
	rpcRes, rpcErr := client.Call("script.GetBlock", trpc.GetBlockArgs{Hash: tcommon.HexToHash(hashStr)})
	return GetIndexedTransactionFromBlock(rpcRes, rpcErr, txIndex)
}

func GetIndexedTransactionFromBlock(rpcRes *rpcc.RPCResponse, rpcErr error, txIndex tcommon.JSONUint64) (result common.EthGetTransactionResult, err error) {
	result = common.EthGetTransactionResult{}
	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := common.ScriptGetBlockResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		if txIndex >= tcommon.JSONUint64(len(trpcResult.Txs)) {
			return result, fmt.Errorf("transaction index out of range")
		}
		result.TransactionIndex = hexutil.Uint64(txIndex)
		var objmap map[string]json.RawMessage
		json.Unmarshal(jsonBytes, &objmap)
		result.BlockHash = trpcResult.Hash
		result.BlockHeight = hexutil.Uint64(trpcResult.Height)
		result.Nonce = hexutil.Uint64(0)

		if objmap["transactions"] != nil {
			var txmaps []map[string]json.RawMessage
			json.Unmarshal(objmap["transactions"], &txmaps)
			indexedTx := trpcResult.Txs[txIndex]
			omap := txmaps[txIndex]
			result.TxHash = indexedTx.Hash
			if types.TxType(indexedTx.Type) == types.TxSmartContract {
				tx := types.SmartContractTx{}
				json.Unmarshal(omap["raw"], &tx)
				result.From = tx.From.Address
				if (tx.To.Address == tcommon.Address{}) {
					result.To = nil // conform to ETH standard
				} else {
					result.To = &tx.To.Address
				}
				result.GasPrice = "0x" + tx.GasPrice.Text(16)
				result.Gas = hexutil.Uint64(tx.GasLimit)
				result.Value = "0x" + tx.From.Coins.SPAYWei.Text(16)
				result.Input = tx.Data.String()
				result.Nonce = hexutil.Uint64(tx.From.Sequence) - 1 // off-by-one: Ethereum's account nonce starts from 0, while Script's account sequnce starts from 1
				data := tx.From.Signature.ToBytes()
				GetRSVfromSignature(data, &result)
			} else if types.TxType(indexedTx.Type) == types.TxSend {
				tx := types.SendTx{}
				json.Unmarshal(omap["raw"], &tx)
				result.From = tx.Inputs[0].Address
				if (tx.Outputs[0].Address == tcommon.Address{}) {
					result.To = nil // conform to ETH standard
				} else {
					result.To = &tx.Outputs[0].Address
				}
				result.Gas = hexutil.Uint64(tx.Fee.SPAYWei.Uint64())
				result.Value = "0x" + tx.Inputs[0].Coins.SPAYWei.Text(16)
				result.Nonce = hexutil.Uint64(tx.Inputs[0].Sequence) - 1 // off-by-one: Ethereum's account nonce starts from 0, while Script's account sequnce starts from 1
				data := tx.Inputs[0].Signature.ToBytes()
				GetRSVfromSignature(data, &result)
			}
		}
		return trpcResult, nil
	}
	_, err = common.HandleScriptRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return result, err
	}
	return result, nil
}
