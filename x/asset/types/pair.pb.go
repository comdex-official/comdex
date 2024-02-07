// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: comdex/asset/v1beta1/pair.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
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

type Pair struct {
	Id       uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	AssetIn  uint64 `protobuf:"varint,2,opt,name=asset_in,json=assetIn,proto3" json:"asset_in,omitempty" yaml:"asset_in"`
	AssetOut uint64 `protobuf:"varint,3,opt,name=asset_out,json=assetOut,proto3" json:"asset_out,omitempty" yaml:"asset_out"`
}

func (m *Pair) Reset()         { *m = Pair{} }
func (m *Pair) String() string { return proto.CompactTextString(m) }
func (*Pair) ProtoMessage()    {}
func (*Pair) Descriptor() ([]byte, []int) {
	return fileDescriptor_65bd24918e5ac160, []int{0}
}
func (m *Pair) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Pair) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Pair.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Pair) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Pair.Merge(m, src)
}
func (m *Pair) XXX_Size() int {
	return m.Size()
}
func (m *Pair) XXX_DiscardUnknown() {
	xxx_messageInfo_Pair.DiscardUnknown(m)
}

var xxx_messageInfo_Pair proto.InternalMessageInfo

type PairInfo struct {
	Id       uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	AssetIn  uint64 `protobuf:"varint,2,opt,name=asset_in,json=assetIn,proto3" json:"asset_in,omitempty" yaml:"asset_in"`
	DenomIn  string `protobuf:"bytes,3,opt,name=denom_in,json=denomIn,proto3" json:"denom_in,omitempty" yaml:"denom"`
	AssetOut uint64 `protobuf:"varint,4,opt,name=asset_out,json=assetOut,proto3" json:"asset_out,omitempty" yaml:"asset_out"`
	DenomOut string `protobuf:"bytes,5,opt,name=denom_out,json=denomOut,proto3" json:"denom_out,omitempty" yaml:"denom"`
}

func (m *PairInfo) Reset()         { *m = PairInfo{} }
func (m *PairInfo) String() string { return proto.CompactTextString(m) }
func (*PairInfo) ProtoMessage()    {}
func (*PairInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_65bd24918e5ac160, []int{1}
}
func (m *PairInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PairInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PairInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PairInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PairInfo.Merge(m, src)
}
func (m *PairInfo) XXX_Size() int {
	return m.Size()
}
func (m *PairInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_PairInfo.DiscardUnknown(m)
}

var xxx_messageInfo_PairInfo proto.InternalMessageInfo

type AssetPair struct {
	Id                    uint64                `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                  string                `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty" yaml:"name"`
	Denom                 string                `protobuf:"bytes,3,opt,name=denom,proto3" json:"denom,omitempty" yaml:"denom"`
	Decimals              cosmossdk_io_math.Int `protobuf:"bytes,4,opt,name=decimals,proto3,customtype=cosmossdk.io/math.Int" json:"decimals" yaml:"decimals"`
	IsOnChain             bool                  `protobuf:"varint,5,opt,name=is_on_chain,json=isOnChain,proto3" json:"is_on_chain,omitempty" yaml:"is_on_chain"`
	IsOraclePriceRequired bool                  `protobuf:"varint,6,opt,name=is_oracle_price_required,json=isOraclePriceRequired,proto3" json:"is_oracle_price_required,omitempty" yaml:"is_oracle_price_required"`
	IsCdpMintable         bool                  `protobuf:"varint,7,opt,name=is_cdp_mintable,json=isCdpMintable,proto3" json:"is_cdp_mintable,omitempty" yaml:"is_cdp_mintable"`
	AssetOut              uint64                `protobuf:"varint,8,opt,name=asset_out,json=assetOut,proto3" json:"asset_out,omitempty" yaml:"asset_out"`
}

func (m *AssetPair) Reset()         { *m = AssetPair{} }
func (m *AssetPair) String() string { return proto.CompactTextString(m) }
func (*AssetPair) ProtoMessage()    {}
func (*AssetPair) Descriptor() ([]byte, []int) {
	return fileDescriptor_65bd24918e5ac160, []int{2}
}
func (m *AssetPair) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AssetPair) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AssetPair.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AssetPair) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AssetPair.Merge(m, src)
}
func (m *AssetPair) XXX_Size() int {
	return m.Size()
}
func (m *AssetPair) XXX_DiscardUnknown() {
	xxx_messageInfo_AssetPair.DiscardUnknown(m)
}

