// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: kyve/multi_coin_rewards/v1beta1/tx.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/cosmos/cosmos-sdk/types/msgservice"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
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

// MsgUpdateParams defines a SDK message for updating the module parameters.
type MsgUpdateParams struct {
	// authority is the address of the governance account.
	Authority string `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority,omitempty"`
	// payload defines the x/multi_coin_rewards parameters to update.
	Payload string `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (m *MsgUpdateParams) Reset()         { *m = MsgUpdateParams{} }
func (m *MsgUpdateParams) String() string { return proto.CompactTextString(m) }
func (*MsgUpdateParams) ProtoMessage()    {}
func (*MsgUpdateParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_702f1149b462214b, []int{0}
}
func (m *MsgUpdateParams) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgUpdateParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgUpdateParams.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgUpdateParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgUpdateParams.Merge(m, src)
}
func (m *MsgUpdateParams) XXX_Size() int {
	return m.Size()
}
func (m *MsgUpdateParams) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgUpdateParams.DiscardUnknown(m)
}

var xxx_messageInfo_MsgUpdateParams proto.InternalMessageInfo

func (m *MsgUpdateParams) GetAuthority() string {
	if m != nil {
		return m.Authority
	}
	return ""
}

func (m *MsgUpdateParams) GetPayload() string {
	if m != nil {
		return m.Payload
	}
	return ""
}

// MsgUpdateParamsResponse defines the Msg/UpdateParams response type.
type MsgUpdateParamsResponse struct {
}

func (m *MsgUpdateParamsResponse) Reset()         { *m = MsgUpdateParamsResponse{} }
func (m *MsgUpdateParamsResponse) String() string { return proto.CompactTextString(m) }
func (*MsgUpdateParamsResponse) ProtoMessage()    {}
func (*MsgUpdateParamsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_702f1149b462214b, []int{1}
}
func (m *MsgUpdateParamsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgUpdateParamsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgUpdateParamsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgUpdateParamsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgUpdateParamsResponse.Merge(m, src)
}
func (m *MsgUpdateParamsResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgUpdateParamsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgUpdateParamsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgUpdateParamsResponse proto.InternalMessageInfo

// MsgEnableMultiCoinReward enables multi-coin rewards for the sender address
// and claims all current pending rewards.
type MsgToggleMultiCoinRewards struct {
	// creator ...
	Creator string `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	// enabled ...
	Enabled bool `protobuf:"varint,2,opt,name=enabled,proto3" json:"enabled,omitempty"`
}

func (m *MsgToggleMultiCoinRewards) Reset()         { *m = MsgToggleMultiCoinRewards{} }
func (m *MsgToggleMultiCoinRewards) String() string { return proto.CompactTextString(m) }
func (*MsgToggleMultiCoinRewards) ProtoMessage()    {}
func (*MsgToggleMultiCoinRewards) Descriptor() ([]byte, []int) {
	return fileDescriptor_702f1149b462214b, []int{2}
}
func (m *MsgToggleMultiCoinRewards) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgToggleMultiCoinRewards) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgToggleMultiCoinRewards.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgToggleMultiCoinRewards) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgToggleMultiCoinRewards.Merge(m, src)
}
func (m *MsgToggleMultiCoinRewards) XXX_Size() int {
	return m.Size()
}
func (m *MsgToggleMultiCoinRewards) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgToggleMultiCoinRewards.DiscardUnknown(m)
}

var xxx_messageInfo_MsgToggleMultiCoinRewards proto.InternalMessageInfo

func (m *MsgToggleMultiCoinRewards) GetCreator() string {
	if m != nil {
		return m.Creator
	}
	return ""
}

func (m *MsgToggleMultiCoinRewards) GetEnabled() bool {
	if m != nil {
		return m.Enabled
	}
	return false
}

// MsgEnableMultiCoinRewardResponse ...
type MsgToggleMultiCoinRewardsResponse struct {
}

