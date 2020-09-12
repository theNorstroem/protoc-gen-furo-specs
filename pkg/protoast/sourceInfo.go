package protoast

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

type SourceInfo struct {
	Messages []MessageInfo
}
type MessageInfo struct {
	Name       string
	Info       *descriptor.SourceCodeInfo_Location
	FieldInfos []FieldInfo
}
type FieldInfo struct {
	Name  string
	Info  *descriptor.SourceCodeInfo_Location
	Field *descriptor.FieldDescriptorProto
}

func GetSourceInfo(descr *descriptor.FileDescriptorProto) SourceInfo {
	SourceInfo := SourceInfo{}

	for _, location := range descr.GetSourceCodeInfo().GetLocation() {
		// 4 111 2 222 => 4 MessageIndex 2 FieldIndex
		// for field with index 222 in message with index 111
		// location info for descriptor.MessageType[111]Field[222]
		if len(location.GetPath()) == 2 && location.Path[0] == 4 {
			msgIndex := location.Path[1]
			SourceInfo.Messages = append(SourceInfo.Messages, MessageInfo{
				Name:       *descr.MessageType[msgIndex].Name,
				Info:       location,
				FieldInfos: []FieldInfo{},
			})
		}

		// 4 111 2 222 =>	 4 MessageIndex 2 FieldIndex
		// for field with index 222 in message with index 111
		// location info for descriptor.MessageType[111]Field[222]
		if len(location.GetPath()) == 4 && location.Path[0] == 4 && location.Path[2] == 2 {
			msgIndex := location.Path[1]
			fieldIndex := location.Path[3]
			fi := FieldInfo{
				Name:  *descr.MessageType[msgIndex].Field[fieldIndex].Name,
				Info:  location,
				Field: descr.MessageType[msgIndex].Field[fieldIndex],
			}
			SourceInfo.Messages[msgIndex].FieldInfos = append(SourceInfo.Messages[msgIndex].FieldInfos, fi)
		}

	}
	return SourceInfo
}
