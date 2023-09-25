package types

import (
	fmt "fmt"
	io "io"
	math "math"
	math_bits "math/bits"

	_ "github.com/cosmos/cosmos-sdk/types"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	proto "github.com/cosmos/gogoproto/proto"
	_ "github.com/gogo/protobuf/gogoproto"
	_ "google.golang.org/protobuf/types/known/timestamppb"
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

type Extended_Pair_Old struct {
	Id              uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	AssetIn         uint64 `protobuf:"varint,2,opt,name=asset_in,json=assetIn,proto3" json:"asset_in,omitempty" yaml:"asset_in"`
	AssetOut        uint64 `protobuf:"varint,3,opt,name=asset_out,json=assetOut,proto3" json:"asset_out,omitempty" yaml:"asset_out"`
	IsInterPool     bool   `protobuf:"varint,4,opt,name=is_inter_pool,json=isInterPool,proto3" json:"is_inter_pool,omitempty" yaml:"is_inter_pool"`
	AssetOutPoolID  uint64 `protobuf:"varint,5,opt,name=asset_out_pool_id,json=assetOutPoolId,proto3" json:"asset_out_pool_id,omitempty" yaml:"asset_out_pool_id"`
	MinUsdValueLeft uint64 `protobuf:"varint,6,opt,name=min_usd_value_left,json=minUsdValueLeft,proto3" json:"min_usd_value_left,omitempty" yaml:"min_usd_value_left"`
}

func (m *Extended_Pair_Old) Reset()         { *m = Extended_Pair_Old{} }
func (m *Extended_Pair_Old) String() string { return proto.CompactTextString(m) }
func (*Extended_Pair_Old) ProtoMessage()    {}
func (*Extended_Pair_Old) Descriptor() ([]byte, []int) {
	return fileDescriptor_b9f686c34547a278, []int{0}
}
func (m *Extended_Pair_Old) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Extended_Pair_Old) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Extended_Pair_Old.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Extended_Pair_Old) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Extended_Pair_Old.Merge(m, src)
}
func (m *Extended_Pair_Old) XXX_Size() int {
	return m.Size()
}
func (m *Extended_Pair_Old) XXX_DiscardUnknown() {
	xxx_messageInfo_Extended_Pair_Old.DiscardUnknown(m)
}

var xxx_messageInfo_Extended_Pair_Old proto.InternalMessageInfo

func (m *Extended_Pair_Old) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Extended_Pair_Old) GetAssetIn() uint64 {
	if m != nil {
		return m.AssetIn
	}
	return 0
}

func (m *Extended_Pair_Old) GetAssetOut() uint64 {
	if m != nil {
		return m.AssetOut
	}
	return 0
}

func (m *Extended_Pair_Old) GetIsInterPool() bool {
	if m != nil {
		return m.IsInterPool
	}
	return false
}

func (m *Extended_Pair_Old) GetAssetOutPoolID() uint64 {
	if m != nil {
		return m.AssetOutPoolID
	}
	return 0
}

func (m *Extended_Pair_Old) GetMinUsdValueLeft() uint64 {
	if m != nil {
		return m.MinUsdValueLeft
	}
	return 0
}

type AssetRatesParams_Old struct {
	AssetID              uint64                                 `protobuf:"varint,1,opt,name=asset_id,json=assetId,proto3" json:"asset_id,omitempty" yaml:"asset_id"`
	UOptimal             github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=u_optimal,json=uOptimal,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"u_optimal" yaml:"u_optimal"`
	Base                 github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,3,opt,name=base,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"base" yaml:"base"`
	Slope1               github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,4,opt,name=slope1,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"slope1" yaml:"slope1"`
	Slope2               github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,5,opt,name=slope2,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"slope2" yaml:"slope2"`
	EnableStableBorrow   bool                                   `protobuf:"varint,6,opt,name=enable_stable_borrow,json=enableStableBorrow,proto3" json:"enable_stable_borrow,omitempty" yaml:"enable_stable_borrow"`
	StableBase           github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,7,opt,name=stable_base,json=stableBase,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"stable_base" yaml:"stable_base"`
	StableSlope1         github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,8,opt,name=stable_slope1,json=stableSlope1,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"stable_slope1" yaml:"stable_slope1"`
	StableSlope2         github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,9,opt,name=stable_slope2,json=stableSlope2,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"stable_slope2" yaml:"stable_slope2"`
	Ltv                  github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,10,opt,name=ltv,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"ltv" yaml:"ltv"`
	LiquidationThreshold github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,11,opt,name=liquidation_threshold,json=liquidationThreshold,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"liquidation_threshold" yaml:"liquidation_threshold"`
	LiquidationPenalty   github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,12,opt,name=liquidation_penalty,json=liquidationPenalty,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"liquidation_penalty" yaml:"liquidation_penalty"`
	LiquidationBonus     github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,13,opt,name=liquidation_bonus,json=liquidationBonus,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"liquidation_bonus" yaml:"liquidation_bonus"`
	ReserveFactor        github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,14,opt,name=reserve_factor,json=reserveFactor,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"reserve_factor" yaml:"reserve_factor"`
	CAssetID             uint64                                 `protobuf:"varint,15,opt,name=c_asset_id,json=cAssetId,proto3" json:"c_asset_id,omitempty" yaml:"c_asset_id"`
}

