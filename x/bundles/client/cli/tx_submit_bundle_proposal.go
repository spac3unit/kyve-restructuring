package cli

import (
	"github.com/KYVENetwork/chain/x/bundles/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdSubmitBundleProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-bundle-proposal [staker] [pool_id] [storage_id] [byte_size] [from_height] [to_height] [from_key] [to_key] [to_value] [bundle_hash]",
		Short: "Broadcast message submit-bundle-proposal",
		Args:  cobra.ExactArgs(10),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argStaker := args[0]

			argPoolId, err := cast.ToUint64E(args[1])
			if err != nil {
				return err
			}

			argStorageId := args[2]

			argByteSize, err := cast.ToUint64E(args[3])
			if err != nil {
				return err
			}

			argFromHeight, err := cast.ToUint64E(args[4])
			if err != nil {
				return err
			}

			argToHeight, err := cast.ToUint64E(args[5])
			if err != nil {
				return err
			}

			argFromKey := args[6]

			argToKey := args[7]

			argToValue := args[8]

			argBundleHash := args[9]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSubmitBundleProposal(
				clientCtx.GetFromAddress().String(),
				argStaker,
				argPoolId,
				argStorageId,
				argByteSize,
				argFromHeight,
				argToHeight,
				argFromKey,
				argToKey,
				argToValue,
				argBundleHash,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
