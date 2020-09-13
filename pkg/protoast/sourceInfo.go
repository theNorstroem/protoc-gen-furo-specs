package protoast

import (
	"google.golang.org/protobuf/types/descriptorpb"
)

type SourceInfo struct {
	Messages []MessageInfo
}
type MessageInfo struct {
	Name       string
	Info       *descriptorpb.SourceCodeInfo_Location
	FieldInfos []FieldInfo
	Message    descriptorpb.DescriptorProto
}
type FieldInfo struct {
	Name    string
	Info    *descriptorpb.SourceCodeInfo_Location
	Field   *descriptorpb.FieldDescriptorProto
	Message descriptorpb.DescriptorProto
}

func GetSourceInfo(descr *descriptorpb.FileDescriptorProto) SourceInfo {
	SourceInfo := SourceInfo{}

	for _, location := range descr.GetSourceCodeInfo().GetLocation() {
		// 4 111 2 222 => 4 MessageIndex 2 FieldIndex
		// for field with index 222 in message with index 111
		// location info for descriptor.MessageType[111]Field[222]
		if len(location.GetPath()) == 2 && location.Path[0] == 4 {
			msgIndex := location.Path[1]
			SourceInfo.Messages = append(SourceInfo.Messages, MessageInfo{
				Name:       *descr.MessageType[msgIndex].Name,
				Message:    *descr.MessageType[msgIndex],
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
				Name:    *descr.MessageType[msgIndex].Field[fieldIndex].Name,
				Info:    location,
				Field:   descr.MessageType[msgIndex].Field[fieldIndex],
				Message: *descr.MessageType[msgIndex],
			}
			SourceInfo.Messages[msgIndex].FieldInfos = append(SourceInfo.Messages[msgIndex].FieldInfos, fi)
		}

		if len(location.GetPath()) == 5 && location.Path[0] == 4 && location.Path[2] == 2 {
			a := location.Path[1]
			a = a
		}

	}
	return SourceInfo
}
