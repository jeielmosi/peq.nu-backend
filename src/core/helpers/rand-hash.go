package helpers

import (
	"math/rand"
	"strings"
	"time"
)

// TODO: verify if could increase alphabet
const (
	FIRST_CHAR_ALPHABET  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ-0123456789abcdefghijklmnopqrstuvwxyz"
	FIRST_CHAR_FORBIDDEN = "_"
	ALPHABET             = FIRST_CHAR_ALPHABET + FIRST_CHAR_FORBIDDEN
	HASH_SIZE            = 7
)

func NewRandomHash(size uint) string {
	if size == 0 {
		return ""
	}

	rnd := rand.New(
		rand.NewSource(
			time.Now().UTC().UnixNano(),
		),
	)

	var sb strings.Builder

	idx := rnd.Intn(len(FIRST_CHAR_ALPHABET))
	r := rune(FIRST_CHAR_ALPHABET[idx])
	sb.WriteRune(r)

	for i := uint(1); i < size; i++ {
		idx = rnd.Intn(len(ALPHABET))
		r = rune(ALPHABET[idx])
		sb.WriteRune(r)
	}

	return sb.String()
}
