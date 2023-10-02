package proxy

import (
	"bytes"
	"encoding/hex"
	"log"
	"net/http"
	"strings"

	"github.com/andreburgaud/crypt2go/ecb"
	"golang.org/x/crypto/blowfish"
)

func decryptSubdomain(subDomain string) (host string, err error) {
	key := []byte("thisissecretkey")
	ciphertext, err := hex.DecodeString(subDomain)
	if err != nil {
		log.Println(err)
	}
	block, err := blowfish.NewCipher(key)
	if err != nil {
		log.Println(err)
	}
	mode := ecb.NewECBDecrypter(block)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)
	plaintext = PKCS5UnPadding(plaintext)
	host = string(plaintext)
	return host, err
}
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func EncryptSubdomain(host string) (subDomain string, err error) {
	key := []byte("thisissecretkey")
	block, err := blowfish.NewCipher(key)
	if err != nil {
		log.Println(err)
	}
	mode := ecb.NewECBEncrypter(block)
	plaintext := []byte(host)
	plaintext = PKCS5Padding(plaintext, block.BlockSize())
	ciphertext := make([]byte, len(plaintext))
	mode.CryptBlocks(ciphertext, plaintext)
	subDomain = hex.EncodeToString(ciphertext)
	return subDomain, err
}

func ParseAppHost(request *http.Request) string {
	host := request.Host
	subDomain := strings.Split(host, ".")[0]
	appHost, _ := decryptSubdomain(subDomain)
	appHost = strings.Split(appHost, string(rune(0)))[0]
	return appHost
}
