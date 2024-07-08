// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.13.0
// source: app/grpc/proto/sample.proto

package pb

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// ListRamEventsRequest
type ListRamEventsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// id
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// key
	Key string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	// status
	Status string `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
	// ref
	Ref string `protobuf:"bytes,4,opt,name=ref,proto3" json:"ref,omitempty"`
	// type
	Type int32 `protobuf:"varint,5,opt,name=type,proto3" json:"type,omitempty"`
	// limit
	Limit int32 `protobuf:"varint,101,opt,name=limit,proto3" json:"limit,omitempty"`
	// offset
	Offset int32 `protobuf:"varint,102,opt,name=offset,proto3" json:"offset,omitempty"`
	// order_direction
	OrderDirection string `protobuf:"bytes,103,opt,name=order_direction,json=orderDirection,proto3" json:"order_direction,omitempty"`
	// order_by
	OrderBy string `protobuf:"bytes,104,opt,name=order_by,json=orderBy,proto3" json:"order_by,omitempty"`
	// created_at_lte
	CreatedAtLte int32 `protobuf:"varint,105,opt,name=created_at_lte,json=createdAtLte,proto3" json:"created_at_lte,omitempty"`
	// created_at_gte
	CreatedAtGte int32 `protobuf:"varint,106,opt,name=created_at_gte,json=createdAtGte,proto3" json:"created_at_gte,omitempty"`
}

func (x *ListRamEventsRequest) Reset() {
	*x = ListRamEventsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_grpc_proto_sample_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListRamEventsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRamEventsRequest) ProtoMessage() {}

func (x *ListRamEventsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_app_grpc_proto_sample_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRamEventsRequest.ProtoReflect.Descriptor instead.
func (*ListRamEventsRequest) Descriptor() ([]byte, []int) {
	return file_app_grpc_proto_sample_proto_rawDescGZIP(), []int{0}
}

func (x *ListRamEventsRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *ListRamEventsRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *ListRamEventsRequest) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *ListRamEventsRequest) GetRef() string {
	if x != nil {
		return x.Ref
	}
	return ""
}

func (x *ListRamEventsRequest) GetType() int32 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *ListRamEventsRequest) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *ListRamEventsRequest) GetOffset() int32 {
	if x != nil {
		return x.Offset
	}
	return 0
}

func (x *ListRamEventsRequest) GetOrderDirection() string {
	if x != nil {
		return x.OrderDirection
	}
	return ""
}

func (x *ListRamEventsRequest) GetOrderBy() string {
	if x != nil {
		return x.OrderBy
	}
	return ""
}

func (x *ListRamEventsRequest) GetCreatedAtLte() int32 {
	if x != nil {
		return x.CreatedAtLte
	}
	return 0
}

func (x *ListRamEventsRequest) GetCreatedAtGte() int32 {
	if x != nil {
		return x.CreatedAtGte
	}
	return 0
}

// ListRamEventsResponse
type ListRamEventsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// events
	Events []*RamEvent `protobuf:"bytes,1,rep,name=events,proto3" json:"events,omitempty"`
}

func (x *ListRamEventsResponse) Reset() {
	*x = ListRamEventsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_grpc_proto_sample_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListRamEventsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRamEventsResponse) ProtoMessage() {}

func (x *ListRamEventsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_app_grpc_proto_sample_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRamEventsResponse.ProtoReflect.Descriptor instead.
func (*ListRamEventsResponse) Descriptor() ([]byte, []int) {
	return file_app_grpc_proto_sample_proto_rawDescGZIP(), []int{1}
}

func (x *ListRamEventsResponse) GetEvents() []*RamEvent {
	if x != nil {
		return x.Events
	}
	return nil
}

// UpsertRamEventRequest
type UpsertRamEventRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// id
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// type
	Type int32 `protobuf:"varint,2,opt,name=type,proto3" json:"type,omitempty"`
	// payload (json string)
	Payload string `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
	// ref
	Ref string `protobuf:"bytes,4,opt,name=ref,proto3" json:"ref,omitempty"`
	// key
	Key string `protobuf:"bytes,5,opt,name=key,proto3" json:"key,omitempty"`
	// status
	Status string `protobuf:"bytes,6,opt,name=status,proto3" json:"status,omitempty"`
	// retry_count
	RetryCount int32 `protobuf:"varint,7,opt,name=retry_count,json=retryCount,proto3" json:"retry_count,omitempty"`
}

