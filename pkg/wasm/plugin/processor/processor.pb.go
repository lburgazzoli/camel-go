// Code generated by protoc-gen-go-plugin. DO NOT EDIT.
// versions:
// 	protoc-gen-go-plugin v0.1.0
// 	protoc               v3.21.12
// source: processor.proto

package processor

import (
	context "context"

	timestamppb "github.com/knqyf263/go-plugin/types/known/timestamppb"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Source        string                 `protobuf:"bytes,2,opt,name=source,proto3" json:"source,omitempty"`
	Type          string                 `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
	Subject       string                 `protobuf:"bytes,4,opt,name=subject,proto3" json:"subject,omitempty"`
	ContentType   string                 `protobuf:"bytes,5,opt,name=content_type,json=contentType,proto3" json:"content_type,omitempty"`
	ContentSchema string                 `protobuf:"bytes,6,opt,name=content_schema,json=contentSchema,proto3" json:"content_schema,omitempty"`
	Time          *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=time,proto3" json:"time,omitempty"`
	Attributes    map[string]string      `protobuf:"bytes,8,rep,name=attributes,proto3" json:"attributes,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Data          []byte                 `protobuf:"bytes,9,opt,name=data,proto3" json:"data,omitempty"`
	Annotations   map[string]string      `protobuf:"bytes,10,rep,name=annotations,proto3" json:"annotations,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Message) ProtoReflect() protoreflect.Message {
	panic(`not implemented`)
}

func (x *Message) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Message) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

func (x *Message) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Message) GetSubject() string {
	if x != nil {
		return x.Subject
	}
	return ""
}

func (x *Message) GetContentType() string {
	if x != nil {
		return x.ContentType
	}
	return ""
}

func (x *Message) GetContentSchema() string {
	if x != nil {
		return x.ContentSchema
	}
	return ""
}

func (x *Message) GetTime() *timestamppb.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

func (x *Message) GetAttributes() map[string]string {
	if x != nil {
		return x.Attributes
	}
	return nil
}

func (x *Message) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Message) GetAnnotations() map[string]string {
	if x != nil {
		return x.Annotations
	}
	return nil
}

// go:plugin type=plugin version=1
type Processors interface {
	Process(context.Context, *Message) (*Message, error)
}
