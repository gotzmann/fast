package gguf

import (
	"encoding/binary"
	"os"
)

type String struct {
	Length uint64
	Str    []byte
}

func (s *String) String() string {
	return string(s.Str)
}

func (s *String) Read(file *os.File) error {
	err := binary.Read(file, binary.LittleEndian, &s.Length)
	if err != nil {
		return err
	}

	s.Str = make([]byte, s.Length)
	_, err = file.Read(s.Str)
	if err != nil {
		return err
	}

	return nil
}