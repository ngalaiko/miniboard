// Code generated by protoc-gen-go. DO NOT EDIT.
// source: articles_service.proto

package articles

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	field_mask "google.golang.org/genproto/protobuf/field_mask"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ArticleView int32

const (
	// Not specified, equivalent to BASIC.
	ArticleView_ARTICLE_VIEW_UNSPECIFIED ArticleView = 0
	// Server response does not include content field.
	// Default value.
	ArticleView_ARTICLE_VIEW_BASIC ArticleView = 1
	// Service response includes article's content.
	ArticleView_ARTICLE_VIEW_FULL ArticleView = 2
)

var ArticleView_name = map[int32]string{
	0: "ARTICLE_VIEW_UNSPECIFIED",
	1: "ARTICLE_VIEW_BASIC",
	2: "ARTICLE_VIEW_FULL",
}

var ArticleView_value = map[string]int32{
	"ARTICLE_VIEW_UNSPECIFIED": 0,
	"ARTICLE_VIEW_BASIC":       1,
	"ARTICLE_VIEW_FULL":        2,
}

func (x ArticleView) String() string {
	return proto.EnumName(ArticleView_name, int32(x))
}

func (ArticleView) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_75f616a2c230f591, []int{0}
}

type ListArticlesRequest struct {
	// The maximum number of articles to return.
	PageSize int64 `protobuf:"varint,1,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	// The next_page_token value returned from a previous List request, if any.
	PageToken string `protobuf:"bytes,2,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
	// Filter by is_read field.
	IsRead *wrappers.BoolValue `protobuf:"bytes,3,opt,name=is_read,json=isRead,proto3" json:"is_read,omitempty"`
	// Filter by is_favorite field.
	IsFavorite *wrappers.BoolValue `protobuf:"bytes,4,opt,name=is_favorite,json=isFavorite,proto3" json:"is_favorite,omitempty"`
	// Filter by title field.
	Title                *wrappers.StringValue `protobuf:"bytes,5,opt,name=title,proto3" json:"title,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *ListArticlesRequest) Reset()         { *m = ListArticlesRequest{} }
func (m *ListArticlesRequest) String() string { return proto.CompactTextString(m) }
func (*ListArticlesRequest) ProtoMessage()    {}
func (*ListArticlesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_75f616a2c230f591, []int{0}
}

func (m *ListArticlesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListArticlesRequest.Unmarshal(m, b)
}
func (m *ListArticlesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListArticlesRequest.Marshal(b, m, deterministic)
}
func (m *ListArticlesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListArticlesRequest.Merge(m, src)
}
func (m *ListArticlesRequest) XXX_Size() int {
	return xxx_messageInfo_ListArticlesRequest.Size(m)
}
func (m *ListArticlesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListArticlesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListArticlesRequest proto.InternalMessageInfo

func (m *ListArticlesRequest) GetPageSize() int64 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

func (m *ListArticlesRequest) GetPageToken() string {
	if m != nil {
		return m.PageToken
	}
	return ""
}

func (m *ListArticlesRequest) GetIsRead() *wrappers.BoolValue {
	if m != nil {
		return m.IsRead
	}
	return nil
}

func (m *ListArticlesRequest) GetIsFavorite() *wrappers.BoolValue {
	if m != nil {
		return m.IsFavorite
	}
	return nil
}

func (m *ListArticlesRequest) GetTitle() *wrappers.StringValue {
	if m != nil {
		return m.Title
	}
	return nil
}

type ListArticlesResponse struct {
	// There will be a maximum number of articles returned based on the page_size field int the request.
	Articles []*Article `protobuf:"bytes,1,rep,name=articles,proto3" json:"articles,omitempty"`
	// Token to retrieve the next page of results, or empty if there are no more results in the list.
	NextPageToken        string   `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListArticlesResponse) Reset()         { *m = ListArticlesResponse{} }
func (m *ListArticlesResponse) String() string { return proto.CompactTextString(m) }
func (*ListArticlesResponse) ProtoMessage()    {}
func (*ListArticlesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_75f616a2c230f591, []int{1}
}

