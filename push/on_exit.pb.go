// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: on_exit.proto

package push

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

type OnExit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PID uint32 `protobuf:"varint,1,opt,name=PID,proto3" json:"PID,omitempty"`
}

func (x *OnExit) Reset() {
	*x = OnExit{}
	if protoimpl.UnsafeEnabled {
		mi := &file_on_exit_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OnExit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OnExit) ProtoMessage() {}

func (x *OnExit) ProtoReflect() protoreflect.Message {
	mi := &file_on_exit_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OnExit.ProtoReflect.Descriptor instead.
func (*OnExit) Descriptor() ([]byte, []int) {
	return file_on_exit_proto_rawDescGZIP(), []int{0}
}

func (x *OnExit) GetPID() uint32 {
	if x != nil {
		return x.PID
	}
	return 0
}

var File_on_exit_proto protoreflect.FileDescriptor

var file_on_exit_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x6f, 0x6e, 0x5f, 0x65, 0x78, 0x69, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x04, 0x70, 0x75, 0x73, 0x68, 0x22, 0x1a, 0x0a, 0x06, 0x4f, 0x6e, 0x45, 0x78, 0x69, 0x74, 0x12,
	0x10, 0x0a, 0x03, 0x50, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x03, 0x50, 0x49,
	0x44, 0x42, 0x08, 0x5a, 0x06, 0x2f, 0x3b, 0x70, 0x75, 0x73, 0x68, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_on_exit_proto_rawDescOnce sync.Once
	file_on_exit_proto_rawDescData = file_on_exit_proto_rawDesc
)

func file_on_exit_proto_rawDescGZIP() []byte {
	file_on_exit_proto_rawDescOnce.Do(func() {
		file_on_exit_proto_rawDescData = protoimpl.X.CompressGZIP(file_on_exit_proto_rawDescData)
	})
	return file_on_exit_proto_rawDescData
}

var file_on_exit_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_on_exit_proto_goTypes = []interface{}{
	(*OnExit)(nil), // 0: push.OnExit
}
var file_on_exit_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_on_exit_proto_init() }
func file_on_exit_proto_init() {
	if File_on_exit_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_on_exit_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OnExit); i {
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
			RawDescriptor: file_on_exit_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_on_exit_proto_goTypes,
		DependencyIndexes: file_on_exit_proto_depIdxs,
		MessageInfos:      file_on_exit_proto_msgTypes,
	}.Build()
	File_on_exit_proto = out.File
	file_on_exit_proto_rawDesc = nil
	file_on_exit_proto_goTypes = nil
	file_on_exit_proto_depIdxs = nil
}