func (m *AssetRatesParams_Old) Reset()         { *m = AssetRatesParams_Old{} }
func (m *AssetRatesParams_Old) String() string { return proto.CompactTextString(m) }
func (*AssetRatesParams_Old) ProtoMessage()    {}
func (*AssetRatesParams_Old) Descriptor() ([]byte, []int) {
	return fileDescriptor_b9f686c34547a278, []int{1}
}
func (m *AssetRatesParams_Old) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AssetRatesParams_Old) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AssetRatesParams_Old.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AssetRatesParams_Old) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AssetRatesParams_Old.Merge(m, src)
}
func (m *AssetRatesParams_Old) XXX_Size() int {
	return m.Size()
}
func (m *AssetRatesParams_Old) XXX_DiscardUnknown() {
	xxx_messageInfo_AssetRatesParams_Old.DiscardUnknown(m)
}

var xxx_messageInfo_AssetRatesParams_Old proto.InternalMessageInfo

func (m *AssetRatesParams_Old) GetAssetID() uint64 {
	if m != nil {
		return m.AssetID
	}
	return 0
}

func (m *AssetRatesParams_Old) GetEnableStableBorrow() bool {
	if m != nil {
		return m.EnableStableBorrow
	}
	return false
}

func (m *AssetRatesParams_Old) GetCAssetID() uint64 {
	if m != nil {
		return m.CAssetID
	}
	return 0
}

func init() {
	proto.RegisterType((*Extended_Pair_Old)(nil), "comdex.lend.v1beta1.Extended_Pair_Old")
	proto.RegisterType((*AssetRatesParams_Old)(nil), "comdex.lend.v1beta1.AssetRatesParams_Old")
}

func init() { proto.RegisterFile("comdex/lend/v1beta1/temp.proto", fileDescriptor_b9f686c34547a278) }

