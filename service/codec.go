package service

import (
	"errors"
	"math"
	"strings"

	"github.com/vladikan/url-shortener/config"
)

// Encode will convert base10 to baseN (configured)
func Encode(num uint64, cfg *config.ServiceSettings) string {
	base := cfg.GetBase()

	b := make([]byte, 0)
	for num > 0 {
		r := math.Mod(float64(num), float64(base))
		num /= uint64(base)

		idx := int(r)
		ch := cfg.Chars[idx]
		b = append([]byte{ch}, b...)
	}

	return string(b)
}

// Decode will convert baseN (configured) to base10
func Decode(s string, cfg *config.ServiceSettings) (uint64, error) {
	base := cfg.GetBase()

	var r uint64
	ln := len(s)

	for i, v := range s {
		pow := ln - (i + 1)
		pos := strings.IndexRune(cfg.Chars, v)

		if pos == -1 {
			return 0, errors.New("invalid character: " + string(v))
		}

		r += uint64(pos) * uint64(math.Pow(float64(base), float64(pow)))
	}

	return uint64(r), nil
}
