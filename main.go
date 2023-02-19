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
	merkelTree.PrintWholeTree()

	// case-1: 验证整颗默克尔树
	if ok, _ := merkelTree.VerifyTree(); ok {
		fmt.Println("MerkelTree verify success!")
	} else {
		fmt.Println("MerkelTree verify failed!")
	}

	// case-2: 验证某个data是否在树中
	testData := []string{"\"test data 1\"", "\"test data NOT-exists\""}
	for _, testCase := range testData {
		if ok, _ := merkelTree.VerifyData([]byte(testCase)); ok {
			fmt.Printf("case:[%s] is in merkel-tree\n", testCase)
		} else {
			fmt.Printf("case:[%s] is not in merkel-tree\n", testCase)
		}
	}
}
