package utils

import (
	"math/rand"
	"time"
)

func RandomString(n int) string {
	var lettters = []byte("qwertyuioplkjhgfdsazxcvbnmQWERTYUIOPLKJHGFDSAZXCVBNM")
	result := make([]byte, n)
	rand.Seed(time.Now().Unix())
	for i, _ := range result {
		result[i] = lettters[rand.Intn(len(lettters))]
	}
	return string(result)
}
