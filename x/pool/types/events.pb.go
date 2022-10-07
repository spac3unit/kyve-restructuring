// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: kyve/pool/v1beta1/events.proto

package types

import (
	fmt "fmt"
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

// EventCreatePool ...
type EventCreatePool struct {
	// id ...
	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// name ...
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// runtime ...
	Runtime string `protobuf:"bytes,3,opt,name=runtime,proto3" json:"runtime,omitempty"`
	// logo ...
	Logo string `protobuf:"bytes,4,opt,name=logo,proto3" json:"logo,omitempty"`
	// config ...
	Config string `protobuf:"bytes,5,opt,name=config,proto3" json:"config,omitempty"`
	// start_key ...
	StartKey string `protobuf:"bytes,6,opt,name=start_key,json=startKey,proto3" json:"start_key,omitempty"`
	// upload_interval ...
	UploadInterval uint64 `protobuf:"varint,7,opt,name=upload_interval,json=uploadInterval,proto3" json:"upload_interval,omitempty"`
	// operating_cost ...
	OperatingCost uint64 `protobuf:"varint,8,opt,name=operating_cost,json=operatingCost,proto3" json:"operating_cost,omitempty"`
	// min_stake ...
	MinStake uint64 `protobuf:"varint,9,opt,name=min_stake,json=minStake,proto3" json:"min_stake,omitempty"`
	// max_bundle_size ...
	MaxBundleSize uint64 `protobuf:"varint,10,opt,name=max_bundle_size,json=maxBundleSize,proto3" json:"max_bundle_size,omitempty"`
	// version ...
	Version string `protobuf:"bytes,11,opt,name=version,proto3" json:"version,omitempty"`
	// binaries ...
	Binaries string `protobuf:"bytes,12,opt,name=binaries,proto3" json:"binaries,omitempty"`
}

func (m *EventCreatePool) Reset()         { *m = EventCreatePool{} }
func (m *EventCreatePool) String() string { return proto.CompactTextString(m) }
func (*EventCreatePool) ProtoMessage()    {}
func (*EventCreatePool) Descriptor() ([]byte, []int) {
	return fileDescriptor_c1828a100d789238, []int{0}
}
func (m *EventCreatePool) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EventCreatePool) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EventCreatePool.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EventCreatePool) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventCreatePool.Merge(m, src)
}
func (m *EventCreatePool) XXX_Size() int {
	return m.Size()
}
func (m *EventCreatePool) XXX_DiscardUnknown() {
	xxx_messageInfo_EventCreatePool.DiscardUnknown(m)
}

var xxx_messageInfo_EventCreatePool proto.InternalMessageInfo

func (m *EventCreatePool) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *EventCreatePool) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *EventCreatePool) GetRuntime() string {
	if m != nil {
		return m.Runtime
	}
	return ""
}

func (m *EventCreatePool) GetLogo() string {
	if m != nil {
		return m.Logo
	}
	return ""
}

func (m *EventCreatePool) GetConfig() string {
	if m != nil {
		return m.Config
	}
	return ""
}

func (m *EventCreatePool) GetStartKey() string {
	if m != nil {
		return m.StartKey
	}
	return ""
}

func (m *EventCreatePool) GetUploadInterval() uint64 {
	if m != nil {
		return m.UploadInterval
	}
	return 0
}

func (m *EventCreatePool) GetOperatingCost() uint64 {
	if m != nil {
		return m.OperatingCost
	}
	return 0
}

func (m *EventCreatePool) GetMinStake() uint64 {
	if m != nil {
		return m.MinStake
	}
	return 0
}

func (m *EventCreatePool) GetMaxBundleSize() uint64 {
	if m != nil {
		return m.MaxBundleSize
	}
	return 0
}

func (m *EventCreatePool) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *EventCreatePool) GetBinaries() string {
	if m != nil {
		return m.Binaries
	}
	return ""
}