func (x *UpsertRamEventRequest) Reset() {
	*x = UpsertRamEventRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_grpc_proto_sample_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpsertRamEventRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpsertRamEventRequest) ProtoMessage() {}

func (x *UpsertRamEventRequest) ProtoReflect() protoreflect.Message {
	mi := &file_app_grpc_proto_sample_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpsertRamEventRequest.ProtoReflect.Descriptor instead.
func (*UpsertRamEventRequest) Descriptor() ([]byte, []int) {
	return file_app_grpc_proto_sample_proto_rawDescGZIP(), []int{2}
}

func (x *UpsertRamEventRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UpsertRamEventRequest) GetType() int32 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *UpsertRamEventRequest) GetPayload() string {
	if x != nil {
		return x.Payload
	}
	return ""
}

func (x *UpsertRamEventRequest) GetRef() string {
	if x != nil {
		return x.Ref
	}
	return ""
}

func (x *UpsertRamEventRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *UpsertRamEventRequest) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *UpsertRamEventRequest) GetRetryCount() int32 {
	if x != nil {
		return x.RetryCount
	}
	return 0
}

// UpsertRamEventResponse
type UpsertRamEventResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// message
	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *UpsertRamEventResponse) Reset() {
	*x = UpsertRamEventResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_grpc_proto_sample_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpsertRamEventResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpsertRamEventResponse) ProtoMessage() {}

func (x *UpsertRamEventResponse) ProtoReflect() protoreflect.Message {
	mi := &file_app_grpc_proto_sample_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpsertRamEventResponse.ProtoReflect.Descriptor instead.
func (*UpsertRamEventResponse) Descriptor() ([]byte, []int) {
	return file_app_grpc_proto_sample_proto_rawDescGZIP(), []int{3}
}

func (x *UpsertRamEventResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type RamEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// id
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// type
	Type int32 `protobuf:"varint,2,opt,name=type,proto3" json:"type,omitempty"`
	// type_name
	TypeName string `protobuf:"bytes,3,opt,name=type_name,json=typeName,proto3" json:"type_name,omitempty"`
	// payload (json string)
	Payload string `protobuf:"bytes,4,opt,name=payload,proto3" json:"payload,omitempty"`
	// ref
	Ref string `protobuf:"bytes,5,opt,name=ref,proto3" json:"ref,omitempty"`
	// key
	Key string `protobuf:"bytes,6,opt,name=key,proto3" json:"key,omitempty"`
	// status
	Status string `protobuf:"bytes,7,opt,name=status,proto3" json:"status,omitempty"`
	// retry_count
	RetryCount int32 `protobuf:"varint,8,opt,name=retry_count,json=retryCount,proto3" json:"retry_count,omitempty"`
	// created_at
	CreatedAt int64 `protobuf:"varint,9,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// updated_at
	UpdatedAt int64 `protobuf:"varint,10,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	// actions
	Actions []*RamEventAction `protobuf:"bytes,11,rep,name=actions,proto3" json:"actions,omitempty"`
}

func (x *RamEvent) Reset() {
	*x = RamEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_grpc_proto_sample_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RamEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RamEvent) ProtoMessage() {}

func (x *RamEvent) ProtoReflect() protoreflect.Message {
	mi := &file_app_grpc_proto_sample_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RamEvent.ProtoReflect.Descriptor instead.
func (*RamEvent) Descriptor() ([]byte, []int) {
	return file_app_grpc_proto_sample_proto_rawDescGZIP(), []int{4}
}

func (x *RamEvent) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *RamEvent) GetType() int32 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *RamEvent) GetTypeName() string {
	if x != nil {
		return x.TypeName
	}
	return ""
}

