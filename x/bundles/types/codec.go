package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// this line is used by starport scaffolding # 1
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgSubmitBundleProposal{}, "bundles/SubmitBundleProposal", nil)
	cdc.RegisterConcrete(&MsgVoteBundleProposal{}, "bundles/VoteBundleProposal", nil)
	cdc.RegisterConcrete(&MsgClaimUploaderRole{}, "bundles/ClaimUploaderRole", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// this line is used by starport scaffolding # 3
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSubmitBundleProposal{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgVoteBundleProposal{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgClaimUploaderRole{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