var xxx_messageInfo_AssetPair proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Pair)(nil), "comdex.asset.v1beta1.Pair")
	proto.RegisterType((*PairInfo)(nil), "comdex.asset.v1beta1.PairInfo")
	proto.RegisterType((*AssetPair)(nil), "comdex.asset.v1beta1.AssetPair")
}

func init() { proto.RegisterFile("comdex/asset/v1beta1/pair.proto", fileDescriptor_65bd24918e5ac160) }

var fileDescriptor_65bd24918e5ac160 = []byte{
	// 516 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x53, 0x4f, 0x6f, 0xd3, 0x3c,
	0x18, 0x4f, 0xba, 0x6e, 0x6d, 0xbc, 0xf7, 0xa5, 0x93, 0xe9, 0xa6, 0x68, 0x87, 0x18, 0x65, 0x12,
	0x42, 0x42, 0x4b, 0xa8, 0x90, 0x38, 0x70, 0x23, 0x3b, 0x55, 0x02, 0x75, 0xe4, 0x88, 0x90, 0x22,
	0x37, 0x71, 0x5b, 0x8b, 0xc6, 0x0e, 0xb1, 0x83, 0xe8, 0xb7, 0xe0, 0x63, 0xf0, 0x3d, 0xb8, 0xf4,
	0xb8, 0x23, 0x70, 0x88, 0xa0, 0xfd, 0x06, 0xf9, 0x04, 0xc8, 0x76, 0xcb, 0x0a, 0xda, 0xa4, 0x49,
	0x9c, 0xf2, 0xe4, 0xf9, 0xfd, 0x79, 0x9e, 0x9f, 0x2d, 0x03, 0x94, 0xf2, 0x3c, 0x23, 0x1f, 0x43,
	0x2c, 0x04, 0x91, 0xe1, 0x87, 0xc1, 0x98, 0x48, 0x3c, 0x08, 0x0b, 0x4c, 0xcb, 0xa0, 0x28, 0xb9,
	0xe4, 0xb0, 0x6f, 0x08, 0x81, 0x26, 0x04, 0x1b, 0xc2, 0x69, 0x7f, 0xca, 0xa7, 0x5c, 0x13, 0x42,
	0x55, 0x19, 0xae, 0xbf, 0x00, 0xed, 0x4b, 0x4c, 0x4b, 0x78, 0x0f, 0xb4, 0x68, 0xe6, 0xda, 0x0f,
	0xec, 0x47, 0xed, 0xb8, 0x45, 0x33, 0x18, 0x80, 0xae, 0x96, 0x27, 0x94, 0xb9, 0x2d, 0xd5, 0x8d,
	0xee, 0x37, 0x35, 0xea, 0x2d, 0x70, 0x3e, 0x7f, 0xee, 0x6f, 0x11, 0x3f, 0xee, 0xe8, 0x72, 0xc8,
	0xe0, 0x00, 0x38, 0xa6, 0xcb, 0x2b, 0xe9, 0xee, 0x69, 0x41, 0xbf, 0xa9, 0xd1, 0xd1, 0xae, 0x80,
	0x57, 0xd2, 0x8f, 0x8d, 0xed, 0xa8, 0x92, 0xfe, 0x37, 0x1b, 0x74, 0xd5, 0xec, 0x21, 0x9b, 0xf0,
	0x7f, 0x9e, 0xff, 0x18, 0x74, 0x33, 0xc2, 0x78, 0xae, 0xf8, 0x6a, 0xbc, 0x13, 0x1d, 0x35, 0x35,
	0xfa, 0xcf, 0xf0, 0x35, 0xe2, 0xc7, 0x1d, 0xfd, 0xfd, 0x7b, 0xd9, 0xf6, 0x5d, 0x96, 0x85, 0xe7,
	0xc0, 0x31, 0xfe, 0x4a, 0xb2, 0x7f, 0xcb, 0x00, 0xb3, 0x82, 0xca, 0xf6, 0x65, 0x0f, 0x38, 0x2f,
	0x94, 0xf6, 0xc6, 0xc3, 0x3d, 0x03, 0x6d, 0x86, 0x73, 0xa2, 0x83, 0x39, 0x51, 0xaf, 0xa9, 0xd1,
	0xa1, 0xf1, 0x51, 0x5d, 0x3f, 0xd6, 0x20, 0x7c, 0x08, 0xf6, 0xb5, 0xdd, 0xad, 0x71, 0x0c, 0x0c,
	0x5f, 0xaa, 0xe4, 0x29, 0xcd, 0xf1, 0x5c, 0xe8, 0x2c, 0x4e, 0xf4, 0x64, 0x59, 0x23, 0xeb, 0x7b,
	0x8d, 0x8e, 0x53, 0x2e, 0x72, 0x2e, 0x44, 0xf6, 0x2e, 0xa0, 0x3c, 0xcc, 0xb1, 0x9c, 0x05, 0x43,
	0x26, 0xaf, 0x8f, 0x71, 0x2b, 0xd3, 0x8b, 0x9b, 0x12, 0x3e, 0x03, 0x87, 0x54, 0x24, 0x9c, 0x25,
	0xe9, 0x0c, 0x53, 0xa6, 0x93, 0x76, 0xa3, 0x93, 0xa6, 0x46, 0xd0, 0x68, 0x76, 0x40, 0x3f, 0x76,
	0xa8, 0x18, 0xb1, 0x0b, 0x55, 0xc3, 0xb7, 0xc0, 0x55, 0x50, 0x89, 0xd3, 0x39, 0x49, 0x8a, 0x92,
	0xa6, 0x24, 0x29, 0xc9, 0xfb, 0x8a, 0x96, 0x24, 0x73, 0x0f, 0xb4, 0xc9, 0x59, 0x53, 0x23, 0x74,
	0x6d, 0x72, 0x13, 0xd3, 0x8f, 0x8f, 0xa9, 0x18, 0x69, 0xe4, 0x52, 0x01, 0xf1, 0xa6, 0x0f, 0x23,
	0xd0, 0xa3, 0x22, 0x49, 0xb3, 0x22, 0xc9, 0x29, 0x93, 0x78, 0x3c, 0x27, 0x6e, 0x47, 0x9b, 0x9e,
	0x36, 0x35, 0x3a, 0xf9, 0x6d, 0xba, 0x4b, 0xf0, 0xe3, 0xff, 0xa9, 0xb8, 0xc8, 0x8a, 0x57, 0x9b,
	0xff, 0x3f, 0x2f, 0xbd, 0x7b, 0x97, 0x4b, 0x8f, 0x5e, 0x2f, 0x7f, 0x7a, 0xd6, 0xe7, 0x95, 0x67,
	0x2d, 0x57, 0x9e, 0x7d, 0xb5, 0xf2, 0xec, 0x1f, 0x2b, 0xcf, 0xfe, 0xb4, 0xf6, 0xac, 0xab, 0xb5,
	0x67, 0x7d, 0x5d, 0x7b, 0xd6, 0x9b, 0x70, 0x4a, 0xe5, 0xac, 0x1a, 0x07, 0x29, 0xcf, 0x43, 0xf3,
	0xea, 0xce, 0xf9, 0x64, 0x42, 0x53, 0x8a, 0xe7, 0x9b, 0xff, 0x70, 0xfb, 0x50, 0xe5, 0xa2, 0x20,
	0x62, 0x7c, 0xa0, 0x9f, 0xdd, 0xd3, 0x5f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x81, 0x40, 0xde, 0xdb,
	0xc5, 0x03, 0x00, 0x00,
}

