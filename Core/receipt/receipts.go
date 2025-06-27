package receipt

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	pb "FichainCore/proto"
)

type Receipts struct {
	receipts []*Receipt
}

// NewReceipts creates a new Receipts collection from a slice of receipts.
func NewReceipts(receipts []*Receipt) *Receipts {
	return &Receipts{receipts: receipts}
}

func (rs *Receipts) Receipts() []*Receipt {
	return rs.receipts
}

// Proto converts the Receipts collection to its protobuf representation.
func (rs *Receipts) Proto() proto.Message {
	pbReceipts := &pb.Receipts{
		Receipts: make([]*pb.Receipt, len(rs.receipts)),
	}
	for i, r := range rs.receipts {
		pbReceipts.Receipts[i] = r.Proto().(*pb.Receipt)
	}
	return pbReceipts
}

// FromProto populates the Receipts collection from its protobuf representation.
func (rs *Receipts) FromProto(p proto.Message) error {
	pbReceipts, ok := p.(*pb.Receipts)
	if !ok {
		return fmt.Errorf("invalid proto message type: %T", p)
	}

	rs.receipts = make([]*Receipt, len(pbReceipts.Receipts))
	for i, pbRcpt := range pbReceipts.Receipts {
		r := &Receipt{}
		if err := r.FromProto(pbRcpt); err != nil {
			return fmt.Errorf("failed to convert receipt %d from proto: %w", i, err)
		}
		rs.receipts[i] = r
	}
	return nil
}

// Marshal serializes the Receipts collection into a byte slice.
func (rs *Receipts) Marshal() ([]byte, error) {
	return proto.Marshal(rs.Proto())
}

// Unmarshal deserializes a byte slice into the Receipts collection.
func (rs *Receipts) Unmarshal(data []byte) error {
	pbReceipts := &pb.Receipts{}
	if err := proto.Unmarshal(data, pbReceipts); err != nil {
		return err
	}
	return rs.FromProto(pbReceipts)
}

func MarshalReceipts(receipts []*Receipt) ([]byte, error) {
	rcptsProto := &pb.Receipts{}
	for _, v := range receipts {
		rcptsProto.Receipts = append(rcptsProto.Receipts, v.Proto().(*pb.Receipt))
	}
	return proto.Marshal(rcptsProto)
}

func UnmarshalReceipts(data []byte) ([]*Receipt, error) {
	rcptsProto := &pb.Receipts{}
	if err := proto.Unmarshal(data, rcptsProto); err != nil {
		return nil, err
	}

	receipts := make([]*Receipt, len(rcptsProto.Receipts))
	for i, pbRcpt := range rcptsProto.Receipts {
		r := &Receipt{}
		if err := r.FromProto(pbRcpt); err != nil {
			return nil, err
		}
		receipts[i] = r
	}
	return receipts, nil
}
