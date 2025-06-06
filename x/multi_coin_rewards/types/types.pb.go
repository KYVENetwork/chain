// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: kyve/multi_coin_rewards/v1beta1/types.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
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

// UnbondingState stores the state for the unbonding of stakes and delegations.
type QueueState struct {
	// low_index is the tail of the queue. It is the
	// oldest entry in the queue. If this entry isn't
	// due, non of the other entries is.
	LowIndex uint64 `protobuf:"varint,1,opt,name=low_index,json=lowIndex,proto3" json:"low_index,omitempty"`
	// high_index is the head of the queue. New entries
	// are added to the top.
	HighIndex uint64 `protobuf:"varint,2,opt,name=high_index,json=highIndex,proto3" json:"high_index,omitempty"`
}

func (m *QueueState) Reset()         { *m = QueueState{} }
func (m *QueueState) String() string { return proto.CompactTextString(m) }
func (*QueueState) ProtoMessage()    {}
func (*QueueState) Descriptor() ([]byte, []int) {
	return fileDescriptor_3f0d9743633e1637, []int{0}
}
func (m *QueueState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueueState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueueState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueueState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueueState.Merge(m, src)
}
func (m *QueueState) XXX_Size() int {
	return m.Size()
}
func (m *QueueState) XXX_DiscardUnknown() {
	xxx_messageInfo_QueueState.DiscardUnknown(m)
}

var xxx_messageInfo_QueueState proto.InternalMessageInfo

func (m *QueueState) GetLowIndex() uint64 {
	if m != nil {
		return m.LowIndex
	}
	return 0
}

func (m *QueueState) GetHighIndex() uint64 {
	if m != nil {
		return m.HighIndex
	}
	return 0
}

// MultiCoinPendingRewardsEntry ...
type MultiCoinPendingRewardsEntry struct {
	// index is needed for the queue-algorithm which
	// processes the commission changes
	Index uint64 `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty"`
	// address ...
	Address string `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	// rewards ...
	Rewards      github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,3,rep,name=rewards,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"rewards"`
	CreationDate int64                                    `protobuf:"varint,4,opt,name=creation_date,json=creationDate,proto3" json:"creation_date,omitempty"`
}

func (m *MultiCoinPendingRewardsEntry) Reset()         { *m = MultiCoinPendingRewardsEntry{} }
func (m *MultiCoinPendingRewardsEntry) String() string { return proto.CompactTextString(m) }
func (*MultiCoinPendingRewardsEntry) ProtoMessage()    {}
func (*MultiCoinPendingRewardsEntry) Descriptor() ([]byte, []int) {
	return fileDescriptor_3f0d9743633e1637, []int{1}
}
func (m *MultiCoinPendingRewardsEntry) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MultiCoinPendingRewardsEntry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MultiCoinPendingRewardsEntry.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MultiCoinPendingRewardsEntry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MultiCoinPendingRewardsEntry.Merge(m, src)
}
func (m *MultiCoinPendingRewardsEntry) XXX_Size() int {
	return m.Size()
}
func (m *MultiCoinPendingRewardsEntry) XXX_DiscardUnknown() {
	xxx_messageInfo_MultiCoinPendingRewardsEntry.DiscardUnknown(m)
}

var xxx_messageInfo_MultiCoinPendingRewardsEntry proto.InternalMessageInfo

func (m *MultiCoinPendingRewardsEntry) GetIndex() uint64 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *MultiCoinPendingRewardsEntry) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *MultiCoinPendingRewardsEntry) GetRewards() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.Rewards
	}
	return nil
}

func (m *MultiCoinPendingRewardsEntry) GetCreationDate() int64 {
	if m != nil {
		return m.CreationDate
	}
	return 0
}

// MultiCoinDistributionPolicy ...
type MultiCoinDistributionPolicy struct {
	Entries []*MultiCoinDistributionDenomEntry `protobuf:"bytes,1,rep,name=entries,proto3" json:"entries,omitempty"`
}