func (m *Pair) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Pair) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Pair) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.AssetOut != 0 {
		i = encodeVarintPair(dAtA, i, uint64(m.AssetOut))
		i--
		dAtA[i] = 0x18
	}
	if m.AssetIn != 0 {
		i = encodeVarintPair(dAtA, i, uint64(m.AssetIn))
		i--
		dAtA[i] = 0x10
	}
	if m.Id != 0 {
		i = encodeVarintPair(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *PairInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PairInfo) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PairInfo) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.DenomOut) > 0 {
		i -= len(m.DenomOut)
		copy(dAtA[i:], m.DenomOut)
		i = encodeVarintPair(dAtA, i, uint64(len(m.DenomOut)))
		i--
		dAtA[i] = 0x2a
	}
	if m.AssetOut != 0 {
		i = encodeVarintPair(dAtA, i, uint64(m.AssetOut))
		i--
		dAtA[i] = 0x20
	}
	if len(m.DenomIn) > 0 {
		i -= len(m.DenomIn)
		copy(dAtA[i:], m.DenomIn)
		i = encodeVarintPair(dAtA, i, uint64(len(m.DenomIn)))
		i--
		dAtA[i] = 0x1a
	}
	if m.AssetIn != 0 {
		i = encodeVarintPair(dAtA, i, uint64(m.AssetIn))
		i--
		dAtA[i] = 0x10
	}
	if m.Id != 0 {
		i = encodeVarintPair(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *AssetPair) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AssetPair) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AssetPair) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.AssetOut != 0 {
		i = encodeVarintPair(dAtA, i, uint64(m.AssetOut))
		i--
		dAtA[i] = 0x40
	}
	if m.IsCdpMintable {
		i--
		if m.IsCdpMintable {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x38
	}
	if m.IsOraclePriceRequired {
		i--
		if m.IsOraclePriceRequired {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x30
	}
	if m.IsOnChain {
		i--
		if m.IsOnChain {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x28
	}
	{
		size := m.Decimals.Size()
		i -= size
		if _, err := m.Decimals.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPair(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintPair(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintPair(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0x12
	}
	if m.Id != 0 {
		i = encodeVarintPair(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintPair(dAtA []byte, offset int, v uint64) int {
	offset -= sovPair(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Pair) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovPair(uint64(m.Id))
	}
	if m.AssetIn != 0 {
		n += 1 + sovPair(uint64(m.AssetIn))
	}
	if m.AssetOut != 0 {
		n += 1 + sovPair(uint64(m.AssetOut))
	}
	return n
}

func (m *PairInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovPair(uint64(m.Id))
	}
	if m.AssetIn != 0 {
		n += 1 + sovPair(uint64(m.AssetIn))
	}
	l = len(m.DenomIn)
	if l > 0 {
		n += 1 + l + sovPair(uint64(l))
	}
	if m.AssetOut != 0 {
		n += 1 + sovPair(uint64(m.AssetOut))
	}
	l = len(m.DenomOut)
	if l > 0 {
		n += 1 + l + sovPair(uint64(l))
	}
	return n
}

func (m *AssetPair) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovPair(uint64(m.Id))
	}
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovPair(uint64(l))
	}
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovPair(uint64(l))
	}
	l = m.Decimals.Size()
	n += 1 + l + sovPair(uint64(l))
	if m.IsOnChain {
		n += 2
	}
	if m.IsOraclePriceRequired {
		n += 2
	}
	if m.IsCdpMintable {
		n += 2
	}
	if m.AssetOut != 0 {
		n += 1 + sovPair(uint64(m.AssetOut))
	}
	return n
}

func sovPair(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozPair(x uint64) (n int) {
	return sovPair(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Pair) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPair
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
			return fmt.Errorf("proto: Pair: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Pair: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPair
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AssetIn", wireType)
			}
			m.AssetIn = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPair
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AssetIn |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AssetOut", wireType)
			}
			m.AssetOut = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPair
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AssetOut |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipPair(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPair
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
func (m *PairInfo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPair
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
			return fmt.Errorf("proto: PairInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PairInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPair
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AssetIn", wireType)
			}
			m.AssetIn = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPair
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AssetIn |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DenomIn", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPair
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPair
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPair
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DenomIn = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AssetOut", wireType)
			}
			m.AssetOut = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPair
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AssetOut |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DenomOut", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPair
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPair
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPair
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DenomOut = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPair(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPair
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
func (m *AssetPair) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPair
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
			return fmt.Errorf("proto: AssetPair: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AssetPair: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPair
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPair
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPair
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPair
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPair
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPair
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPair
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Decimals", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPair
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPair
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPair
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Decimals.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsOnChain", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPair
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.IsOnChain = bool(v != 0)
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsOraclePriceRequired", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPair
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.IsOraclePriceRequired = bool(v != 0)
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsCdpMintable", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPair
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.IsCdpMintable = bool(v != 0)
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AssetOut", wireType)
			}
			m.AssetOut = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPair
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AssetOut |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipPair(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPair
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
func skipPair(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowPair
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
					return 0, ErrIntOverflowPair
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
					return 0, ErrIntOverflowPair
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
				return 0, ErrInvalidLengthPair
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupPair
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthPair
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthPair        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowPair          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupPair = fmt.Errorf("proto: unexpected end of group")
)
