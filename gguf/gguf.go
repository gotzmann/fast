package gguf


import (
	"fmt"
	"log"
	"os"
)

// DebugReadGGUF reads a GGUF file and prints its header information.
func DebugReadGGUF(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Printf("Error closing file: %v\n", err)
		}
	}(f)

	// Read the header
	var header Header
	err = header.Read(f)
	if err != nil {
		return err
	}

	// Print header information
	fmt.Printf("GGUF File Information:\n")
	fmt.Println(header.String())

	return nil
}