var fileDescriptor_b9f686c34547a278 = []byte{
	// 834 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x95, 0xcf, 0x6f, 0xdb, 0x36,
	0x14, 0xc7, 0xe3, 0x34, 0x4b, 0x6c, 0xba, 0x76, 0x1a, 0xc6, 0x05, 0xb4, 0x6c, 0x93, 0x3a, 0x1e,
	0x86, 0x5e, 0x2a, 0x21, 0xde, 0x69, 0x43, 0x87, 0x62, 0x5a, 0xd6, 0x2d, 0x45, 0xd1, 0xa4, 0xec,
	0x7e, 0x60, 0xbb, 0x08, 0x94, 0x48, 0x3b, 0x44, 0x29, 0xd1, 0x13, 0x29, 0xb7, 0x39, 0xec, 0xb4,
	0x7f, 0x60, 0x7f, 0x56, 0xb1, 0x53, 0x8f, 0xc3, 0x0e, 0xc2, 0xe0, 0x1c, 0x77, 0xf3, 0x5f, 0x30,
	0x88, 0x94, 0x7f, 0x35, 0xb9, 0x08, 0xdb, 0x49, 0x7a, 0xdf, 0xf7, 0xf8, 0xf9, 0x3e, 0x3c, 0x12,
	0x24, 0x70, 0x13, 0x99, 0x52, 0xf6, 0x3a, 0x10, 0x2c, 0xa3, 0xc1, 0xf4, 0x38, 0x66, 0x9a, 0x1c,
	0x07, 0x9a, 0xa5, 0x13, 0x7f, 0x92, 0x4b, 0x2d, 0xe1, 0xa1, 0xcd, 0xfb, 0x55, 0xde, 0xaf, 0xf3,
	0x47, 0x83, 0xb1, 0x1c, 0x4b, 0x93, 0x0f, 0xaa, 0x3f, 0x5b, 0x7a, 0xe4, 0x8d, 0xa5, 0x1c, 0x0b,
	0x16, 0x98, 0x28, 0x2e, 0x46, 0x81, 0xe6, 0x29, 0x53, 0x9a, 0x2c, 0x58, 0x47, 0x6e, 0x22, 0x55,
	0x2a, 0x55, 0x10, 0x13, 0xc5, 0x96, 0x5e, 0x89, 0xe4, 0x99, 0xcd, 0xa3, 0x7f, 0xb6, 0xc1, 0xc1,
	0xd7, 0xaf, 0x35, 0xcb, 0x28, 0xa3, 0xd1, 0x39, 0xe1, 0x79, 0x74, 0x26, 0x28, 0xec, 0x83, 0x6d,
	0x4e, 0x9d, 0xd6, 0xbd, 0xd6, 0xfd, 0x1d, 0xbc, 0xcd, 0x29, 0xf4, 0x41, 0x9b, 0x28, 0xc5, 0x74,
	0xc4, 0x33, 0x67, 0xbb, 0x52, 0xc3, 0xc3, 0x79, 0xe9, 0xed, 0x5f, 0x92, 0x54, 0x7c, 0x8e, 0x16,
	0x19, 0x84, 0xf7, 0xcc, 0xef, 0x69, 0x06, 0x8f, 0x41, 0xc7, 0xaa, 0xb2, 0xd0, 0xce, 0x2d, 0xb3,
	0x60, 0x30, 0x2f, 0xbd, 0x3b, 0xeb, 0x0b, 0x64, 0xa1, 0x11, 0xb6, 0xd8, 0xb3, 0x42, 0xc3, 0x87,
	0xa0, 0xc7, 0x55, 0xc4, 0x33, 0xcd, 0xf2, 0x68, 0x22, 0xa5, 0x70, 0x76, 0xee, 0xb5, 0xee, 0xb7,
	0x43, 0x67, 0x5e, 0x7a, 0x03, 0xbb, 0x6c, 0x23, 0x8d, 0x70, 0x97, 0xab, 0xd3, 0x2a, 0x3c, 0x97,
	0x52, 0xc0, 0x9f, 0xc0, 0xc1, 0x92, 0x6a, 0xf2, 0x11, 0xa7, 0xce, 0x7b, 0xc6, 0xd8, 0x9f, 0x95,
	0x5e, 0xff, 0xcb, 0xda, 0xa6, 0x2a, 0x3e, 0x3d, 0x99, 0x97, 0x9e, 0xf3, 0x4e, 0x2b, 0x8b, 0x45,
	0x08, 0xf7, 0xc9, 0x7a, 0x2d, 0x85, 0x4f, 0x00, 0x4c, 0x79, 0x16, 0x15, 0x8a, 0x46, 0x53, 0x22,
	0x0a, 0x16, 0x09, 0x36, 0xd2, 0xce, 0xae, 0x61, 0x7f, 0x34, 0x2f, 0xbd, 0xf7, 0x2d, 0xe9, 0x7a,
	0x0d, 0xc2, 0xfb, 0x29, 0xcf, 0xbe, 0x57, 0xf4, 0x87, 0x4a, 0x7a, 0x5a, 0x29, 0x7f, 0x74, 0xc1,
	0xc0, 0xb4, 0x82, 0x89, 0x66, 0xea, 0x9c, 0xe4, 0x24, 0x55, 0x66, 0xe0, 0x9f, 0x2d, 0x07, 0x5c,
	0x8f, 0x3d, 0x74, 0x67, 0xa5, 0xb7, 0x67, 0x6a, 0x4d, 0xbf, 0x9b, 0xb3, 0xa6, 0xcb, 0x59, 0x53,
	0x18, 0x81, 0x4e, 0x11, 0xc9, 0x89, 0xe6, 0x29, 0x11, 0x66, 0x73, 0x3a, 0x61, 0xf8, 0xa6, 0xf4,
	0xb6, 0xfe, 0x2a, 0xbd, 0x4f, 0xc6, 0x5c, 0x5f, 0x14, 0xb1, 0x9f, 0xc8, 0x34, 0xa8, 0xcf, 0x81,
	0xfd, 0x3c, 0x50, 0xf4, 0x65, 0xa0, 0x2f, 0x27, 0x4c, 0xf9, 0x27, 0x2c, 0x59, 0xed, 0xcc, 0x12,
	0x84, 0x70, 0xbb, 0x38, 0xb3, 0xbf, 0xf0, 0x39, 0xd8, 0xa9, 0x4e, 0x8f, 0xd9, 0xc7, 0x4e, 0xf8,
	0x45, 0x63, 0x76, 0xd7, 0xb2, 0x2b, 0x06, 0xc2, 0x06, 0x05, 0x7f, 0x04, 0xbb, 0x4a, 0xc8, 0x09,
	0x3b, 0x36, 0xbb, 0xdc, 0x09, 0x1f, 0x35, 0x86, 0xf6, 0x2c, 0xd4, 0x52, 0x10, 0xae, 0x71, 0x4b,
	0xf0, 0xd0, 0x6c, 0xfe, 0x7f, 0x05, 0x0f, 0x17, 0xe0, 0x21, 0x7c, 0x0e, 0x06, 0x2c, 0x23, 0xb1,
	0x60, 0x91, 0xd2, 0xe6, 0x13, 0xcb, 0x3c, 0x97, 0xaf, 0xcc, 0x39, 0x68, 0x87, 0xde, 0xbc, 0xf4,
	0x3e, 0xb0, 0x0b, 0x6f, 0xaa, 0x42, 0x18, 0x5a, 0xf9, 0x85, 0x51, 0x43, 0x23, 0x42, 0x06, 0xba,
	0x8b, 0xaa, 0x6a, 0xbc, 0x7b, 0xa6, 0xe1, 0x93, 0xc6, 0x0d, 0xc3, 0xba, 0xe1, 0x15, 0x0a, 0x61,
	0x60, 0xa3, 0xb0, 0x9a, 0xf5, 0x4b, 0xd0, 0xab, 0x73, 0xf5, 0xc8, 0xdb, 0xc6, 0xe8, 0x71, 0x63,
	0xa3, 0xc1, 0x86, 0xd1, 0x62, 0xf2, 0xb7, 0x6d, 0xfc, 0xc2, 0xce, 0xff, 0x1d, 0xb3, 0xa1, 0xd3,
	0xf9, 0xff, 0xcc, 0x86, 0x9b, 0x66, 0x43, 0xf8, 0x0c, 0xdc, 0x12, 0x7a, 0xea, 0x00, 0x63, 0xf1,
	0xb0, 0xb1, 0x05, 0xb0, 0x16, 0x42, 0x4f, 0x11, 0xae, 0x40, 0xf0, 0xb7, 0x16, 0xb8, 0x2b, 0xf8,
	0x2f, 0x05, 0xa7, 0x44, 0x73, 0x99, 0x45, 0xfa, 0x22, 0x67, 0xea, 0x42, 0x0a, 0xea, 0x74, 0x8d,
	0xc5, 0xb3, 0xc6, 0x16, 0x1f, 0xd6, 0x16, 0x37, 0x41, 0x11, 0x1e, 0xac, 0xe9, 0xdf, 0x2d, 0x64,
	0xf8, 0x2b, 0x38, 0x5c, 0xaf, 0x9f, 0xb0, 0x8c, 0x08, 0x7d, 0xe9, 0xdc, 0x36, 0x2d, 0x3c, 0x6d,
	0xdc, 0xc2, 0xd1, 0xf5, 0x16, 0x6a, 0x24, 0xc2, 0x70, 0x4d, 0x3d, 0xb7, 0x22, 0x7c, 0x05, 0x0e,
	0xd6, 0x6b, 0x63, 0x99, 0x15, 0xca, 0xe9, 0x19, 0xf3, 0x27, 0x8d, 0xcd, 0x9d, 0xeb, 0xe6, 0x06,
	0x88, 0xf0, 0x9d, 0x35, 0x2d, 0xac, 0x24, 0x98, 0x81, 0x7e, 0xce, 0x14, 0xcb, 0xa7, 0x2c, 0x1a,
	0x91, 0x44, 0xcb, 0xdc, 0xe9, 0x1b, 0xd7, 0x6f, 0x1a, 0xbb, 0xde, 0xb5, 0xae, 0x9b, 0x34, 0x84,
	0x7b, 0xb5, 0xf0, 0xd8, 0xc4, 0xf0, 0x11, 0x00, 0x49, 0xb4, 0xbc, 0x74, 0xf7, 0xcd, 0xa5, 0xfb,
	0xf1, 0xac, 0xf4, 0xda, 0x5f, 0xad, 0x6e, 0xdd, 0x03, 0x4b, 0x5a, 0xd5, 0x21, 0xdc, 0x4e, 0x6c,
	0x9a, 0x86, 0xdf, 0xbe, 0x99, 0xb9, 0xad, 0xb7, 0x33, 0xb7, 0xf5, 0xf7, 0xcc, 0x6d, 0xfd, 0x7e,
	0xe5, 0x6e, 0xbd, 0xbd, 0x72, 0xb7, 0xfe, 0xbc, 0x72, 0xb7, 0x7e, 0xf6, 0x37, 0x5a, 0xad, 0xde,
	0xf2, 0x07, 0x72, 0x34, 0xe2, 0x09, 0x27, 0xa2, 0x8e, 0x83, 0xfa, 0xf5, 0x37, 0x6d, 0xc7, 0xbb,
	0xe6, 0x2d, 0xfe, 0xf4, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x9e, 0xe2, 0xfe, 0xc7, 0x19, 0x08,
	0x00, 0x00,
}

