package main

import (
	"bytes"
	"encoding/ascii85"
)

func decode85(cookieValue string) string {
	cookieEncodedBytes := []byte(cookieValue)
	cookieDecodedBytes := make([]byte, len(cookieEncodedBytes))
	nCookieDecodedBytes, _, _ := ascii85.Decode(cookieDecodedBytes, cookieEncodedBytes, true)
	cookieDecodedBytes = cookieDecodedBytes[:nCookieDecodedBytes]

	//ascii85 adds /x00 null bytes at the end
	cookieDecodedBytes = bytes.Trim(cookieDecodedBytes, "\x00")
	return string(cookieDecodedBytes)
}

func encode85(cookieValue string) string {
	cookieBytes := []byte(cookieValue)
	cookieEncodedb85Bytes := make([]byte, ascii85.MaxEncodedLen(len(cookieBytes)))
	_ = ascii85.Encode(cookieEncodedb85Bytes, cookieBytes)
	cookieEncodedString := string(cookieEncodedb85Bytes)
	return cookieEncodedString
}
