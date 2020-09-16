// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.10.0
// source: feeds_service.proto

package feeds

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type GetFeedRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetFeedRequest) Reset() {
	*x = GetFeedRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_feeds_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFeedRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFeedRequest) ProtoMessage() {}

func (x *GetFeedRequest) ProtoReflect() protoreflect.Message {
	mi := &file_feeds_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFeedRequest.ProtoReflect.Descriptor instead.
func (*GetFeedRequest) Descriptor() ([]byte, []int) {
	return file_feeds_service_proto_rawDescGZIP(), []int{0}
}

func (x *GetFeedRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type ListFeedsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The maximum number of feeds to return.
	PageSize int64 `protobuf:"varint,1,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	// The next_page_token value returned from a previous List request, if any.
	PageToken string `protobuf:"bytes,2,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
}

func (x *ListFeedsRequest) Reset() {
	*x = ListFeedsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_feeds_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListFeedsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListFeedsRequest) ProtoMessage() {}

func (x *ListFeedsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_feeds_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListFeedsRequest.ProtoReflect.Descriptor instead.
func (*ListFeedsRequest) Descriptor() ([]byte, []int) {
	return file_feeds_service_proto_rawDescGZIP(), []int{1}
}

func (x *ListFeedsRequest) GetPageSize() int64 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListFeedsRequest) GetPageToken() string {
	if x != nil {
		return x.PageToken
	}
	return ""
}

type ListFeedsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// There will be a maximum number of feeds returned based on the page_size field int the request.
	Feeds []*Feed `protobuf:"bytes,1,rep,name=feeds,proto3" json:"feeds,omitempty"`
	// Token to retrieve the next page of results, or empty if there are no more results in the list.
	NextPageToken string `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
}

func (x *ListFeedsResponse) Reset() {
	*x = ListFeedsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_feeds_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListFeedsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListFeedsResponse) ProtoMessage() {}

func (x *ListFeedsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_feeds_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListFeedsResponse.ProtoReflect.Descriptor instead.
func (*ListFeedsResponse) Descriptor() ([]byte, []int) {
	return file_feeds_service_proto_rawDescGZIP(), []int{2}
}

func (x *ListFeedsResponse) GetFeeds() []*Feed {
	if x != nil {
		return x.Feeds
	}
	return nil
}

func (x *ListFeedsResponse) GetNextPageToken() string {
	if x != nil {
		return x.NextPageToken
	}
	return ""
}

type Feed struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Id of the resource.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Id of user who owns this feed.
	UserId string `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	// Time of the last feed fetch by the reader.
	LastFetched *timestamp.Timestamp `protobuf:"bytes,3,opt,name=last_fetched,json=lastFetched,proto3" json:"last_fetched,omitempty"`
	// Url is the feed's url.
	Url string `protobuf:"bytes,4,opt,name=url,proto3" json:"url,omitempty"`
	// Title is the feed's title.
	Title string `protobuf:"bytes,5,opt,name=title,proto3" json:"title,omitempty"`
}

func (x *Feed) Reset() {
	*x = Feed{}
	if protoimpl.UnsafeEnabled {
		mi := &file_feeds_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Feed) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Feed) ProtoMessage() {}

func (x *Feed) ProtoReflect() protoreflect.Message {
	mi := &file_feeds_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Feed.ProtoReflect.Descriptor instead.
func (*Feed) Descriptor() ([]byte, []int) {
	return file_feeds_service_proto_rawDescGZIP(), []int{3}
}

func (x *Feed) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Feed) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *Feed) GetLastFetched() *timestamp.Timestamp {
	if x != nil {
		return x.LastFetched
	}
	return nil
}

func (x *Feed) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *Feed) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

var File_feeds_service_proto protoreflect.FileDescriptor

