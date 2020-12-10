// Copyright 2020 Syntio Inc.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package impl represents implementation of message validation process.
package impl

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/syntio/central-consumer/registry"
)

var tmpFilePath = registry.Cfg.Protoparam.TmpFilePath
var tmpFileName = registry.Cfg.Protoparam.TmpFileName
var mode os.FileMode = registry.Cfg.FileMode

const msgDescriptorLimit = 1

// ProtobufValidator is a validator structure for protobuf format.
type ProtobufValidator struct{}

// Validate validates a protobuf message with a schema.
//
// Function returns the validation boolean result. An error is returned if any errors occur during the
// function execution.
func (proto *ProtobufValidator) Validate(message, schema []byte) (bool, error) {
	var valid bool = false

	completeFilePath := filepath.Join(tmpFilePath, tmpFileName)
	file, err := os.OpenFile(completeFilePath, os.O_CREATE|os.O_WRONLY, mode)
	if err != nil {
		return valid, err
	}
	_, err = file.WriteString(string(schema))
	if err != nil {
		return valid, err
	}

	mainFileDescriptor, err := parseMainFileDescriptor(tmpFilePath, tmpFileName)
	if err != nil {
		return valid, err
	}

	mainMessage, err := parseMainMessageDescriptor(mainFileDescriptor)
	if err != nil {
		return valid, err
	}

	if err = mainMessage.Unmarshal(message); err != nil {
		return valid, err
	}

	if err := mainMessage.ValidateRecursive(); err == nil {
		valid = true
	}
	//file.Close()
	return valid, err
}

// parseMainFileDescriptor retrieves a file-descriptor structure of the protobuf schema which is used to
// programmatically setup a protobuf message structure. The schema has to be imported from a file.
//
// An error is returned if any errors occur during the function execution.
func parseMainFileDescriptor(importPath, fileName string) (*desc.FileDescriptor, error) {
	var mainFileDescriptor *desc.FileDescriptor = nil

	parser := protoparse.Parser{
		ImportPaths: []string{importPath},
	}

	fileDescriptors, err := parser.ParseFiles(fileName)
	if err != nil {
		return mainFileDescriptor, err
	}

	mainFileDescriptor = fileDescriptors[0]
	return mainFileDescriptor, err
}

// parseMainMessageDescriptor retrieves a protobuf message structure from the file-descriptor structure.
// The protobuf message structure is used for protobuf message validation.
//
// An error is returned if any errors occur during the function execution.
func parseMainMessageDescriptor(mainFileDescriptor *desc.FileDescriptor) (*dynamic.Message, error) {
	var mainMessageDescriptor *dynamic.Message
	var err error

	messageDescriptors := mainFileDescriptor.GetMessageTypes()
	if len(messageDescriptors) < msgDescriptorLimit {
		err = fmt.Errorf("ERROR: main file descriptor has less than 1 message descriptors.\n")
	}
	mainMessageDescriptor = dynamic.NewMessage(messageDescriptors[0])

	return mainMessageDescriptor, err
}
