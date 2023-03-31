// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.4
// source: token.proto

package JwtGererater

import (
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

type Requset struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccessToken string `protobuf:"bytes,1,opt,name=accessToken,proto3" json:"accessToken,omitempty"`
}

func (x *Requset) Reset() {
	*x = Requset{}
	if protoimpl.UnsafeEnabled {
		mi := &file_token_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Requset) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Requset) ProtoMessage() {}

func (x *Requset) ProtoReflect() protoreflect.Message {
	mi := &file_token_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Requset.ProtoReflect.Descriptor instead.
func (*Requset) Descriptor() ([]byte, []int) {
	return file_token_proto_rawDescGZIP(), []int{0}
}

func (x *Requset) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

type Respond struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode   int32   `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	ErrorMessage string  `protobuf:"bytes,2,opt,name=error_message,json=errorMessage,proto3" json:"error_message,omitempty"`
	Data         *Tokens `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *Respond) Reset() {
	*x = Respond{}
	if protoimpl.UnsafeEnabled {
		mi := &file_token_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Respond) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Respond) ProtoMessage() {}

func (x *Respond) ProtoReflect() protoreflect.Message {
	mi := &file_token_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Respond.ProtoReflect.Descriptor instead.
func (*Respond) Descriptor() ([]byte, []int) {
	return file_token_proto_rawDescGZIP(), []int{1}
}

func (x *Respond) GetStatusCode() int32 {
	if x != nil {
		return x.StatusCode
	}
	return 0
}

func (x *Respond) GetErrorMessage() string {
	if x != nil {
		return x.ErrorMessage
	}
	return ""
}

func (x *Respond) GetData() *Tokens {
	if x != nil {
		return x.Data
	}
	return nil
}

type Tokens struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RefreshToken string `protobuf:"bytes,2,opt,name=refreshToken,proto3" json:"refreshToken,omitempty"`
}

func (x *Tokens) Reset() {
	*x = Tokens{}
	if protoimpl.UnsafeEnabled {
		mi := &file_token_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Tokens) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tokens) ProtoMessage() {}

func (x *Tokens) ProtoReflect() protoreflect.Message {
	mi := &file_token_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tokens.ProtoReflect.Descriptor instead.
func (*Tokens) Descriptor() ([]byte, []int) {
	return file_token_proto_rawDescGZIP(), []int{2}
}

func (x *Tokens) GetRefreshToken() string {
	if x != nil {
		return x.RefreshToken
	}
	return ""
}

var File_token_proto protoreflect.FileDescriptor

var file_token_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x4a,
	0x77, 0x74, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x2b, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x73,
	0x65, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x75, 0x0a, 0x07, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x64, 0x12,
	0x1f, 0x0a, 0x0b, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64, 0x65,
	0x12, 0x23, 0x0a, 0x0d, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x24, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x4a, 0x77, 0x74, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x2e, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x2c, 0x0a, 0x06, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x12, 0x22, 0x0a, 0x0c, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x66,
	0x72, 0x65, 0x73, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x32, 0x42, 0x0a, 0x08, 0x4a, 0x77, 0x74,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x36, 0x0a, 0x0c, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x11, 0x2e, 0x4a, 0x77, 0x74, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x2e, 0x52, 0x65, 0x71, 0x75, 0x73, 0x65, 0x74, 0x1a, 0x11, 0x2e, 0x4a, 0x77, 0x74, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x64, 0x22, 0x00, 0x42, 0x11, 0x5a,
	0x0f, 0x2e, 0x2f, 0x3b, 0x4a, 0x77, 0x74, 0x47, 0x65, 0x72, 0x65, 0x72, 0x61, 0x74, 0x65, 0x72,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_token_proto_rawDescOnce sync.Once
	file_token_proto_rawDescData = file_token_proto_rawDesc
)

func file_token_proto_rawDescGZIP() []byte {
	file_token_proto_rawDescOnce.Do(func() {
		file_token_proto_rawDescData = protoimpl.X.CompressGZIP(file_token_proto_rawDescData)
	})
	return file_token_proto_rawDescData
}

var file_token_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_token_proto_goTypes = []interface{}{
	(*Requset)(nil), // 0: JwtToken.Requset
	(*Respond)(nil), // 1: JwtToken.Respond
	(*Tokens)(nil),  // 2: JwtToken.tokens
}
var file_token_proto_depIdxs = []int32{
	2, // 0: JwtToken.Respond.data:type_name -> JwtToken.tokens
	0, // 1: JwtToken.JwtToken.refreshToken:input_type -> JwtToken.Requset
	1, // 2: JwtToken.JwtToken.refreshToken:output_type -> JwtToken.Respond
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_token_proto_init() }
func file_token_proto_init() {
	if File_token_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_token_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Requset); i {
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
		file_token_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Respond); i {
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
		file_token_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Tokens); i {
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
			RawDescriptor: file_token_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_token_proto_goTypes,
		DependencyIndexes: file_token_proto_depIdxs,
		MessageInfos:      file_token_proto_msgTypes,
	}.Build()
	File_token_proto = out.File
	file_token_proto_rawDesc = nil
	file_token_proto_goTypes = nil
	file_token_proto_depIdxs = nil
}