func (m *ListArticlesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListArticlesResponse.Unmarshal(m, b)
}
func (m *ListArticlesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListArticlesResponse.Marshal(b, m, deterministic)
}
func (m *ListArticlesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListArticlesResponse.Merge(m, src)
}
func (m *ListArticlesResponse) XXX_Size() int {
	return xxx_messageInfo_ListArticlesResponse.Size(m)
}
func (m *ListArticlesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListArticlesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListArticlesResponse proto.InternalMessageInfo

func (m *ListArticlesResponse) GetArticles() []*Article {
	if m != nil {
		return m.Articles
	}
	return nil
}

func (m *ListArticlesResponse) GetNextPageToken() string {
	if m != nil {
		return m.NextPageToken
	}
	return ""
}

type UpdateArticleRequest struct {
	// The article resource which replaces the resource on the server.
	Article *Article `protobuf:"bytes,1,opt,name=article,proto3" json:"article,omitempty"`
	// FieldMask of fields to update.
	// see https://developers.google.com/protocol-buffers/docs/reference/google.protobuf#fieldmask
	UpdateMask           *field_mask.FieldMask `protobuf:"bytes,2,opt,name=update_mask,json=updateMask,proto3" json:"update_mask,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *UpdateArticleRequest) Reset()         { *m = UpdateArticleRequest{} }
func (m *UpdateArticleRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateArticleRequest) ProtoMessage()    {}
func (*UpdateArticleRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_75f616a2c230f591, []int{2}
}

func (m *UpdateArticleRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateArticleRequest.Unmarshal(m, b)
}
func (m *UpdateArticleRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateArticleRequest.Marshal(b, m, deterministic)
}
func (m *UpdateArticleRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateArticleRequest.Merge(m, src)
}
func (m *UpdateArticleRequest) XXX_Size() int {
	return xxx_messageInfo_UpdateArticleRequest.Size(m)
}
func (m *UpdateArticleRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateArticleRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateArticleRequest proto.InternalMessageInfo

func (m *UpdateArticleRequest) GetArticle() *Article {
	if m != nil {
		return m.Article
	}
	return nil
}

func (m *UpdateArticleRequest) GetUpdateMask() *field_mask.FieldMask {
	if m != nil {
		return m.UpdateMask
	}
	return nil
}

type GetArticleRequest struct {
	// Name of the resource.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Specifies which parts of the article resource should be returned
	// in the response.
	View                 ArticleView `protobuf:"varint,2,opt,name=view,proto3,enum=app.miniboard.users.articles.v1.ArticleView" json:"view,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *GetArticleRequest) Reset()         { *m = GetArticleRequest{} }
func (m *GetArticleRequest) String() string { return proto.CompactTextString(m) }
func (*GetArticleRequest) ProtoMessage()    {}
func (*GetArticleRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_75f616a2c230f591, []int{3}
}

func (m *GetArticleRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetArticleRequest.Unmarshal(m, b)
}
func (m *GetArticleRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetArticleRequest.Marshal(b, m, deterministic)
}
func (m *GetArticleRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetArticleRequest.Merge(m, src)
}
func (m *GetArticleRequest) XXX_Size() int {
	return xxx_messageInfo_GetArticleRequest.Size(m)
}
func (m *GetArticleRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetArticleRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetArticleRequest proto.InternalMessageInfo

func (m *GetArticleRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *GetArticleRequest) GetView() ArticleView {
	if m != nil {
		return m.View
	}
	return ArticleView_ARTICLE_VIEW_UNSPECIFIED
}

type DeleteArticleRequest struct {
	// Name of the resource to delete.
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteArticleRequest) Reset()         { *m = DeleteArticleRequest{} }
func (m *DeleteArticleRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteArticleRequest) ProtoMessage()    {}
func (*DeleteArticleRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_75f616a2c230f591, []int{4}
}

func (m *DeleteArticleRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteArticleRequest.Unmarshal(m, b)
}
func (m *DeleteArticleRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteArticleRequest.Marshal(b, m, deterministic)
}
func (m *DeleteArticleRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteArticleRequest.Merge(m, src)
}
func (m *DeleteArticleRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteArticleRequest.Size(m)
}
func (m *DeleteArticleRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteArticleRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteArticleRequest proto.InternalMessageInfo

func (m *DeleteArticleRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type Article struct {
	// Name of the resource, for example "users/user1/articles/article1".
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Url of the article.
	Url string `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
	// Title of the article.
	Title string `protobuf:"bytes,3,opt,name=title,proto3" json:"title,omitempty"`
	// Icon url.
	IconUrl string `protobuf:"bytes,4,opt,name=icon_url,json=iconUrl,proto3" json:"icon_url,omitempty"`
	// Time when article was added.
	CreateTime *timestamp.Timestamp `protobuf:"bytes,5,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"`
	// Readable content of the article.
	Content []byte `protobuf:"bytes,6,opt,name=content,proto3" json:"content,omitempty"`
	// Mark the article as read.
	IsRead bool `protobuf:"varint,7,opt,name=is_read,json=isRead,proto3" json:"is_read,omitempty"`
	// Mark the article as favorite.
	IsFavorite bool `protobuf:"varint,8,opt,name=is_favorite,json=isFavorite,proto3" json:"is_favorite,omitempty"`
	// Name of source website.
	SiteName string `protobuf:"bytes,9,opt,name=site_name,json=siteName,proto3" json:"site_name,omitempty"`
	// SHA256 sum of the content.
	ContentSha256Sum     string   `protobuf:"bytes,10,opt,name=content_sha256_sum,json=contentSha256Sum,proto3" json:"content_sha256_sum,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Article) Reset()         { *m = Article{} }
func (m *Article) String() string { return proto.CompactTextString(m) }
func (*Article) ProtoMessage()    {}
func (*Article) Descriptor() ([]byte, []int) {
	return fileDescriptor_75f616a2c230f591, []int{5}
}

func (m *Article) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Article.Unmarshal(m, b)
}
func (m *Article) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Article.Marshal(b, m, deterministic)
}
func (m *Article) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Article.Merge(m, src)
}
func (m *Article) XXX_Size() int {
	return xxx_messageInfo_Article.Size(m)
}
func (m *Article) XXX_DiscardUnknown() {
	xxx_messageInfo_Article.DiscardUnknown(m)
}

var xxx_messageInfo_Article proto.InternalMessageInfo

func (m *Article) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Article) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *Article) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Article) GetIconUrl() string {
	if m != nil {
		return m.IconUrl
	}
	return ""
}

