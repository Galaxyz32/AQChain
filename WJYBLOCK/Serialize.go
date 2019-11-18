package main

import (
	msg "AQChain/WJYBLOCK/pb"
	"AQChain/models"
	"bytes"
	"encoding/gob"
	"github.com/golang/protobuf/proto"
	"log"
)

//mn已用protobuf重写
func SerializeMerkleNode(m models.MerkleNode) []byte {

	mn := msg.MerkleNode{
		UserName:           m.UserName,
		User2Name:          m.User2name,
		ValueOfMerkleNode:  m.ValueOfMerkleNode,
		ContentHash:        m.ContentHash,
		TimeStampOfContent: m.TimeStampOfContent,
		TypeOfThis:         int32(m.TypeOfThis),
	}

	data, err := proto.Marshal(&mn)
	if err != nil {
		log.Fatalln("Marshal data error: ", err)
	}
	return data
}
func DeserializeMerkleNode(d []byte) *models.MerkleNode {
	var m models.MerkleNode
	var target msg.MerkleNode

	err = proto.Unmarshal(d, &target)
	if err != nil {
		log.Fatalln("Unmarshal data error: ", err)
	}
	m.TypeOfThis = int(target.TypeOfThis)
	m.TimeStampOfContent = target.TimeStampOfContent
	m.ContentHash = target.ContentHash
	m.ValueOfMerkleNode = target.ValueOfMerkleNode
	m.User2name = target.User2Name
	m.UserName = target.UserName

	//decoder := gob.NewDecoder(bytes.NewReader(d))
	//err := decoder.Decode(&m)
	//if err != nil {
	//	log.Panic(err)
	//}

	return &m
}

//File暂未用到
func SerializeFile(f models.FileInBlockchain) []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(f)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}
func DeserializeFile(d []byte) *models.FileInBlockchain {
	var f models.FileInBlockchain

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&f)
	if err != nil {
		log.Panic(err)
	}

	return &f
}

//匿名字段不能序列化

func SerializeBlockBody(b []models.MerkleNode) []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

func SerializeBlock(b models.Block) []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}
func DeserializeBlock(d []byte) *models.Block {
	var b models.Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&b)
	if err != nil {
		log.Panic(err)
	}
	return &b
}

//User已用Protobuf重写
func SerializeUser(u models.User) []byte {
	user := msg.User{
		UserName:      u.UserName,
		UserAddTime:   u.UserAddTime,
		UserFileNume:  int32(u.UserFileNume),
		Contribution:  u.Contribution,
		Contribution1: u.Contribution1,
		Contribution2: u.Contribution2,
		Contribution3: u.Contribution3,
		Useronlion:    u.Useronlion,
	}

	//var result bytes.Buffer

	result, err := proto.Marshal(&user)
	if err != nil {
		log.Fatalln("Marshal data error: ", err)
	}
	//encoder := gob.NewEncoder(&result)
	//
	//err := encoder.Encode(u)
	//if err != nil {
	//	log.Panic(err)
	//}
	return result
}
func DeserializeUser(d []byte) *models.User {
	var u models.User
	var target msg.User

	err = proto.Unmarshal(d, &target)
	if err != nil {
		log.Fatalln("Unmarshal data error: ", err)
	}

	u.Useronlion = target.Useronlion
	u.Contribution3 = target.Contribution3
	u.Contribution = target.Contribution
	u.Contribution2 = target.Contribution2
	u.Contribution1 = target.Contribution1
	u.UserFileNume = int(target.UserFileNume)
	u.UserAddTime = target.UserAddTime
	u.UserName = target.UserName

	//decoder := gob.NewDecoder(bytes.NewReader(d))
	//err := decoder.Decode(&u)
	//if err != nil {
	//	log.Panic(err)
	//}

	return &u
}
