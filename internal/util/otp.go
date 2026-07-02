package util

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
)

func GenerateOTP(digits int) (string, error) {
	bi, err := rand.Int(
		rand.Reader,
		big.NewInt(int64(math.Pow(10, float64(digits)))),
	)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%0*d", digits, bi), nil
}