func (m *MsgToggleMultiCoinRewardsResponse) Reset()         { *m = MsgToggleMultiCoinRewardsResponse{} }
func (m *MsgToggleMultiCoinRewardsResponse) String() string { return proto.CompactTextString(m) }
func (*MsgToggleMultiCoinRewardsResponse) ProtoMessage()    {}
func (*MsgToggleMultiCoinRewardsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_702f1149b462214b, []int{3}
}
func (m *MsgToggleMultiCoinRewardsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgToggleMultiCoinRewardsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgToggleMultiCoinRewardsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgToggleMultiCoinRewardsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgToggleMultiCoinRewardsResponse.Merge(m, src)
}
func (m *MsgToggleMultiCoinRewardsResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgToggleMultiCoinRewardsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgToggleMultiCoinRewardsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgToggleMultiCoinRewardsResponse proto.InternalMessageInfo

// MsgEnableMultiCoinReward enables multi-coin rewards for the sender address
// and claims all current pending rewards.
type MsgSetMultiCoinRewardsDistributionPolicy struct {
	// creator ...
	Creator string `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	// policy ...
	Policy *MultiCoinDistributionPolicy `protobuf:"bytes,2,opt,name=policy,proto3" json:"policy,omitempty"`
}

func (m *MsgSetMultiCoinRewardsDistributionPolicy) Reset() {
	*m = MsgSetMultiCoinRewardsDistributionPolicy{}
}
func (m *MsgSetMultiCoinRewardsDistributionPolicy) String() string { return proto.CompactTextString(m) }
func (*MsgSetMultiCoinRewardsDistributionPolicy) ProtoMessage()    {}
func (*MsgSetMultiCoinRewardsDistributionPolicy) Descriptor() ([]byte, []int) {
	return fileDescriptor_702f1149b462214b, []int{4}
}
func (m *MsgSetMultiCoinRewardsDistributionPolicy) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgSetMultiCoinRewardsDistributionPolicy) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSetMultiCoinRewardsDistributionPolicy.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgSetMultiCoinRewardsDistributionPolicy) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSetMultiCoinRewardsDistributionPolicy.Merge(m, src)
}
func (m *MsgSetMultiCoinRewardsDistributionPolicy) XXX_Size() int {
	return m.Size()
}
func (m *MsgSetMultiCoinRewardsDistributionPolicy) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgSetMultiCoinRewardsDistributionPolicy.DiscardUnknown(m)
}

var xxx_messageInfo_MsgSetMultiCoinRewardsDistributionPolicy proto.InternalMessageInfo

func (m *MsgSetMultiCoinRewardsDistributionPolicy) GetCreator() string {
	if m != nil {
		return m.Creator
	}
	return ""
}

func (m *MsgSetMultiCoinRewardsDistributionPolicy) GetPolicy() *MultiCoinDistributionPolicy {
	if m != nil {
		return m.Policy
	}
	return nil
}

// MsgEnableMultiCoinRewardResponse ...
type MsgSetMultiCoinRewardsDistributionPolicyResponse struct {
}

func (m *MsgSetMultiCoinRewardsDistributionPolicyResponse) Reset() {
	*m = MsgSetMultiCoinRewardsDistributionPolicyResponse{}
}
func (m *MsgSetMultiCoinRewardsDistributionPolicyResponse) String() string {
	return proto.CompactTextString(m)
}
func (*MsgSetMultiCoinRewardsDistributionPolicyResponse) ProtoMessage() {}
func (*MsgSetMultiCoinRewardsDistributionPolicyResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_702f1149b462214b, []int{5}
}
func (m *MsgSetMultiCoinRewardsDistributionPolicyResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgSetMultiCoinRewardsDistributionPolicyResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSetMultiCoinRewardsDistributionPolicyResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgSetMultiCoinRewardsDistributionPolicyResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSetMultiCoinRewardsDistributionPolicyResponse.Merge(m, src)
}
func (m *MsgSetMultiCoinRewardsDistributionPolicyResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgSetMultiCoinRewardsDistributionPolicyResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgSetMultiCoinRewardsDistributionPolicyResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgSetMultiCoinRewardsDistributionPolicyResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgUpdateParams)(nil), "kyve.multi_coin_rewards.v1beta1.MsgUpdateParams")
	proto.RegisterType((*MsgUpdateParamsResponse)(nil), "kyve.multi_coin_rewards.v1beta1.MsgUpdateParamsResponse")
	proto.RegisterType((*MsgToggleMultiCoinRewards)(nil), "kyve.multi_coin_rewards.v1beta1.MsgToggleMultiCoinRewards")
	proto.RegisterType((*MsgToggleMultiCoinRewardsResponse)(nil), "kyve.multi_coin_rewards.v1beta1.MsgToggleMultiCoinRewardsResponse")
	proto.RegisterType((*MsgSetMultiCoinRewardsDistributionPolicy)(nil), "kyve.multi_coin_rewards.v1beta1.MsgSetMultiCoinRewardsDistributionPolicy")
	proto.RegisterType((*MsgSetMultiCoinRewardsDistributionPolicyResponse)(nil), "kyve.multi_coin_rewards.v1beta1.MsgSetMultiCoinRewardsDistributionPolicyResponse")
}

func init() {
	proto.RegisterFile("kyve/multi_coin_rewards/v1beta1/tx.proto", fileDescriptor_702f1149b462214b)
}

var fileDescriptor_702f1149b462214b = []byte{
	// 496 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x54, 0x4f, 0x6b, 0x13, 0x41,
	0x1c, 0xcd, 0x58, 0xac, 0x76, 0x2c, 0x0a, 0x8b, 0xd8, 0x64, 0x0f, 0xab, 0x46, 0x0f, 0xa1, 0xe2,
	0x4e, 0x13, 0x41, 0x25, 0x78, 0x31, 0xea, 0x41, 0x64, 0xa5, 0xdd, 0x56, 0x41, 0x41, 0xc2, 0xec,
	0xee, 0x30, 0x19, 0xba, 0xbb, 0xb3, 0xcc, 0xcc, 0xa6, 0x5d, 0x4f, 0xe2, 0x27, 0xf0, 0xe8, 0x47,
	0xf0, 0xd8, 0x83, 0x1f, 0xc1, 0x83, 0xc7, 0xe2, 0xc9, 0xa3, 0x24, 0x87, 0x82, 0x9f, 0x42, 0xf6,
	0x5f, 0x4a, 0x93, 0xd4, 0xac, 0xd2, 0xd3, 0xf2, 0x63, 0xde, 0xfb, 0xbd, 0xf7, 0xf6, 0xc1, 0x0f,
	0xb6, 0x76, 0x93, 0x21, 0x41, 0x41, 0xec, 0x2b, 0xd6, 0x77, 0x39, 0x0b, 0xfb, 0x82, 0xec, 0x61,
	0xe1, 0x49, 0x34, 0x6c, 0x3b, 0x44, 0xe1, 0x36, 0x52, 0xfb, 0x66, 0x24, 0xb8, 0xe2, 0xda, 0xf5,
	0x14, 0x69, 0xce, 0x22, 0xcd, 0x02, 0xa9, 0xaf, 0xb9, 0x5c, 0x06, 0x5c, 0xa2, 0x40, 0x52, 0x34,
	0x6c, 0xa7, 0x9f, 0x9c, 0xa9, 0x37, 0xf2, 0x87, 0x7e, 0x36, 0xa1, 0x7c, 0x28, 0x9e, 0xee, 0x2c,
	0x94, 0x4f, 0x22, 0x52, 0x80, 0x9b, 0x12, 0x5e, 0xb1, 0x24, 0x7d, 0x15, 0x79, 0x58, 0x91, 0x4d,
	0x2c, 0x70, 0x20, 0xb5, 0xfb, 0x70, 0x05, 0xc7, 0x6a, 0xc0, 0x05, 0x53, 0x49, 0x1d, 0xdc, 0x00,
	0xad, 0x95, 0x5e, 0xfd, 0xc7, 0xd7, 0xbb, 0x57, 0x0b, 0x91, 0xc7, 0x9e, 0x27, 0x88, 0x94, 0xdb,
	0x4a, 0xb0, 0x90, 0xda, 0xc7, 0x50, 0xad, 0x0e, 0x2f, 0x44, 0x38, 0xf1, 0x39, 0xf6, 0xea, 0xe7,
	0x52, 0x96, 0x5d, 0x8e, 0xdd, 0xcb, 0x1f, 0x8f, 0x0e, 0xd6, 0x8f, 0x91, 0xcd, 0x06, 0x5c, 0x9b,
	0x12, 0xb5, 0x89, 0x8c, 0x78, 0x28, 0x49, 0xf3, 0x1d, 0x6c, 0x58, 0x92, 0xee, 0x70, 0x4a, 0x7d,
	0x62, 0xa5, 0x11, 0x9e, 0x70, 0x16, 0xda, 0x79, 0x80, 0x54, 0xc1, 0x15, 0x04, 0x2b, 0x2e, 0x72,
	0x5f, 0x76, 0x39, 0xa6, 0x2f, 0x24, 0xc4, 0x8e, 0x4f, 0x72, 0xed, 0x8b, 0x76, 0x39, 0x76, 0x57,
	0x53, 0xed, 0x12, 0xd7, 0xbc, 0x05, 0x6f, 0x9e, 0xba, 0x7e, 0xe2, 0xe1, 0x0b, 0x80, 0x2d, 0x4b,
	0xd2, 0x6d, 0xa2, 0xa6, 0x21, 0x4f, 0x99, 0x54, 0x82, 0x39, 0xb1, 0x62, 0x3c, 0xdc, 0xe4, 0x3e,
	0x73, 0x93, 0xbf, 0x78, 0xda, 0x81, 0xcb, 0x51, 0x86, 0xc9, 0x2c, 0x5d, 0xea, 0x3c, 0x32, 0x17,
	0xb4, 0x6d, 0x4e, 0xe4, 0x66, 0x75, 0xec, 0x62, 0xd7, 0x54, 0x9e, 0x0e, 0xdc, 0xa8, 0xea, 0xb4,
	0x8c, 0xd7, 0xf9, 0xbd, 0x04, 0x97, 0x2c, 0x49, 0xb5, 0xf7, 0x70, 0xf5, 0x44, 0xef, 0x1b, 0x8b,
	0xfd, 0x9d, 0x2c, 0x4d, 0x7f, 0xf8, 0xaf, 0x8c, 0xd2, 0x83, 0xf6, 0x19, 0xc0, 0x6b, 0xa7, 0x94,
	0xdc, 0xad, 0xb2, 0x74, 0x3e, 0x57, 0xef, 0xfd, 0x3f, 0x77, 0x62, 0xed, 0x1b, 0x80, 0xb7, 0x67,
	0x7f, 0xe8, 0x9c, 0xe6, 0x9f, 0x57, 0x11, 0xab, 0x54, 0x8d, 0xbe, 0x75, 0x66, 0xab, 0xca, 0x18,
	0xfa, 0xf9, 0x0f, 0x47, 0x07, 0xeb, 0xa0, 0xb7, 0xf5, 0x7d, 0x64, 0x80, 0xc3, 0x91, 0x01, 0x7e,
	0x8d, 0x0c, 0xf0, 0x69, 0x6c, 0xd4, 0x0e, 0xc7, 0x46, 0xed, 0xe7, 0xd8, 0xa8, 0xbd, 0x7d, 0x40,
	0x99, 0x1a, 0xc4, 0x8e, 0xe9, 0xf2, 0x00, 0xbd, 0x78, 0xf3, 0xfa, 0xd9, 0x4b, 0xa2, 0xf6, 0xb8,
	0xd8, 0x45, 0xee, 0x00, 0xb3, 0x10, 0xed, 0xcf, 0x3b, 0x20, 0xd9, 0xe1, 0x70, 0x96, 0xb3, 0xcb,
	0x71, 0xef, 0x4f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x73, 0x23, 0x38, 0x46, 0xe7, 0x04, 0x00, 0x00,
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
	// UpdateParams defines a governance operation for updating the x/multi_coin_rewards module
	// parameters. The authority is hard-coded to the x/gov module account.
	UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error)
	// ToggleMultiCoinRewards ...
	ToggleMultiCoinRewards(ctx context.Context, in *MsgToggleMultiCoinRewards, opts ...grpc.CallOption) (*MsgToggleMultiCoinRewardsResponse, error)
	// SetMultiCoinRewardDistributionPolicy ...
	SetMultiCoinRewardDistributionPolicy(ctx context.Context, in *MsgSetMultiCoinRewardsDistributionPolicy, opts ...grpc.CallOption) (*MsgSetMultiCoinRewardsDistributionPolicyResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error) {
	out := new(MsgUpdateParamsResponse)
	err := c.cc.Invoke(ctx, "/kyve.multi_coin_rewards.v1beta1.Msg/UpdateParams", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) ToggleMultiCoinRewards(ctx context.Context, in *MsgToggleMultiCoinRewards, opts ...grpc.CallOption) (*MsgToggleMultiCoinRewardsResponse, error) {
	out := new(MsgToggleMultiCoinRewardsResponse)
	err := c.cc.Invoke(ctx, "/kyve.multi_coin_rewards.v1beta1.Msg/ToggleMultiCoinRewards", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) SetMultiCoinRewardDistributionPolicy(ctx context.Context, in *MsgSetMultiCoinRewardsDistributionPolicy, opts ...grpc.CallOption) (*MsgSetMultiCoinRewardsDistributionPolicyResponse, error) {
	out := new(MsgSetMultiCoinRewardsDistributionPolicyResponse)
	err := c.cc.Invoke(ctx, "/kyve.multi_coin_rewards.v1beta1.Msg/SetMultiCoinRewardDistributionPolicy", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	// UpdateParams defines a governance operation for updating the x/multi_coin_rewards module
	// parameters. The authority is hard-coded to the x/gov module account.
	UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error)
	// ToggleMultiCoinRewards ...
	ToggleMultiCoinRewards(context.Context, *MsgToggleMultiCoinRewards) (*MsgToggleMultiCoinRewardsResponse, error)
	// SetMultiCoinRewardDistributionPolicy ...
	SetMultiCoinRewardDistributionPolicy(context.Context, *MsgSetMultiCoinRewardsDistributionPolicy) (*MsgSetMultiCoinRewardsDistributionPolicyResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) UpdateParams(ctx context.Context, req *MsgUpdateParams) (*MsgUpdateParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateParams not implemented")
}
func (*UnimplementedMsgServer) ToggleMultiCoinRewards(ctx context.Context, req *MsgToggleMultiCoinRewards) (*MsgToggleMultiCoinRewardsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ToggleMultiCoinRewards not implemented")
}
func (*UnimplementedMsgServer) SetMultiCoinRewardDistributionPolicy(ctx context.Context, req *MsgSetMultiCoinRewardsDistributionPolicy) (*MsgSetMultiCoinRewardsDistributionPolicyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetMultiCoinRewardDistributionPolicy not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_UpdateParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kyve.multi_coin_rewards.v1beta1.Msg/UpdateParams",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateParams(ctx, req.(*MsgUpdateParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_ToggleMultiCoinRewards_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgToggleMultiCoinRewards)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ToggleMultiCoinRewards(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kyve.multi_coin_rewards.v1beta1.Msg/ToggleMultiCoinRewards",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ToggleMultiCoinRewards(ctx, req.(*MsgToggleMultiCoinRewards))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_SetMultiCoinRewardDistributionPolicy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSetMultiCoinRewardsDistributionPolicy)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SetMultiCoinRewardDistributionPolicy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kyve.multi_coin_rewards.v1beta1.Msg/SetMultiCoinRewardDistributionPolicy",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SetMultiCoinRewardDistributionPolicy(ctx, req.(*MsgSetMultiCoinRewardsDistributionPolicy))
	}
	return interceptor(ctx, in, info, handler)
}

var Msg_serviceDesc = _Msg_serviceDesc
var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "kyve.multi_coin_rewards.v1beta1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdateParams",
			Handler:    _Msg_UpdateParams_Handler,
		},
		{
			MethodName: "ToggleMultiCoinRewards",
			Handler:    _Msg_ToggleMultiCoinRewards_Handler,
		},
		{
			MethodName: "SetMultiCoinRewardDistributionPolicy",
			Handler:    _Msg_SetMultiCoinRewardDistributionPolicy_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "kyve/multi_coin_rewards/v1beta1/tx.proto",
}

func (m *MsgUpdateParams) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgUpdateParams) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgUpdateParams) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Payload) > 0 {
		i -= len(m.Payload)
		copy(dAtA[i:], m.Payload)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Payload)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Authority) > 0 {
		i -= len(m.Authority)
		copy(dAtA[i:], m.Authority)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Authority)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgUpdateParamsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgUpdateParamsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgUpdateParamsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgToggleMultiCoinRewards) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgToggleMultiCoinRewards) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgToggleMultiCoinRewards) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Enabled {
		i--
		if m.Enabled {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
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

func (m *MsgToggleMultiCoinRewardsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgToggleMultiCoinRewardsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgToggleMultiCoinRewardsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgSetMultiCoinRewardsDistributionPolicy) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgSetMultiCoinRewardsDistributionPolicy) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgSetMultiCoinRewardsDistributionPolicy) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Policy != nil {
		{
			size, err := m.Policy.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTx(dAtA, i, uint64(size))
		}
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

func (m *MsgSetMultiCoinRewardsDistributionPolicyResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgSetMultiCoinRewardsDistributionPolicyResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgSetMultiCoinRewardsDistributionPolicyResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
func (m *MsgUpdateParams) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Authority)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Payload)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgUpdateParamsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgToggleMultiCoinRewards) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if m.Enabled {
		n += 2
	}
	return n
}

func (m *MsgToggleMultiCoinRewardsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgSetMultiCoinRewardsDistributionPolicy) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if m.Policy != nil {
		l = m.Policy.Size()
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgSetMultiCoinRewardsDistributionPolicyResponse) Size() (n int) {
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
func (m *MsgUpdateParams) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgUpdateParams: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgUpdateParams: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Authority", wireType)
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
			m.Authority = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Payload", wireType)
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
			m.Payload = string(dAtA[iNdEx:postIndex])
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
func (m *MsgUpdateParamsResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgUpdateParamsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgUpdateParamsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *MsgToggleMultiCoinRewards) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgToggleMultiCoinRewards: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgToggleMultiCoinRewards: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field Enabled", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
			m.Enabled = bool(v != 0)
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
func (m *MsgToggleMultiCoinRewardsResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgToggleMultiCoinRewardsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgToggleMultiCoinRewardsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *MsgSetMultiCoinRewardsDistributionPolicy) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgSetMultiCoinRewardsDistributionPolicy: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgSetMultiCoinRewardsDistributionPolicy: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field Policy", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Policy == nil {
				m.Policy = &MultiCoinDistributionPolicy{}
			}
			if err := m.Policy.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
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
func (m *MsgSetMultiCoinRewardsDistributionPolicyResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgSetMultiCoinRewardsDistributionPolicyResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgSetMultiCoinRewardsDistributionPolicyResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
