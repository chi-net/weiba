package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/chi-net/weiba/core/types"
	"math"
	"math/rand"
	"strconv"
	"time"
)

func Decode(from string) uint64 {
	result := uint64(0)
	strlist := "abcdefghijklmnopqrstuvwxyz!@#$ABCDEFGHIJKLMNOPQRSTUVWXYZ%^&*1234567890()-=_+[]{}|\\:;<>?,./`~"

	// Create a map for quick lookup of character indices
	indexMap := make(map[rune]int)
	for i, c := range strlist {
		indexMap[c] = i
	}

	for _, char := range from {
		index, exists := indexMap[char]
		if !exists {
			// fmt.Printf("Character '%c' not found in strlist\n", char)
			return 0 // or handle the error appropriately
		}
		// fmt.Print(index, " ")                                // Print the index for debugging
		result = result*uint64(len(strlist)) + uint64(index) // Update the result
	}
	// fmt.Println(result)
	// fmt.Println() // Print a newline after indices
	return result
}

func Check(msg string, encoded string, config types.YmlConfigurationData) bool {
	// Convert string to a slice of runes
	var runes []rune
	for _, r := range msg {
		runes = append(runes, r)
	}
	msg = string(runes)
	// sha1
	h := sha1.New()
	h.Write([]byte(msg))
	hash := h.Sum(nil)
	hashHex := hex.EncodeToString(hash)
	//fmt.Printf("SHA-1 Hash (hex): %s\n", hashHex)
	hashInt, _ := strconv.ParseUint(hashHex[:16], 16, 64)
	modulus := uint64(math.Pow(92, float64(config.TransformDigits)))
	result := hashInt % modulus
	//fmt.Println(hashInt)
	//fmt.Println(result)
	return result == Decode(encoded)
}

func Generate(data types.ImportedGICAuthData) (int64, int64) {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source) // Create a new Rand instance
	num := r.Int63n(int64(len(data.Data)))

	source = rand.NewSource(time.Now().UnixNano())
	r = rand.New(source) // Create a new Rand instance
	// fmt.Println(len(data.Data[num].Data))
	num2 := r.Int63n(int64(len(data.Data[num].Data)))

	return num, num2
}
