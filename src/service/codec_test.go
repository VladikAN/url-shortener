package service

import (
	"testing"

	"github.com/vladikan/url-shortener/config"
)

const expectedIdx = 1048585
const expectedCode = "eyWP"
const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func TestEncodeToString(t *testing.T) {
	cfg := config.ServiceSettings{Chars: chars}
	rst := Encode(expectedIdx, &cfg)
	if rst != expectedCode {
		t.Errorf("Encode gives `%s` and expected is `%s`", rst, expectedCode)
	}
}

func TestDecodeToString(t *testing.T) {
	cfg := config.ServiceSettings{Chars: chars}
	rst, _ := Decode(expectedCode, &cfg)
	if rst != expectedIdx {
		t.Errorf("Decode gives `%d` and expected is `%d`", rst, expectedIdx)
	}
}
