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

const (
	chunkSize = 1 << 20 // 1 MB chunks, adjust as needed
	numChunks = 10
)

func main() {
	inputFileName := "input.pcap"
	outputFileName := "output.pcap"

	// Create channels for communication between stages
	chunkChannel := make(chan *chunkData, numChunks)
	encryptedChunkChannel := make(chan *chunkData, numChunks)

	// Wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Stage 1: Split the file into chunks and send them to the encryption stage
	wg.Add(1)
	go splitAndSend(inputFileName, chunkChannel, &wg)

	// Stage 2: Encrypt each chunk concurrently and send them to the combining stage
	for i := 0; i < numChunks; i++ {
		wg.Add(1)
		go encryptAndSend(chunkChannel, encryptedChunkChannel, &wg)
	}

	// Stage 3: Combine the encrypted chunks and write to the output file
	go combineAndWrite(outputFileName, encryptedChunkChannel, &wg)

	// Wait for all goroutines to finish
	wg.Wait()

	fmt.Println("Encryption complete.")
}

type chunkData struct {
	index    int
	fileName string
}

func splitAndSend(inputFileName string, chunkChannel chan<- *chunkData, wg *sync.WaitGroup) {
	defer close(chunkChannel)
	defer wg.Done()

	file, err := os.Open(inputFileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}

	chunkSize := (fileInfo.Size() + int64(numChunks) - 1) / int64(numChunks)

	for i := 0; i < numChunks; i++ {
		chunkFileName := fmt.Sprintf("chunk%d.pcap", i)
		chunkFile, err := os.Create(chunkFileName)
		if err != nil {
			fmt.Println("Error creating chunk file:", err)
			return
		}

		bytesToCopy := chunkSize
		if i == numChunks-1 {
			bytesToCopy = fileInfo.Size() - int64(i)*chunkSize
		}

		_, err = io.CopyN(chunkFile, file, bytesToCopy)
		if err != nil {
			chunkFile.Close()
			fmt.Println("Error copying chunk:", err)
			return
		}

		chunkFile.Close()

		chunkChannel <- &chunkData{index: i, fileName: chunkFileName}
	}
}

func encryptAndSend(chunkChannel <-chan *chunkData, encryptedChunkChannel chan<- *chunkData, wg *sync.WaitGroup) {
	defer wg.Done()

	for chunk := range chunkChannel {
		encryptedChunkFileName := fmt.Sprintf("encrypted_chunk%d.pcap", chunk.index)

		// Simulate encryption, replace with your encryption logic
		err := simulateEncryption(chunk.fileName, encryptedChunkFileName)
		if err != nil {
			fmt.Printf("Error encrypting chunk %d: %v\n", chunk.index, err)
			return
		}

		encryptedChunkChannel <- &chunkData{index: chunk.index, fileName: encryptedChunkFileName}
	}
}

func combineAndWrite(outputFileName string, encryptedChunkChannel <-chan *chunkData, wg *sync.WaitGroup) {
	defer wg.Done()

	outputFile, err := os.Create(outputFileName)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	// Map to store encrypted chunk data temporarily
	encryptedChunks := make(map[int]string)

	// Receive encrypted chunks and store in the map
	for i := 0; i < numChunks; i++ {
		encryptedChunk := <-encryptedChunkChannel
		encryptedChunks[encryptedChunk.index] = encryptedChunk.fileName
	}

	// Combine encrypted chunks and write to the output file
	for i := 0; i < numChunks; i++ {
		encryptedChunkFileName := encryptedChunks[i]
		encryptedChunkFile, err := os.Open(encryptedChunkFileName)
		if err != nil {
			fmt.Printf("Error opening encrypted chunk file %d: %v\n", i, err)
			return
		}
		defer encryptedChunkFile.Close()

		_, err = io.Copy(outputFile, encryptedChunkFile)
		if err != nil {
			fmt.Printf("Error copying encrypted chunk %d: %v\n", i, err)
			return
		}
	}
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
