package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/amery/protogen/pkg/protogen"
)

func saveRawMessage(m protoreflect.ProtoMessage, name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := proto.Marshal(m)
	if err != nil {
		return err
	}

	_, err = bytes.NewBuffer(b).WriteTo(f)
	return err
}

func saveRawRequest(gen *protogen.Plugin, name string) error {
	dirName, fileName := filepath.Split(filepath.Clean(name))
	switch {
	case fileName == "":
		return protogen.ErrInvalidName
	case !strings.HasSuffix(fileName, ".pb"):
		fileName += ".req.pb"
	}

	name = dirName + fileName
	gen.Printf("writing binary request to %q", name)
	return saveRawMessage(gen.Request(), name)
}
