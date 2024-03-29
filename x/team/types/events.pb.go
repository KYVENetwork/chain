// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: kyve/team/v1beta1/events.proto

package types

import (
	fmt "fmt"
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

// MsgCreateTeamVestingAccount is an event emitted when a new team vesting account gets created.
// emitted_by: MsgCreateTeamVestingAccount
type EventCreateTeamVestingAccount struct {
	// authority which initiated this action
	Authority string `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority,omitempty"`
	// id is a unique identify for each vesting account, tied to a single team member.
	Id uint64 `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	// total_allocation is the number of tokens reserved for this team member.
	TotalAllocation uint64 `protobuf:"varint,3,opt,name=total_allocation,json=totalAllocation,proto3" json:"total_allocation,omitempty"`
	// commencement is the unix timestamp of the member's official start date.
	Commencement uint64 `protobuf:"varint,4,opt,name=commencement,proto3" json:"commencement,omitempty"`
}

func (m *EventCreateTeamVestingAccount) Reset()         { *m = EventCreateTeamVestingAccount{} }
func (m *EventCreateTeamVestingAccount) String() string { return proto.CompactTextString(m) }
func (*EventCreateTeamVestingAccount) ProtoMessage()    {}
func (*EventCreateTeamVestingAccount) Descriptor() ([]byte, []int) {
	return fileDescriptor_198acea0777f469a, []int{0}
}
func (m *EventCreateTeamVestingAccount) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EventCreateTeamVestingAccount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EventCreateTeamVestingAccount.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EventCreateTeamVestingAccount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventCreateTeamVestingAccount.Merge(m, src)
}
func (m *EventCreateTeamVestingAccount) XXX_Size() int {
	return m.Size()
}
func (m *EventCreateTeamVestingAccount) XXX_DiscardUnknown() {
	xxx_messageInfo_EventCreateTeamVestingAccount.DiscardUnknown(m)
}

var xxx_messageInfo_EventCreateTeamVestingAccount proto.InternalMessageInfo

func (m *EventCreateTeamVestingAccount) GetAuthority() string {
	if m != nil {
		return m.Authority
	}
	return ""
}

func (m *EventCreateTeamVestingAccount) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *EventCreateTeamVestingAccount) GetTotalAllocation() uint64 {
	if m != nil {
		return m.TotalAllocation
	}
	return 0
}

func (m *EventCreateTeamVestingAccount) GetCommencement() uint64 {
	if m != nil {
		return m.Commencement
	}
	return 0
}

// EventClawback is an event emitted when the authority claws back tokens from a team vesting account.
// emitted_by: MsgClawback
type EventClawback struct {
	// authority which initiated this action
	Authority string `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority,omitempty"`
	// id is a unique identify for each vesting account, tied to a single team member.
	Id uint64 `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	// clawback is a unix timestamp of a clawback. If timestamp is zero
	// it means that the account has not received a clawback
	Clawback uint64 `protobuf:"varint,3,opt,name=clawback,proto3" json:"clawback,omitempty"`
	// amount which got clawed back.
	Amount uint64 `protobuf:"varint,4,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (m *EventClawback) Reset()         { *m = EventClawback{} }
func (m *EventClawback) String() string { return proto.CompactTextString(m) }
func (*EventClawback) ProtoMessage()    {}
func (*EventClawback) Descriptor() ([]byte, []int) {
	return fileDescriptor_198acea0777f469a, []int{1}
}
func (m *EventClawback) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EventClawback) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EventClawback.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EventClawback) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventClawback.Merge(m, src)
}
func (m *EventClawback) XXX_Size() int {
	return m.Size()
}
func (m *EventClawback) XXX_DiscardUnknown() {
	xxx_messageInfo_EventClawback.DiscardUnknown(m)
}

var xxx_messageInfo_EventClawback proto.InternalMessageInfo

func (m *EventClawback) GetAuthority() string {
	if m != nil {
		return m.Authority
	}
	return ""
}

func (m *EventClawback) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *EventClawback) GetClawback() uint64 {
	if m != nil {
		return m.Clawback
	}
	return 0
}

func (m *EventClawback) GetAmount() uint64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

