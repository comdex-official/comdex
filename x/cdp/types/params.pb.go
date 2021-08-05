// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: comdex/cdp/v1alpha1/params.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
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

type Params struct {
	CollateralParams []CollateralParam `protobuf:"bytes,1,rep,name=collateral_params,json=collateralParams,proto3" json:"collateral_params" yaml:"collateral_params"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_bf9893c9be4c20a1, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetCollateralParams() []CollateralParam {
	if m != nil {
		return m.CollateralParams
	}
	return nil
}

type CollateralParam struct {
	CollateralDenom  string                                 `protobuf:"bytes,1,opt,name=collateral_denom,json=collateralDenom,proto3" json:"collateral_denom,omitempty"`
	DebtDenom        string                                 `protobuf:"bytes,2,opt,name=debt_denom,json=debtDenom,proto3" json:"debt_denom,omitempty"`
	Type             string                                 `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
	LiquidationRatio github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,4,opt,name=liquidation_ratio,json=liquidationRatio,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"liquidation_ratio" yaml:"liquidation_ratio"`
}

func (m *CollateralParam) Reset()         { *m = CollateralParam{} }
func (m *CollateralParam) String() string { return proto.CompactTextString(m) }
func (*CollateralParam) ProtoMessage()    {}
func (*CollateralParam) Descriptor() ([]byte, []int) {
	return fileDescriptor_bf9893c9be4c20a1, []int{1}
}
func (m *CollateralParam) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *CollateralParam) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_CollateralParam.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *CollateralParam) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CollateralParam.Merge(m, src)
}
func (m *CollateralParam) XXX_Size() int {
	return m.Size()
}
func (m *CollateralParam) XXX_DiscardUnknown() {
	xxx_messageInfo_CollateralParam.DiscardUnknown(m)
}

var xxx_messageInfo_CollateralParam proto.InternalMessageInfo

func (m *CollateralParam) GetCollateralDenom() string {
	if m != nil {
		return m.CollateralDenom
	}
	return ""
}

func (m *CollateralParam) GetDebtDenom() string {
	if m != nil {
		return m.DebtDenom
	}
	return ""
}

func (m *CollateralParam) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func init() {
	proto.RegisterType((*Params)(nil), "comdex.cdp.v1alpha1.Params")
	proto.RegisterType((*CollateralParam)(nil), "comdex.cdp.v1alpha1.CollateralParam")
}

func init() { proto.RegisterFile("comdex/cdp/v1alpha1/params.proto", fileDescriptor_bf9893c9be4c20a1) }

var fileDescriptor_bf9893c9be4c20a1 = []byte{
	// 358 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x92, 0x41, 0x4f, 0xf2, 0x30,
	0x1c, 0xc6, 0xb7, 0x17, 0x42, 0x42, 0xdf, 0x03, 0xb0, 0xf7, 0x3d, 0x2c, 0x24, 0x8e, 0x65, 0x31,
	0x06, 0x0f, 0xac, 0x99, 0xde, 0x3c, 0x22, 0x89, 0x89, 0x27, 0xb3, 0xa3, 0x17, 0xd2, 0x75, 0x05,
	0x1a, 0x3b, 0x3a, 0xd7, 0x82, 0x72, 0xf0, 0x3b, 0xf8, 0xb1, 0x38, 0x72, 0x34, 0x1e, 0x88, 0x81,
	0x8b, 0x67, 0x3f, 0x81, 0x69, 0x8b, 0x02, 0xe2, 0x65, 0xed, 0x9e, 0xe7, 0xb7, 0xe7, 0xc9, 0xbf,
	0x2b, 0xf0, 0x31, 0xcf, 0x52, 0xf2, 0x08, 0x71, 0x9a, 0xc3, 0x69, 0x84, 0x58, 0x3e, 0x42, 0x11,
	0xcc, 0x51, 0x81, 0x32, 0x11, 0xe6, 0x05, 0x97, 0xdc, 0xf9, 0x67, 0x88, 0x10, 0xa7, 0x79, 0xf8,
	0x45, 0x34, 0xff, 0x0f, 0xf9, 0x90, 0x6b, 0x1f, 0xaa, 0x9d, 0x41, 0x9b, 0x1e, 0xe6, 0x22, 0xe3,
	0x02, 0x26, 0x48, 0x10, 0x38, 0x8d, 0x12, 0x22, 0x51, 0x04, 0x31, 0xa7, 0x63, 0xe3, 0x07, 0x4f,
	0xa0, 0x72, 0xa3, 0xa3, 0x1d, 0x01, 0x1a, 0x98, 0x33, 0x86, 0x24, 0x29, 0x10, 0xeb, 0x9b, 0x3e,
	0xd7, 0xf6, 0x4b, 0xed, 0xbf, 0x67, 0xc7, 0xe1, 0x2f, 0x85, 0xe1, 0xe5, 0x37, 0xad, 0x13, 0xba,
	0xfe, 0x7c, 0xd9, 0xb2, 0x3e, 0x96, 0x2d, 0x77, 0x86, 0x32, 0x76, 0x11, 0x1c, 0x84, 0x05, 0x71,
	0x1d, 0xef, 0x7f, 0x22, 0x82, 0x77, 0x1b, 0xd4, 0x7e, 0xe4, 0x38, 0xa7, 0x60, 0x87, 0xeb, 0xa7,
	0x64, 0xcc, 0x33, 0xd7, 0xf6, 0xed, 0x76, 0x35, 0xae, 0x6d, 0xf5, 0x9e, 0x92, 0x9d, 0x23, 0x00,
	0x52, 0x92, 0xc8, 0x0d, 0xf4, 0x47, 0x43, 0x55, 0xa5, 0x18, 0xdb, 0x01, 0x65, 0x39, 0xcb, 0x89,
	0x5b, 0xd2, 0x86, 0xde, 0x3b, 0x0f, 0xa0, 0xc1, 0xe8, 0xfd, 0x84, 0xa6, 0x48, 0x52, 0x3e, 0xee,
	0x17, 0x6a, 0x71, 0xcb, 0x0a, 0xe8, 0x5e, 0xab, 0x01, 0x5e, 0x97, 0xad, 0x93, 0x21, 0x95, 0xa3,
	0x49, 0xa2, 0x86, 0x86, 0x9b, 0xe3, 0x33, 0x4b, 0x47, 0xa4, 0x77, 0x50, 0xa5, 0x88, 0xb0, 0x47,
	0xf0, 0x76, 0xd4, 0x83, 0xc0, 0x20, 0xae, 0xef, 0x68, 0xb1, 0x7a, 0x76, 0xaf, 0xe6, 0x2b, 0xcf,
	0x5e, 0xac, 0x3c, 0xfb, 0x6d, 0xe5, 0xd9, 0xcf, 0x6b, 0xcf, 0x5a, 0xac, 0x3d, 0xeb, 0x65, 0xed,
	0x59, 0xb7, 0x9d, 0xbd, 0x3e, 0x75, 0xd0, 0x1d, 0x3e, 0x18, 0x50, 0x4c, 0x11, 0xdb, 0xbc, 0x43,
	0x73, 0x1b, 0x74, 0x75, 0x52, 0xd1, 0x7f, 0xee, 0xfc, 0x33, 0x00, 0x00, 0xff, 0xff, 0xf8, 0x4d,
	0x6d, 0xed, 0x28, 0x02, 0x00, 0x00,
}

func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.CollateralParams) > 0 {
		for iNdEx := len(m.CollateralParams) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.CollateralParams[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintParams(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *CollateralParam) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *CollateralParam) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *CollateralParam) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.LiquidationRatio.Size()
		i -= size
		if _, err := m.LiquidationRatio.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	if len(m.Type) > 0 {
		i -= len(m.Type)
		copy(dAtA[i:], m.Type)
		i = encodeVarintParams(dAtA, i, uint64(len(m.Type)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.DebtDenom) > 0 {
		i -= len(m.DebtDenom)
		copy(dAtA[i:], m.DebtDenom)
		i = encodeVarintParams(dAtA, i, uint64(len(m.DebtDenom)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.CollateralDenom) > 0 {
		i -= len(m.CollateralDenom)
		copy(dAtA[i:], m.CollateralDenom)
		i = encodeVarintParams(dAtA, i, uint64(len(m.CollateralDenom)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintParams(dAtA []byte, offset int, v uint64) int {
	offset -= sovParams(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.CollateralParams) > 0 {
		for _, e := range m.CollateralParams {
			l = e.Size()
			n += 1 + l + sovParams(uint64(l))
		}
	}
	return n
}

func (m *CollateralParam) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.CollateralDenom)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	l = len(m.DebtDenom)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	l = len(m.Type)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	l = m.LiquidationRatio.Size()
	n += 1 + l + sovParams(uint64(l))
	return n
}

func sovParams(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozParams(x uint64) (n int) {
	return sovParams(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParams
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
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CollateralParams", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CollateralParams = append(m.CollateralParams, CollateralParam{})
			if err := m.CollateralParams[len(m.CollateralParams)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipParams(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthParams
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthParams
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
func (m *CollateralParam) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParams
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
			return fmt.Errorf("proto: CollateralParam: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CollateralParam: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CollateralDenom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CollateralDenom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DebtDenom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DebtDenom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Type = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LiquidationRatio", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.LiquidationRatio.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipParams(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthParams
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthParams
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
func skipParams(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
				return 0, ErrInvalidLengthParams
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupParams
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthParams
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthParams        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowParams          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupParams = fmt.Errorf("proto: unexpected end of group")
)