func (x *RamEvent) GetPayload() string {
	if x != nil {
		return x.Payload
	}
	return ""
}

func (x *RamEvent) GetRef() string {
	if x != nil {
		return x.Ref
	}
	return ""
}

func (x *RamEvent) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *RamEvent) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *RamEvent) GetRetryCount() int32 {
	if x != nil {
		return x.RetryCount
	}
	return 0
}

func (x *RamEvent) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *RamEvent) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

func (x *RamEvent) GetActions() []*RamEventAction {
	if x != nil {
		return x.Actions
	}
	return nil
}

type RamEventAction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// id
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// event_id
	EventId int64 `protobuf:"varint,2,opt,name=event_id,json=eventId,proto3" json:"event_id,omitempty"`
	// retry_id
	RetryId int32 `protobuf:"varint,3,opt,name=retry_id,json=retryId,proto3" json:"retry_id,omitempty"`
	// status
	Status string `protobuf:"bytes,4,opt,name=status,proto3" json:"status,omitempty"`
	// error
	Error string `protobuf:"bytes,5,opt,name=error,proto3" json:"error,omitempty"`
	// created_at
	CreatedAt int64 `protobuf:"varint,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// updated_at
	UpdatedAt int64 `protobuf:"varint,7,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *RamEventAction) Reset() {
	*x = RamEventAction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_grpc_proto_sample_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RamEventAction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RamEventAction) ProtoMessage() {}

func (x *RamEventAction) ProtoReflect() protoreflect.Message {
	mi := &file_app_grpc_proto_sample_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RamEventAction.ProtoReflect.Descriptor instead.
func (*RamEventAction) Descriptor() ([]byte, []int) {
	return file_app_grpc_proto_sample_proto_rawDescGZIP(), []int{5}
}

func (x *RamEventAction) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *RamEventAction) GetEventId() int64 {
	if x != nil {
		return x.EventId
	}
	return 0
}

func (x *RamEventAction) GetRetryId() int32 {
	if x != nil {
		return x.RetryId
	}
	return 0
}

func (x *RamEventAction) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *RamEventAction) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

func (x *RamEventAction) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *RamEventAction) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

var File_app_grpc_proto_sample_proto protoreflect.FileDescriptor