// EventFundPool is an event emitted when a pool is funded.
type EventFundPool struct {
	// pool_id is the unique ID of the pool.
	PoolId uint64 `protobuf:"varint,1,opt,name=pool_id,json=poolId,proto3" json:"pool_id,omitempty"`
	// address is the account address of the pool funder.
	Address string `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	// amount ...
	Amount uint64 `protobuf:"varint,3,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (m *EventFundPool) Reset()         { *m = EventFundPool{} }
func (m *EventFundPool) String() string { return proto.CompactTextString(m) }
func (*EventFundPool) ProtoMessage()    {}
func (*EventFundPool) Descriptor() ([]byte, []int) {
	return fileDescriptor_c1828a100d789238, []int{1}
}
func (m *EventFundPool) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EventFundPool) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EventFundPool.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EventFundPool) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventFundPool.Merge(m, src)
}
func (m *EventFundPool) XXX_Size() int {
	return m.Size()
}
func (m *EventFundPool) XXX_DiscardUnknown() {
	xxx_messageInfo_EventFundPool.DiscardUnknown(m)
}

var xxx_messageInfo_EventFundPool proto.InternalMessageInfo

func (m *EventFundPool) GetPoolId() uint64 {
	if m != nil {
		return m.PoolId
	}
	return 0
}

func (m *EventFundPool) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *EventFundPool) GetAmount() uint64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

// EventDefundPool is an event emitted when a pool is defunded.
type EventDefundPool struct {
	// pool_id is the unique ID of the pool.
	PoolId uint64 `protobuf:"varint,1,opt,name=pool_id,json=poolId,proto3" json:"pool_id,omitempty"`
	// address is the account address of the pool funder.
	Address string `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	// amount ...
	Amount uint64 `protobuf:"varint,3,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (m *EventDefundPool) Reset()         { *m = EventDefundPool{} }
func (m *EventDefundPool) String() string { return proto.CompactTextString(m) }
func (*EventDefundPool) ProtoMessage()    {}
func (*EventDefundPool) Descriptor() ([]byte, []int) {
	return fileDescriptor_c1828a100d789238, []int{2}
}
func (m *EventDefundPool) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EventDefundPool) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EventDefundPool.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EventDefundPool) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventDefundPool.Merge(m, src)
}
func (m *EventDefundPool) XXX_Size() int {
	return m.Size()
}
func (m *EventDefundPool) XXX_DiscardUnknown() {
	xxx_messageInfo_EventDefundPool.DiscardUnknown(m)
}

var xxx_messageInfo_EventDefundPool proto.InternalMessageInfo

func (m *EventDefundPool) GetPoolId() uint64 {
	if m != nil {
		return m.PoolId
	}
	return 0
}

func (m *EventDefundPool) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *EventDefundPool) GetAmount() uint64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

// EventPoolOutOfFunds is an event emitted when a pool has run out of funds
type EventPoolOutOfFunds struct {
	// pool_id is the unique ID of the pool.
	PoolId uint64 `protobuf:"varint,1,opt,name=pool_id,json=poolId,proto3" json:"pool_id,omitempty"`
}

func (m *EventPoolOutOfFunds) Reset()         { *m = EventPoolOutOfFunds{} }
func (m *EventPoolOutOfFunds) String() string { return proto.CompactTextString(m) }
func (*EventPoolOutOfFunds) ProtoMessage()    {}
func (*EventPoolOutOfFunds) Descriptor() ([]byte, []int) {
	return fileDescriptor_c1828a100d789238, []int{3}
}
func (m *EventPoolOutOfFunds) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EventPoolOutOfFunds) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EventPoolOutOfFunds.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EventPoolOutOfFunds) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventPoolOutOfFunds.Merge(m, src)
}
func (m *EventPoolOutOfFunds) XXX_Size() int {
	return m.Size()
}
func (m *EventPoolOutOfFunds) XXX_DiscardUnknown() {
	xxx_messageInfo_EventPoolOutOfFunds.DiscardUnknown(m)
}

var xxx_messageInfo_EventPoolOutOfFunds proto.InternalMessageInfo

func (m *EventPoolOutOfFunds) GetPoolId() uint64 {
	if m != nil {
		return m.PoolId
	}
	return 0
}

