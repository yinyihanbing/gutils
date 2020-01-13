package goutils

import (
	"crypto/md5"
	"encoding/hex"
	"crypto/sha256"
	"hash"
	"crypto/hmac"
	"crypto"
	"encoding/pem"
	"errors"
	"crypto/rsa"
	"crypto/x509"
	"crypto/rand"
	"encoding/base64"
)

// MD5编码
func EncodeMd5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

// 加密密码
func EncodePassword(rawPwd string, salt string) string {
	pwd := PBKDF2([]byte(rawPwd), []byte(salt), 10000, 50, sha256.New)
	return hex.EncodeToString(pwd)
}

// http://code.google.com/p/go/source/browse/pbkdf2/pbkdf2.go?repo=crypto
func PBKDF2(password, salt []byte, iter, keyLen int, h func() hash.Hash) []byte {
	prf := hmac.New(h, password)
	hashLen := prf.Size()
	numBlocks := (keyLen + hashLen - 1) / hashLen

	var buf [4]byte
	dk := make([]byte, 0, numBlocks*hashLen)
	U := make([]byte, hashLen)
	for block := 1; block <= numBlocks; block++ {
		// N.B.: || means concatenation, ^ means XOR
		// for each block T_i = U_1 ^ U_2 ^ ... ^ U_iter
		// U_1 = PRF(password, salt || uint(i))
		prf.Reset()
		prf.Write(salt)
		buf[0] = byte(block >> 24)
		buf[1] = byte(block >> 16)
		buf[2] = byte(block >> 8)
		buf[3] = byte(block)
		prf.Write(buf[:4])
		dk = prf.Sum(dk)
		T := dk[len(dk)-hashLen:]
		copy(U, T)

		// U_n = PRF(password, U_(n-1))
		for n := 2; n <= iter; n++ {
			prf.Reset()
			prf.Write(U)
			U = U[:0]
			U = prf.Sum(U)
			for x := range U {
				T[x] ^= U[x]
			}
		}
	}
	return dk[:keyLen]
}

// 获取RSA签名值, privateKey=私钥, hash=哈希算法(crypto.MD5等), content=要进行签名计算的内容
func RsaGetSignByPKCS1Key(privateKey string, hash crypto.Hash, content string) (string, error) {
	blockPri, _ := pem.Decode([]byte(privateKey))
	if blockPri == nil {
		return "", errors.New("privateKey error")
	}

	var priKey *rsa.PrivateKey
	priKey, err := x509.ParsePKCS1PrivateKey(blockPri.Bytes)
	if err != nil {
		return "", err
	}

	h := hash.New()
	h.Write([]byte(content))
	hashed := h.Sum(nil)
	b, err := rsa.SignPKCS1v15(rand.Reader, priKey, hash, hashed)
	if err != nil {
		return "", err
	}

	sign := base64.StdEncoding.EncodeToString(b)

	return sign, nil
}

// 获取RSA签名值, privateKey=私钥, hash=哈希算法(crypto.MD5等), content=要进行签名计算的内容
func RsaGetSignByPKCS8Key(privateKey string, hash crypto.Hash, content string) (string, error) {
	blockPri, _ := pem.Decode([]byte(privateKey))
	if blockPri == nil {
		return "", errors.New("privateKey error")
	}

	prkI, err := x509.ParsePKCS8PrivateKey(blockPri.Bytes)
	if err != nil {
		return "", err
	}
	priKey := prkI.(*rsa.PrivateKey)

	h := hash.New()
	h.Write([]byte(content))
	hashed := h.Sum(nil)
	b, err := rsa.SignPKCS1v15(rand.Reader, priKey, hash, hashed)
	if err != nil {
		return "", err
	}

	sign := base64.StdEncoding.EncodeToString(b)

	return sign, nil
}

// RSA验证签名
func RsaCheckSign(publicKey string, hash crypto.Hash, content, sign string) error {
	signData, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return err
	}

	blockPri, _ := pem.Decode([]byte(publicKey))
	if blockPri == nil {
		return errors.New("publicKey error")
	}

	pub, err := x509.ParsePKIXPublicKey(blockPri.Bytes)
	if err != nil {
		return err
	}
	pubKey := pub.(*rsa.PublicKey)

	h := hash.New()
	h.Write([]byte(content))
	hashed := h.Sum(nil)
	err = rsa.VerifyPKCS1v15(pubKey, hash, hashed, signData)

	return err
}
