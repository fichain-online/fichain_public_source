package tracing

import (
	"strconv"
)

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[BalanceChangeUnspecified-0]
	_ = x[BalanceIncreaseRewardMineUncle-1]
	_ = x[BalanceIncreaseRewardMineBlock-2]
	_ = x[BalanceIncreaseWithdrawal-3]
	_ = x[BalanceIncreaseGenesisBalance-4]
	_ = x[BalanceIncreaseRewardTransactionFee-5]
	_ = x[BalanceDecreaseGasBuy-6]
	_ = x[BalanceIncreaseGasReturn-7]
	_ = x[BalanceIncreaseDaoContract-8]
	_ = x[BalanceDecreaseDaoAccount-9]
	_ = x[BalanceChangeTransfer-10]
	_ = x[BalanceChangeTouchAccount-11]
	_ = x[BalanceIncreaseSelfdestruct-12]
	_ = x[BalanceDecreaseSelfdestruct-13]
	_ = x[BalanceDecreaseSelfdestructBurn-14]
	_ = x[BalanceChangeRevert-15]
}

const _BalanceChangeReason_name = "UnspecifiedBalanceIncreaseRewardMineUncleBalanceIncreaseRewardMineBlockBalanceIncreaseWithdrawalBalanceIncreaseGenesisBalanceBalanceIncreaseRewardTransactionFeeBalanceDecreaseGasBuyBalanceIncreaseGasReturnBalanceIncreaseDaoContractBalanceDecreaseDaoAccountTransferTouchAccountBalanceIncreaseSelfdestructBalanceDecreaseSelfdestructBalanceDecreaseSelfdestructBurnRevert"

var _BalanceChangeReason_index = [...]uint16{
	0,
	11,
	41,
	71,
	96,
	125,
	160,
	181,
	205,
	231,
	256,
	264,
	276,
	303,
	330,
	361,
	367,
}

func (i BalanceChangeReason) String() string {
	if i >= BalanceChangeReason(len(_BalanceChangeReason_index)-1) {
		return "BalanceChangeReason(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _BalanceChangeReason_name[_BalanceChangeReason_index[i]:_BalanceChangeReason_index[i+1]]
}