// EventClaimedUnlocked is an event emitted when the authority claims unlocked $KYVE for a recipient.
// emitted_by: MsgClaimUnlocked
type EventClaimedUnlocked struct {
	// authority which initiated this action
	Authority string `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority,omitempty"`
	// id is a unique identify for each vesting account, tied to a single team member.
	Id uint64 `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	// amount is the number of tokens claimed from the unlocked amount.
	Amount uint64 `protobuf:"varint,3,opt,name=amount,proto3" json:"amount,omitempty"`
	// recipient is the receiver address of the claim.
	Recipient string `protobuf:"bytes,4,opt,name=recipient,proto3" json:"recipient,omitempty"`
}

func (m *EventClaimedUnlocked) Reset()         { *m = EventClaimedUnlocked{} }
func (m *EventClaimedUnlocked) String() string { return proto.CompactTextString(m) }
func (*EventClaimedUnlocked) ProtoMessage()    {}
func (*EventClaimedUnlocked) Descriptor() ([]byte, []int) {
	return fileDescriptor_198acea0777f469a, []int{2}
}
func (m *EventClaimedUnlocked) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EventClaimedUnlocked) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EventClaimedUnlocked.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EventClaimedUnlocked) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventClaimedUnlocked.Merge(m, src)
}
func (m *EventClaimedUnlocked) XXX_Size() int {
	return m.Size()
}
func (m *EventClaimedUnlocked) XXX_DiscardUnknown() {
	xxx_messageInfo_EventClaimedUnlocked.DiscardUnknown(m)
}

var xxx_messageInfo_EventClaimedUnlocked proto.InternalMessageInfo

func (m *EventClaimedUnlocked) GetAuthority() string {
	if m != nil {
		return m.Authority
	}
	return ""
}

func (m *EventClaimedUnlocked) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *EventClaimedUnlocked) GetAmount() uint64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *EventClaimedUnlocked) GetRecipient() string {
	if m != nil {
		return m.Recipient
	}
	return ""
}

// EventClaimInflationRewards is an event emitted when the authority claims inflation rewards for a recipient.
// emitted_by: MsgClaimInflationRewards
type EventClaimInflationRewards struct {
	// authority which initiated this action
	Authority string `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority,omitempty"`
	// id is a unique identify for each vesting account, tied to a single team member.
	Id uint64 `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	// amount is the amount of inflation rewards the authority should claim for the account holder
	Amount uint64 `protobuf:"varint,3,opt,name=amount,proto3" json:"amount,omitempty"`
	// recipient is the receiver address of the claim.
	Recipient string `protobuf:"bytes,4,opt,name=recipient,proto3" json:"recipient,omitempty"`
}

func (m *EventClaimInflationRewards) Reset()         { *m = EventClaimInflationRewards{} }
func (m *EventClaimInflationRewards) String() string { return proto.CompactTextString(m) }
func (*EventClaimInflationRewards) ProtoMessage()    {}
func (*EventClaimInflationRewards) Descriptor() ([]byte, []int) {
	return fileDescriptor_198acea0777f469a, []int{3}
}
func (m *EventClaimInflationRewards) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EventClaimInflationRewards) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EventClaimInflationRewards.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EventClaimInflationRewards) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventClaimInflationRewards.Merge(m, src)
}
func (m *EventClaimInflationRewards) XXX_Size() int {
	return m.Size()
}
func (m *EventClaimInflationRewards) XXX_DiscardUnknown() {
	xxx_messageInfo_EventClaimInflationRewards.DiscardUnknown(m)
}

var xxx_messageInfo_EventClaimInflationRewards proto.InternalMessageInfo

func (m *EventClaimInflationRewards) GetAuthority() string {
	if m != nil {
		return m.Authority
	}
	return ""
}

func (m *EventClaimInflationRewards) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *EventClaimInflationRewards) GetAmount() uint64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *EventClaimInflationRewards) GetRecipient() string {
	if m != nil {
		return m.Recipient
	}
	return ""
}

// EventClaimAuthorityRewards is an event emitted when the authority claims its inflation rewards for a recipient.
// emitted_by: MsgClaimAuthorityRewards
type EventClaimAuthorityRewards struct {
	// authority which initiated this action
	Authority string `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority,omitempty"`
	// amount is the amount of inflation rewards the authority should claim for the account holder
	Amount uint64 `protobuf:"varint,2,opt,name=amount,proto3" json:"amount,omitempty"`
	// recipient is the receiver address of the claim.
	Recipient string `protobuf:"bytes,3,opt,name=recipient,proto3" json:"recipient,omitempty"`
}

