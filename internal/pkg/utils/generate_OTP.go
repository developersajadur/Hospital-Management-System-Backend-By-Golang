package utils


import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

func GenerateOTP() string {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return fmt.Sprintf("%06d", 100000+randInt()%900000)
	}
	return fmt.Sprintf("%06d", n.Int64())
}

func randInt() int {
	return int(time.Now().UnixNano() % 1000000)
}