func init() {
	proto.RegisterType((*EventCreatePool)(nil), "kyve.pool.v1beta1.EventCreatePool")
	proto.RegisterType((*EventFundPool)(nil), "kyve.pool.v1beta1.EventFundPool")
	proto.RegisterType((*EventDefundPool)(nil), "kyve.pool.v1beta1.EventDefundPool")
	proto.RegisterType((*EventPoolOutOfFunds)(nil), "kyve.pool.v1beta1.EventPoolOutOfFunds")
}

func init() { proto.RegisterFile("kyve/pool/v1beta1/events.proto", fileDescriptor_c1828a100d789238) }

var fileDescriptor_c1828a100d789238 = []byte{
	// 443 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x52, 0x5d, 0x6b, 0x13, 0x41,
	0x14, 0xcd, 0xc6, 0x98, 0x8f, 0xd1, 0x26, 0x38, 0x82, 0x0e, 0x0a, 0x4b, 0x09, 0xa8, 0xf5, 0x65,
	0x97, 0xe2, 0x3f, 0x68, 0xac, 0x50, 0x0a, 0x56, 0x52, 0x10, 0x2c, 0xc2, 0x32, 0x9b, 0xbd, 0x49,
	0x87, 0xec, 0xce, 0x5d, 0x66, 0xee, 0xae, 0x49, 0x7f, 0x85, 0x3f, 0xcb, 0xc7, 0x3e, 0xfa, 0x28,
	0xc9, 0x1f, 0x91, 0x99, 0xdd, 0xf4, 0xcd, 0x47, 0xdf, 0xe6, 0x9c, 0x7b, 0x38, 0x33, 0xe7, 0xcc,
	0x65, 0xe1, 0x7a, 0x5b, 0x43, 0x5c, 0x22, 0xe6, 0x71, 0x7d, 0x9a, 0x02, 0xc9, 0xd3, 0x18, 0x6a,
	0xd0, 0x64, 0xa3, 0xd2, 0x20, 0x21, 0x7f, 0xe6, 0xe6, 0x91, 0x9b, 0x47, 0xed, 0x7c, 0xba, 0xef,
	0xb2, 0xc9, 0xb9, 0xd3, 0xcc, 0x0c, 0x48, 0x82, 0x2f, 0x88, 0x39, 0x1f, 0xb3, 0xae, 0xca, 0x44,
	0x70, 0x1c, 0x9c, 0xf4, 0xe6, 0x5d, 0x95, 0x71, 0xce, 0x7a, 0x5a, 0x16, 0x20, 0xba, 0xc7, 0xc1,
	0xc9, 0x68, 0xee, 0xcf, 0x5c, 0xb0, 0x81, 0xa9, 0x34, 0xa9, 0x02, 0xc4, 0x23, 0x4f, 0x1f, 0xa0,
	0x53, 0xe7, 0xb8, 0x42, 0xd1, 0x6b, 0xd4, 0xee, 0xcc, 0x5f, 0xb0, 0xfe, 0x02, 0xf5, 0x52, 0xad,
	0xc4, 0x63, 0xcf, 0xb6, 0x88, 0xbf, 0x66, 0x23, 0x4b, 0xd2, 0x50, 0xb2, 0x86, 0xad, 0xe8, 0xfb,
	0xd1, 0xd0, 0x13, 0x97, 0xb0, 0xe5, 0xef, 0xd8, 0xa4, 0x2a, 0x73, 0x94, 0x59, 0xa2, 0x34, 0x81,
	0xa9, 0x65, 0x2e, 0x06, 0xfe, 0x4d, 0xe3, 0x86, 0xbe, 0x68, 0x59, 0xfe, 0x86, 0x8d, 0xb1, 0x04,
	0x23, 0x49, 0xe9, 0x55, 0xb2, 0x40, 0x4b, 0x62, 0xe8, 0x75, 0x47, 0x0f, 0xec, 0x0c, 0x2d, 0xb9,
	0xcb, 0x0a, 0xa5, 0x13, 0x4b, 0x72, 0x0d, 0x62, 0xe4, 0x15, 0xc3, 0x42, 0xe9, 0x6b, 0x87, 0xf9,
	0x5b, 0x36, 0x29, 0xe4, 0x26, 0x49, 0x2b, 0x9d, 0xe5, 0x90, 0x58, 0x75, 0x07, 0x82, 0x35, 0x26,
	0x85, 0xdc, 0x9c, 0x79, 0xf6, 0x5a, 0xdd, 0xf9, 0xdc, 0x35, 0x18, 0xab, 0x50, 0x8b, 0x27, 0x4d,
	0xee, 0x16, 0xf2, 0x57, 0x6c, 0x98, 0x2a, 0x2d, 0x8d, 0x02, 0x2b, 0x9e, 0x36, 0x51, 0x0e, 0x78,
	0x7a, 0xc3, 0x8e, 0x7c, 0xc9, 0x9f, 0x2a, 0x9d, 0xf9, 0x8a, 0x5f, 0xb2, 0x81, 0xfb, 0x86, 0xe4,
	0xa1, 0xe7, 0xbe, 0x83, 0x17, 0x99, 0xf3, 0x97, 0x59, 0x66, 0xc0, 0xda, 0xb6, 0xee, 0x03, 0x74,
	0x1d, 0xca, 0x02, 0x2b, 0x4d, 0xbe, 0xf0, 0xde, 0xbc, 0x45, 0xd3, 0xef, 0xed, 0x07, 0x7e, 0x84,
	0xe5, 0x7f, 0x70, 0x8f, 0xd8, 0x73, 0xef, 0xee, 0x7c, 0xaf, 0x2a, 0xba, 0x5a, 0xba, 0x08, 0xf6,
	0x9f, 0x37, 0x9c, 0xcd, 0x7e, 0xed, 0xc2, 0xe0, 0x7e, 0x17, 0x06, 0x7f, 0x76, 0x61, 0xf0, 0x73,
	0x1f, 0x76, 0xee, 0xf7, 0x61, 0xe7, 0xf7, 0x3e, 0xec, 0xdc, 0xbc, 0x5f, 0x29, 0xba, 0xad, 0xd2,
	0x68, 0x81, 0x45, 0x7c, 0xf9, 0xed, 0xeb, 0xf9, 0x67, 0xa0, 0x1f, 0x68, 0xd6, 0xf1, 0xe2, 0x56,
	0x2a, 0x1d, 0x6f, 0x9a, 0xb5, 0xa5, 0x6d, 0x09, 0x36, 0xed, 0xfb, 0x75, 0xfd, 0xf0, 0x37, 0x00,
	0x00, 0xff, 0xff, 0x07, 0x71, 0xa1, 0x47, 0xd0, 0x02, 0x00, 0x00,
}