func (m *EventClaimAuthorityRewards) Reset()         { *m = EventClaimAuthorityRewards{} }
func (m *EventClaimAuthorityRewards) String() string { return proto.CompactTextString(m) }
func (*EventClaimAuthorityRewards) ProtoMessage()    {}
func (*EventClaimAuthorityRewards) Descriptor() ([]byte, []int) {
	return fileDescriptor_198acea0777f469a, []int{4}
}
func (m *EventClaimAuthorityRewards) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EventClaimAuthorityRewards) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EventClaimAuthorityRewards.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EventClaimAuthorityRewards) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventClaimAuthorityRewards.Merge(m, src)
}
func (m *EventClaimAuthorityRewards) XXX_Size() int {
	return m.Size()
}
func (m *EventClaimAuthorityRewards) XXX_DiscardUnknown() {
	xxx_messageInfo_EventClaimAuthorityRewards.DiscardUnknown(m)
}

var xxx_messageInfo_EventClaimAuthorityRewards proto.InternalMessageInfo

func (m *EventClaimAuthorityRewards) GetAuthority() string {
	if m != nil {
		return m.Authority
	}
	return ""
}

func (m *EventClaimAuthorityRewards) GetAmount() uint64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *EventClaimAuthorityRewards) GetRecipient() string {
	if m != nil {
		return m.Recipient
	}
	return ""
}

func init() {
	proto.RegisterType((*EventCreateTeamVestingAccount)(nil), "kyve.team.v1beta1.EventCreateTeamVestingAccount")
	proto.RegisterType((*EventClawback)(nil), "kyve.team.v1beta1.EventClawback")
	proto.RegisterType((*EventClaimedUnlocked)(nil), "kyve.team.v1beta1.EventClaimedUnlocked")
	proto.RegisterType((*EventClaimInflationRewards)(nil), "kyve.team.v1beta1.EventClaimInflationRewards")
	proto.RegisterType((*EventClaimAuthorityRewards)(nil), "kyve.team.v1beta1.EventClaimAuthorityRewards")
}

func init() { proto.RegisterFile("kyve/team/v1beta1/events.proto", fileDescriptor_198acea0777f469a) }

var fileDescriptor_198acea0777f469a = []byte{
	// 372 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x92, 0xcf, 0x6a, 0xe2, 0x50,
	0x14, 0xc6, 0xbd, 0x51, 0x64, 0xbc, 0xcc, 0xdf, 0x30, 0x0c, 0x41, 0x66, 0x82, 0x64, 0xa5, 0x9b,
	0x04, 0x99, 0x27, 0x70, 0xc4, 0xc5, 0x30, 0x30, 0x8b, 0x30, 0x23, 0xb4, 0x9b, 0x72, 0x73, 0x73,
	0xaa, 0x97, 0xe4, 0xde, 0x9b, 0x26, 0x47, 0xad, 0x5d, 0xf5, 0x11, 0xfa, 0x00, 0x7d, 0xa0, 0x2e,
	0x5d, 0x76, 0x59, 0xf4, 0x45, 0x4a, 0x62, 0x34, 0xb6, 0x50, 0xa8, 0x9b, 0x2e, 0xef, 0xf7, 0x9d,
	0x7c, 0xbf, 0xef, 0x84, 0x43, 0xed, 0x68, 0x39, 0x07, 0x0f, 0x81, 0x49, 0x6f, 0xde, 0x0f, 0x00,
	0x59, 0xdf, 0x83, 0x39, 0x28, 0xcc, 0xdc, 0x24, 0xd5, 0xa8, 0xcd, 0x2f, 0xb9, 0xef, 0xe6, 0xbe,
	0x5b, 0xfa, 0xce, 0x2d, 0xa1, 0x3f, 0x46, 0xf9, 0xcc, 0x30, 0x05, 0x86, 0xf0, 0x0f, 0x98, 0x1c,
	0x43, 0x86, 0x42, 0x4d, 0x06, 0x9c, 0xeb, 0x99, 0x42, 0xf3, 0x3b, 0x6d, 0xb1, 0x19, 0x4e, 0x75,
	0x2a, 0x70, 0x69, 0x91, 0x0e, 0xe9, 0xb6, 0xfc, 0x4a, 0x30, 0x3f, 0x52, 0x43, 0x84, 0x96, 0xd1,
	0x21, 0xdd, 0x86, 0x6f, 0x88, 0xd0, 0xec, 0xd1, 0xcf, 0xa8, 0x91, 0xc5, 0x67, 0x2c, 0x8e, 0x35,
	0x67, 0x28, 0xb4, 0xb2, 0xea, 0x85, 0xfb, 0xa9, 0xd0, 0x07, 0x7b, 0xd9, 0x74, 0xe8, 0x7b, 0xae,
	0xa5, 0x04, 0xc5, 0x41, 0x82, 0x42, 0xab, 0x51, 0x8c, 0x3d, 0xd1, 0x9c, 0x0b, 0xfa, 0x61, 0xdb,
	0x2e, 0x66, 0x8b, 0x80, 0xf1, 0xe8, 0xc8, 0x36, 0x6d, 0xfa, 0x8e, 0x97, 0x5f, 0x96, 0x2d, 0xf6,
	0x6f, 0xf3, 0x1b, 0x6d, 0x32, 0x99, 0x6f, 0x58, 0x82, 0xcb, 0x97, 0x73, 0x45, 0xbf, 0xee, 0x90,
	0x42, 0x42, 0xf8, 0x5f, 0xc5, 0x9a, 0x47, 0x10, 0x1e, 0x49, 0xae, 0xd2, 0xeb, 0x87, 0xe9, 0x79,
	0x4a, 0x0a, 0x5c, 0x24, 0x62, 0xb7, 0x71, 0xcb, 0xaf, 0x04, 0xe7, 0x9a, 0xd0, 0x76, 0x05, 0xff,
	0xad, 0xce, 0xe3, 0xe2, 0x57, 0xf9, 0xb0, 0x60, 0x69, 0x98, 0xbd, 0x49, 0x85, 0xe4, 0xb0, 0xc1,
	0x60, 0x17, 0xfe, 0xba, 0x06, 0x15, 0xd1, 0x78, 0x99, 0x58, 0x7f, 0x46, 0xfc, 0x35, 0xbc, 0x5b,
	0xdb, 0x64, 0xb5, 0xb6, 0xc9, 0xc3, 0xda, 0x26, 0x37, 0x1b, 0xbb, 0xb6, 0xda, 0xd8, 0xb5, 0xfb,
	0x8d, 0x5d, 0x3b, 0xed, 0x4d, 0x04, 0x4e, 0x67, 0x81, 0xcb, 0xb5, 0xf4, 0xfe, 0x9c, 0x8c, 0x47,
	0x7f, 0x01, 0x17, 0x3a, 0x8d, 0x3c, 0x3e, 0x65, 0x42, 0x79, 0x97, 0xdb, 0x4b, 0xc7, 0x65, 0x02,
	0x59, 0xd0, 0x2c, 0x2e, 0xfc, 0xe7, 0x63, 0x00, 0x00, 0x00, 0xff, 0xff, 0xfe, 0x6f, 0xd8, 0x00,
	0x03, 0x03, 0x00, 0x00,
}

