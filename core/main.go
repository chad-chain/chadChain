package main

import (
	"fmt"

	n "github.com/malay44/chadChain/core/network"
	t "github.com/malay44/chadChain/core/types"
)

func main() {
	// n.Http()
	n.Rpc()

	bbc := t.Block{}
	bbch := t.Header{}
	bbct := t.Transaction{}
	bbca := t.Account{}
	bbch.CreateHeader([32]byte{}, [20]byte{}, [32]byte{}, [32]byte{}, 0, 0, 0, 0, []byte{}, 0)
	bbct.CreateTransaction([20]byte{}, 0, 0, nil, nil, nil)
	bbca.CreateAccount([20]byte{}, 0, 0)

	fmt.Println(bbch)
	fmt.Println(bbct)
	fmt.Println(bbca)

	bbc.CreateBlock(bbch, []t.Transaction{bbct})
	fmt.Println(bbc)

	fmt.Println("Hello, world!")

}