func (m *MultiCoinDistributionPolicy) Reset()         { *m = MultiCoinDistributionPolicy{} }
func (m *MultiCoinDistributionPolicy) String() string { return proto.CompactTextString(m) }
func (*MultiCoinDistributionPolicy) ProtoMessage()    {}
func (*MultiCoinDistributionPolicy) Descriptor() ([]byte, []int) {
	return fileDescriptor_3f0d9743633e1637, []int{2}
}
func (m *MultiCoinDistributionPolicy) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MultiCoinDistributionPolicy) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MultiCoinDistributionPolicy.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MultiCoinDistributionPolicy) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MultiCoinDistributionPolicy.Merge(m, src)
}
func (m *MultiCoinDistributionPolicy) XXX_Size() int {
	return m.Size()
}
func (m *MultiCoinDistributionPolicy) XXX_DiscardUnknown() {
	xxx_messageInfo_MultiCoinDistributionPolicy.DiscardUnknown(m)
}

var xxx_messageInfo_MultiCoinDistributionPolicy proto.InternalMessageInfo

func (m *MultiCoinDistributionPolicy) GetEntries() []*MultiCoinDistributionDenomEntry {
	if m != nil {
		return m.Entries
	}
	return nil
}

// MultiCoinDistributionDenomEntry ...
type MultiCoinDistributionDenomEntry struct {
	Denom       string                                  `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty"`
	PoolWeights []*MultiCoinDistributionPoolWeightEntry `protobuf:"bytes,2,rep,name=pool_weights,json=poolWeights,proto3" json:"pool_weights,omitempty"`
}

func (m *MultiCoinDistributionDenomEntry) Reset()         { *m = MultiCoinDistributionDenomEntry{} }
func (m *MultiCoinDistributionDenomEntry) String() string { return proto.CompactTextString(m) }
func (*MultiCoinDistributionDenomEntry) ProtoMessage()    {}
func (*MultiCoinDistributionDenomEntry) Descriptor() ([]byte, []int) {
	return fileDescriptor_3f0d9743633e1637, []int{3}
}
func (m *MultiCoinDistributionDenomEntry) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MultiCoinDistributionDenomEntry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MultiCoinDistributionDenomEntry.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MultiCoinDistributionDenomEntry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MultiCoinDistributionDenomEntry.Merge(m, src)
}
func (m *MultiCoinDistributionDenomEntry) XXX_Size() int {
	return m.Size()
}
func (m *MultiCoinDistributionDenomEntry) XXX_DiscardUnknown() {
	xxx_messageInfo_MultiCoinDistributionDenomEntry.DiscardUnknown(m)
}

var xxx_messageInfo_MultiCoinDistributionDenomEntry proto.InternalMessageInfo

func (m *MultiCoinDistributionDenomEntry) GetDenom() string {
	if m != nil {
		return m.Denom
	}
	return ""
}

func (m *MultiCoinDistributionDenomEntry) GetPoolWeights() []*MultiCoinDistributionPoolWeightEntry {
	if m != nil {
		return m.PoolWeights
	}
	return nil
}

// MultiCoinDistributionPoolWeightEntry ...
type MultiCoinDistributionPoolWeightEntry struct {
	PoolId uint64                      `protobuf:"varint,1,opt,name=pool_id,json=poolId,proto3" json:"pool_id,omitempty"`
	Weight cosmossdk_io_math.LegacyDec `protobuf:"bytes,2,opt,name=weight,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"weight"`
}

func (m *MultiCoinDistributionPoolWeightEntry) Reset()         { *m = MultiCoinDistributionPoolWeightEntry{} }
func (m *MultiCoinDistributionPoolWeightEntry) String() string { return proto.CompactTextString(m) }
func (*MultiCoinDistributionPoolWeightEntry) ProtoMessage()    {}
func (*MultiCoinDistributionPoolWeightEntry) Descriptor() ([]byte, []int) {
	return fileDescriptor_3f0d9743633e1637, []int{4}
}
func (m *MultiCoinDistributionPoolWeightEntry) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MultiCoinDistributionPoolWeightEntry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MultiCoinDistributionPoolWeightEntry.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MultiCoinDistributionPoolWeightEntry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MultiCoinDistributionPoolWeightEntry.Merge(m, src)
}
func (m *MultiCoinDistributionPoolWeightEntry) XXX_Size() int {
	return m.Size()
}
func (m *MultiCoinDistributionPoolWeightEntry) XXX_DiscardUnknown() {
	xxx_messageInfo_MultiCoinDistributionPoolWeightEntry.DiscardUnknown(m)
}

var xxx_messageInfo_MultiCoinDistributionPoolWeightEntry proto.InternalMessageInfo

func (m *MultiCoinDistributionPoolWeightEntry) GetPoolId() uint64 {
	if m != nil {
		return m.PoolId
	}
	return 0
}

