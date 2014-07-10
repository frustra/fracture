package protocol

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func GenerateKey(size int) []byte {
	key := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		return nil
	} else {
		return key
	}
}

type KeyPair struct {
	private *rsa.PrivateKey
	public  *rsa.PublicKey
}

func (kp KeyPair) Serialize() []byte {
	buf, err := x509.MarshalPKIXPublicKey(kp.public)
	if err != nil {
		return make([]byte, 0)
	} else {
		return buf
	}
}

func GenerateKeyPair(size int) (*KeyPair, error) {
	priv, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return nil, err
	} else {
		return &KeyPair{priv, &priv.PublicKey}, nil
	}
}

func DecryptRSABytes(buf []byte, keyPair *KeyPair) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, keyPair.private, buf)
}

func EncryptRSABytes(buf []byte, keyPair *KeyPair) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, keyPair.public, buf)
}

func ParsePublicKey(buf []byte) (*KeyPair, error) {
	pub, err := x509.ParsePKIXPublicKey(buf)
	if err != nil {
		return nil, err
	}
	pubKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("RSA Key is not valid")
	} else {
		return &KeyPair{nil, pubKey}, nil
	}
}

func CheckAuth(username, serverId string, keyPair *KeyPair, secret []byte) (string, error) {
	h := sha1.New()
	io.WriteString(h, serverId)
	h.Write(secret)
	h.Write(keyPair.Serialize())
	hash := h.Sum(nil)
	negative := (hash[0] & 0x80) == 0x80
	if negative {
		carry := true
		for i := len(hash) - 1; i >= 0; i-- {
			hash[i] = byte(^hash[i])
			if carry {
				carry = hash[i] == 0xff
				hash[i]++
			}
		}
	}

	token := strings.TrimLeft(hex.EncodeToString(hash), "0")
	if negative {
		token = "-" + token
	}

	resp, err := http.Get("https://sessionserver.mojang.com/session/minecraft/hasJoined?username=" + url.QueryEscape(username) + "&serverId=" + url.QueryEscape(string(token)))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	} else {
		bodyJson := make(map[string]interface{})
		err = json.Unmarshal(body, &bodyJson)
		if err != nil {
			return "", err
		} else if id, ok := bodyJson["id"].(string); ok && len(id) == 32 {
			return id[0:8] + "-" + id[8:12] + "-" + id[12:16] + "-" + id[16:20] + "-" + id[20:32], nil
		} else {
			return "", err
		}
	}
}
