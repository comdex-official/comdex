// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: comdex/asset/v1beta1/extendedPairVault.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	github_com_cosmos_gogoproto_types "github.com/cosmos/gogoproto/types"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type ExtendedPairVault struct {
	Id                  uint64                      `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	AppId               uint64                      `protobuf:"varint,2,opt,name=app_id,json=appId,proto3" json:"app_id,omitempty" yaml:"app_id"`
	PairId              uint64                      `protobuf:"varint,3,opt,name=pair_id,json=pairId,proto3" json:"pair_id,omitempty" yaml:"pair_id"`
	StabilityFee        cosmossdk_io_math.LegacyDec `protobuf:"bytes,4,opt,name=stability_fee,json=stabilityFee,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"stability_fee" yaml:"stability_fee"`
	ClosingFee          cosmossdk_io_math.LegacyDec `protobuf:"bytes,5,opt,name=closing_fee,json=closingFee,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"closing_fee" yaml:"closing_fee"`
	LiquidationPenalty  cosmossdk_io_math.LegacyDec `protobuf:"bytes,6,opt,name=liquidation_penalty,json=liquidationPenalty,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"liquidation_penalty" yaml:"liquidation_penalty"`
	DrawDownFee         cosmossdk_io_math.LegacyDec `protobuf:"bytes,7,opt,name=draw_down_fee,json=drawDownFee,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"draw_down_fee" yaml:"draw_down_fee"`
	IsVaultActive       bool                        `protobuf:"varint,8,opt,name=is_vault_active,json=isVaultActive,proto3" json:"is_vault_active,omitempty" yaml:"active_flag"`
	DebtCeiling         cosmossdk_io_math.Int       `protobuf:"bytes,9,opt,name=debt_ceiling,json=debtCeiling,proto3,customtype=cosmossdk.io/math.Int" json:"debt_ceiling" yaml:"debt_ceiling"`
	DebtFloor           cosmossdk_io_math.Int       `protobuf:"bytes,10,opt,name=debt_floor,json=debtFloor,proto3,customtype=cosmossdk.io/math.Int" json:"debt_floor" yaml:"debt_floor"`
	IsStableMintVault   bool                        `protobuf:"varint,11,opt,name=is_stable_mint_vault,json=isStableMintVault,proto3" json:"is_stable_mint_vault,omitempty" yaml:"is_stable_mint_vault"`
	MinCr               cosmossdk_io_math.LegacyDec `protobuf:"bytes,12,opt,name=min_cr,json=minCr,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"min_cr" yaml:"min_cr"`
	PairName            string                      `protobuf:"bytes,13,opt,name=pair_name,json=pairName,proto3" json:"pair_name,omitempty" yaml:"pair_name"`
	AssetOutOraclePrice bool                        `protobuf:"varint,14,opt,name=asset_out_oracle_price,json=assetOutOraclePrice,proto3" json:"asset_out_oracle_price,omitempty" yaml:"asset_out_oracle_price"`
	AssetOutPrice       uint64                      `protobuf:"varint,15,opt,name=asset_out_price,json=assetOutPrice,proto3" json:"asset_out_price,omitempty" yaml:"asset_out_price"`
	MinUsdValueLeft     uint64                      `protobuf:"varint,16,opt,name=min_usd_value_left,json=minUsdValueLeft,proto3" json:"min_usd_value_left,omitempty" yaml:"min_usd_value_left"`
	BlockHeight         int64                       `protobuf:"varint,17,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty" yaml:"block_height"`
	BlockTime           time.Time                   `protobuf:"bytes,18,opt,name=block_time,json=blockTime,proto3,stdtime" json:"block_time" yaml:"block_time"`
}

