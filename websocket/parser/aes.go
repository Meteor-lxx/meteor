package parser

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
	"wss/helper"
)

type AesPack struct {
	Data *helper.Package
	AesData string
}

func (aes AesPack)PackDecode() *helper.Package {
	data, _ := base64.StdEncoding.DecodeString(aes.AesData)
	salt := data[8:16]
	hashStr := []byte("DlClientPost2019")
	hash := md5.Sum(hashStr)
	key,iv := GetKeyIv(hash[:],salt)
	var pack helper.Package
	msg := AesDecryptCBC(data[16:],key,iv)
	_ = json.Unmarshal([]byte(msg), &pack)
	return &pack
}

func (aes AesPack)PackEncode() []byte  {
	data := []byte("DlClientPost2019")
	hash := md5.Sum(data)
	salt := make([]byte, 8)
	rand.Seed(time.Now().Unix())
	rand.Read(salt)
	key,iv := GetKeyIv(hash[:],salt)
	msg, _ := json.Marshal(aes.Data)
	encryptedData := AesEncryptCBC(msg,key,iv)
	var buffer bytes.Buffer
	buffer.Write([]byte("Salted__"))
	buffer.Write(salt)
	buffer.Write(encryptedData)
	msgStr := base64.StdEncoding.EncodeToString(buffer.Bytes())
	buffer.Reset()
	lenStr := len(msgStr)
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(lenStr))
	buffer.Write(buf)
	buffer.Write([]byte(msgStr))
	sendPack := buffer.Bytes()
	buffer.Reset()
	return sendPack
}

func AesEncryptCBC(origData []byte, key []byte,iv []byte) (encrypted []byte) {
	block, _ := aes.NewCipher(key)
	origData=Padding(origData,block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, iv)
	encrypted = make([]byte, len(origData))
	blockMode.CryptBlocks(encrypted, origData)
	base64.StdEncoding.EncodeToString(encrypted)
	return encrypted
}

func Padding(plainText []byte,blockSize int) []byte{
	//计算要填充的长度
	n:= blockSize-len(plainText)%blockSize
	//对原来的明文填充n个n
	temp:=bytes.Repeat([]byte{byte(n)},n)
	plainText=append(plainText,temp...)
	return plainText
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unhanding := int(src[length-1])
	if length -unhanding < 0 {
		return []byte("")
	}
	return src[:(length - unhanding)]
}

func AesDecryptCBC(encrypted []byte, key []byte,iv []byte) string {
	cipherBlock, err := aes.NewCipher(key)
	if err != nil{
		fmt.Println(err)
	}
	cipher.NewCBCDecrypter(cipherBlock, iv).CryptBlocks(encrypted, encrypted)
	return string(PKCS5UnPadding(encrypted))
}

func GetKeyIv(hash []byte,salt []byte) (key []byte, iv []byte) {
	var salted,dx []byte
	for len(salted) < 48 {
		h := md5.New()
		if len(dx)!=0 {
			h.Write(dx)
		}
		h.Write([]byte(fmt.Sprintf("%x",hash)))
		h.Write(salt)
		dx = h.Sum(nil)
		salted = append(salted, dx...)
	}
	key = salted[:32]
	iv = salted[32:48]
	return key,iv
}