func init() {
	proto.RegisterType((*QueueState)(nil), "kyve.multi_coin_rewards.v1beta1.QueueState")
	proto.RegisterType((*MultiCoinPendingRewardsEntry)(nil), "kyve.multi_coin_rewards.v1beta1.MultiCoinPendingRewardsEntry")
	proto.RegisterType((*MultiCoinDistributionPolicy)(nil), "kyve.multi_coin_rewards.v1beta1.MultiCoinDistributionPolicy")
	proto.RegisterType((*MultiCoinDistributionDenomEntry)(nil), "kyve.multi_coin_rewards.v1beta1.MultiCoinDistributionDenomEntry")
	proto.RegisterType((*MultiCoinDistributionPoolWeightEntry)(nil), "kyve.multi_coin_rewards.v1beta1.MultiCoinDistributionPoolWeightEntry")
}

func init() {
	proto.RegisterFile("kyve/multi_coin_rewards/v1beta1/types.proto", fileDescriptor_3f0d9743633e1637)
}

var fileDescriptor_3f0d9743633e1637 = []byte{
	// 543 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x53, 0x4f, 0x6b, 0x13, 0x41,
	0x14, 0xcf, 0xda, 0xda, 0x98, 0x69, 0x3d, 0x38, 0x14, 0x8c, 0xad, 0x6e, 0x42, 0xea, 0x21, 0x28,
	0xce, 0x52, 0x45, 0x3c, 0x78, 0x91, 0x98, 0x80, 0xc5, 0x3f, 0xa4, 0x2b, 0x28, 0xf6, 0x12, 0x26,
	0x3b, 0xc3, 0xee, 0x98, 0xdd, 0x79, 0x61, 0x67, 0xd2, 0x74, 0xc1, 0x0f, 0xe1, 0x57, 0xf0, 0x26,
	0x9e, 0xfc, 0x18, 0x3d, 0xf6, 0x28, 0x1e, 0x5a, 0x49, 0x0e, 0x7e, 0x0d, 0x99, 0x99, 0x4d, 0x28,
	0x58, 0xa8, 0x78, 0xd9, 0x9d, 0xf7, 0xef, 0x37, 0xbf, 0xf7, 0x7b, 0x6f, 0xd0, 0xfd, 0x51, 0x71,
	0xc8, 0x83, 0x6c, 0x92, 0x6a, 0x31, 0x88, 0x40, 0xc8, 0x41, 0xce, 0xa7, 0x34, 0x67, 0x2a, 0x38,
	0xdc, 0x1d, 0x72, 0x4d, 0x77, 0x03, 0x5d, 0x8c, 0xb9, 0x22, 0xe3, 0x1c, 0x34, 0xe0, 0x86, 0x49,
	0x26, 0x7f, 0x27, 0x93, 0x32, 0x79, 0xeb, 0x06, 0xcd, 0x84, 0x84, 0xc0, 0x7e, 0x5d, 0xcd, 0x96,
	0x1f, 0x81, 0xca, 0x40, 0x05, 0x43, 0xaa, 0xf8, 0x12, 0xd4, 0x14, 0x97, 0xf1, 0xcd, 0x18, 0x62,
	0xb0, 0xc7, 0xc0, 0x9c, 0x9c, 0xb7, 0xf5, 0x02, 0xa1, 0xfd, 0x09, 0x9f, 0xf0, 0xb7, 0x9a, 0x6a,
	0x8e, 0xb7, 0x51, 0x2d, 0x85, 0xe9, 0x40, 0x48, 0xc6, 0x8f, 0xea, 0x5e, 0xd3, 0x6b, 0xaf, 0x86,
	0xd7, 0x52, 0x98, 0xee, 0x19, 0x1b, 0xdf, 0x41, 0x28, 0x11, 0x71, 0x52, 0x46, 0xaf, 0xd8, 0x68,
	0xcd, 0x78, 0x6c, 0xb8, 0x75, 0xe6, 0xa1, 0xdb, 0xaf, 0x0d, 0xe3, 0xe7, 0x20, 0x64, 0x9f, 0x4b,
	0x26, 0x64, 0x1c, 0x3a, 0xda, 0x3d, 0xa9, 0xf3, 0x02, 0x6f, 0xa2, 0xab, 0xe7, 0x81, 0x9d, 0x81,
	0xeb, 0xa8, 0x4a, 0x19, 0xcb, 0xb9, 0x52, 0x16, 0xb2, 0x16, 0x2e, 0x4c, 0xfc, 0x11, 0x55, 0xcb,
	0xb6, 0xeb, 0x2b, 0xcd, 0x95, 0xf6, 0xfa, 0xc3, 0x5b, 0xc4, 0xb5, 0x48, 0x4c, 0x8b, 0x0b, 0x29,
	0x88, 0xb9, 0xae, 0xf3, 0xf8, 0xf8, 0xb4, 0x51, 0xf9, 0x76, 0xd6, 0x68, 0xc7, 0x42, 0x27, 0x93,
	0x21, 0x89, 0x20, 0x0b, 0x4a, 0x3d, 0xdc, 0xef, 0x81, 0x62, 0xa3, 0x52, 0x62, 0x53, 0xa0, 0xbe,
	0xfe, 0xfe, 0x7e, 0xcf, 0x0b, 0x17, 0x17, 0xe0, 0x1d, 0x74, 0x3d, 0xca, 0x39, 0xd5, 0x02, 0xe4,
	0x80, 0x51, 0xcd, 0xeb, 0xab, 0x4d, 0xaf, 0xbd, 0x12, 0x6e, 0x2c, 0x9c, 0x5d, 0xaa, 0x79, 0xab,
	0x40, 0xdb, 0xcb, 0x06, 0xbb, 0x42, 0xe9, 0x5c, 0x0c, 0x27, 0x26, 0xd8, 0x87, 0x54, 0x44, 0x05,
	0x3e, 0x40, 0x55, 0x2e, 0x75, 0x2e, 0xb8, 0xaa, 0x7b, 0x96, 0xef, 0x33, 0x72, 0xc9, 0x18, 0xc9,
	0x85, 0x70, 0x5d, 0x2e, 0x21, 0xb3, 0x92, 0x85, 0x0b, 0xc0, 0xd6, 0x17, 0x0f, 0x35, 0x2e, 0x49,
	0x36, 0xfa, 0x32, 0x63, 0x59, 0x7d, 0x6b, 0xa1, 0x33, 0x70, 0x82, 0x36, 0xc6, 0x00, 0xe9, 0x60,
	0xca, 0x45, 0x9c, 0x68, 0x23, 0xb2, 0xa1, 0xd6, 0xfb, 0x3f, 0x6a, 0x7d, 0x80, 0xf4, 0xbd, 0x05,
	0x72, 0xfc, 0xd6, 0xc7, 0x4b, 0x87, 0x6a, 0x7d, 0x42, 0x77, 0xff, 0xa5, 0x08, 0xdf, 0x44, 0x55,
	0xcb, 0x48, 0xb0, 0x72, 0x13, 0xd6, 0x8c, 0xb9, 0xc7, 0xf0, 0x53, 0xb4, 0xe6, 0x58, 0xba, 0x4d,
	0xe8, 0xec, 0x98, 0xa1, 0xfe, 0x3c, 0x6d, 0x6c, 0xbb, 0x11, 0x2a, 0x36, 0x22, 0x02, 0x82, 0x8c,
	0xea, 0x84, 0xbc, 0xe2, 0x31, 0x8d, 0x8a, 0x2e, 0x8f, 0xc2, 0xb2, 0xa4, 0xb3, 0x7f, 0x3c, 0xf3,
	0xbd, 0x93, 0x99, 0xef, 0xfd, 0x9a, 0xf9, 0xde, 0xe7, 0xb9, 0x5f, 0x39, 0x99, 0xfb, 0x95, 0x1f,
	0x73, 0xbf, 0x72, 0xf0, 0xe4, 0xdc, 0x4e, 0xbc, 0xfc, 0xf0, 0xae, 0xf7, 0x86, 0xeb, 0x29, 0xe4,
	0xa3, 0x20, 0x4a, 0xa8, 0x90, 0xc1, 0xd1, 0x45, 0x6f, 0xd2, 0x2e, 0xca, 0x70, 0xcd, 0x3e, 0x91,
	0x47, 0x7f, 0x02, 0x00, 0x00, 0xff, 0xff, 0x06, 0xf2, 0x6f, 0x51, 0xbb, 0x03, 0x00, 0x00,
}

