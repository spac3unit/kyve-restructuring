package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSubmitBundleProposal = "submit_bundle_proposal"

var _ sdk.Msg = &MsgSubmitBundleProposal{}

func NewMsgSubmitBundleProposal(creator string, staker string, poolId uint64, storageId string, byteSize uint64, fromHeight uint64, toHeight uint64, fromKey string, toKey string, toValue string, bundleHash string) *MsgSubmitBundleProposal {
	return &MsgSubmitBundleProposal{
		Creator: creator,
		Staker:      staker,
		PoolId:  poolId,
		StorageId: storageId,
		ByteSize: byteSize,
		FromHeight: fromHeight,
		ToHeight: toHeight,
		FromKey: fromKey,
		ToKey: toKey,
		ToValue: toValue,
		BundleHash: bundleHash,
	}
}

func (msg *MsgSubmitBundleProposal) Route() string {
	return RouterKey
}

func (msg *MsgSubmitBundleProposal) Type() string {
	return TypeMsgSubmitBundleProposal
}

func (msg *MsgSubmitBundleProposal) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSubmitBundleProposal) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSubmitBundleProposal) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
