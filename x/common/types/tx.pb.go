// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: comdex/common/v1beta1/tx.proto

package types

import (
	context "context"
	fmt "fmt"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type MsgRegisterContract struct {
	Creator      string `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	GameName     string `protobuf:"bytes,2,opt,name=game_name,json=gameName,proto3" json:"game_name,omitempty"`
	GameId       uint64 `protobuf:"varint,3,opt,name=game_id,json=gameId,proto3" json:"game_id,omitempty"`
	ContractAddr string `protobuf:"bytes,4,opt,name=contractAddr,proto3" json:"contractAddr,omitempty"`
}

func (m *MsgRegisterContract) Reset()         { *m = MsgRegisterContract{} }
func (m *MsgRegisterContract) String() string { return proto.CompactTextString(m) }
func (*MsgRegisterContract) ProtoMessage()    {}
func (*MsgRegisterContract) Descriptor() ([]byte, []int) {
	return fileDescriptor_63826287044af113, []int{0}
}
func (m *MsgRegisterContract) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgRegisterContract) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgRegisterContract.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgRegisterContract) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgRegisterContract.Merge(m, src)
}
func (m *MsgRegisterContract) XXX_Size() int {
	return m.Size()
}
func (m *MsgRegisterContract) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgRegisterContract.DiscardUnknown(m)
}

var xxx_messageInfo_MsgRegisterContract proto.InternalMessageInfo

func (m *MsgRegisterContract) GetCreator() string {
	if m != nil {
		return m.Creator
	}
	return ""
}

func (m *MsgRegisterContract) GetGameName() string {
	if m != nil {
		return m.GameName
	}
	return ""
}

func (m *MsgRegisterContract) GetGameId() uint64 {
	if m != nil {
		return m.GameId
	}
	return 0
}

func (m *MsgRegisterContract) GetContractAddr() string {
	if m != nil {
		return m.ContractAddr
	}
	return ""
}

type MsgRegisterContractResponse struct {
}

func (m *MsgRegisterContractResponse) Reset()         { *m = MsgRegisterContractResponse{} }
func (m *MsgRegisterContractResponse) String() string { return proto.CompactTextString(m) }
func (*MsgRegisterContractResponse) ProtoMessage()    {}
func (*MsgRegisterContractResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_63826287044af113, []int{1}
}
func (m *MsgRegisterContractResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgRegisterContractResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgRegisterContractResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgRegisterContractResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgRegisterContractResponse.Merge(m, src)
}
func (m *MsgRegisterContractResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgRegisterContractResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgRegisterContractResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgRegisterContractResponse proto.InternalMessageInfo

type MsgDeRegisterContract struct {
	Creator string `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	GameId  uint64 `protobuf:"varint,2,opt,name=game_id,json=gameId,proto3" json:"game_id,omitempty"`
}

func (m *MsgDeRegisterContract) Reset()         { *m = MsgDeRegisterContract{} }
func (m *MsgDeRegisterContract) String() string { return proto.CompactTextString(m) }
func (*MsgDeRegisterContract) ProtoMessage()    {}
func (*MsgDeRegisterContract) Descriptor() ([]byte, []int) {
	return fileDescriptor_63826287044af113, []int{2}
}
func (m *MsgDeRegisterContract) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgDeRegisterContract) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgDeRegisterContract.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgDeRegisterContract) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgDeRegisterContract.Merge(m, src)
}
func (m *MsgDeRegisterContract) XXX_Size() int {
	return m.Size()
}
func (m *MsgDeRegisterContract) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgDeRegisterContract.DiscardUnknown(m)
}

var xxx_messageInfo_MsgDeRegisterContract proto.InternalMessageInfo

func (m *MsgDeRegisterContract) GetCreator() string {
	if m != nil {
		return m.Creator
	}
	return ""
}

func (m *MsgDeRegisterContract) GetGameId() uint64 {
	if m != nil {
		return m.GameId
	}
	return 0
}

type MsgDeRegisterContractResponse struct {
}

func (m *MsgDeRegisterContractResponse) Reset()         { *m = MsgDeRegisterContractResponse{} }
func (m *MsgDeRegisterContractResponse) String() string { return proto.CompactTextString(m) }
func (*MsgDeRegisterContractResponse) ProtoMessage()    {}
func (*MsgDeRegisterContractResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_63826287044af113, []int{3}
}
func (m *MsgDeRegisterContractResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgDeRegisterContractResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgDeRegisterContractResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgDeRegisterContractResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgDeRegisterContractResponse.Merge(m, src)
}
func (m *MsgDeRegisterContractResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgDeRegisterContractResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgDeRegisterContractResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgDeRegisterContractResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgRegisterContract)(nil), "comdex.common.v1beta1.MsgRegisterContract")
	proto.RegisterType((*MsgRegisterContractResponse)(nil), "comdex.common.v1beta1.MsgRegisterContractResponse")
	proto.RegisterType((*MsgDeRegisterContract)(nil), "comdex.common.v1beta1.MsgDeRegisterContract")
	proto.RegisterType((*MsgDeRegisterContractResponse)(nil), "comdex.common.v1beta1.MsgDeRegisterContractResponse")
}

