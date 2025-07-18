package tracing

import (
	"strconv"
)

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[NonceChangeUnspecified-0]
	_ = x[NonceChangeGenesis-1]
	_ = x[NonceChangeEoACall-2]
	_ = x[NonceChangeContractCreator-3]
	_ = x[NonceChangeNewContract-4]
	_ = x[NonceChangeAuthorization-5]
	_ = x[NonceChangeRevert-6]
}

const _NonceChangeReason_name = "UnspecifiedGenesisEoACallContractCreatorNewContractAuthorizationRevert"

var _NonceChangeReason_index = [...]uint8{0, 11, 18, 25, 40, 51, 64, 70}

func (i NonceChangeReason) String() string {
	if i >= NonceChangeReason(len(_NonceChangeReason_index)-1) {
		return "NonceChangeReason(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _NonceChangeReason_name[_NonceChangeReason_index[i]:_NonceChangeReason_index[i+1]]
}
