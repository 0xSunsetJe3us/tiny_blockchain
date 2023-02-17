package merkel_tree

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"hash"
	"math"
	"strings"
)

type MerkelTree struct {
	root        *node
	leaves      []*node
	hashHandler func() hash.Hash
}

type node struct {
	parent    *node
	left      *node
	right     *node
	leaf      bool // ?
	single    bool // ?
	hashValue []byte
	data      []byte
}

//@brief: 创建一个新的默克尔树
func NewMerkelTree(hashType string, data [][]byte) (*MerkelTree, error) {
	var err error
	merkelTree := &MerkelTree{}
	merkelTree.hashHandler = merkelTree.buildHash(hashType)
	merkelTree.root, merkelTree.leaves, err = merkelTree.buildMerkelTreeRoot(data)
	if err != nil {
		return nil, err
	}
	return merkelTree, nil
}

//@brief: 构造默克尔树根节点
func (m *MerkelTree) buildMerkelTreeRoot(data [][]byte) (*node, []*node, error) {
	if len(data) <= 0 {
		return nil, nil, errors.New("empty data")
	}
	leaves := m.buildMerkelTreeLeaves(data)
	root, err := m.buildMerkelTreeNode(leaves)
	if err != nil {
		return nil, nil, err
	}
	return root, leaves, nil
}

// @brief: 构造默克尔树非根父节点, 递归
func (m *MerkelTree) buildMerkelTreeNode(leaves []*node) (*node, error) {
	mergeNodesCnt := int(math.Ceil(float64(len(leaves) / 2))) // 要合并的个数
	// 两两hash成根节点, 如果只有一个，就增加一个节点, hash自己
	for i := 0; i < mergeNodesCnt; i++ {
		leftNode := leaves[i]
		rightNode := leaves[i+1]
	}

	return nil, nil
}

//@brief: 构造默克尔树叶节点, 将所有交易info构造成叶子
func (m *MerkelTree) buildMerkelTreeLeaves(data [][]byte) []*node {
	var leaves []*node
	for _, item := range data {
		n := &node{parent: nil, right: nil, left: nil, leaf: true, single: false, hashValue: m.callHashHandler(item), data: item}
		leaves = append(leaves, n)
	}
	return leaves
}

//@brief: 获取默克尔根
func (m *MerkelTree) GetMerkelRootHashValue() []byte {
	return m.root.hashValue
}

//@brief: 构造加密函数
func (m *MerkelTree) buildHash(hashType string) func() hash.Hash {
	switch strings.ToLower(hashType) {
	case "md5":
		return md5.New
	case "sha1":
		return sha1.New
	case "sha256":
		return sha256.New
	case "sha512":
		return sha512.New
	default:
		return sha1.New
	}
}

//@brief: 调用crypto函数，返回结果
func (m *MerkelTree) callHashHandler(data []byte) []byte {
	handler := m.hashHandler()
	handler.Write(data)
	return handler.Sum(nil)
}
