// Copyright (c) 2023 Kevin Nguyen <kevin.nguyen.ai@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"fmt"
	"io"
	"os"
	"sync"
)

const chunkSize = 1 << 20 // 1 MB chunks, adjust as needed

func main() {
	inputFileName := "input.pcapng"
	outputFileName := "output.pcapng"
	numChunks := 10

	// Split the file into smaller chunks
	err := splitFile(inputFileName, numChunks)
	if err != nil {
		fmt.Println("Error splitting file:", err)
		return
	}

	// Process each chunk concurrently
	var wg sync.WaitGroup
	for i := 0; i < numChunks; i++ {
		wg.Add(1)
		go func(chunkIndex int) {
			defer wg.Done()

			chunkFileName := fmt.Sprintf("chunk%d.pcapng", chunkIndex)
			encryptedChunkFileName := fmt.Sprintf("encrypted_chunk%d.pcapng", chunkIndex)

			// Simulate encryption, replace with your encryption logic
			// For a real implementation, you would read the chunk, encrypt it, and then write the encrypted data
			err := simulateEncryption(chunkFileName, encryptedChunkFileName)
			if err != nil {
				fmt.Printf("Error encrypting chunk %d: %v\n", chunkIndex, err)
				return
			}
		}(i)
	}

	// Wait for all encryption goroutines to finish
	wg.Wait()

	// Combine the encrypted chunks into a single file
	err = combineChunks(outputFileName, numChunks)
	if err != nil {
		fmt.Println("Error combining chunks:", err)
		return
	}

	fmt.Println("Encryption complete.")
}

func splitFile(inputFileName string, numChunks int) error {
	file, err := os.Open(inputFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	chunkSize := (fileInfo.Size() + int64(numChunks) - 1) / int64(numChunks)

	for i := 0; i < numChunks; i++ {
		chunkFileName := fmt.Sprintf("chunk%d.pcapng", i)
		chunkFile, err := os.Create(chunkFileName)
		if err != nil {
			return err
		}

		bytesToCopy := chunkSize
		if i == numChunks-1 {
			bytesToCopy = fileInfo.Size() - int64(i)*chunkSize
		}

		_, err = io.CopyN(chunkFile, file, bytesToCopy)
		if err != nil {
			chunkFile.Close()
			return err
		}

		chunkFile.Close()
	}

	return nil
}

func simulateEncryption(inputFileName, outputFileName string) error {
	inputFile, err := os.Open(inputFileName)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputFileName)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Simulate encryption by copying data from input to output
	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		return err
	}

	return nil
}

func combineChunks(outputFileName string, numChunks int) error {
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	for i := 0; i < numChunks; i++ {
		chunkFileName := fmt.Sprintf("encrypted_chunk%d.pcapng", i)
		chunkFile, err := os.Open(chunkFileName)
		if err != nil {
			return err
		}
		defer chunkFile.Close()

		_, err = io.Copy(outputFile, chunkFile)
		if err != nil {
			return err
		}
	}

	return nil
}
