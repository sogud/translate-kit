package translatekit

import (
	"encoding/base64"
	"strings"
	"testing"
	"unicode/utf16"
	"unicode/utf8"
)

func TestYoudaoDecryptRejectsInvalidCiphertextAndPadding(t *testing.T) {
	session := youdaoSession{Key: "fixture-aes", IV: "fixture-iv"}
	for _, data := range [][]byte{nil, []byte("bad"), []byte(base64.URLEncoding.EncodeToString([]byte("short"))), []byte(base64.URLEncoding.EncodeToString(make([]byte, 16)))} {
		if _, err := decodeYoudao(data, session); err == nil {
			t.Fatal("Invalid ciphertext accepted")
		}
	}
}

func TestYoudaoLongTextSplitsWithoutLosingCharacters(t *testing.T) {
	for _, text := range []string{
		strings.Repeat("Read this sentence carefully. ", 6000),
		strings.Repeat("这是一段很长的说明。", 300),
		strings.Repeat("🙂", 1300),
		strings.Repeat("x", 2500),
	} {
		chunks := splitYoudaoText(text)
		if len(chunks) < 2 {
			t.Fatal("Long text was not split")
		}
		for _, chunk := range chunks {
			if !utf8.ValidString(chunk) || len(utf16.Encode([]rune(chunk))) > 1000 {
				t.Fatal("Chunk is invalid or exceeds the web text limit")
			}
		}
		withoutSpaces := func(value string) string {
			return strings.Join(strings.Fields(value), "")
		}
		if withoutSpaces(strings.Join(chunks, "")) != withoutSpaces(text) {
			t.Fatal("Text was lost while splitting")
		}
	}
	chunks := splitYoudaoText(strings.Repeat("a", 800) + ". " + strings.Repeat("b", 500))
	if !strings.HasSuffix(chunks[0], ".") || strings.HasPrefix(chunks[1], "a") {
		t.Fatal("Available sentence boundary was ignored")
	}
}
