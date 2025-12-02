package ascii

import (
	"crypto/sha256"
	"fmt"
	"os"
)

func CheckHash(fileName string) error {
	fileHash := "c3ec7584fb7ecfbd739e6b3f6f63fd1fe557d2ae3e24f870730d9cf8b2559e94"
	winHash := "56d0071a1d7439793953dae6ab3086e1ba4f2947028bc3d6ac4ec475956dff62"
	if fileName == "banners\\shadow.txt" || fileName == "banners/shadow.txt" {
		fileHash = "78ccd616680eb9068fe1465db1c852ceaffd8c0f318e3aa0414e1635508e85bf"
		winHash = "eb4b49a4abe7496fe728829d37a453c4148bac6f5ce771b181fd46407fb35077"
	} else if fileName == "banners\\thinkertoy.txt" || fileName == "banners/thinkertoy.txt" {
		fileHash = "e3c7a11f41a473d9b0d3bf2132a8f6dabb754bd16efa3897fa835a432d3b9caa"
		winHash = "242fdef5fad0fe9302bad1e38f0af4b0f83d086e49a4a99cdf0e765785640666"
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("could not read %s: %w", fileName, err)
	}

	newHash := sha256.Sum256(data)
	newHashStr := fmt.Sprintf("%x", newHash)
	if newHashStr == fileHash || newHashStr == winHash {
		return nil
	}
	return fmt.Errorf("invalid banner hash for %s", fileName)
}