func (m *ExtendedPairVault) Reset()         { *m = ExtendedPairVault{} }
func (m *ExtendedPairVault) String() string { return proto.CompactTextString(m) }
func (*ExtendedPairVault) ProtoMessage()    {}
func (*ExtendedPairVault) Descriptor() ([]byte, []int) {
	return fileDescriptor_23dd38fcddb231cd, []int{0}
}
func (m *ExtendedPairVault) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ExtendedPairVault) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ExtendedPairVault.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ExtendedPairVault) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExtendedPairVault.Merge(m, src)
}
func (m *ExtendedPairVault) XXX_Size() int {
	return m.Size()
}
func (m *ExtendedPairVault) XXX_DiscardUnknown() {
	xxx_messageInfo_ExtendedPairVault.DiscardUnknown(m)
}

var xxx_messageInfo_ExtendedPairVault proto.InternalMessageInfo

func init() {
	proto.RegisterType((*ExtendedPairVault)(nil), "comdex.asset.v1beta1.ExtendedPairVault")
}

func init() {
	proto.RegisterFile("comdex/asset/v1beta1/extendedPairVault.proto", fileDescriptor_23dd38fcddb231cd)
}

var fileDescriptor_23dd38fcddb231cd = []byte{
	// 824 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x55, 0x51, 0x73, 0xdb, 0x44,
	0x10, 0xb6, 0xd2, 0xc6, 0x8d, 0xcf, 0x71, 0x52, 0x5f, 0x4c, 0x10, 0xee, 0xc4, 0x32, 0xf7, 0xe4,
	0x19, 0x40, 0x9a, 0x14, 0x5e, 0x28, 0x33, 0xcc, 0xd4, 0x2d, 0x1d, 0x02, 0x85, 0xa6, 0x07, 0x04,
	0xa6, 0x2f, 0xc7, 0x59, 0x3a, 0x2b, 0x37, 0x3d, 0xe9, 0x84, 0x74, 0x4a, 0xea, 0x7f, 0xd1, 0x9f,
	0xc1, 0x4f, 0x09, 0x6f, 0x7d, 0x64, 0x78, 0x10, 0x90, 0xfc, 0x03, 0xfd, 0x02, 0xe6, 0xee, 0x14,
	0xec, 0x94, 0xcc, 0xb4, 0xc3, 0xdb, 0xed, 0xb7, 0xdf, 0x7e, 0xb7, 0xab, 0xdb, 0x5d, 0x81, 0x0f,
	0x43, 0x99, 0x44, 0xec, 0x45, 0x40, 0x8b, 0x82, 0xa9, 0xe0, 0x64, 0x7f, 0xc6, 0x14, 0xdd, 0x0f,
	0xd8, 0x0b, 0xc5, 0xd2, 0x88, 0x45, 0x87, 0x94, 0xe7, 0x47, 0xb4, 0x14, 0xca, 0xcf, 0x72, 0xa9,
	0x24, 0x1c, 0x58, 0xb6, 0x6f, 0xd8, 0x7e, 0xc3, 0x1e, 0x0e, 0x62, 0x19, 0x4b, 0x43, 0x08, 0xf4,
	0xc9, 0x72, 0x87, 0x5e, 0x2c, 0x65, 0x2c, 0x58, 0x60, 0xac, 0x59, 0x39, 0x0f, 0x14, 0x4f, 0x58,
	0xa1, 0x68, 0x92, 0x59, 0x02, 0xfa, 0x0d, 0x80, 0xfe, 0x17, 0xaf, 0x5f, 0x04, 0xb7, 0xc0, 0x1a,
	0x8f, 0x5c, 0x67, 0xec, 0x4c, 0x6e, 0xe2, 0x35, 0x1e, 0xc1, 0x09, 0x68, 0xd3, 0x2c, 0x23, 0x3c,
	0x72, 0xd7, 0x34, 0x36, 0xed, 0xd7, 0x95, 0xd7, 0x5b, 0xd0, 0x44, 0xdc, 0x43, 0x16, 0x47, 0x78,
	0x9d, 0x66, 0xd9, 0x41, 0x04, 0x3f, 0x00, 0xb7, 0x32, 0xca, 0x73, 0x4d, 0xbd, 0x61, 0xa8, 0xb0,
	0xae, 0xbc, 0x2d, 0x4b, 0x6d, 0x1c, 0x08, 0xb7, 0xf5, 0xe9, 0x20, 0x82, 0x3f, 0x83, 0x5e, 0xa1,
	0xe8, 0x8c, 0x0b, 0xae, 0x16, 0x64, 0xce, 0x98, 0x7b, 0x73, 0xec, 0x4c, 0x3a, 0xd3, 0xcf, 0xce,
	0x2a, 0xaf, 0xf5, 0x47, 0xe5, 0xdd, 0x09, 0x65, 0x91, 0xc8, 0xa2, 0x88, 0x9e, 0xfb, 0x5c, 0x06,
	0x09, 0x55, 0xc7, 0xfe, 0x63, 0x16, 0xd3, 0x70, 0xf1, 0x90, 0x85, 0x75, 0xe5, 0x0d, 0xac, 0xea,
	0x15, 0x05, 0x84, 0x37, 0xff, 0xb5, 0x1f, 0x31, 0x06, 0x9f, 0x81, 0x6e, 0x28, 0x64, 0xc1, 0xd3,
	0xd8, 0xe8, 0xaf, 0x1b, 0xfd, 0x4f, 0xdf, 0x4e, 0x1f, 0x5a, 0xfd, 0x95, 0x78, 0x84, 0x41, 0x63,
	0x69, 0xed, 0x1c, 0xec, 0x08, 0xfe, 0x4b, 0xc9, 0x23, 0xaa, 0xb8, 0x4c, 0x49, 0xc6, 0x52, 0x2a,
	0xd4, 0xc2, 0x6d, 0x9b, 0x3b, 0xee, 0xbf, 0xdd, 0x1d, 0x43, 0x7b, 0xc7, 0x35, 0x3a, 0x08, 0xc3,
	0x15, 0xf4, 0xd0, 0x82, 0x90, 0x80, 0x5e, 0x94, 0xd3, 0x53, 0x12, 0xc9, 0xd3, 0xd4, 0x54, 0x74,
	0xeb, 0x7f, 0x7c, 0xb1, 0x2b, 0x0a, 0x08, 0x77, 0xb5, 0xfd, 0x50, 0x9e, 0xa6, 0xba, 0xa8, 0xcf,
	0xc1, 0x36, 0x2f, 0xc8, 0x89, 0xee, 0x02, 0x42, 0x43, 0xc5, 0x4f, 0x98, 0xbb, 0x31, 0x76, 0x26,
	0x1b, 0xd3, 0xdd, 0xe5, 0x17, 0xb1, 0x38, 0x99, 0x0b, 0x1a, 0x23, 0xdc, 0xe3, 0x85, 0xe9, 0x99,
	0xfb, 0x06, 0x84, 0x3f, 0x82, 0xcd, 0x88, 0xcd, 0x14, 0x09, 0x19, 0x17, 0x3c, 0x8d, 0xdd, 0x8e,
	0xc9, 0xef, 0x93, 0x26, 0xbf, 0x77, 0xfe, 0x9b, 0xdf, 0x41, 0xaa, 0xea, 0xca, 0xdb, 0x69, 0x32,
	0x5b, 0x09, 0xd5, 0x89, 0xb1, 0x99, 0x7a, 0x60, 0x2d, 0xf8, 0x14, 0x00, 0xe3, 0x9d, 0x0b, 0x29,
	0x73, 0x17, 0x18, 0xd9, 0xbb, 0x6f, 0x92, 0xed, 0xaf, 0xc8, 0x9a, 0x40, 0x84, 0x3b, 0xda, 0x78,
	0xa4, 0xcf, 0xf0, 0x10, 0x0c, 0x78, 0x41, 0x74, 0xbf, 0x08, 0x46, 0x12, 0x9e, 0x2a, 0x5b, 0xb7,
	0xdb, 0x35, 0x05, 0x7b, 0x75, 0xe5, 0xdd, 0xb1, 0xf1, 0xd7, 0xb1, 0x10, 0xee, 0xf3, 0xe2, 0x3b,
	0x83, 0x7e, 0xc3, 0x53, 0x65, 0xe7, 0xe6, 0x6b, 0xd0, 0x4e, 0x78, 0x4a, 0xc2, 0xdc, 0xdd, 0xbc,
	0x52, 0xf7, 0x1b, 0xde, 0xa5, 0x19, 0x25, 0x1b, 0x8a, 0xf0, 0x7a, 0xc2, 0xd3, 0x07, 0x39, 0xdc,
	0x07, 0x1d, 0x33, 0x31, 0x29, 0x4d, 0x98, 0xdb, 0x33, 0x7a, 0x83, 0xba, 0xf2, 0x6e, 0xaf, 0x0c,
	0x93, 0x76, 0x21, 0xbc, 0xa1, 0xcf, 0xdf, 0xd2, 0x84, 0xc1, 0x23, 0xb0, 0x6b, 0xb6, 0x02, 0x91,
	0xa5, 0x22, 0x32, 0xa7, 0xa1, 0x60, 0x24, 0xcb, 0x79, 0xc8, 0xdc, 0x2d, 0x53, 0xd3, 0xfb, 0x75,
	0xe5, 0xed, 0x35, 0x8f, 0x78, 0x2d, 0x0f, 0xe1, 0x1d, 0xe3, 0x78, 0x52, 0xaa, 0x27, 0x06, 0x3e,
	0xd4, 0x28, 0x9c, 0x82, 0xed, 0x25, 0xdf, 0x0a, 0x6e, 0x9b, 0xe9, 0x1e, 0xd6, 0x95, 0xb7, 0xfb,
	0xba, 0x60, 0xa3, 0xd4, 0xbb, 0x54, 0xb2, 0x1a, 0x5f, 0x01, 0xa8, 0x0b, 0x2c, 0x8b, 0x88, 0x9c,
	0x50, 0x51, 0x32, 0x22, 0xd8, 0x5c, 0xb9, 0xb7, 0x8d, 0xcc, 0x5e, 0x5d, 0x79, 0xef, 0x2d, 0x3f,
	0xc2, 0x55, 0x0e, 0xc2, 0xdb, 0x09, 0x4f, 0x7f, 0x28, 0xa2, 0x23, 0x0d, 0x3d, 0x66, 0x73, 0x05,
	0xef, 0x81, 0xcd, 0x99, 0x90, 0xe1, 0x73, 0x72, 0xcc, 0x78, 0x7c, 0xac, 0xdc, 0xfe, 0xd8, 0x99,
	0xdc, 0x98, 0xbe, 0xbb, 0x6c, 0xa4, 0x55, 0x2f, 0xc2, 0x5d, 0x63, 0x7e, 0x69, 0x2c, 0xf8, 0x13,
	0x00, 0xd6, 0xab, 0x57, 0xa1, 0x0b, 0xc7, 0xce, 0xa4, 0x7b, 0x77, 0xe8, 0xdb, 0x3d, 0xe9, 0x5f,
	0xee, 0x49, 0xff, 0xfb, 0xcb, 0x3d, 0x39, 0xdd, 0xd3, 0x6f, 0xb8, 0xec, 0xa5, 0x65, 0x2c, 0x7a,
	0xf9, 0xa7, 0xe7, 0xe0, 0x8e, 0x01, 0x34, 0x7d, 0xfa, 0xf4, 0xec, 0xef, 0x51, 0xeb, 0xd7, 0xf3,
	0x51, 0xeb, 0xec, 0x7c, 0xe4, 0xbc, 0x3a, 0x1f, 0x39, 0x7f, 0x9d, 0x8f, 0x9c, 0x97, 0x17, 0xa3,
	0xd6, 0xab, 0x8b, 0x51, 0xeb, 0xf7, 0x8b, 0x51, 0xeb, 0x59, 0x10, 0x73, 0x75, 0x5c, 0xce, 0xfc,
	0x50, 0x26, 0x81, 0xdd, 0xe2, 0x1f, 0xc9, 0xf9, 0x9c, 0x87, 0x9c, 0x8a, 0xc6, 0x0e, 0x2e, 0xff,
	0x02, 0x6a, 0x91, 0xb1, 0x62, 0xd6, 0x36, 0x09, 0x7d, 0xfc, 0x4f, 0x00, 0x00, 0x00, 0xff, 0xff,
	0xaa, 0x0d, 0x8b, 0x6c, 0x22, 0x06, 0x00, 0x00,
}

