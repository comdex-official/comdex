// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: comdex/auction/v1beta1/genesis.proto

package types

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type GenesisState struct {
	SurplusAuction     []SurplusAuction     `protobuf:"bytes,1,rep,name=surplusAuction,proto3" json:"surplusAuction" yaml:"surplusAuction"`
	DebtAuction        []DebtAuction        `protobuf:"bytes,2,rep,name=debtAuction,proto3" json:"debtAuction" yaml:"debtAuction"`
	DutchAuction       []DutchAuction       `protobuf:"bytes,3,rep,name=dutchAuction,proto3" json:"dutchAuction" yaml:"dutchAuction"`
	ProtocolStatistics []ProtocolStatistics `protobuf:"bytes,4,rep,name=protocolStatistics,proto3" json:"protocolStatistics" yaml:"protocolStatistics"`
	AuctionParams      []AuctionParams      `protobuf:"bytes,5,rep,name=auctionParams,proto3" json:"auctionParams" yaml:"auctionParams"`
	Params             Params               `protobuf:"bytes,6,opt,name=params,proto3" json:"params"`
	UserBiddingID      uint64               `protobuf:"varint,7,opt,name=UserBiddingID,proto3" json:"UserBiddingID,omitempty"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_49088f171dd3086d, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetSurplusAuction() []SurplusAuction {
	if m != nil {
		return m.SurplusAuction
	}
	return nil
}

func (m *GenesisState) GetDebtAuction() []DebtAuction {
	if m != nil {
		return m.DebtAuction
	}
	return nil
}

func (m *GenesisState) GetDutchAuction() []DutchAuction {
	if m != nil {
		return m.DutchAuction
	}
	return nil
}

func (m *GenesisState) GetProtocolStatistics() []ProtocolStatistics {
	if m != nil {
		return m.ProtocolStatistics
	}
	return nil
}

func (m *GenesisState) GetAuctionParams() []AuctionParams {
	if m != nil {
		return m.AuctionParams
	}
	return nil
}

func (m *GenesisState) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

func (m *GenesisState) GetUserBiddingID() uint64 {
	if m != nil {
		return m.UserBiddingID
	}
	return 0
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "comdex.auction.v1beta1.GenesisState")
}

func init() {
	proto.RegisterFile("comdex/auction/v1beta1/genesis.proto", fileDescriptor_49088f171dd3086d)
}

var fileDescriptor_49088f171dd3086d = []byte{
	// 425 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0xb1, 0x8e, 0xd3, 0x30,
	0x18, 0x80, 0x13, 0xae, 0x57, 0x24, 0xf7, 0x8e, 0xc1, 0x1c, 0x28, 0x04, 0xc8, 0x95, 0x5c, 0x0f,
	0x55, 0x48, 0x24, 0xea, 0xb1, 0x21, 0x96, 0x8b, 0x2a, 0x21, 0xc4, 0x52, 0xa5, 0x62, 0x61, 0x73,
	0x12, 0x37, 0xb5, 0x94, 0xd4, 0x51, 0xec, 0x20, 0x3a, 0xf0, 0x0e, 0x2c, 0xbc, 0x53, 0xc7, 0x8e,
	0x4c, 0x15, 0x6a, 0xdf, 0x80, 0x27, 0x40, 0xb1, 0x5d, 0x35, 0xa1, 0xf5, 0xd6, 0x26, 0x9f, 0xbf,
	0x2f, 0xfe, 0xf5, 0x83, 0x41, 0x4c, 0xf3, 0x04, 0x7f, 0xf7, 0x51, 0x15, 0x73, 0x42, 0x17, 0xfe,
	0xb7, 0x51, 0x84, 0x39, 0x1a, 0xf9, 0x29, 0x5e, 0x60, 0x46, 0x98, 0x57, 0x94, 0x94, 0x53, 0xf8,
	0x54, 0x52, 0x9e, 0xa2, 0x3c, 0x45, 0xd9, 0x57, 0x29, 0x4d, 0xa9, 0x40, 0xfc, 0xfa, 0x97, 0xa4,
	0xed, 0x1b, 0x8d, 0xb3, 0x40, 0x25, 0xca, 0x95, 0xd2, 0xd6, 0x85, 0xf7, 0x09, 0x49, 0xdd, 0x6a,
	0xa8, 0x88, 0x24, 0x09, 0x59, 0xa4, 0x4a, 0xe6, 0xfe, 0x3a, 0x07, 0x17, 0x1f, 0xe5, 0x17, 0x4f,
	0x39, 0xe2, 0x18, 0xe6, 0xe0, 0x11, 0xab, 0xca, 0x22, 0xab, 0xd8, 0xbd, 0x3c, 0x69, 0x99, 0xfd,
	0xb3, 0x61, 0xef, 0xee, 0xb5, 0x77, 0xfa, 0x26, 0xde, 0xb4, 0x45, 0x07, 0x2f, 0x57, 0x9b, 0x6b,
	0xe3, 0xef, 0xe6, 0xfa, 0xc9, 0x12, 0xe5, 0xd9, 0x7b, 0xb7, 0xed, 0x72, 0xc3, 0xff, 0xe4, 0x10,
	0x81, 0x5e, 0x82, 0x23, 0xbe, 0x6f, 0x3d, 0x10, 0xad, 0x1b, 0x5d, 0x6b, 0x7c, 0x40, 0x03, 0x5b,
	0x85, 0xa0, 0x0c, 0x35, 0x2c, 0x6e, 0xd8, 0x74, 0x42, 0x0c, 0x2e, 0x92, 0x8a, 0xc7, 0xf3, 0x7d,
	0xe3, 0x4c, 0x34, 0x06, 0xda, 0x46, 0x83, 0x0d, 0x9e, 0xab, 0xc8, 0x63, 0x15, 0x69, 0xbc, 0x73,
	0xc3, 0x96, 0x16, 0xfe, 0x00, 0x50, 0x8c, 0x34, 0xa6, 0x59, 0x3d, 0x49, 0xc2, 0x38, 0x89, 0x99,
	0xd5, 0x11, 0xb1, 0x37, 0xba, 0xd8, 0xe4, 0xe8, 0x44, 0xf0, 0x4a, 0x25, 0x9f, 0xc9, 0xe4, 0xb1,
	0xd3, 0x0d, 0x4f, 0x84, 0x20, 0x01, 0x97, 0x4a, 0x3e, 0x11, 0xcb, 0x62, 0x9d, 0x8b, 0xf2, 0xad,
	0xae, 0x7c, 0xdf, 0x84, 0x83, 0x17, 0x2a, 0x7a, 0x25, 0xa3, 0x2d, 0x93, 0x1b, 0xb6, 0xcd, 0xf0,
	0x03, 0xe8, 0xca, 0x85, 0xb4, 0xba, 0x7d, 0x73, 0xd8, 0xbb, 0x73, 0xb4, 0xb7, 0x93, 0xf2, 0x4e,
	0x2d, 0x0f, 0xd5, 0x19, 0x38, 0x00, 0x97, 0x5f, 0x18, 0x2e, 0x03, 0xb9, 0x87, 0x9f, 0xc6, 0xd6,
	0xc3, 0xbe, 0x39, 0xec, 0x84, 0xed, 0x87, 0xc1, 0xe7, 0xd5, 0xd6, 0x31, 0xd7, 0x5b, 0xc7, 0xfc,
	0xb3, 0x75, 0xcc, 0x9f, 0x3b, 0xc7, 0x58, 0xef, 0x1c, 0xe3, 0xf7, 0xce, 0x31, 0xbe, 0x8e, 0x52,
	0xc2, 0xe7, 0x55, 0x54, 0x37, 0x7d, 0xd9, 0x7d, 0x4b, 0x67, 0x33, 0x12, 0x13, 0x94, 0xa9, 0xff,
	0xfe, 0x61, 0xeb, 0xf9, 0xb2, 0xc0, 0x2c, 0xea, 0x8a, 0x79, 0xbd, 0xfb, 0x17, 0x00, 0x00, 0xff,
	0xff, 0x51, 0x21, 0x0b, 0xb7, 0xb3, 0x03, 0x00, 0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.UserBiddingID != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.UserBiddingID))
		i--
		dAtA[i] = 0x38
	}
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	if len(m.AuctionParams) > 0 {
		for iNdEx := len(m.AuctionParams) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.AuctionParams[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x2a
		}
	}
	if len(m.ProtocolStatistics) > 0 {
		for iNdEx := len(m.ProtocolStatistics) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.ProtocolStatistics[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.DutchAuction) > 0 {
		for iNdEx := len(m.DutchAuction) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.DutchAuction[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.DebtAuction) > 0 {
		for iNdEx := len(m.DebtAuction) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.DebtAuction[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.SurplusAuction) > 0 {
		for iNdEx := len(m.SurplusAuction) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.SurplusAuction[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.SurplusAuction) > 0 {
		for _, e := range m.SurplusAuction {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.DebtAuction) > 0 {
		for _, e := range m.DebtAuction {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.DutchAuction) > 0 {
		for _, e := range m.DutchAuction {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.ProtocolStatistics) > 0 {
		for _, e := range m.ProtocolStatistics {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.AuctionParams) > 0 {
		for _, e := range m.AuctionParams {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	l = m.Params.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if m.UserBiddingID != 0 {
		n += 1 + sovGenesis(uint64(m.UserBiddingID))
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SurplusAuction", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SurplusAuction = append(m.SurplusAuction, SurplusAuction{})
			if err := m.SurplusAuction[len(m.SurplusAuction)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DebtAuction", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DebtAuction = append(m.DebtAuction, DebtAuction{})
			if err := m.DebtAuction[len(m.DebtAuction)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DutchAuction", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DutchAuction = append(m.DutchAuction, DutchAuction{})
			if err := m.DutchAuction[len(m.DutchAuction)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProtocolStatistics", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ProtocolStatistics = append(m.ProtocolStatistics, ProtocolStatistics{})
			if err := m.ProtocolStatistics[len(m.ProtocolStatistics)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AuctionParams", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AuctionParams = append(m.AuctionParams, AuctionParams{})
			if err := m.AuctionParams[len(m.AuctionParams)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field UserBiddingID", wireType)
			}
			m.UserBiddingID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.UserBiddingID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