func (m *QueueState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueueState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueueState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.HighIndex != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.HighIndex))
		i--
		dAtA[i] = 0x10
	}
	if m.LowIndex != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.LowIndex))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *MultiCoinPendingRewardsEntry) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MultiCoinPendingRewardsEntry) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MultiCoinPendingRewardsEntry) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.CreationDate != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.CreationDate))
		i--
		dAtA[i] = 0x20
	}
	if len(m.Rewards) > 0 {
		for iNdEx := len(m.Rewards) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Rewards[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTypes(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0x12
	}
	if m.Index != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Index))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *MultiCoinDistributionPolicy) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MultiCoinDistributionPolicy) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MultiCoinDistributionPolicy) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Entries) > 0 {
		for iNdEx := len(m.Entries) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Entries[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTypes(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *MultiCoinDistributionDenomEntry) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MultiCoinDistributionDenomEntry) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MultiCoinDistributionDenomEntry) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.PoolWeights) > 0 {
		for iNdEx := len(m.PoolWeights) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.PoolWeights[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTypes(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MultiCoinDistributionPoolWeightEntry) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MultiCoinDistributionPoolWeightEntry) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MultiCoinDistributionPoolWeightEntry) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.Weight.Size()
		i -= size
		if _, err := m.Weight.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if m.PoolId != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.PoolId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintTypes(dAtA []byte, offset int, v uint64) int {
	offset -= sovTypes(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueueState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.LowIndex != 0 {
		n += 1 + sovTypes(uint64(m.LowIndex))
	}
	if m.HighIndex != 0 {
		n += 1 + sovTypes(uint64(m.HighIndex))
	}
	return n
}

func (m *MultiCoinPendingRewardsEntry) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Index != 0 {
		n += 1 + sovTypes(uint64(m.Index))
	}
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	if len(m.Rewards) > 0 {
		for _, e := range m.Rewards {
			l = e.Size()
			n += 1 + l + sovTypes(uint64(l))
		}
	}
	if m.CreationDate != 0 {
		n += 1 + sovTypes(uint64(m.CreationDate))
	}
	return n
}

