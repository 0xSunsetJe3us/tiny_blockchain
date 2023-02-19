package main

import (
	"encoding/json"
	"fmt"
	"tiny_blockchain/merkel_tree"
)

func initData() [][]byte {
	var data [][]byte
	for i := 0; i < 4; i++ {
		str := fmt.Sprintf("test data %d", i)
		bz, _ := json.Marshal(str)
		data = append(data, bz)
	}
	return data
}

func main() {
	data := initData()                                      // 1. init data
	merkelTree, _ := merkel_tree.NewMerkelTree("md5", data) // 2. build a merkel tree
	//fmt.Printf("%x", merkel_tree.GetMerkelRootHashValue())
	merkelTree.PrintWholeTree()
}