func (m *EventCreateTeamVestingAccount) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EventCreateTeamVestingAccount) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EventCreateTeamVestingAccount) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Commencement != 0 {
		i = encodeVarintEvents(dAtA, i, uint64(m.Commencement))
		i--
		dAtA[i] = 0x20
	}
	if m.TotalAllocation != 0 {
		i = encodeVarintEvents(dAtA, i, uint64(m.TotalAllocation))
		i--
		dAtA[i] = 0x18
	}
	if m.Id != 0 {
		i = encodeVarintEvents(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Authority) > 0 {
		i -= len(m.Authority)
		copy(dAtA[i:], m.Authority)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Authority)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *EventClawback) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EventClawback) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EventClawback) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Amount != 0 {
		i = encodeVarintEvents(dAtA, i, uint64(m.Amount))
		i--
		dAtA[i] = 0x20
	}
	if m.Clawback != 0 {
		i = encodeVarintEvents(dAtA, i, uint64(m.Clawback))
		i--
		dAtA[i] = 0x18
	}
	if m.Id != 0 {
		i = encodeVarintEvents(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Authority) > 0 {
		i -= len(m.Authority)
		copy(dAtA[i:], m.Authority)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Authority)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *EventClaimedUnlocked) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EventClaimedUnlocked) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EventClaimedUnlocked) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Recipient) > 0 {
		i -= len(m.Recipient)
		copy(dAtA[i:], m.Recipient)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Recipient)))
		i--
		dAtA[i] = 0x22
	}
	if m.Amount != 0 {
		i = encodeVarintEvents(dAtA, i, uint64(m.Amount))
		i--
		dAtA[i] = 0x18
	}
	if m.Id != 0 {
		i = encodeVarintEvents(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Authority) > 0 {
		i -= len(m.Authority)
		copy(dAtA[i:], m.Authority)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Authority)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *EventClaimInflationRewards) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EventClaimInflationRewards) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EventClaimInflationRewards) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Recipient) > 0 {
		i -= len(m.Recipient)
		copy(dAtA[i:], m.Recipient)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Recipient)))
		i--
		dAtA[i] = 0x22
	}
	if m.Amount != 0 {
		i = encodeVarintEvents(dAtA, i, uint64(m.Amount))
		i--
		dAtA[i] = 0x18
	}
	if m.Id != 0 {
		i = encodeVarintEvents(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Authority) > 0 {
		i -= len(m.Authority)
		copy(dAtA[i:], m.Authority)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Authority)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *EventClaimAuthorityRewards) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EventClaimAuthorityRewards) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EventClaimAuthorityRewards) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Recipient) > 0 {
		i -= len(m.Recipient)
		copy(dAtA[i:], m.Recipient)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Recipient)))
		i--
		dAtA[i] = 0x1a
	}
	if m.Amount != 0 {
		i = encodeVarintEvents(dAtA, i, uint64(m.Amount))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Authority) > 0 {
		i -= len(m.Authority)
		copy(dAtA[i:], m.Authority)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Authority)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintEvents(dAtA []byte, offset int, v uint64) int {
	offset -= sovEvents(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *EventCreateTeamVestingAccount) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Authority)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	if m.Id != 0 {
		n += 1 + sovEvents(uint64(m.Id))
	}
	if m.TotalAllocation != 0 {
		n += 1 + sovEvents(uint64(m.TotalAllocation))
	}
	if m.Commencement != 0 {
		n += 1 + sovEvents(uint64(m.Commencement))
	}
	return n
}