func init() { proto.RegisterFile("comdex/common/v1beta1/tx.proto", fileDescriptor_63826287044af113) }

var fileDescriptor_63826287044af113 = []byte{
	// 323 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4b, 0xce, 0xcf, 0x4d,
	0x49, 0xad, 0xd0, 0x4f, 0xce, 0xcf, 0xcd, 0xcd, 0xcf, 0xd3, 0x2f, 0x33, 0x4c, 0x4a, 0x2d, 0x49,
	0x34, 0xd4, 0x2f, 0xa9, 0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0x85, 0xc8, 0xeb, 0x41,
	0xe4, 0xf5, 0xa0, 0xf2, 0x4a, 0x9d, 0x8c, 0x5c, 0xc2, 0xbe, 0xc5, 0xe9, 0x41, 0xa9, 0xe9, 0x99,
	0xc5, 0x25, 0xa9, 0x45, 0xce, 0xf9, 0x79, 0x25, 0x45, 0x89, 0xc9, 0x25, 0x42, 0x12, 0x5c, 0xec,
	0xc9, 0x45, 0xa9, 0x89, 0x25, 0xf9, 0x45, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x30, 0xae,
	0x90, 0x34, 0x17, 0x67, 0x7a, 0x62, 0x6e, 0x6a, 0x7c, 0x5e, 0x62, 0x6e, 0xaa, 0x04, 0x13, 0x58,
	0x8e, 0x03, 0x24, 0xe0, 0x97, 0x98, 0x9b, 0x2a, 0x24, 0xce, 0xc5, 0x0e, 0x96, 0xcc, 0x4c, 0x91,
	0x60, 0x56, 0x60, 0xd4, 0x60, 0x09, 0x62, 0x03, 0x71, 0x3d, 0x53, 0x84, 0x94, 0xb8, 0x78, 0x92,
	0xa1, 0x66, 0x3b, 0xa6, 0xa4, 0x14, 0x49, 0xb0, 0x80, 0x35, 0xa2, 0x88, 0x29, 0xc9, 0x72, 0x49,
	0x63, 0x71, 0x4a, 0x50, 0x6a, 0x71, 0x41, 0x7e, 0x5e, 0x71, 0xaa, 0x92, 0x17, 0x97, 0xa8, 0x6f,
	0x71, 0xba, 0x4b, 0x2a, 0x09, 0x6e, 0x45, 0x72, 0x0e, 0x13, 0xb2, 0x73, 0x94, 0xe4, 0xb9, 0x64,
	0xb1, 0x9a, 0x05, 0xb3, 0xcc, 0xe8, 0x33, 0x23, 0x17, 0xb3, 0x6f, 0x71, 0xba, 0x50, 0x11, 0x97,
	0x00, 0x86, 0x7d, 0x5a, 0x7a, 0x58, 0xc3, 0x52, 0x0f, 0x8b, 0xe3, 0xa5, 0x8c, 0x88, 0x57, 0x0b,
	0xb3, 0x5b, 0xa8, 0x82, 0x4b, 0x08, 0x8b, 0x2f, 0x75, 0x70, 0x9b, 0x84, 0xa9, 0x5a, 0xca, 0x84,
	0x14, 0xd5, 0x30, 0x9b, 0x9d, 0xbc, 0x4e, 0x3c, 0x92, 0x63, 0xbc, 0xf0, 0x48, 0x8e, 0xf1, 0xc1,
	0x23, 0x39, 0xc6, 0x09, 0x8f, 0xe5, 0x18, 0x2e, 0x3c, 0x96, 0x63, 0xb8, 0xf1, 0x58, 0x8e, 0x21,
	0xca, 0x20, 0x3d, 0xb3, 0x24, 0xa3, 0x34, 0x09, 0x64, 0x9c, 0x3e, 0xc4, 0x64, 0xdd, 0xfc, 0xb4,
	0xb4, 0xcc, 0xe4, 0xcc, 0xc4, 0x1c, 0x28, 0x5f, 0x1f, 0x9e, 0xf6, 0x4a, 0x2a, 0x0b, 0x52, 0x8b,
	0x93, 0xd8, 0xc0, 0xe9, 0xce, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x15, 0xc0, 0x4f, 0x20, 0x99,
	0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MsgClient interface {
	RegisterContract(ctx context.Context, in *MsgRegisterContract, opts ...grpc.CallOption) (*MsgRegisterContractResponse, error)
	DeRegisterContract(ctx context.Context, in *MsgDeRegisterContract, opts ...grpc.CallOption) (*MsgDeRegisterContractResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) RegisterContract(ctx context.Context, in *MsgRegisterContract, opts ...grpc.CallOption) (*MsgRegisterContractResponse, error) {
	out := new(MsgRegisterContractResponse)
	err := c.cc.Invoke(ctx, "/comdex.common.v1beta1.Msg/RegisterContract", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) DeRegisterContract(ctx context.Context, in *MsgDeRegisterContract, opts ...grpc.CallOption) (*MsgDeRegisterContractResponse, error) {
	out := new(MsgDeRegisterContractResponse)
	err := c.cc.Invoke(ctx, "/comdex.common.v1beta1.Msg/DeRegisterContract", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	RegisterContract(context.Context, *MsgRegisterContract) (*MsgRegisterContractResponse, error)
	DeRegisterContract(context.Context, *MsgDeRegisterContract) (*MsgDeRegisterContractResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) RegisterContract(ctx context.Context, req *MsgRegisterContract) (*MsgRegisterContractResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterContract not implemented")
}
func (*UnimplementedMsgServer) DeRegisterContract(ctx context.Context, req *MsgDeRegisterContract) (*MsgDeRegisterContractResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeRegisterContract not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_RegisterContract_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRegisterContract)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RegisterContract(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/comdex.common.v1beta1.Msg/RegisterContract",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RegisterContract(ctx, req.(*MsgRegisterContract))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_DeRegisterContract_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgDeRegisterContract)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).DeRegisterContract(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/comdex.common.v1beta1.Msg/DeRegisterContract",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).DeRegisterContract(ctx, req.(*MsgDeRegisterContract))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "comdex.common.v1beta1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterContract",
			Handler:    _Msg_RegisterContract_Handler,
		},
		{
			MethodName: "DeRegisterContract",
			Handler:    _Msg_DeRegisterContract_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "comdex/common/v1beta1/tx.proto",
}

func (m *MsgRegisterContract) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgRegisterContract) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgRegisterContract) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ContractAddr) > 0 {
		i -= len(m.ContractAddr)
		copy(dAtA[i:], m.ContractAddr)
		i = encodeVarintTx(dAtA, i, uint64(len(m.ContractAddr)))
		i--
		dAtA[i] = 0x22
	}
	if m.GameId != 0 {
		i = encodeVarintTx(dAtA, i, uint64(m.GameId))
		i--
		dAtA[i] = 0x18
	}
	if len(m.GameName) > 0 {
		i -= len(m.GameName)
		copy(dAtA[i:], m.GameName)
		i = encodeVarintTx(dAtA, i, uint64(len(m.GameName)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgRegisterContractResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgRegisterContractResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgRegisterContractResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgDeRegisterContract) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgDeRegisterContract) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgDeRegisterContract) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.GameId != 0 {
		i = encodeVarintTx(dAtA, i, uint64(m.GameId))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgDeRegisterContractResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgDeRegisterContractResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgDeRegisterContractResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func encodeVarintTx(dAtA []byte, offset int, v uint64) int {
	offset -= sovTx(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MsgRegisterContract) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.GameName)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if m.GameId != 0 {
		n += 1 + sovTx(uint64(m.GameId))
	}
	l = len(m.ContractAddr)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgRegisterContractResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgDeRegisterContract) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if m.GameId != 0 {
		n += 1 + sovTx(uint64(m.GameId))
	}
	return n
}

func (m *MsgDeRegisterContractResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func sovTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTx(x uint64) (n int) {
	return sovTx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgRegisterContract) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgRegisterContract: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgRegisterContract: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Creator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GameName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.GameName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GameId", wireType)
			}
			m.GameId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.GameId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ContractAddr", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ContractAddr = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (m *MsgRegisterContractResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgRegisterContractResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgRegisterContractResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (m *MsgDeRegisterContract) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgDeRegisterContract: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgDeRegisterContract: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Creator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GameId", wireType)
			}
			m.GameId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.GameId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (m *MsgDeRegisterContractResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgDeRegisterContractResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgDeRegisterContractResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func skipTx(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTx
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
					return 0, ErrIntOverflowTx
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
					return 0, ErrIntOverflowTx
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
				return 0, ErrInvalidLengthTx
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTx
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTx
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTx        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTx          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTx = fmt.Errorf("proto: unexpected end of group")
)
