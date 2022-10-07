package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// this line is used by starport scaffolding # 1
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	// this line is used by starport scaffolding # 2
	cdc.RegisterConcrete(&MsgCreateStaker{}, "registry/Stake", nil)
	cdc.RegisterConcrete(&MsgUpdateMetadata{}, "registry/UpdateMetadata", nil)
	cdc.RegisterConcrete(&MsgUpdateCommission{}, "registry/UpdateCommission", nil)
	cdc.RegisterConcrete(&MsgJoinPool{}, "registry/MsgJoinPool", nil)
	cdc.RegisterConcrete(&MsgLeavePool{}, "registry/MsgLeavePool", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {

	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgCreateStaker{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgUpdateCommission{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgUpdateMetadata{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgJoinPool{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgLeavePool{})

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