func (m *EventClawback) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Authority)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	if m.Id != 0 {
		n += 1 + sovEvents(uint64(m.Id))
	}
	if m.Clawback != 0 {
		n += 1 + sovEvents(uint64(m.Clawback))
	}
	if m.Amount != 0 {
		n += 1 + sovEvents(uint64(m.Amount))
	}
	return n
}

func (m *EventClaimedUnlocked) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Authority)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	if m.Id != 0 {
		n += 1 + sovEvents(uint64(m.Id))
	}
	if m.Amount != 0 {
		n += 1 + sovEvents(uint64(m.Amount))
	}
	l = len(m.Recipient)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	return n
}

func (m *EventClaimInflationRewards) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Authority)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	if m.Id != 0 {
		n += 1 + sovEvents(uint64(m.Id))
	}
	if m.Amount != 0 {
		n += 1 + sovEvents(uint64(m.Amount))
	}
	l = len(m.Recipient)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	return n
}

func (m *EventClaimAuthorityRewards) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Authority)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	if m.Amount != 0 {
		n += 1 + sovEvents(uint64(m.Amount))
	}
	l = len(m.Recipient)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	return n
}

func sovEvents(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozEvents(x uint64) (n int) {
	return sovEvents(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *EventCreateTeamVestingAccount) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvents
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
			return fmt.Errorf("proto: EventCreateTeamVestingAccount: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EventCreateTeamVestingAccount: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Authority", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Authority = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalAllocation", wireType)
			}
			m.TotalAllocation = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TotalAllocation |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Commencement", wireType)
			}
			m.Commencement = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Commencement |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipEvents(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvents
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
func (m *EventClawback) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvents
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
			return fmt.Errorf("proto: EventClawback: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EventClawback: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Authority", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Authority = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Clawback", wireType)
			}
			m.Clawback = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Clawback |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			m.Amount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Amount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipEvents(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvents
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
func (m *EventClaimedUnlocked) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvents
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
			return fmt.Errorf("proto: EventClaimedUnlocked: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EventClaimedUnlocked: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Authority", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Authority = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			m.Amount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Amount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Recipient", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Recipient = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEvents(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvents
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
func (m *EventClaimInflationRewards) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvents
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
			return fmt.Errorf("proto: EventClaimInflationRewards: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EventClaimInflationRewards: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Authority", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Authority = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			m.Amount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Amount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Recipient", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Recipient = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEvents(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvents
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
func (m *EventClaimAuthorityRewards) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvents
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
			return fmt.Errorf("proto: EventClaimAuthorityRewards: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EventClaimAuthorityRewards: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Authority", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Authority = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			m.Amount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Amount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Recipient", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Recipient = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEvents(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvents
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
func skipEvents(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowEvents
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
					return 0, ErrIntOverflowEvents
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
					return 0, ErrIntOverflowEvents
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
				return 0, ErrInvalidLengthEvents
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupEvents
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthEvents
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthEvents        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowEvents          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupEvents = fmt.Errorf("proto: unexpected end of group")
)