func (m *Article) GetCreateTime() *timestamp.Timestamp {
	if m != nil {
		return m.CreateTime
	}
	return nil
}

func (m *Article) GetContent() []byte {
	if m != nil {
		return m.Content
	}
	return nil
}

func (m *Article) GetIsRead() bool {
	if m != nil {
		return m.IsRead
	}
	return false
}

func (m *Article) GetIsFavorite() bool {
	if m != nil {
		return m.IsFavorite
	}
	return false
}

func (m *Article) GetSiteName() string {
	if m != nil {
		return m.SiteName
	}
	return ""
}

func (m *Article) GetContentSha256Sum() string {
	if m != nil {
		return m.ContentSha256Sum
	}
	return ""
}

func init() {
	proto.RegisterEnum("app.miniboard.users.articles.v1.ArticleView", ArticleView_name, ArticleView_value)
	proto.RegisterType((*ListArticlesRequest)(nil), "app.miniboard.users.articles.v1.ListArticlesRequest")
	proto.RegisterType((*ListArticlesResponse)(nil), "app.miniboard.users.articles.v1.ListArticlesResponse")
	proto.RegisterType((*UpdateArticleRequest)(nil), "app.miniboard.users.articles.v1.UpdateArticleRequest")
	proto.RegisterType((*GetArticleRequest)(nil), "app.miniboard.users.articles.v1.GetArticleRequest")
	proto.RegisterType((*DeleteArticleRequest)(nil), "app.miniboard.users.articles.v1.DeleteArticleRequest")
	proto.RegisterType((*Article)(nil), "app.miniboard.users.articles.v1.Article")
}

func init() { proto.RegisterFile("articles_service.proto", fileDescriptor_75f616a2c230f591) }