var file_app_grpc_proto_sample_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x61, 0x70, 0x70, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x73, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x67,
	0x6f, 0x69, 0x74, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x62, 0x1a,
	0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x62,
	0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xee, 0x02, 0x0a, 0x14, 0x4c,
	0x69, 0x73, 0x74, 0x52, 0x61, 0x6d, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x1f, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x72, 0x65, 0x66, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x72, 0x65, 0x66, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x65, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x6c, 0x69, 0x6d,
	0x69, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x66, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x12, 0x39, 0x0a, 0x0f, 0x6f, 0x72,
	0x64, 0x65, 0x72, 0x5f, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x67, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x10, 0xba, 0x48, 0x0d, 0x72, 0x0b, 0x52, 0x03, 0x41, 0x53, 0x43, 0x52,
	0x04, 0x44, 0x45, 0x53, 0x43, 0x52, 0x0e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x44, 0x69, 0x72, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x38, 0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x62,
	0x79, 0x18, 0x68, 0x20, 0x01, 0x28, 0x09, 0x42, 0x1d, 0xba, 0x48, 0x1a, 0x72, 0x18, 0x52, 0x0a,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x52, 0x0a, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x42, 0x79, 0x12,
	0x24, 0x0a, 0x0e, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x5f, 0x6c, 0x74,
	0x65, 0x18, 0x69, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64,
	0x41, 0x74, 0x4c, 0x74, 0x65, 0x12, 0x24, 0x0a, 0x0e, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64,
	0x5f, 0x61, 0x74, 0x5f, 0x67, 0x74, 0x65, 0x18, 0x6a, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x47, 0x74, 0x65, 0x22, 0x4b, 0x0a, 0x15, 0x4c,
	0x69, 0x73, 0x74, 0x52, 0x61, 0x6d, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x32, 0x0a, 0x06, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x69, 0x74, 0x2e, 0x61, 0x70, 0x70, 0x2e,
	0x67, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x62, 0x2e, 0x52, 0x61, 0x6d, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x52, 0x06, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x22, 0xe1, 0x01, 0x0a, 0x15, 0x55, 0x70, 0x73,
	0x65, 0x72, 0x74, 0x52, 0x61, 0x6d, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64,
	0x12, 0x10, 0x0a, 0x03, 0x72, 0x65, 0x66, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x72,
	0x65, 0x66, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x45, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x2d, 0xba, 0x48, 0x2a, 0x72, 0x28, 0x52, 0x07, 0x43, 0x52, 0x45,
	0x41, 0x54, 0x45, 0x44, 0x52, 0x08, 0x48, 0x41, 0x4e, 0x44, 0x4c, 0x49, 0x4e, 0x47, 0x52, 0x04,
	0x44, 0x4f, 0x4e, 0x45, 0x52, 0x0d, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x5f, 0x48, 0x41, 0x4e,
	0x44, 0x4c, 0x45, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x72,
	0x65, 0x74, 0x72, 0x79, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0a, 0x72, 0x65, 0x74, 0x72, 0x79, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x32, 0x0a, 0x16,
	0x55, 0x70, 0x73, 0x65, 0x72, 0x74, 0x52, 0x61, 0x6d, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x22, 0xbc, 0x02, 0x0a, 0x08, 0x52, 0x61, 0x6d, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x79, 0x70, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x79, 0x70, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x72, 0x65, 0x66, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x72, 0x65, 0x66, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x16, 0x0a, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x74, 0x72, 0x79, 0x5f, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x72, 0x65, 0x74, 0x72, 0x79,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64,
	0x5f, 0x61, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f,
	0x61, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x64, 0x41, 0x74, 0x12, 0x3a, 0x0a, 0x07, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x0b,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x67, 0x6f, 0x69, 0x74, 0x2e, 0x61, 0x70, 0x70, 0x2e,
	0x67, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x62, 0x2e, 0x52, 0x61, 0x6d, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22,
	0xc2, 0x01, 0x0a, 0x0e, 0x52, 0x61, 0x6d, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x41, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x19, 0x0a,
	0x08, 0x72, 0x65, 0x74, 0x72, 0x79, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x07, 0x72, 0x65, 0x74, 0x72, 0x79, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x5f, 0x61, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64,
	0x5f, 0x61, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x64, 0x41, 0x74, 0x32, 0xa3, 0x02, 0x0a, 0x0c, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x82, 0x01, 0x0a, 0x0d, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x61,
	0x6d, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x26, 0x2e, 0x67, 0x6f, 0x69, 0x74, 0x2e, 0x61,
	0x70, 0x70, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x62, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52,
	0x61, 0x6d, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x27, 0x2e, 0x67, 0x6f, 0x69, 0x74, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e,
	0x70, 0x62, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x61, 0x6d, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x20, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1a,
	0x12, 0x18, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f,
	0x72, 0x61, 0x6d, 0x2d, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x8d, 0x01, 0x0a, 0x0e, 0x55,
	0x70, 0x73, 0x65, 0x72, 0x74, 0x52, 0x61, 0x6d, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x27, 0x2e,
	0x67, 0x6f, 0x69, 0x74, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x62,
	0x2e, 0x55, 0x70, 0x73, 0x65, 0x72, 0x74, 0x52, 0x61, 0x6d, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x28, 0x2e, 0x67, 0x6f, 0x69, 0x74, 0x2e, 0x61, 0x70,
	0x70, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x62, 0x2e, 0x55, 0x70, 0x73, 0x65, 0x72, 0x74,
	0x52, 0x61, 0x6d, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x28, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x22, 0x3a, 0x01, 0x2a, 0x22, 0x1d, 0x2f, 0x61, 0x64,
	0x6d, 0x69, 0x6e, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x61, 0x6d, 0x2d, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x42, 0x28, 0x5a, 0x26, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x6f, 0x6b, 0x79, 0x61, 0x6e, 0x31,
	0x39, 0x38, 0x39, 0x2f, 0x67, 0x6f, 0x69, 0x74, 0x2f, 0x61, 0x70, 0x70, 0x2f, 0x67, 0x72, 0x70,
	0x63, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_app_grpc_proto_sample_proto_rawDescOnce sync.Once
	file_app_grpc_proto_sample_proto_rawDescData = file_app_grpc_proto_sample_proto_rawDesc
)

