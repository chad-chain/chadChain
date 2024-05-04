package initialize

import (
	db "github.com/malay44/chadChain/core/storage"
	t "github.com/malay44/chadChain/core/types"
	"log"
)

func GlobalDBVar() {
	err := db.BadgerDB.View(db.Get([]byte("stateRootHash"), &t.StateRootHash))
	if err != nil {
		log.Default().Printf("Error in init StateRootHash\n")
		log.Default().Println(err.Error())
	} else {
		log.Default().Printf("StateRootHash: %x\n", t.StateRootHash)
	}
	err = db.BadgerDB.View(db.Get([]byte("latestBlock"), &t.LatestBlock))
	if err != nil {
		log.Default().Printf("Error in init latestBlock: \n")
		log.Default().Println(err.Error())
	} else {
		log.Default().Printf("LatestBlock: %x\n", t.LatestBlock)
	}
}