var file_feeds_service_proto_rawDesc = []byte{
	0x0a, 0x13, 0x66, 0x65, 0x65, 0x64, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1c, 0x61, 0x70, 0x70, 0x2e, 0x6d, 0x69, 0x6e, 0x69, 0x62,
	0x6f, 0x61, 0x72, 0x64, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x66, 0x65, 0x65, 0x64, 0x73,
	0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x20, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x46, 0x65, 0x65, 0x64, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x22, 0x4e, 0x0a, 0x10, 0x4c, 0x69, 0x73, 0x74, 0x46, 0x65, 0x65, 0x64,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65,
	0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x70, 0x61, 0x67,
	0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f,
	0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x61, 0x67, 0x65, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x75, 0x0a, 0x11, 0x4c, 0x69, 0x73, 0x74, 0x46, 0x65, 0x65, 0x64,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x38, 0x0a, 0x05, 0x66, 0x65, 0x65,
	0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x6d,
	0x69, 0x6e, 0x69, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x66,
	0x65, 0x65, 0x64, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x65, 0x65, 0x64, 0x52, 0x05, 0x66, 0x65,
	0x65, 0x64, 0x73, 0x12, 0x26, 0x0a, 0x0f, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x70, 0x61, 0x67, 0x65,
	0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6e, 0x65,
	0x78, 0x74, 0x50, 0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x96, 0x01, 0x0a, 0x04,
	0x46, 0x65, 0x65, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x3d, 0x0a,
	0x0c, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x66, 0x65, 0x74, 0x63, 0x68, 0x65, 0x64, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x0b, 0x6c, 0x61, 0x73, 0x74, 0x46, 0x65, 0x74, 0x63, 0x68, 0x65, 0x64, 0x12, 0x10, 0x0a, 0x03,
	0x75, 0x72, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x14,
	0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74,
	0x69, 0x74, 0x6c, 0x65, 0x32, 0x8f, 0x02, 0x0a, 0x0c, 0x46, 0x65, 0x65, 0x64, 0x73, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x79, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x46, 0x65, 0x65, 0x64,
	0x12, 0x2c, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x6d, 0x69, 0x6e, 0x69, 0x62, 0x6f, 0x61, 0x72, 0x64,
	0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x66, 0x65, 0x65, 0x64, 0x73, 0x2e, 0x76, 0x31, 0x2e,
	0x47, 0x65, 0x74, 0x46, 0x65, 0x65, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22,
	0x2e, 0x61, 0x70, 0x70, 0x2e, 0x6d, 0x69, 0x6e, 0x69, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x75,
	0x73, 0x65, 0x72, 0x73, 0x2e, 0x66, 0x65, 0x65, 0x64, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x65,
	0x65, 0x64, 0x22, 0x1c, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x16, 0x12, 0x14, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x76, 0x31, 0x2f, 0x66, 0x65, 0x65, 0x64, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x3d, 0x2a, 0x7d,
	0x12, 0x83, 0x01, 0x0a, 0x09, 0x4c, 0x69, 0x73, 0x74, 0x46, 0x65, 0x65, 0x64, 0x73, 0x12, 0x2e,
	0x2e, 0x61, 0x70, 0x70, 0x2e, 0x6d, 0x69, 0x6e, 0x69, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x75,
	0x73, 0x65, 0x72, 0x73, 0x2e, 0x66, 0x65, 0x65, 0x64, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69,
	0x73, 0x74, 0x46, 0x65, 0x65, 0x64, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2f,
	0x2e, 0x61, 0x70, 0x70, 0x2e, 0x6d, 0x69, 0x6e, 0x69, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x75,
	0x73, 0x65, 0x72, 0x73, 0x2e, 0x66, 0x65, 0x65, 0x64, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69,
	0x73, 0x74, 0x46, 0x65, 0x65, 0x64, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x15, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0f, 0x12, 0x0d, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31,
	0x2f, 0x66, 0x65, 0x65, 0x64, 0x73, 0x42, 0x07, 0x5a, 0x05, 0x66, 0x65, 0x65, 0x64, 0x73, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_feeds_service_proto_rawDescOnce sync.Once
	file_feeds_service_proto_rawDescData = file_feeds_service_proto_rawDesc
)

func file_feeds_service_proto_rawDescGZIP() []byte {
	file_feeds_service_proto_rawDescOnce.Do(func() {
		file_feeds_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_feeds_service_proto_rawDescData)
	})
	return file_feeds_service_proto_rawDescData
}