func (m *EventCreatePool) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EventCreatePool) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EventCreatePool) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Binaries) > 0 {
		i -= len(m.Binaries)
		copy(dAtA[i:], m.Binaries)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Binaries)))
		i--
		dAtA[i] = 0x62
	}
	if len(m.Version) > 0 {
		i -= len(m.Version)
		copy(dAtA[i:], m.Version)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Version)))
		i--
		dAtA[i] = 0x5a
	}
	if m.MaxBundleSize != 0 {
		i = encodeVarintEvents(dAtA, i, uint64(m.MaxBundleSize))
		i--
		dAtA[i] = 0x50
	}
	if m.MinStake != 0 {
		i = encodeVarintEvents(dAtA, i, uint64(m.MinStake))
		i--
		dAtA[i] = 0x48
	}
	if m.OperatingCost != 0 {
		i = encodeVarintEvents(dAtA, i, uint64(m.OperatingCost))
		i--
		dAtA[i] = 0x40
	}
	if m.UploadInterval != 0 {
		i = encodeVarintEvents(dAtA, i, uint64(m.UploadInterval))
		i--
		dAtA[i] = 0x38
	}
	if len(m.StartKey) > 0 {
		i -= len(m.StartKey)
		copy(dAtA[i:], m.StartKey)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.StartKey)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.Config) > 0 {
		i -= len(m.Config)
		copy(dAtA[i:], m.Config)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Config)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.Logo) > 0 {
		i -= len(m.Logo)
		copy(dAtA[i:], m.Logo)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Logo)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Runtime) > 0 {
		i -= len(m.Runtime)
		copy(dAtA[i:], m.Runtime)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Runtime)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0x12
	}
	if m.Id != 0 {
		i = encodeVarintEvents(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *EventFundPool) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EventFundPool) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EventFundPool) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Amount != 0 {
		i = encodeVarintEvents(dAtA, i, uint64(m.Amount))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0x12
	}
	if m.PoolId != 0 {
		i = encodeVarintEvents(dAtA, i, uint64(m.PoolId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *EventDefundPool) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EventDefundPool) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EventDefundPool) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Amount != 0 {
		i = encodeVarintEvents(dAtA, i, uint64(m.Amount))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0x12
	}
	if m.PoolId != 0 {
		i = encodeVarintEvents(dAtA, i, uint64(m.PoolId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *EventPoolOutOfFunds) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EventPoolOutOfFunds) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EventPoolOutOfFunds) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.PoolId != 0 {
		i = encodeVarintEvents(dAtA, i, uint64(m.PoolId))
		i--
		dAtA[i] = 0x8
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
func (m *EventCreatePool) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovEvents(uint64(m.Id))
	}
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	l = len(m.Runtime)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	l = len(m.Logo)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	l = len(m.Config)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	l = len(m.StartKey)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	if m.UploadInterval != 0 {
		n += 1 + sovEvents(uint64(m.UploadInterval))
	}
	if m.OperatingCost != 0 {
		n += 1 + sovEvents(uint64(m.OperatingCost))
	}
	if m.MinStake != 0 {
		n += 1 + sovEvents(uint64(m.MinStake))
	}
	if m.MaxBundleSize != 0 {
		n += 1 + sovEvents(uint64(m.MaxBundleSize))
	}
	l = len(m.Version)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	l = len(m.Binaries)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	return n
}

