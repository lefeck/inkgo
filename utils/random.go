package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06d", r.Intn(1000000))
}
