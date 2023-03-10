// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: type.proto

package pb

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

//packet type
type Type int32

const (
	Type_Handshake_    Type = 0
	Type_HandshakeAck_ Type = 1
	Type_Heartbeat_    Type = 2
	Type_Request_      Type = 3
	Type_Response_     Type = 4
	Type_Push_         Type = 5
	Type_Kick_         Type = 6
	Type_Join          Type = 7
	Type_Leave         Type = 8
)

// Enum value maps for Type.
var (
	Type_name = map[int32]string{
		0: "Handshake_",
		1: "HandshakeAck_",
		2: "Heartbeat_",
		3: "Request_",
		4: "Response_",
		5: "Push_",
		6: "Kick_",
		7: "Join",
		8: "Leave",
	}
	Type_value = map[string]int32{
		"Handshake_":    0,
		"HandshakeAck_": 1,
		"Heartbeat_":    2,
		"Request_":      3,
		"Response_":     4,
		"Push_":         5,
		"Kick_":         6,
		"Join":          7,
		"Leave":         8,
	}
)

func (x Type) Enum() *Type {
	p := new(Type)
	*p = x
	return p
}

func (x Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Type) Descriptor() protoreflect.EnumDescriptor {
	return file_type_proto_enumTypes[0].Descriptor()
}

func (Type) Type() protoreflect.EnumType {
	return &file_type_proto_enumTypes[0]
}

func (x Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Type.Descriptor instead.
func (Type) EnumDescriptor() ([]byte, []int) {
	return file_type_proto_rawDescGZIP(), []int{0}
}

var File_type_proto protoreflect.FileDescriptor

var file_type_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x74, 0x79, 0x70, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62,
	0x2a, 0x81, 0x01, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0e, 0x0a, 0x0a, 0x48, 0x61, 0x6e,
	0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x5f, 0x10, 0x00, 0x12, 0x11, 0x0a, 0x0d, 0x48, 0x61, 0x6e,
	0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x41, 0x63, 0x6b, 0x5f, 0x10, 0x01, 0x12, 0x0e, 0x0a, 0x0a,
	0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x5f, 0x10, 0x02, 0x12, 0x0c, 0x0a, 0x08,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x10, 0x03, 0x12, 0x0d, 0x0a, 0x09, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x5f, 0x10, 0x04, 0x12, 0x09, 0x0a, 0x05, 0x50, 0x75, 0x73,
	0x68, 0x5f, 0x10, 0x05, 0x12, 0x09, 0x0a, 0x05, 0x4b, 0x69, 0x63, 0x6b, 0x5f, 0x10, 0x06, 0x12,
	0x08, 0x0a, 0x04, 0x4a, 0x6f, 0x69, 0x6e, 0x10, 0x07, 0x12, 0x09, 0x0a, 0x05, 0x4c, 0x65, 0x61,
	0x76, 0x65, 0x10, 0x08, 0x42, 0x06, 0x5a, 0x04, 0x2f, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_type_proto_rawDescOnce sync.Once
	file_type_proto_rawDescData = file_type_proto_rawDesc
)

func file_type_proto_rawDescGZIP() []byte {
	file_type_proto_rawDescOnce.Do(func() {
		file_type_proto_rawDescData = protoimpl.X.CompressGZIP(file_type_proto_rawDescData)
	})
	return file_type_proto_rawDescData
}

var file_type_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_type_proto_goTypes = []interface{}{
	(Type)(0), // 0: pb.Type
}
var file_type_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_type_proto_init() }
func file_type_proto_init() {
	if File_type_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_type_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_type_proto_goTypes,
		DependencyIndexes: file_type_proto_depIdxs,
		EnumInfos:         file_type_proto_enumTypes,
	}.Build()
	File_type_proto = out.File
	file_type_proto_rawDesc = nil
	file_type_proto_goTypes = nil
	file_type_proto_depIdxs = nil
}
