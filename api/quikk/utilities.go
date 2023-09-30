package quikk

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
)

func encrypt(key, secret, noww string) string {
	to_encode := fmt.Sprintf("date: %s", noww)
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(to_encode))
	buf := hash.Sum(nil)
	encoded := base64.StdEncoding.Strict().EncodeToString(buf)
	url_encoded := url.QueryEscape(encoded)
	return fmt.Sprintf(`keyId="%s",algorithm="hmac-sha256",signature="%s"`, key, url_encoded)
}