var fileDescriptor_75f616a2c230f591 = []byte{
	// 811 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x55, 0xdf, 0x4e, 0x3b, 0x45,
	0x14, 0x76, 0xdb, 0x42, 0xdb, 0xd3, 0x1f, 0x52, 0xc6, 0x8a, 0x6b, 0x41, 0x59, 0xd7, 0xc4, 0x94,
	0x86, 0x6c, 0xc3, 0x02, 0x5e, 0x60, 0x4c, 0xa4, 0xd0, 0x9a, 0x26, 0x95, 0x90, 0x2d, 0xc5, 0x84,
	0x9b, 0xcd, 0xd0, 0x0e, 0x75, 0xc2, 0xfe, 0x73, 0x67, 0xb6, 0x28, 0xc6, 0x1b, 0xc3, 0x8d, 0xd7,
	0xc6, 0x78, 0xe5, 0x1b, 0xf8, 0x08, 0xbe, 0x85, 0xaf, 0xe0, 0x3b, 0x78, 0x6b, 0x66, 0x76, 0xb7,
	0xd0, 0x3f, 0xbf, 0x14, 0xee, 0x76, 0xce, 0xf9, 0xce, 0x99, 0x6f, 0xbe, 0x39, 0xdf, 0x2c, 0x6c,
	0xe2, 0x90, 0xd3, 0x81, 0x43, 0x98, 0xcd, 0x48, 0x38, 0xa6, 0x03, 0x62, 0x04, 0xa1, 0xcf, 0x7d,
	0xb4, 0x83, 0x83, 0xc0, 0x70, 0xa9, 0x47, 0x6f, 0x7c, 0x1c, 0x0e, 0x8d, 0x88, 0x91, 0x90, 0x19,
	0x29, 0xd6, 0x18, 0xef, 0x57, 0xb7, 0x47, 0xbe, 0x3f, 0x72, 0x48, 0x03, 0x07, 0xb4, 0x81, 0x3d,
	0xcf, 0xe7, 0x98, 0x53, 0xdf, 0x63, 0x71, 0x79, 0x75, 0x2b, 0xc9, 0xca, 0xd5, 0x4d, 0x74, 0xdb,
	0x20, 0x6e, 0xc0, 0x7f, 0x4c, 0x92, 0x3b, 0xb3, 0x49, 0x4e, 0x5d, 0xc2, 0x38, 0x76, 0x83, 0x04,
	0xa0, 0xcd, 0x02, 0x6e, 0x29, 0x71, 0x86, 0xb6, 0x8b, 0xd9, 0x5d, 0x82, 0xf8, 0x78, 0x16, 0x71,
	0x1f, 0xe2, 0x20, 0x10, 0x1c, 0x65, 0x44, 0xff, 0x4f, 0x81, 0xf7, 0xba, 0x94, 0xf1, 0x93, 0x84,
	0xb1, 0x45, 0xbe, 0x8f, 0x08, 0xe3, 0x68, 0x0b, 0x8a, 0x01, 0x1e, 0x11, 0x9b, 0xd1, 0x07, 0xa2,
	0x2a, 0x9a, 0x52, 0xcb, 0x5a, 0x05, 0x11, 0xe8, 0xd1, 0x07, 0x82, 0x3e, 0x02, 0x90, 0x49, 0xee,
	0xdf, 0x11, 0x4f, 0xcd, 0x68, 0x4a, 0xad, 0x68, 0x49, 0xf8, 0xa5, 0x08, 0xa0, 0x03, 0xc8, 0x53,
	0x66, 0x87, 0x04, 0x0f, 0xd5, 0xac, 0xa6, 0xd4, 0x4a, 0x66, 0xd5, 0x88, 0x59, 0x18, 0x29, 0x0b,
	0xa3, 0xe9, 0xfb, 0xce, 0x15, 0x76, 0x22, 0x62, 0xad, 0x52, 0x66, 0x11, 0x3c, 0x44, 0x5f, 0x40,
	0x89, 0x32, 0xfb, 0x16, 0x8f, 0xfd, 0x90, 0x72, 0xa2, 0xe6, 0x96, 0x16, 0x02, 0x65, 0xed, 0x04,
	0x8d, 0x4c, 0x58, 0xe1, 0x94, 0x3b, 0x44, 0x5d, 0x91, 0x65, 0xdb, 0x73, 0x65, 0x3d, 0x1e, 0x52,
	0x6f, 0x14, 0x17, 0xc6, 0x50, 0xfd, 0x51, 0x81, 0xca, 0xf4, 0xc9, 0x59, 0xe0, 0x7b, 0x8c, 0xa0,
	0x33, 0x28, 0xa4, 0xf7, 0xa7, 0x2a, 0x5a, 0xb6, 0x56, 0x32, 0x6b, 0xc6, 0x92, 0x4b, 0x36, 0x92,
	0x26, 0xd6, 0xa4, 0x12, 0x7d, 0x06, 0xeb, 0x1e, 0xf9, 0x81, 0xdb, 0x73, 0x42, 0xad, 0x89, 0xf0,
	0x45, 0x2a, 0x96, 0xfe, 0x87, 0x02, 0x95, 0x7e, 0x30, 0xc4, 0x9c, 0xa4, 0x3d, 0x92, 0x1b, 0x68,
	0x42, 0x3e, 0x69, 0x26, 0xf5, 0x7f, 0x0d, 0x8b, 0xb4, 0x50, 0x88, 0x1a, 0xc9, 0xde, 0x72, 0x24,
	0x24, 0x81, 0x45, 0xa2, 0xb6, 0xc5, 0xd4, 0x7c, 0x83, 0xd9, 0x9d, 0x05, 0x31, 0x5c, 0x7c, 0xeb,
	0x14, 0x36, 0xbe, 0x26, 0x7c, 0x86, 0x15, 0x82, 0x9c, 0x87, 0xdd, 0x98, 0x52, 0xd1, 0x92, 0xdf,
	0xe8, 0x2b, 0xc8, 0x8d, 0x29, 0xb9, 0x97, 0xed, 0xdf, 0x35, 0xf7, 0x5e, 0x4a, 0xf3, 0x8a, 0x92,
	0x7b, 0x4b, 0x56, 0xea, 0x75, 0xa8, 0x9c, 0x11, 0x87, 0xcc, 0x69, 0xb0, 0x60, 0x37, 0xfd, 0xef,
	0x0c, 0xe4, 0x13, 0xd8, 0x42, 0x36, 0x65, 0xc8, 0x46, 0xa1, 0x93, 0x88, 0x2d, 0x3e, 0x51, 0x25,
	0x9d, 0x8e, 0xac, 0x8c, 0xc5, 0x0b, 0xf4, 0x21, 0x14, 0xe8, 0xc0, 0xf7, 0x6c, 0x01, 0xce, 0xc9,
	0x44, 0x5e, 0xac, 0xfb, 0xa1, 0x23, 0x64, 0x1b, 0x84, 0x44, 0xc8, 0x26, 0x0c, 0x97, 0x0c, 0xd5,
	0xbc, 0x6c, 0x97, 0xa9, 0x1b, 0x2d, 0x88, 0xe1, 0x22, 0x80, 0x54, 0xc8, 0x0f, 0x7c, 0x8f, 0x13,
	0x8f, 0xab, 0xab, 0x9a, 0x52, 0x7b, 0x63, 0xa5, 0x4b, 0xf4, 0xc1, 0x93, 0x2f, 0xf2, 0x9a, 0x52,
	0x2b, 0x4c, 0x66, 0x7f, 0x67, 0x7a, 0xf6, 0x0b, 0x32, 0xf9, 0x7c, 0xbe, 0xb7, 0xa0, 0xc8, 0x28,
	0x27, 0xb6, 0x3c, 0x6c, 0x51, 0x92, 0x2d, 0x88, 0xc0, 0xb9, 0x38, 0xf0, 0x1e, 0xa0, 0x64, 0x07,
	0x9b, 0x7d, 0x87, 0xcd, 0xa3, 0xcf, 0x6d, 0x16, 0xb9, 0x2a, 0x48, 0x54, 0x39, 0xc9, 0xf4, 0x64,
	0xa2, 0x17, 0xb9, 0xf5, 0x6b, 0x28, 0x3d, 0xd3, 0x1f, 0x6d, 0x83, 0x7a, 0x62, 0x5d, 0x76, 0x4e,
	0xbb, 0x2d, 0xfb, 0xaa, 0xd3, 0xfa, 0xd6, 0xee, 0x9f, 0xf7, 0x2e, 0x5a, 0xa7, 0x9d, 0x76, 0xa7,
	0x75, 0x56, 0x7e, 0x07, 0x6d, 0x02, 0x9a, 0xca, 0x36, 0x4f, 0x7a, 0x9d, 0xd3, 0xb2, 0x82, 0xde,
	0x87, 0x8d, 0xa9, 0x78, 0xbb, 0xdf, 0xed, 0x96, 0x33, 0xe6, 0xe3, 0x0a, 0xac, 0xa7, 0x76, 0xea,
	0xc5, 0xaf, 0x24, 0xfa, 0x53, 0x81, 0x37, 0xcf, 0x6d, 0x86, 0x0e, 0x97, 0xce, 0xc7, 0x82, 0xf7,
	0xa8, 0x7a, 0xf4, 0xca, 0xaa, 0xd8, 0xcb, 0xba, 0xf6, 0xcb, 0x3f, 0xff, 0xfe, 0x96, 0xa9, 0x22,
	0x55, 0x3e, 0xbf, 0xe3, 0xfd, 0x86, 0xac, 0x6b, 0xd4, 0x1b, 0x13, 0x9f, 0xfe, 0xa5, 0xc0, 0xda,
	0x94, 0xff, 0xd0, 0xf2, 0xad, 0x16, 0xf9, 0xb5, 0xfa, 0x62, 0x7b, 0xea, 0xc7, 0x92, 0xd4, 0xa1,
	0xb9, 0x9b, 0x92, 0xfa, 0x29, 0x41, 0x19, 0xe2, 0x7e, 0xbf, 0x9c, 0xa5, 0xd8, 0xa8, 0xff, 0x7c,
	0x3c, 0x71, 0xf4, 0xef, 0x0a, 0xc0, 0x93, 0x2b, 0x91, 0xb9, 0x74, 0xd3, 0x39, 0x0b, 0xbf, 0x82,
	0xe8, 0xae, 0x24, 0xfa, 0x29, 0xfa, 0x64, 0x42, 0xf4, 0x6d, 0x04, 0xd1, 0xaf, 0x0a, 0xac, 0x4d,
	0x59, 0xf8, 0x05, 0x32, 0x2e, 0xb2, 0x7c, 0x75, 0x73, 0xce, 0x66, 0x2d, 0xf1, 0x47, 0x4c, 0xb9,
	0xd4, 0x97, 0x73, 0x69, 0xc2, 0xf5, 0xe4, 0x19, 0xbe, 0x59, 0x95, 0x6d, 0x0e, 0xfe, 0x0f, 0x00,
	0x00, 0xff, 0xff, 0x13, 0xd4, 0xec, 0x9f, 0xbf, 0x07, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ArticlesServiceClient is the client API for ArticlesService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ArticlesServiceClient interface {
	// List articles
	//
	// An endpoint to list the user's articles.
	ListArticles(ctx context.Context, in *ListArticlesRequest, opts ...grpc.CallOption) (*ListArticlesResponse, error)
	// Update article
	//
	// An endpoint to update the article.
	UpdateArticle(ctx context.Context, in *UpdateArticleRequest, opts ...grpc.CallOption) (*Article, error)
	// Get article
	//
	// An endpoint to get an article by id.
	GetArticle(ctx context.Context, in *GetArticleRequest, opts ...grpc.CallOption) (*Article, error)
	// Delete article
	//
	// An endpoint to delete an article by id.
	DeleteArticle(ctx context.Context, in *DeleteArticleRequest, opts ...grpc.CallOption) (*empty.Empty, error)
}

type articlesServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewArticlesServiceClient(cc grpc.ClientConnInterface) ArticlesServiceClient {
	return &articlesServiceClient{cc}
}

func (c *articlesServiceClient) ListArticles(ctx context.Context, in *ListArticlesRequest, opts ...grpc.CallOption) (*ListArticlesResponse, error) {
	out := new(ListArticlesResponse)
	err := c.cc.Invoke(ctx, "/app.miniboard.users.articles.v1.ArticlesService/ListArticles", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *articlesServiceClient) UpdateArticle(ctx context.Context, in *UpdateArticleRequest, opts ...grpc.CallOption) (*Article, error) {
	out := new(Article)
	err := c.cc.Invoke(ctx, "/app.miniboard.users.articles.v1.ArticlesService/UpdateArticle", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *articlesServiceClient) GetArticle(ctx context.Context, in *GetArticleRequest, opts ...grpc.CallOption) (*Article, error) {
	out := new(Article)
	err := c.cc.Invoke(ctx, "/app.miniboard.users.articles.v1.ArticlesService/GetArticle", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *articlesServiceClient) DeleteArticle(ctx context.Context, in *DeleteArticleRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/app.miniboard.users.articles.v1.ArticlesService/DeleteArticle", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ArticlesServiceServer is the server API for ArticlesService service.
type ArticlesServiceServer interface {
	// List articles
	//
	// An endpoint to list the user's articles.
	ListArticles(context.Context, *ListArticlesRequest) (*ListArticlesResponse, error)
	// Update article
	//
	// An endpoint to update the article.
	UpdateArticle(context.Context, *UpdateArticleRequest) (*Article, error)
	// Get article
	//
	// An endpoint to get an article by id.
	GetArticle(context.Context, *GetArticleRequest) (*Article, error)
	// Delete article
	//
	// An endpoint to delete an article by id.
	DeleteArticle(context.Context, *DeleteArticleRequest) (*empty.Empty, error)
}

// UnimplementedArticlesServiceServer can be embedded to have forward compatible implementations.
type UnimplementedArticlesServiceServer struct {
}

func (*UnimplementedArticlesServiceServer) ListArticles(ctx context.Context, req *ListArticlesRequest) (*ListArticlesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListArticles not implemented")
}
func (*UnimplementedArticlesServiceServer) UpdateArticle(ctx context.Context, req *UpdateArticleRequest) (*Article, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateArticle not implemented")
}
func (*UnimplementedArticlesServiceServer) GetArticle(ctx context.Context, req *GetArticleRequest) (*Article, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetArticle not implemented")
}
func (*UnimplementedArticlesServiceServer) DeleteArticle(ctx context.Context, req *DeleteArticleRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteArticle not implemented")
}

func RegisterArticlesServiceServer(s *grpc.Server, srv ArticlesServiceServer) {
	s.RegisterService(&_ArticlesService_serviceDesc, srv)
}

func _ArticlesService_ListArticles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListArticlesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArticlesServiceServer).ListArticles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.miniboard.users.articles.v1.ArticlesService/ListArticles",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArticlesServiceServer).ListArticles(ctx, req.(*ListArticlesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ArticlesService_UpdateArticle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateArticleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArticlesServiceServer).UpdateArticle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.miniboard.users.articles.v1.ArticlesService/UpdateArticle",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArticlesServiceServer).UpdateArticle(ctx, req.(*UpdateArticleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ArticlesService_GetArticle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetArticleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArticlesServiceServer).GetArticle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.miniboard.users.articles.v1.ArticlesService/GetArticle",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArticlesServiceServer).GetArticle(ctx, req.(*GetArticleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ArticlesService_DeleteArticle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteArticleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArticlesServiceServer).DeleteArticle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.miniboard.users.articles.v1.ArticlesService/DeleteArticle",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArticlesServiceServer).DeleteArticle(ctx, req.(*DeleteArticleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ArticlesService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "app.miniboard.users.articles.v1.ArticlesService",
	HandlerType: (*ArticlesServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListArticles",
			Handler:    _ArticlesService_ListArticles_Handler,
		},
		{
			MethodName: "UpdateArticle",
			Handler:    _ArticlesService_UpdateArticle_Handler,
		},
		{
			MethodName: "GetArticle",
			Handler:    _ArticlesService_GetArticle_Handler,
		},
		{
			MethodName: "DeleteArticle",
			Handler:    _ArticlesService_DeleteArticle_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "articles_service.proto",
}
