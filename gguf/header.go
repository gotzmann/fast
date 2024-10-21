package gguf

import (
	"encoding/binary"
	"fmt"
	"os"
)

// Header represents the header information in a GGUF file.
type Header struct {
	MagicNumber     uint32
	Version         uint32
	TensorCount     uint64
	MetadataKVCount uint64
	Metadata        []MetadataKV
}

func (h *Header) String() string {
	magicToString := func(magic uint32) string {
		return fmt.Sprintf("%c%c%c%c", byte(magic), byte(magic>>8), byte(magic>>16), byte(magic>>24))
	}

	str := fmt.Sprintf("Magic Number: %x (%s)\n", h.MagicNumber, magicToString(h.MagicNumber))
	str += fmt.Sprintf("Version: %d\n", h.Version)
	str += fmt.Sprintf("Tensor Count: %d\n", h.TensorCount)
	str += fmt.Sprintf("Metadata KV Count: %d\n", h.MetadataKVCount)
	str += fmt.Sprintf("Metadata:\n")
	for i, kv := range h.Metadata {
		str += fmt.Sprintf("\tMetadata KV %d:\n", i)
		str += fmt.Sprintf("\t\t%s\n", kv.String())
	}
	return str
}

func (h *Header) Read(file *os.File) error {
	err := binary.Read(file, binary.LittleEndian, &h.MagicNumber)
	if err != nil {
		return err
	}
	err = binary.Read(file, binary.LittleEndian, &h.Version)
	if err != nil {
		return err
	}
	err = binary.Read(file, binary.LittleEndian, &h.TensorCount)
	if err != nil {
		return err
	}
	err = binary.Read(file, binary.LittleEndian, &h.MetadataKVCount)
	if err != nil {
		return err
	}

	h.Metadata = make([]MetadataKV, h.MetadataKVCount)
	for i := range h.Metadata {
		err = h.Metadata[i].Read(file)
		if err != nil {
			return err
		}
	}

	return nil
}
