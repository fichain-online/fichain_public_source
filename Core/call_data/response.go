package call_data

import (
	"encoding/hex"
	"fmt"

	"google.golang.org/protobuf/proto"

	"FichainCore/common"
	pb "FichainCore/proto" // make sure this matches your proto import
)

type CallSmartContractResponse struct {
	Hash common.Hash
	Data []byte
}

// Serialize to protobuf
func (r *CallSmartContractResponse) Proto() proto.Message {
	return &pb.CallSmartContractResponse{
		Hash: r.Hash.Bytes(),
		Data: r.Data,
	}
}

// Marshal to bytes
func (r *CallSmartContractResponse) Marshal() ([]byte, error) {
	return proto.Marshal(r.Proto())
}

// Unmarshal from bytes
func (r *CallSmartContractResponse) Unmarshal(bz []byte) error {
	pbResp := &pb.CallSmartContractResponse{}
	if err := proto.Unmarshal(bz, pbResp); err != nil {
		return err
	}
	return r.FromProto(pbResp)
}

// Convert from protobuf to struct
func (r *CallSmartContractResponse) FromProto(pbResp *pb.CallSmartContractResponse) error {
	r.Hash = common.BytesToHash(pbResp.Hash)
	r.Data = pbResp.Data
	return nil
}

// String returns a human-readable representation of the response
func (r *CallSmartContractResponse) String() string {
	return fmt.Sprintf(
		"CallSmartContractResponse{Hash: %s, Data: %s}",
		r.Hash.Hex(),
		hex.EncodeToString(r.Data),
	)
}
