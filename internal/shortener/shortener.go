package shortener

import (
	"crypto/rand"
	"math/big"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const codeLength = 6

func GenerateCode() (string, error) {
	code := make([]byte, codeLength)
	max := big.NewInt(int64(len(alphabet)))

	for i := 0; i < codeLength; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		code[i] = alphabet[n.Int64()]
	}

	return string(code), nil
}