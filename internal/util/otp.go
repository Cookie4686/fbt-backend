package util

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
)

const digitLen = 10

func GenerateOTP(digits int) (string, error) {
	bi, err := rand.Int(
		rand.Reader,
		big.NewInt(int64(math.Pow(digitLen, float64(digits)))),
	)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%0*d", digits, bi), nil
}
