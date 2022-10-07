package types

import (
	"fmt"
)

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any failure.
func (gs GenesisState) Validate() error {

	// Staker
	stakerLeaving := make(map[string]bool)

	// Valaccounts
	valaccountMap := make(map[string]struct{})
	for _, elem := range gs.ValaccountList {
		index := string(ValaccountKey(elem.PoolId, elem.Staker))
		if _, ok := valaccountMap[index]; ok {
			return fmt.Errorf("duplicated index for valaccount %v", elem)
		}
		valaccountMap[index] = struct{}{}
		stakerLeaving[index] = elem.IsLeaving
	}

	// Commission Change
	commissionChangeMap := make(map[string]struct{})

	for _, elem := range gs.CommissionChangeEntries {
		index := string(CommissionChangeEntryKey(elem.Index))
		if _, ok := commissionChangeMap[index]; ok {
			return fmt.Errorf("duplicated index for commission change entry %v", elem)
		}
		if elem.Index > gs.QueueStateCommission.HighIndex {
			return fmt.Errorf("commission change entry index too high: %v", elem)
		}
		if elem.Index < gs.QueueStateCommission.LowIndex {
			return fmt.Errorf("commission change entry index too low: %v", elem)
		}

		commissionChangeMap[index] = struct{}{}
	}

	// Leave Pool
	leavePoolMap := make(map[string]struct{})

	for _, elem := range gs.LeavePoolEntries {
		index := string(UnbondingStakeEntryKey(elem.Index))
		if _, ok := leavePoolMap[index]; ok {
			return fmt.Errorf("duplicated index for unbonding stake entry %v", elem)
		}
		if elem.Index > gs.QueueStateLeave.HighIndex {
			return fmt.Errorf("unbonding stake entry index too high: %v", elem)
		}
		if elem.Index < gs.QueueStateLeave.LowIndex {
			return fmt.Errorf("unbonding stake entry index too low: %v", elem)
		}
		if stakerLeaving[string(ValaccountKey(elem.PoolId, elem.Staker))] != true {
			return fmt.Errorf("inconsistent staker leave: %v", elem)
		}
		stakerLeaving[string(ValaccountKey(elem.PoolId, elem.Staker))] = false

		leavePoolMap[index] = struct{}{}
	}

	for staker, isLeaving := range stakerLeaving {
		if isLeaving != false {
			return fmt.Errorf("inconsistent staker leave: %v", staker)

		}
	}

	return gs.Params.Validate()
}