func (m *EventFundPool) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.PoolId != 0 {
		n += 1 + sovEvents(uint64(m.PoolId))
	}
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	if m.Amount != 0 {
		n += 1 + sovEvents(uint64(m.Amount))
	}
	return n
}

func (m *EventDefundPool) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.PoolId != 0 {
		n += 1 + sovEvents(uint64(m.PoolId))
	}
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	if m.Amount != 0 {
		n += 1 + sovEvents(uint64(m.Amount))
	}
	return n
}

func (m *EventPoolOutOfFunds) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.PoolId != 0 {
		n += 1 + sovEvents(uint64(m.PoolId))
	}
	return n
}

func sovEvents(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozEvents(x uint64) (n int) {
	return sovEvents(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *EventCreatePool) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: EventCreatePool: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EventCreatePool: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
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
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
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
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Runtime", wireType)
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
			m.Runtime = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Logo", wireType)
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
			m.Logo = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Config", wireType)
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
			m.Config = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StartKey", wireType)
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
			m.StartKey = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field UploadInterval", wireType)
			}
			m.UploadInterval = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.UploadInterval |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field OperatingCost", wireType)
			}
			m.OperatingCost = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.OperatingCost |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 9:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinStake", wireType)
			}
			m.MinStake = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MinStake |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 10:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxBundleSize", wireType)
			}
			m.MaxBundleSize = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxBundleSize |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Version", wireType)
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
			m.Version = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 12:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Binaries", wireType)
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
			m.Binaries = string(dAtA[iNdEx:postIndex])
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
func (m *EventFundPool) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: EventFundPool: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EventFundPool: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolId", wireType)
			}
			m.PoolId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
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
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
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
func (m *EventDefundPool) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: EventDefundPool: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EventDefundPool: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolId", wireType)
			}
			m.PoolId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
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
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
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
func (m *EventPoolOutOfFunds) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: EventPoolOutOfFunds: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EventPoolOutOfFunds: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolId", wireType)
			}
			m.PoolId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
