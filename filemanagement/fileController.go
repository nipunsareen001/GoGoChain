package filemanagement

import (
	"encoding/csv"
	"fmt"
	"os"
)

// OpenFileForAppending checks if the transaction file exists (creating it if necessary)
// and then opens it for appending. It returns the file and any encountered error.
func OpenFileForAppending(filename string) (*os.File, error) {
	// Check if the file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// File does not exist, attempt to create it
		if _, err := os.Create(filename); err != nil {
			return nil, fmt.Errorf("unable to create file: %v", err)
		}
	}

	// Open the file for appending
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	return file, nil
}

// AppendDataInFile writes the provided content to the given file and appends a newline.
func AppendDataInFile(file *os.File, content []byte) error {
	// Write the content to the file
	if _, err := file.Write(content); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	// Write a newline to the file
	if _, err := file.WriteString("\n"); err != nil {
		return fmt.Errorf("error writing newline to file: %v", err)
	}

	return nil
}

// IMPROVEMENTS, this func is currently used else where, in can be used from here for better project flow
// // ReadAllFiles reads all files from the keystore directory and matches by address.
// func readFilesFromKeyStore(addr string) (string, error) {
// 	files, err := os.ReadDir(string(enum.KeystorePath))
// 	if err != nil {
// 		return "", fmt.Errorf("error reading directory: %v", err)
// 	}
// 	for _, file := range files {
// 		if strings.Contains(file.Name(), addr) {
// 			return file.Name(), nil
// 		}
// 	}
// 	return "", fmt.Errorf("no file found for address: %s", addr)
// }

// IsAlreadyInFile checks if the given multiaddr string is already present in the file.
func IsAlreadyInFile(filename string, multiaddr string) (bool, error) {
	f, err := os.Open(filename)
	if err != nil {
		return false, fmt.Errorf("Error opening file %s: %v", filename, err)
	}
	defer f.Close()

	r := csv.NewReader(f)

	records, err := r.ReadAll()
	if err != nil {
		return false, fmt.Errorf("Error reading from file %s: %v", filename, err)
	}

	for _, record := range records {
		if len(record) > 0 && record[0] == multiaddr {
			return true, nil
		}
	}

	return false, nil
}
