package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

var PubKey []byte
var PriKey []byte

var PubKeyMap map[string]string

func RSAKeyMap() {

	//写上自己的公钥
	PubKeyMap[strconv.Itoa(IDname)] = string(PubKey)
	//观察PubKeyMap
	//go func() {
	//	for {
	//		fmt.Printf("\x1b[32m")
	//		spew.Dump(PubKeyMap)
	//		fmt.Printf("\x1b[0m")
	//		time.Sleep(10e9)
	//
	//	}
	//}()
}

func RSAKeyGenInit() {
	//生成RSA密钥
	err := GenRsaKey(1024)
	if err != nil {
		panic(err)
	}
	PriKey = InitKey(WorkSpace + strconv.Itoa(IDname) + " S.pem")
	PubKey = InitKey(WorkSpace + strconv.Itoa(IDname) + " P.pem")
}

func rsamain() {
	//str:="helloworld"

	IDname = 1
	//IDname=2
	WorkSpace = CaculatePath()

	err := GenRsaKey(1024)
	if err != nil {
		panic(err)
	}

	PriKey := InitKey(WorkSpace + strconv.Itoa(IDname) + " S.pem")
	PubKey := InitKey(WorkSpace + strconv.Itoa(IDname) + " P.pem")

	var theMsg = "the message you want to encode 你好 世界"
	fmt.Println("Source:", theMsg)
	//私钥签名
	sig, _ := RsaSign([]byte(theMsg), PriKey)

	RSAVerify(theMsg, sig, PubKey)

}

//RSA公钥私钥产生
func GenRsaKey(bits int) error {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	file, err := os.Create(WorkSpace + strconv.Itoa(IDname) + " S.pem")
	if err != nil {
		return err
	}

	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}

	file, err = os.Create(WorkSpace + strconv.Itoa(IDname) + " P.pem")
	if err != nil {
		return err
	}

	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	return nil
}

func RSAVerify(theMsg string, sig []byte, PubKey []byte) bool {
	//var theMsg = "the message you want to encode 你好 世界"
	//fmt.Println("Source:", theMsg)
	//私钥签名
	//sig, _ := RsaSign([]byte(theMsg))

	fmt.Println(string(sig))
	//公钥验证
	if RsaSignVer([]byte(theMsg), sig, PubKey) != nil {
		fmt.Println("验证失败")
		return false
	} else {
		fmt.Println("验证通过")
		return true

	}

	////公钥加密
	// enc, _ := RsaEncrypt([]byte(theMsg))
	//  fmt.Println("Encrypted:", string(enc))
	//  //私钥解密
	//  decstr, _ := RsaDecrypt(enc)
	//  fmt.Println("Decrypted:", string(decstr))
}

func InitKey(path string) []byte {
	PubK, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	//num,_:=PriK.Read(PriKey)
	PubKey, err := ioutil.ReadAll(PubK)
	//fmt.Println(num)
	PubK.Close()
	return PubKey
}

//私钥签名
func RsaSign(data []byte, PriKey []byte) ([]byte, error) {
	h := sha256.New()
	h.Write(data)
	hashed := h.Sum(nil)
	//获取私钥
	block, _ := pem.Decode(PriKey)
	if block == nil {
		return nil, errors.New("private key error")
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hashed)
}

//公钥验证
func RsaSignVer(data []byte, signature []byte, PubKey []byte) error {
	hashed := sha256.Sum256(data)
	block, _ := pem.Decode(PubKey)
	if block == nil {
		return errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//验证签名
	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashed[:], signature)
}

// 公钥加密
func RsaEncrypt(data []byte, PubKey []byte) ([]byte, error) {
	//解密pem格式的公钥
	block, _ := pem.Decode(PubKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	return rsa.EncryptPKCS1v15(rand.Reader, pub, data)
}

// 私钥解密
func RsaDecrypt(ciphertext []byte, PriKey []byte) ([]byte, error) {
	//获取私钥
	block, _ := pem.Decode(PriKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 解密
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}
