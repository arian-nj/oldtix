package randutils

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"
)

func GenerateRandomString(length int) string {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.EncodeToString(randomBytes)[:length]
}

func GenerateRandomNumber(max_num int) int {
	max := big.NewInt(int64(max_num))
	random_index, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic(err)
	}
	return int(random_index.Int64())

}
