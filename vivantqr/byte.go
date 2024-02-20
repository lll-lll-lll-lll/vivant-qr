package main

import "fmt"

func bytesToOctalString(bytes []byte) string {
	var octalString string

	for _, b := range bytes {
		// 8進数で3桁になるようにフォーマット
		octalString += fmt.Sprintf("%03o", b)
	}

	return octalString
}

func octalStringToBytes(octalString string) ([]byte, error) {
	var bytes []byte

	for i := 0; i < len(octalString); i += 3 {
		// 残りの文字が3文字未満の場合の対処
		endIndex := i + 3
		if endIndex > len(octalString) {
			endIndex = len(octalString)
		}

		octalByte := octalString[i:endIndex]
		var b int
		_, err := fmt.Sscanf(octalByte, "%03o", &b)
		if err != nil {
			return nil, err
		}
		bytes = append(bytes, byte(b))
	}

	return bytes, nil
}