var file_feeds_service_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_feeds_service_proto_goTypes = []interface{}{
	(*GetFeedRequest)(nil),      // 0: app.miniboard.users.feeds.v1.GetFeedRequest
	(*ListFeedsRequest)(nil),    // 1: app.miniboard.users.feeds.v1.ListFeedsRequest
	(*ListFeedsResponse)(nil),   // 2: app.miniboard.users.feeds.v1.ListFeedsResponse
	(*Feed)(nil),                // 3: app.miniboard.users.feeds.v1.Feed
	(*timestamp.Timestamp)(nil), // 4: google.protobuf.Timestamp
}
var file_feeds_service_proto_depIdxs = []int32{
	3, // 0: app.miniboard.users.feeds.v1.ListFeedsResponse.feeds:type_name -> app.miniboard.users.feeds.v1.Feed
	4, // 1: app.miniboard.users.feeds.v1.Feed.last_fetched:type_name -> google.protobuf.Timestamp
	0, // 2: app.miniboard.users.feeds.v1.FeedsService.GetFeed:input_type -> app.miniboard.users.feeds.v1.GetFeedRequest
	1, // 3: app.miniboard.users.feeds.v1.FeedsService.ListFeeds:input_type -> app.miniboard.users.feeds.v1.ListFeedsRequest
	3, // 4: app.miniboard.users.feeds.v1.FeedsService.GetFeed:output_type -> app.miniboard.users.feeds.v1.Feed
	2, // 5: app.miniboard.users.feeds.v1.FeedsService.ListFeeds:output_type -> app.miniboard.users.feeds.v1.ListFeedsResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_feeds_service_proto_init() }
func file_feeds_service_proto_init() {
	if File_feeds_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_feeds_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFeedRequest); i {
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
		file_feeds_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListFeedsRequest); i {
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
		file_feeds_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListFeedsResponse); i {
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
		file_feeds_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Feed); i {
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
			RawDescriptor: file_feeds_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_feeds_service_proto_goTypes,
		DependencyIndexes: file_feeds_service_proto_depIdxs,
		MessageInfos:      file_feeds_service_proto_msgTypes,
	}.Build()
	File_feeds_service_proto = out.File
	file_feeds_service_proto_rawDesc = nil
	file_feeds_service_proto_goTypes = nil
	file_feeds_service_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// FeedsServiceClient is the client API for FeedsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type FeedsServiceClient interface {
	GetFeed(ctx context.Context, in *GetFeedRequest, opts ...grpc.CallOption) (*Feed, error)
	ListFeeds(ctx context.Context, in *ListFeedsRequest, opts ...grpc.CallOption) (*ListFeedsResponse, error)
}

type feedsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFeedsServiceClient(cc grpc.ClientConnInterface) FeedsServiceClient {
	return &feedsServiceClient{cc}
}

func (c *feedsServiceClient) GetFeed(ctx context.Context, in *GetFeedRequest, opts ...grpc.CallOption) (*Feed, error) {
	out := new(Feed)
	err := c.cc.Invoke(ctx, "/app.miniboard.users.feeds.v1.FeedsService/GetFeed", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *feedsServiceClient) ListFeeds(ctx context.Context, in *ListFeedsRequest, opts ...grpc.CallOption) (*ListFeedsResponse, error) {
	out := new(ListFeedsResponse)
	err := c.cc.Invoke(ctx, "/app.miniboard.users.feeds.v1.FeedsService/ListFeeds", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FeedsServiceServer is the server API for FeedsService service.
type FeedsServiceServer interface {
	GetFeed(context.Context, *GetFeedRequest) (*Feed, error)
	ListFeeds(context.Context, *ListFeedsRequest) (*ListFeedsResponse, error)
}

// UnimplementedFeedsServiceServer can be embedded to have forward compatible implementations.
type UnimplementedFeedsServiceServer struct {
}

func (*UnimplementedFeedsServiceServer) GetFeed(context.Context, *GetFeedRequest) (*Feed, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFeed not implemented")
}
func (*UnimplementedFeedsServiceServer) ListFeeds(context.Context, *ListFeedsRequest) (*ListFeedsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListFeeds not implemented")
}

func RegisterFeedsServiceServer(s *grpc.Server, srv FeedsServiceServer) {
	s.RegisterService(&_FeedsService_serviceDesc, srv)
}

func _FeedsService_GetFeed_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFeedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FeedsServiceServer).GetFeed(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.miniboard.users.feeds.v1.FeedsService/GetFeed",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FeedsServiceServer).GetFeed(ctx, req.(*GetFeedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FeedsService_ListFeeds_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListFeedsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FeedsServiceServer).ListFeeds(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.miniboard.users.feeds.v1.FeedsService/ListFeeds",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FeedsServiceServer).ListFeeds(ctx, req.(*ListFeedsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _FeedsService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "app.miniboard.users.feeds.v1.FeedsService",
	HandlerType: (*FeedsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetFeed",
			Handler:    _FeedsService_GetFeed_Handler,
		},
		{
			MethodName: "ListFeeds",
			Handler:    _FeedsService_ListFeeds_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "feeds_service.proto",
}
