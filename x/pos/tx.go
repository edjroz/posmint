package pos

import (
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/auth/util"
	"github.com/pokt-network/posmint/x/pos/types"
)

func (am AppModule) StakeTx(cdc *codec.Codec, txBuilder auth.TxBuilder, address sdk.ValAddress, passphrase string, amount sdk.Int) (*sdk.TxResponse, error) {
	cliCtx := util.NewCLIContext(am.GetTendermintNode(), sdk.AccAddress(address), passphrase).WithCodec(cdc)
	msg := types.MsgStake{
		Address: address,
		PubKey:  am.node.PrivValidator().GetPubKey(),
		Value:   amount,
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

func (am AppModule) UnstakeTx(cdc *codec.Codec, txBuilder auth.TxBuilder, address sdk.ValAddress, passphrase string) (*sdk.TxResponse, error) {
	cliCtx := util.NewCLIContext(am.GetTendermintNode(), sdk.AccAddress(address), passphrase).WithCodec(cdc)
	msg := types.MsgBeginUnstake{Address: address}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

func (am AppModule) UnjailTx(cdc *codec.Codec, txBuilder auth.TxBuilder, address sdk.ValAddress, passphrase string) (*sdk.TxResponse, error) {
	cliCtx := util.NewCLIContext(am.GetTendermintNode(), sdk.AccAddress(address), passphrase).WithCodec(cdc)
	msg := types.MsgUnjail{ValidatorAddr: address}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

func (am AppModule) Send(cdc *codec.Codec, fromAddr, toAddr sdk.ValAddress, txBuilder auth.TxBuilder, passphrase string, amount sdk.Int) (*sdk.TxResponse, error) {
	cliCtx := util.NewCLIContext(am.GetTendermintNode(), sdk.AccAddress(fromAddr), passphrase).WithCodec(cdc)
	msg := types.MsgSend{
		FromAddress: fromAddr,
		ToAddress:   toAddr,
		Amount:      amount,
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}
