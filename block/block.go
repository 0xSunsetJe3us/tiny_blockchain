// Package block provides struct of a block of the blockchains.
package block

import (
	"hash"
	"time"
	"tiny_blockchain/consensus"
	"tiny_blockchain/crypto"
)

type Block struct {
	Head        *blockHead
	Body        *blockBody
	HashHandler func() hash.Hash
}

type blockHead struct {
	timeStamp      int64  // timestamp 		时间戳
	merkelRootHash []byte // merkelRootHash 默克尔树根hash值, 32bits
	prevBlockHash  []byte // prevBlockHash 	前一个区块Hash值, 32bits
	bits           []byte // bits 			目标targetHash的压缩信息, 32bits
	nonce          int32  // nonce 			随机数(或者叫bits), 挖矿找的值, 32bits
}

type blockBody struct {
	data [][]byte // data 区块内所有的数据
}

// NewSingleBlock 构造一个新的区块
func NewSingleBlock(hashType string, data [][]byte, prevBlockHash []byte) (*Block, error) {
	b := new(Block)
	b.HashHandler = crypto.NewHashHandler(hashType)
	b.Head = b.buildBlockHead(prevBlockHash, data)
	b.Body = b.buildBlockBody(data)
	return nil, nil
}

func (b *Block) buildBlockHead(prevBlockHash []byte, data [][]byte) *blockHead {
	b.Head.timeStamp = time.Now().Unix()
	b.Head.merkelRootHash = consensus.NewMerkelTree(b.HashHandler, data)
	b.Head.prevBlockHash = prevBlockHash
	// TODO: bits and nonce
	return nil
}

func (b *Block) buildBlockBody(data [][]byte) *blockBody {
	return nil
}