func file_app_grpc_proto_sample_proto_rawDescGZIP() []byte {
	file_app_grpc_proto_sample_proto_rawDescOnce.Do(func() {
		file_app_grpc_proto_sample_proto_rawDescData = protoimpl.X.CompressGZIP(file_app_grpc_proto_sample_proto_rawDescData)
	})
	return file_app_grpc_proto_sample_proto_rawDescData
}

var file_app_grpc_proto_sample_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_app_grpc_proto_sample_proto_goTypes = []any{
	(*ListRamEventsRequest)(nil),   // 0: goit.app.grpc.pb.ListRamEventsRequest
	(*ListRamEventsResponse)(nil),  // 1: goit.app.grpc.pb.ListRamEventsResponse
	(*UpsertRamEventRequest)(nil),  // 2: goit.app.grpc.pb.UpsertRamEventRequest
	(*UpsertRamEventResponse)(nil), // 3: goit.app.grpc.pb.UpsertRamEventResponse
	(*RamEvent)(nil),               // 4: goit.app.grpc.pb.RamEvent
	(*RamEventAction)(nil),         // 5: goit.app.grpc.pb.RamEventAction
}
var file_app_grpc_proto_sample_proto_depIdxs = []int32{
	4, // 0: goit.app.grpc.pb.ListRamEventsResponse.events:type_name -> goit.app.grpc.pb.RamEvent
	5, // 1: goit.app.grpc.pb.RamEvent.actions:type_name -> goit.app.grpc.pb.RamEventAction
	0, // 2: goit.app.grpc.pb.AdminService.ListRamEvents:input_type -> goit.app.grpc.pb.ListRamEventsRequest
	2, // 3: goit.app.grpc.pb.AdminService.UpsertRamEvent:input_type -> goit.app.grpc.pb.UpsertRamEventRequest
	1, // 4: goit.app.grpc.pb.AdminService.ListRamEvents:output_type -> goit.app.grpc.pb.ListRamEventsResponse
	3, // 5: goit.app.grpc.pb.AdminService.UpsertRamEvent:output_type -> goit.app.grpc.pb.UpsertRamEventResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_app_grpc_proto_sample_proto_init() }
func file_app_grpc_proto_sample_proto_init() {
	if File_app_grpc_proto_sample_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_app_grpc_proto_sample_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*ListRamEventsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_app_grpc_proto_sample_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*ListRamEventsResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_app_grpc_proto_sample_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*UpsertRamEventRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_app_grpc_proto_sample_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*UpsertRamEventResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_app_grpc_proto_sample_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*RamEvent); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_app_grpc_proto_sample_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*RamEventAction); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_app_grpc_proto_sample_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_app_grpc_proto_sample_proto_goTypes,
		DependencyIndexes: file_app_grpc_proto_sample_proto_depIdxs,
		MessageInfos:      file_app_grpc_proto_sample_proto_msgTypes,
	}.Build()
	File_app_grpc_proto_sample_proto = out.File
	file_app_grpc_proto_sample_proto_rawDesc = nil
	file_app_grpc_proto_sample_proto_goTypes = nil
	file_app_grpc_proto_sample_proto_depIdxs = nil
}
