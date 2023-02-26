// Package consensus provides consensus algorithm.
package consensus

import (
	"bytes"
	"container/list"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"math"
	"strings"
)

type MerkelTree struct {
	root        *node            // root 默克尔根
	leaves      []*node          // leaves 默克尔叶子节点
	hashHandler func() hash.Hash // hashHandler 加密函数
}

type node struct {
	isLeaf    bool   // isLeaf 叶子节点标签
	isSingle  bool   // isSingle 是否只有一个叶子
	left      *node  // left 左叶子
	right     *node  // right 右叶子
	parent    *node  // parent 父节点
	data      []byte // data 如果是叶子就是叶子内容；否则为nil
	hashValue []byte // hashValue 该节点的hash值
}

// NewMerkelTree 创建一个新的默克尔树
func NewMerkelTree(hashFunc func() hash.Hash, data [][]byte) []byte {
	merkelTree := new(MerkelTree)
	merkelTree.hashHandler = hashFunc
	merkelTree.root, merkelTree.leaves = merkelTree.buildMerkelTree(data)
	return merkelTree.root.hashValue
}

// GetMerkelRootHashValue 获取默克尔根
func (m *MerkelTree) GetMerkelRootHashValue() []byte {
	return m.root.hashValue
}

// VerifyTree  验证默克尔树
func VerifyTree() bool {
	return true
}

// VerifyData 验证数据是否默克尔树中
func (m *MerkelTree) VerifyData(data []byte) (bool, error) {
	dataHash := m.callHashHandler(data)
	for _, leaf := range m.leaves {
		if bytes.Compare(dataHash, leaf.hashValue) == 0 {
			return true, nil
		}
	}
	return false, nil
}

func (m *MerkelTree) buildMerkelTree(data [][]byte) (*node, []*node) {
	leaves := m.buildTreeLeaves(data)
	root := m.buildMerkelTreeNode(leaves)
	return root, leaves
}

func (m *MerkelTree) buildTreeLeaves(data [][]byte) []*node {
	var leaves []*node
	for _, d := range data {
		n := &node{
			isLeaf:    true,
			isSingle:  true,
			left:      nil,
			right:     nil,
			parent:    nil,
			data:      d,
			hashValue: nil,
		}
		leaves = append(leaves, n)
	}
	return leaves
}

func (m *MerkelTree) buildMerkelTreeNode(leaves []*node) *node {
	// 默克尔树是二叉平衡树, 两两节点hash, 如果只有一个节点, 就hash自己即可
	var parentsNodes []*node
	hashPairCnt := int(math.Ceil(float64(len(leaves) / 2)))
	for i := 0; i < hashPairCnt; i++ {
		singleTag := false
		leftNode := leaves[2*i]
		rightNode := leftNode
		hashVal := leftNode.data
		if len(leaves) < 2*i+1 { // skills: 二叉树用 2idx/2idx+1 这种方式来便捷取到左右孩子
			singleTag = true
			rightNode = leaves[2*i+1]
		}
		if singleTag {
			hashVal = append(hashVal, rightNode.data...)
		}
		n := &node{
			parent:    nil,
			left:      leftNode,
			right:     rightNode,
			isSingle:  singleTag,
			isLeaf:    false,
			data:      nil,
			hashValue: m.calHash(hashVal),
		}
		parentsNodes = append(parentsNodes, n)
	}
	// recursive
	if len(parentsNodes) > 1 {
		return m.buildMerkelTreeNode(parentsNodes)
	}
	return parentsNodes[0]
}

func (m *MerkelTree) calHash(data []byte) []byte {
	handler := m.hashHandler()
	handler.Write(data)
	return handler.Sum(nil)
}

func (m *MerkelTree) buildMerkelTreeLeaves(data [][]byte) []*node {
	var leaves []*node
	for _, item := range data {
		n := &node{parent: nil, right: nil, left: nil, isLeaf: true, isSingle: false, hashValue: m.callHashHandler(item), data: item}
		leaves = append(leaves, n)
	}
	return leaves
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

// PrintWholeTree @brief: 打印整颗树的info, post-loop
func (m *MerkelTree) PrintWholeTree() {
	cnt := 0
	nextCnt := 1
	height := 1
	// use bfs
	queue := list.New()
	queue.PushBack(m.root)
	for queue.Len() != 0 {
		e := queue.Front()
		queue.Remove(e)
		n, _ := e.Value.(*node)

		cnt += 1
		if cnt%nextCnt == 0 {
			fmt.Printf("-- The-%d-level --\n", height)
			nextCnt = int(math.Exp2(float64(height)))
			height += 1
		}

		if n.left != nil || n.right != nil {
			fmt.Printf("[Parents] data: %s, hash: %x, left: %x, right: %x\n", string(n.data), n.hashValue, n.left.hashValue, n.right.hashValue)
		} else {
			fmt.Printf("[Leaf] data: %s, hash: %x, left: %s, right: %s\n", string(n.data), n.hashValue, "null", "null")
		}

		if n.left != nil {
			queue.PushBack(n.left)
		}
		if n.right != nil {
			queue.PushBack(n.right)
		}
	}
}

//@brief: 重新计算每一个节点的hash, 后续遍历
func (n *node) verifyNode(m *MerkelTree) ([]byte, error) {
	if n.isLeaf {
		return m.callHashHandler(n.data), nil
	}
	leftHash, _ := n.left.verifyNode(m)
	rightHash, _ := n.right.verifyNode(m)
	return m.callHashHandler(append(leftHash, rightHash...)), nil
}
