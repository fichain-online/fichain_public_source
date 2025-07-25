package tracing

import (
	"strconv"
)

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[GasChangeUnspecified-0]
	_ = x[GasChangeTxInitialBalance-1]
	_ = x[GasChangeTxIntrinsicGas-2]
	_ = x[GasChangeTxRefunds-3]
	_ = x[GasChangeTxLeftOverReturned-4]
	_ = x[GasChangeCallInitialBalance-5]
	_ = x[GasChangeCallLeftOverReturned-6]
	_ = x[GasChangeCallLeftOverRefunded-7]
	_ = x[GasChangeCallContractCreation-8]
	_ = x[GasChangeCallContractCreation2-9]
	_ = x[GasChangeCallCodeStorage-10]
	_ = x[GasChangeCallOpCode-11]
	_ = x[GasChangeCallPrecompiledContract-12]
	_ = x[GasChangeCallStorageColdAccess-13]
	_ = x[GasChangeCallFailedExecution-14]
	_ = x[GasChangeWitnessContractInit-15]
	_ = x[GasChangeWitnessContractCreation-16]
	_ = x[GasChangeWitnessCodeChunk-17]
	_ = x[GasChangeWitnessContractCollisionCheck-18]
	_ = x[GasChangeTxDataFloor-19]
	_ = x[GasChangeIgnored-255]
}

const (
	_GasChangeReason_name_0 = "UnspecifiedTxInitialBalanceTxIntrinsicGasTxRefundsTxLeftOverReturnedCallInitialBalanceCallLeftOverReturnedCallLeftOverRefundedCallContractCreationCallContractCreation2CallCodeStorageCallOpCodeCallPrecompiledContractCallStorageColdAccessCallFailedExecutionWitnessContractInitWitnessContractCreationWitnessCodeChunkWitnessContractCollisionCheckTxDataFloor"
	_GasChangeReason_name_1 = "Ignored"
)

var _GasChangeReason_index_0 = [...]uint16{
	0,
	11,
	27,
	41,
	50,
	68,
	86,
	106,
	126,
	146,
	167,
	182,
	192,
	215,
	236,
	255,
	274,
	297,
	313,
	342,
	353,
}

func (i GasChangeReason) String() string {
	switch {
	case i <= 19:
		return _GasChangeReason_name_0[_GasChangeReason_index_0[i]:_GasChangeReason_index_0[i+1]]
	case i == 255:
		return _GasChangeReason_name_1
	default:
		return "GasChangeReason(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