func (m *MultiCoinDistributionPolicy) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Entries) > 0 {
		for _, e := range m.Entries {
			l = e.Size()
			n += 1 + l + sovTypes(uint64(l))
		}
	}
	return n
}

func (m *MultiCoinDistributionDenomEntry) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	if len(m.PoolWeights) > 0 {
		for _, e := range m.PoolWeights {
			l = e.Size()
			n += 1 + l + sovTypes(uint64(l))
		}
	}
	return n
}

func (m *MultiCoinDistributionPoolWeightEntry) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.PoolId != 0 {
		n += 1 + sovTypes(uint64(m.PoolId))
	}
	l = m.Weight.Size()
	n += 1 + l + sovTypes(uint64(l))
	return n
}

func sovTypes(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTypes(x uint64) (n int) {
	return sovTypes(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueueState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: QueueState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueueState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LowIndex", wireType)
			}
			m.LowIndex = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LowIndex |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field HighIndex", wireType)
			}
			m.HighIndex = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.HighIndex |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
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
func (m *MultiCoinPendingRewardsEntry) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: MultiCoinPendingRewardsEntry: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MultiCoinPendingRewardsEntry: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Index", wireType)
			}
			m.Index = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Index |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Rewards", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Rewards = append(m.Rewards, types.Coin{})
			if err := m.Rewards[len(m.Rewards)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CreationDate", wireType)
			}
			m.CreationDate = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CreationDate |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
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
func (m *MultiCoinDistributionPolicy) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: MultiCoinDistributionPolicy: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MultiCoinDistributionPolicy: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Entries", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Entries = append(m.Entries, &MultiCoinDistributionDenomEntry{})
			if err := m.Entries[len(m.Entries)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
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
func (m *MultiCoinDistributionDenomEntry) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: MultiCoinDistributionDenomEntry: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MultiCoinDistributionDenomEntry: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolWeights", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PoolWeights = append(m.PoolWeights, &MultiCoinDistributionPoolWeightEntry{})
			if err := m.PoolWeights[len(m.PoolWeights)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
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
func (m *MultiCoinDistributionPoolWeightEntry) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: MultiCoinDistributionPoolWeightEntry: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MultiCoinDistributionPoolWeightEntry: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolId", wireType)
			}
			m.PoolId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PoolId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Weight", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Weight.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
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
func skipTypes(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTypes
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
					return 0, ErrIntOverflowTypes
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
					return 0, ErrIntOverflowTypes
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
				return 0, ErrInvalidLengthTypes
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTypes
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTypes
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTypes        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTypes          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTypes = fmt.Errorf("proto: unexpected end of group")
)
