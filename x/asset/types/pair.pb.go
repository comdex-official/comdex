// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: petri/asset/v1beta1/pair.proto

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

func init() {
	proto.RegisterType((*Pair)(nil), "petri.asset.v1beta1.Pair")
	proto.RegisterType((*PairInfo)(nil), "petri.asset.v1beta1.PairInfo")
}

func init() { proto.RegisterFile("petri/asset/v1beta1/pair.proto", fileDescriptor_65bd24918e5ac160) }

var fileDescriptor_65bd24918e5ac160 = []byte{
	// 308 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4f, 0xce, 0xcf, 0x4d,
	0x49, 0xad, 0xd0, 0x4f, 0x2c, 0x2e, 0x4e, 0x2d, 0xd1, 0x2f, 0x33, 0x4c, 0x4a, 0x2d, 0x49, 0x34,
	0xd4, 0x2f, 0x48, 0xcc, 0x2c, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0x81, 0x28, 0xd0,
	0x03, 0x2b, 0xd0, 0x83, 0x2a, 0x90, 0x12, 0x49, 0xcf, 0x4f, 0xcf, 0x07, 0x2b, 0xd0, 0x07, 0xb1,
	0x20, 0x6a, 0x95, 0x2a, 0xb9, 0x58, 0x02, 0x12, 0x33, 0x8b, 0x84, 0xf8, 0xb8, 0x98, 0x32, 0x53,
	0x24, 0x18, 0x15, 0x18, 0x35, 0x58, 0x82, 0x98, 0x32, 0x53, 0x84, 0xf4, 0xb8, 0x38, 0xc0, 0xda,
	0xe3, 0x33, 0xf3, 0x24, 0x98, 0x40, 0xa2, 0x4e, 0xc2, 0x9f, 0xee, 0xc9, 0xf3, 0x57, 0x26, 0xe6,
	0xe6, 0x58, 0x29, 0xc1, 0x64, 0x94, 0x82, 0xd8, 0xc1, 0x4c, 0xcf, 0x3c, 0x21, 0x43, 0x2e, 0x4e,
	0x88, 0x68, 0x7e, 0x69, 0x89, 0x04, 0x33, 0x58, 0x83, 0xc8, 0xa7, 0x7b, 0xf2, 0x02, 0xc8, 0x1a,
	0xf2, 0x4b, 0x4b, 0x94, 0x82, 0x20, 0xc6, 0xfa, 0x97, 0x96, 0x28, 0xdd, 0x64, 0xe4, 0xe2, 0x00,
	0xd9, 0xed, 0x99, 0x97, 0x96, 0x4f, 0xb1, 0xfd, 0xda, 0x5c, 0x1c, 0x29, 0xa9, 0x79, 0xf9, 0xb9,
	0x20, 0xf5, 0x20, 0xeb, 0x39, 0x9d, 0x04, 0x3e, 0xdd, 0x93, 0xe7, 0x81, 0xa8, 0x07, 0xcb, 0x28,
	0x05, 0xb1, 0x83, 0x69, 0x74, 0xc7, 0xb2, 0x10, 0xe3, 0x58, 0x21, 0x5d, 0x2e, 0x4e, 0x88, 0xf9,
	0x20, 0x2d, 0xac, 0x38, 0x2c, 0x80, 0x38, 0xc1, 0xbf, 0xb4, 0xc4, 0x29, 0xf0, 0xc4, 0x43, 0x39,
	0x86, 0x15, 0x8f, 0xe4, 0x18, 0x4e, 0x3c, 0x92, 0x63, 0xbc, 0xf0, 0x48, 0x8e, 0xf1, 0xc1, 0x23,
	0x39, 0xc6, 0x09, 0x8f, 0xe5, 0x18, 0x2e, 0x3c, 0x96, 0x63, 0xb8, 0xf1, 0x58, 0x8e, 0x21, 0x4a,
	0x3f, 0x3d, 0xb3, 0x24, 0xa3, 0x34, 0x49, 0x2f, 0x39, 0x3f, 0x57, 0x1f, 0x12, 0x5f, 0xba, 0xf9,
	0x69, 0x69, 0x99, 0xc9, 0x99, 0x89, 0x39, 0x50, 0xbe, 0x3e, 0x2c, 0x8a, 0x4b, 0x2a, 0x0b, 0x52,
	0x8b, 0x93, 0xd8, 0xc0, 0x11, 0x66, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0xb1, 0xe6, 0xaf, 0xf4,
	0xff, 0x01, 0x00, 0x00,
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
