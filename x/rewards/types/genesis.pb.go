// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: comdex/rewards/v1beta1/genesis.proto

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

// GenesisState defines the rewards module's genesis state.
type GenesisState struct {
	InternalRewards        []InternalRewards        `protobuf:"bytes,1,rep,name=internal_rewards,json=internalRewards,proto3" json:"internal_rewards" yaml:"internal_rewards"`
	LockerRewardsTracker   []LockerRewardsTracker   `protobuf:"bytes,2,rep,name=locker_rewards_tracker,json=lockerRewardsTracker,proto3" json:"locker_rewards_tracker" yaml:"locker_rewards_tracker"`
	VaultInterestTracker   []VaultInterestTracker   `protobuf:"bytes,3,rep,name=vault_interest_tracker,json=vaultInterestTracker,proto3" json:"vault_interest_tracker" yaml:"vault_interest_tracker"`
	LockerExternalRewards  []LockerExternalRewards  `protobuf:"bytes,4,rep,name=locker_external_rewards,json=lockerExternalRewards,proto3" json:"locker_external_rewards" yaml:"locker_external_rewards"`
	VaultExternalRewards   []VaultExternalRewards   `protobuf:"bytes,5,rep,name=vault_external_rewards,json=vaultExternalRewards,proto3" json:"vault_external_rewards" yaml:"vault_external_rewards"`
	AppIDs                 []uint64                 `protobuf:"varint,6,rep,packed,name=appIDs,proto3" json:"appIDs,omitempty" yaml:"vault_external_rewards"`
	EpochInfo              []EpochInfo              `protobuf:"bytes,7,rep,name=epochInfo,proto3" json:"epochInfo" yaml:"epochInfo"`
	Gauge                  []Gauge                  `protobuf:"bytes,8,rep,name=gauge,proto3" json:"gauge" yaml:"gauge"`
	GaugeByTriggerDuration []GaugeByTriggerDuration `protobuf:"bytes,9,rep,name=gaugeByTriggerDuration,proto3" json:"gaugeByTriggerDuration" yaml:"gaugeByTriggerDuration"`
	Params                 Params                   `protobuf:"bytes,10,opt,name=params,proto3" json:"params" yaml:"params"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_cdfc05d0f3c33bb6, []int{0}
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

func (m *GenesisState) GetInternalRewards() []InternalRewards {
	if m != nil {
		return m.InternalRewards
	}
	return nil
}

func (m *GenesisState) GetLockerRewardsTracker() []LockerRewardsTracker {
	if m != nil {
		return m.LockerRewardsTracker
	}
	return nil
}

func (m *GenesisState) GetVaultInterestTracker() []VaultInterestTracker {
	if m != nil {
		return m.VaultInterestTracker
	}
	return nil
}

func (m *GenesisState) GetLockerExternalRewards() []LockerExternalRewards {
	if m != nil {
		return m.LockerExternalRewards
	}
	return nil
}

func (m *GenesisState) GetVaultExternalRewards() []VaultExternalRewards {
	if m != nil {
		return m.VaultExternalRewards
	}
	return nil
}

func (m *GenesisState) GetAppIDs() []uint64 {
	if m != nil {
		return m.AppIDs
	}
	return nil
}

func (m *GenesisState) GetEpochInfo() []EpochInfo {
	if m != nil {
		return m.EpochInfo
	}
	return nil
}

func (m *GenesisState) GetGauge() []Gauge {
	if m != nil {
		return m.Gauge
	}
	return nil
}

func (m *GenesisState) GetGaugeByTriggerDuration() []GaugeByTriggerDuration {
	if m != nil {
		return m.GaugeByTriggerDuration
	}
	return nil
}

func (m *GenesisState) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "comdex.rewards.v1beta1.GenesisState")
}

func init() {
	proto.RegisterFile("comdex/rewards/v1beta1/genesis.proto", fileDescriptor_cdfc05d0f3c33bb6)
}

var fileDescriptor_cdfc05d0f3c33bb6 = []byte{
	// 542 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x94, 0xcf, 0x6f, 0xd3, 0x30,
	0x1c, 0xc5, 0x1b, 0xba, 0x95, 0xcd, 0x1b, 0x62, 0x8a, 0xba, 0x2e, 0xaa, 0xb4, 0xb4, 0x33, 0xbf,
	0x7a, 0x60, 0x89, 0x3a, 0x4e, 0x70, 0x8c, 0x36, 0x4d, 0x11, 0x20, 0x21, 0x33, 0x21, 0xc1, 0xa5,
	0x72, 0x3b, 0x37, 0x8b, 0x48, 0xe3, 0xc8, 0x71, 0x4b, 0xfb, 0x27, 0x70, 0x40, 0xe2, 0xcf, 0xda,
	0x71, 0x47, 0x4e, 0x13, 0x6a, 0xff, 0x03, 0x2e, 0x5c, 0x51, 0x6c, 0x77, 0xa3, 0x59, 0x1c, 0xb8,
	0xc5, 0xd1, 0x7b, 0x9f, 0xf7, 0xfc, 0xb5, 0x6c, 0xf0, 0x78, 0x40, 0x47, 0xe7, 0x64, 0xea, 0x32,
	0xf2, 0x05, 0xb3, 0xf3, 0xd4, 0x9d, 0x74, 0xfb, 0x84, 0xe3, 0xae, 0x1b, 0x90, 0x98, 0xa4, 0x61,
	0xea, 0x24, 0x8c, 0x72, 0x6a, 0x36, 0xa4, 0xca, 0x51, 0x2a, 0x47, 0xa9, 0x9a, 0xf5, 0x80, 0x06,
	0x54, 0x48, 0xdc, 0xec, 0x4b, 0xaa, 0x9b, 0x8f, 0x34, 0xcc, 0x04, 0x33, 0x3c, 0x52, 0xc8, 0xa6,
	0x2e, 0x78, 0x19, 0x51, 0x8e, 0x22, 0x09, 0x1d, 0x5c, 0x2c, 0x45, 0x50, 0xb7, 0x07, 0x3c, 0x0e,
	0x88, 0xd4, 0xc0, 0xdf, 0x1b, 0x60, 0xfb, 0x54, 0xee, 0xe9, 0x3d, 0xc7, 0x9c, 0x98, 0x29, 0xd8,
	0x09, 0x63, 0x4e, 0x58, 0x8c, 0xa3, 0x9e, 0x32, 0x5a, 0x46, 0xbb, 0xda, 0xd9, 0x3a, 0x7a, 0xe6,
	0x14, 0xef, 0xd6, 0xf1, 0x95, 0x1e, 0xc9, 0xff, 0x5e, 0xeb, 0xf2, 0xba, 0x55, 0xf9, 0x75, 0xdd,
	0xda, 0x9b, 0xe1, 0x51, 0xf4, 0x0a, 0xe6, 0x71, 0x10, 0x3d, 0x0c, 0x57, 0x1d, 0xe6, 0x57, 0x03,
	0x34, 0x22, 0x3a, 0xf8, 0x4c, 0xd8, 0x52, 0xd4, 0xe3, 0x0c, 0x67, 0x6b, 0xeb, 0x9e, 0xc8, 0x7e,
	0xae, 0xcb, 0x7e, 0x23, 0x5c, 0x8a, 0x73, 0x26, 0x3d, 0xde, 0x13, 0x55, 0x60, 0x5f, 0x16, 0x28,
	0x26, 0x43, 0x54, 0x8f, 0x0a, 0xcc, 0xa2, 0xcb, 0x04, 0x8f, 0x23, 0xde, 0x13, 0x2d, 0x49, 0xca,
	0x6f, 0xba, 0x54, 0xcb, 0xbb, 0x7c, 0xc8, 0x5c, 0xbe, 0x32, 0x69, 0xba, 0x14, 0x93, 0x21, 0xaa,
	0x4f, 0x0a, 0xcc, 0xe6, 0x37, 0x03, 0xec, 0xa9, 0xf6, 0x64, 0x9a, 0x3b, 0x94, 0x35, 0x51, 0xe6,
	0xb0, 0x7c, 0x30, 0x27, 0xd3, 0xd5, 0xa3, 0x79, 0xaa, 0xda, 0xd8, 0x2b, 0x93, 0xc9, 0xb3, 0x21,
	0xda, 0x8d, 0x8a, 0xec, 0x7f, 0xcd, 0xe6, 0x4e, 0x9d, 0xf5, 0xff, 0x98, 0x4d, 0xbe, 0x4d, 0xe1,
	0x6c, 0xee, 0x96, 0x91, 0xb3, 0xc9, 0x77, 0x79, 0x09, 0x6a, 0x38, 0x49, 0xfc, 0xe3, 0xd4, 0xaa,
	0xb5, 0xab, 0x9d, 0x35, 0xef, 0xe0, 0xdf, 0x20, 0x65, 0x30, 0x3f, 0x82, 0x4d, 0x71, 0x51, 0xfc,
	0x78, 0x48, 0xad, 0xfb, 0xa2, 0xf8, 0x81, 0xae, 0xf8, 0xc9, 0x52, 0xe8, 0x59, 0xaa, 0xed, 0x8e,
	0x0c, 0xb9, 0x21, 0x40, 0x74, 0x4b, 0x33, 0x7d, 0xb0, 0x2e, 0xae, 0x97, 0xb5, 0x21, 0xb0, 0xfb,
	0x3a, 0xec, 0x69, 0x26, 0xf2, 0xea, 0x0a, 0xb9, 0x2d, 0x91, 0xc2, 0x09, 0x91, 0x24, 0x64, 0x87,
	0xdf, 0x10, 0x5f, 0xde, 0xec, 0x8c, 0x85, 0x41, 0x40, 0xd8, 0xf1, 0x98, 0x61, 0x1e, 0xd2, 0xd8,
	0xda, 0x14, 0x70, 0xa7, 0x1c, 0x9e, 0x77, 0xe5, 0xc7, 0x5d, 0xcc, 0x86, 0x48, 0x13, 0x6a, 0xbe,
	0x05, 0x35, 0xf9, 0x52, 0x59, 0xa0, 0x6d, 0x74, 0xb6, 0x8e, 0x6c, 0x5d, 0xfc, 0x3b, 0xa1, 0xf2,
	0x76, 0x55, 0xdc, 0x03, 0x19, 0x27, 0xbd, 0x10, 0x29, 0x88, 0xf7, 0xfa, 0x72, 0x6e, 0x1b, 0x57,
	0x73, 0xdb, 0xf8, 0x39, 0xb7, 0x8d, 0xef, 0x0b, 0xbb, 0x72, 0xb5, 0xb0, 0x2b, 0x3f, 0x16, 0x76,
	0xe5, 0x53, 0x37, 0x08, 0xf9, 0xc5, 0xb8, 0x9f, 0xe1, 0x5d, 0x19, 0x71, 0x48, 0x87, 0xc3, 0x70,
	0x10, 0xe2, 0x48, 0xad, 0xdd, 0xdb, 0x47, 0x8d, 0xcf, 0x12, 0x92, 0xf6, 0x6b, 0xe2, 0x35, 0x7b,
	0xf1, 0x27, 0x00, 0x00, 0xff, 0xff, 0xda, 0x39, 0x90, 0x60, 0xb7, 0x05, 0x00, 0x00,
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
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x52
	if len(m.GaugeByTriggerDuration) > 0 {
		for iNdEx := len(m.GaugeByTriggerDuration) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.GaugeByTriggerDuration[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x4a
		}
	}
	if len(m.Gauge) > 0 {
		for iNdEx := len(m.Gauge) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Gauge[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x42
		}
	}
	if len(m.EpochInfo) > 0 {
		for iNdEx := len(m.EpochInfo) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.EpochInfo[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x3a
		}
	}
	if len(m.AppIDs) > 0 {
		dAtA3 := make([]byte, len(m.AppIDs)*10)
		var j2 int
		for _, num := range m.AppIDs {
			for num >= 1<<7 {
				dAtA3[j2] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j2++
			}
			dAtA3[j2] = uint8(num)
			j2++
		}
		i -= j2
		copy(dAtA[i:], dAtA3[:j2])
		i = encodeVarintGenesis(dAtA, i, uint64(j2))
		i--
		dAtA[i] = 0x32
	}
	if len(m.VaultExternalRewards) > 0 {
		for iNdEx := len(m.VaultExternalRewards) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.VaultExternalRewards[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	if len(m.LockerExternalRewards) > 0 {
		for iNdEx := len(m.LockerExternalRewards) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.LockerExternalRewards[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	if len(m.VaultInterestTracker) > 0 {
		for iNdEx := len(m.VaultInterestTracker) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.VaultInterestTracker[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	if len(m.LockerRewardsTracker) > 0 {
		for iNdEx := len(m.LockerRewardsTracker) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.LockerRewardsTracker[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	if len(m.InternalRewards) > 0 {
		for iNdEx := len(m.InternalRewards) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.InternalRewards[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	if len(m.InternalRewards) > 0 {
		for _, e := range m.InternalRewards {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.LockerRewardsTracker) > 0 {
		for _, e := range m.LockerRewardsTracker {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.VaultInterestTracker) > 0 {
		for _, e := range m.VaultInterestTracker {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.LockerExternalRewards) > 0 {
		for _, e := range m.LockerExternalRewards {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.VaultExternalRewards) > 0 {
		for _, e := range m.VaultExternalRewards {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.AppIDs) > 0 {
		l = 0
		for _, e := range m.AppIDs {
			l += sovGenesis(uint64(e))
		}
		n += 1 + sovGenesis(uint64(l)) + l
	}
	if len(m.EpochInfo) > 0 {
		for _, e := range m.EpochInfo {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Gauge) > 0 {
		for _, e := range m.Gauge {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.GaugeByTriggerDuration) > 0 {
		for _, e := range m.GaugeByTriggerDuration {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	l = m.Params.Size()
	n += 1 + l + sovGenesis(uint64(l))
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
				return fmt.Errorf("proto: wrong wireType = %d for field InternalRewards", wireType)
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
			m.InternalRewards = append(m.InternalRewards, InternalRewards{})
			if err := m.InternalRewards[len(m.InternalRewards)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LockerRewardsTracker", wireType)
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
			m.LockerRewardsTracker = append(m.LockerRewardsTracker, LockerRewardsTracker{})
			if err := m.LockerRewardsTracker[len(m.LockerRewardsTracker)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field VaultInterestTracker", wireType)
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
			m.VaultInterestTracker = append(m.VaultInterestTracker, VaultInterestTracker{})
			if err := m.VaultInterestTracker[len(m.VaultInterestTracker)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LockerExternalRewards", wireType)
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
			m.LockerExternalRewards = append(m.LockerExternalRewards, LockerExternalRewards{})
			if err := m.LockerExternalRewards[len(m.LockerExternalRewards)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field VaultExternalRewards", wireType)
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
			m.VaultExternalRewards = append(m.VaultExternalRewards, VaultExternalRewards{})
			if err := m.VaultExternalRewards[len(m.VaultExternalRewards)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType == 0 {
				var v uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowGenesis
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.AppIDs = append(m.AppIDs, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowGenesis
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthGenesis
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthGenesis
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				var count int
				for _, integer := range dAtA[iNdEx:postIndex] {
					if integer < 128 {
						count++
					}
				}
				elementCount = count
				if elementCount != 0 && len(m.AppIDs) == 0 {
					m.AppIDs = make([]uint64, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowGenesis
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.AppIDs = append(m.AppIDs, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field AppIDs", wireType)
			}
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EpochInfo", wireType)
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
			m.EpochInfo = append(m.EpochInfo, EpochInfo{})
			if err := m.EpochInfo[len(m.EpochInfo)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Gauge", wireType)
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
			m.Gauge = append(m.Gauge, Gauge{})
			if err := m.Gauge[len(m.Gauge)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GaugeByTriggerDuration", wireType)
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
			m.GaugeByTriggerDuration = append(m.GaugeByTriggerDuration, GaugeByTriggerDuration{})
			if err := m.GaugeByTriggerDuration[len(m.GaugeByTriggerDuration)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 10:
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
