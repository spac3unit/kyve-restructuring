package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgDelegate{}, "registry/Delegate", nil)
	cdc.RegisterConcrete(&MsgWithdrawRewards{}, "registry/WithdrawRewards", nil)
	cdc.RegisterConcrete(&MsgUndelegate{}, "registry/Undelegate", nil)
	cdc.RegisterConcrete(&MsgRedelegate{}, "registry/Redelegate", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgDelegate{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgWithdrawRewards{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgUndelegate{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgRedelegate{})

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