func (m *ExtendedPairVault) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ExtendedPairVault) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ExtendedPairVault) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	n1, err1 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(m.BlockTime, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.BlockTime):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintExtendedPairVault(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x1
	i--
	dAtA[i] = 0x92
	if m.BlockHeight != 0 {
		i = encodeVarintExtendedPairVault(dAtA, i, uint64(m.BlockHeight))
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x88
	}
	if m.MinUsdValueLeft != 0 {
		i = encodeVarintExtendedPairVault(dAtA, i, uint64(m.MinUsdValueLeft))
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x80
	}
	if m.AssetOutPrice != 0 {
		i = encodeVarintExtendedPairVault(dAtA, i, uint64(m.AssetOutPrice))
		i--
		dAtA[i] = 0x78
	}
	if m.AssetOutOraclePrice {
		i--
		if m.AssetOutOraclePrice {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x70
	}
	if len(m.PairName) > 0 {
		i -= len(m.PairName)
		copy(dAtA[i:], m.PairName)
		i = encodeVarintExtendedPairVault(dAtA, i, uint64(len(m.PairName)))
		i--
		dAtA[i] = 0x6a
	}
	{
		size := m.MinCr.Size()
		i -= size
		if _, err := m.MinCr.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintExtendedPairVault(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x62
	if m.IsStableMintVault {
		i--
		if m.IsStableMintVault {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x58
	}
	{
		size := m.DebtFloor.Size()
		i -= size
		if _, err := m.DebtFloor.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintExtendedPairVault(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x52
	{
		size := m.DebtCeiling.Size()
		i -= size
		if _, err := m.DebtCeiling.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintExtendedPairVault(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x4a
	if m.IsVaultActive {
		i--
		if m.IsVaultActive {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x40
	}
	{
		size := m.DrawDownFee.Size()
		i -= size
		if _, err := m.DrawDownFee.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintExtendedPairVault(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x3a
	{
		size := m.LiquidationPenalty.Size()
		i -= size
		if _, err := m.LiquidationPenalty.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintExtendedPairVault(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	{
		size := m.ClosingFee.Size()
		i -= size
		if _, err := m.ClosingFee.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintExtendedPairVault(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	{
		size := m.StabilityFee.Size()
		i -= size
		if _, err := m.StabilityFee.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintExtendedPairVault(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	if m.PairId != 0 {
		i = encodeVarintExtendedPairVault(dAtA, i, uint64(m.PairId))
		i--
		dAtA[i] = 0x18
	}
	if m.AppId != 0 {
		i = encodeVarintExtendedPairVault(dAtA, i, uint64(m.AppId))
		i--
		dAtA[i] = 0x10
	}
	if m.Id != 0 {
		i = encodeVarintExtendedPairVault(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintExtendedPairVault(dAtA []byte, offset int, v uint64) int {
	offset -= sovExtendedPairVault(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ExtendedPairVault) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovExtendedPairVault(uint64(m.Id))
	}
	if m.AppId != 0 {
		n += 1 + sovExtendedPairVault(uint64(m.AppId))
	}
	if m.PairId != 0 {
		n += 1 + sovExtendedPairVault(uint64(m.PairId))
	}
	l = m.StabilityFee.Size()
	n += 1 + l + sovExtendedPairVault(uint64(l))
	l = m.ClosingFee.Size()
	n += 1 + l + sovExtendedPairVault(uint64(l))
	l = m.LiquidationPenalty.Size()
	n += 1 + l + sovExtendedPairVault(uint64(l))
	l = m.DrawDownFee.Size()
	n += 1 + l + sovExtendedPairVault(uint64(l))
	if m.IsVaultActive {
		n += 2
	}
	l = m.DebtCeiling.Size()
	n += 1 + l + sovExtendedPairVault(uint64(l))
	l = m.DebtFloor.Size()
	n += 1 + l + sovExtendedPairVault(uint64(l))
	if m.IsStableMintVault {
		n += 2
	}
	l = m.MinCr.Size()
	n += 1 + l + sovExtendedPairVault(uint64(l))
	l = len(m.PairName)
	if l > 0 {
		n += 1 + l + sovExtendedPairVault(uint64(l))
	}
	if m.AssetOutOraclePrice {
		n += 2
	}
	if m.AssetOutPrice != 0 {
		n += 1 + sovExtendedPairVault(uint64(m.AssetOutPrice))
	}
	if m.MinUsdValueLeft != 0 {
		n += 2 + sovExtendedPairVault(uint64(m.MinUsdValueLeft))
	}
	if m.BlockHeight != 0 {
		n += 2 + sovExtendedPairVault(uint64(m.BlockHeight))
	}
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.BlockTime)
	n += 2 + l + sovExtendedPairVault(uint64(l))
	return n
}

func sovExtendedPairVault(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozExtendedPairVault(x uint64) (n int) {
	return sovExtendedPairVault(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ExtendedPairVault) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowExtendedPairVault
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
			return fmt.Errorf("proto: ExtendedPairVault: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ExtendedPairVault: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExtendedPairVault
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
				return fmt.Errorf("proto: wrong wireType = %d for field AppId", wireType)
			}
			m.AppId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExtendedPairVault
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AppId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PairId", wireType)
			}
			m.PairId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExtendedPairVault
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PairId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StabilityFee", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExtendedPairVault
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
				return ErrInvalidLengthExtendedPairVault
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthExtendedPairVault
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.StabilityFee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClosingFee", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExtendedPairVault
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
				return ErrInvalidLengthExtendedPairVault
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthExtendedPairVault
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ClosingFee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LiquidationPenalty", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExtendedPairVault
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
				return ErrInvalidLengthExtendedPairVault
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthExtendedPairVault
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.LiquidationPenalty.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DrawDownFee", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExtendedPairVault
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
				return ErrInvalidLengthExtendedPairVault
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthExtendedPairVault
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.DrawDownFee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsVaultActive", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExtendedPairVault
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
			m.IsVaultActive = bool(v != 0)
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DebtCeiling", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExtendedPairVault
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
				return ErrInvalidLengthExtendedPairVault
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthExtendedPairVault
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.DebtCeiling.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DebtFloor", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExtendedPairVault
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
				return ErrInvalidLengthExtendedPairVault
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthExtendedPairVault
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.DebtFloor.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 11:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsStableMintVault", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExtendedPairVault
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
			m.IsStableMintVault = bool(v != 0)
		case 12:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinCr", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExtendedPairVault
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
				return ErrInvalidLengthExtendedPairVault
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthExtendedPairVault
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MinCr.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 13:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PairName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExtendedPairVault
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
				return ErrInvalidLengthExtendedPairVault
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthExtendedPairVault
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PairName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 14:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AssetOutOraclePrice", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExtendedPairVault
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
			m.AssetOutOraclePrice = bool(v != 0)
		case 15:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AssetOutPrice", wireType)
			}
			m.AssetOutPrice = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExtendedPairVault
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AssetOutPrice |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 16:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinUsdValueLeft", wireType)
			}
			m.MinUsdValueLeft = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExtendedPairVault
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MinUsdValueLeft |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 17:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockHeight", wireType)
			}
			m.BlockHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExtendedPairVault
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BlockHeight |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 18:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExtendedPairVault
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
				return ErrInvalidLengthExtendedPairVault
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthExtendedPairVault
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(&m.BlockTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipExtendedPairVault(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthExtendedPairVault
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
func skipExtendedPairVault(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowExtendedPairVault
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
					return 0, ErrIntOverflowExtendedPairVault
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
					return 0, ErrIntOverflowExtendedPairVault
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
				return 0, ErrInvalidLengthExtendedPairVault
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupExtendedPairVault
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthExtendedPairVault
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthExtendedPairVault        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowExtendedPairVault          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupExtendedPairVault = fmt.Errorf("proto: unexpected end of group")
)