func (m *Extended_Pair_Old) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Extended_Pair_Old) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Extended_Pair_Old) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.MinUsdValueLeft != 0 {
		i = encodeVarintTemp(dAtA, i, uint64(m.MinUsdValueLeft))
		i--
		dAtA[i] = 0x30
	}
	if m.AssetOutPoolID != 0 {
		i = encodeVarintTemp(dAtA, i, uint64(m.AssetOutPoolID))
		i--
		dAtA[i] = 0x28
	}
	if m.IsInterPool {
		i--
		if m.IsInterPool {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x20
	}
	if m.AssetOut != 0 {
		i = encodeVarintTemp(dAtA, i, uint64(m.AssetOut))
		i--
		dAtA[i] = 0x18
	}
	if m.AssetIn != 0 {
		i = encodeVarintTemp(dAtA, i, uint64(m.AssetIn))
		i--
		dAtA[i] = 0x10
	}
	if m.Id != 0 {
		i = encodeVarintTemp(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *AssetRatesParams_Old) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AssetRatesParams_Old) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AssetRatesParams_Old) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.CAssetID != 0 {
		i = encodeVarintTemp(dAtA, i, uint64(m.CAssetID))
		i--
		dAtA[i] = 0x78
	}
	{
		size := m.ReserveFactor.Size()
		i -= size
		if _, err := m.ReserveFactor.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTemp(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x72
	{
		size := m.LiquidationBonus.Size()
		i -= size
		if _, err := m.LiquidationBonus.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTemp(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x6a
	{
		size := m.LiquidationPenalty.Size()
		i -= size
		if _, err := m.LiquidationPenalty.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTemp(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x62
	{
		size := m.LiquidationThreshold.Size()
		i -= size
		if _, err := m.LiquidationThreshold.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTemp(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x5a
	{
		size := m.Ltv.Size()
		i -= size
		if _, err := m.Ltv.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTemp(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x52
	{
		size := m.StableSlope2.Size()
		i -= size
		if _, err := m.StableSlope2.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTemp(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x4a
	{
		size := m.StableSlope1.Size()
		i -= size
		if _, err := m.StableSlope1.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTemp(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x42
	{
		size := m.StableBase.Size()
		i -= size
		if _, err := m.StableBase.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTemp(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x3a
	if m.EnableStableBorrow {
		i--
		if m.EnableStableBorrow {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x30
	}
	{
		size := m.Slope2.Size()
		i -= size
		if _, err := m.Slope2.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTemp(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	{
		size := m.Slope1.Size()
		i -= size
		if _, err := m.Slope1.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTemp(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size := m.Base.Size()
		i -= size
		if _, err := m.Base.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTemp(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.UOptimal.Size()
		i -= size
		if _, err := m.UOptimal.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTemp(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if m.AssetID != 0 {
		i = encodeVarintTemp(dAtA, i, uint64(m.AssetID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintTemp(dAtA []byte, offset int, v uint64) int {
	offset -= sovTemp(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Extended_Pair_Old) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovTemp(uint64(m.Id))
	}
	if m.AssetIn != 0 {
		n += 1 + sovTemp(uint64(m.AssetIn))
	}
	if m.AssetOut != 0 {
		n += 1 + sovTemp(uint64(m.AssetOut))
	}
	if m.IsInterPool {
		n += 2
	}
	if m.AssetOutPoolID != 0 {
		n += 1 + sovTemp(uint64(m.AssetOutPoolID))
	}
	if m.MinUsdValueLeft != 0 {
		n += 1 + sovTemp(uint64(m.MinUsdValueLeft))
	}
	return n
}

func (m *AssetRatesParams_Old) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.AssetID != 0 {
		n += 1 + sovTemp(uint64(m.AssetID))
	}
	l = m.UOptimal.Size()
	n += 1 + l + sovTemp(uint64(l))
	l = m.Base.Size()
	n += 1 + l + sovTemp(uint64(l))
	l = m.Slope1.Size()
	n += 1 + l + sovTemp(uint64(l))
	l = m.Slope2.Size()
	n += 1 + l + sovTemp(uint64(l))
	if m.EnableStableBorrow {
		n += 2
	}
	l = m.StableBase.Size()
	n += 1 + l + sovTemp(uint64(l))
	l = m.StableSlope1.Size()
	n += 1 + l + sovTemp(uint64(l))
	l = m.StableSlope2.Size()
	n += 1 + l + sovTemp(uint64(l))
	l = m.Ltv.Size()
	n += 1 + l + sovTemp(uint64(l))
	l = m.LiquidationThreshold.Size()
	n += 1 + l + sovTemp(uint64(l))
	l = m.LiquidationPenalty.Size()
	n += 1 + l + sovTemp(uint64(l))
	l = m.LiquidationBonus.Size()
	n += 1 + l + sovTemp(uint64(l))
	l = m.ReserveFactor.Size()
	n += 1 + l + sovTemp(uint64(l))
	if m.CAssetID != 0 {
		n += 1 + sovTemp(uint64(m.CAssetID))
	}
	return n
}

func sovTemp(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTemp(x uint64) (n int) {
	return sovTemp(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Extended_Pair_Old) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTemp
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
			return fmt.Errorf("proto: Extended_Pair_Old: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Extended_Pair_Old: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTemp
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
					return ErrIntOverflowTemp
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
					return ErrIntOverflowTemp
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
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsInterPool", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTemp
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
			m.IsInterPool = bool(v != 0)
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AssetOutPoolID", wireType)
			}
			m.AssetOutPoolID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTemp
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AssetOutPoolID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinUsdValueLeft", wireType)
			}
			m.MinUsdValueLeft = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTemp
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
		default:
			iNdEx = preIndex
			skippy, err := skipTemp(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTemp
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
func (m *AssetRatesParams_Old) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTemp
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
			return fmt.Errorf("proto: AssetRatesParams_Old: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AssetRatesParams_Old: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AssetID", wireType)
			}
			m.AssetID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTemp
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AssetID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UOptimal", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTemp
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
				return ErrInvalidLengthTemp
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTemp
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.UOptimal.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Base", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTemp
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
				return ErrInvalidLengthTemp
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTemp
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Base.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Slope1", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTemp
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
				return ErrInvalidLengthTemp
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTemp
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Slope1.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Slope2", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTemp
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
				return ErrInvalidLengthTemp
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTemp
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Slope2.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EnableStableBorrow", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTemp
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
			m.EnableStableBorrow = bool(v != 0)
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StableBase", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTemp
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
				return ErrInvalidLengthTemp
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTemp
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.StableBase.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StableSlope1", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTemp
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
				return ErrInvalidLengthTemp
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTemp
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.StableSlope1.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StableSlope2", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTemp
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
				return ErrInvalidLengthTemp
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTemp
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.StableSlope2.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ltv", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTemp
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
				return ErrInvalidLengthTemp
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTemp
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Ltv.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LiquidationThreshold", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTemp
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
				return ErrInvalidLengthTemp
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTemp
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.LiquidationThreshold.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 12:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LiquidationPenalty", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTemp
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
				return ErrInvalidLengthTemp
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTemp
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.LiquidationPenalty.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 13:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LiquidationBonus", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTemp
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
				return ErrInvalidLengthTemp
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTemp
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.LiquidationBonus.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 14:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReserveFactor", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTemp
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
				return ErrInvalidLengthTemp
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTemp
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ReserveFactor.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 15:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CAssetID", wireType)
			}
			m.CAssetID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTemp
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CAssetID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTemp(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTemp
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
func skipTemp(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTemp
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
					return 0, ErrIntOverflowTemp
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
					return 0, ErrIntOverflowTemp
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
				return 0, ErrInvalidLengthTemp
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTemp
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTemp
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTemp        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTemp          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTemp = fmt.Errorf("proto: unexpected end of group")
)
