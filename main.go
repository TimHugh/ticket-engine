package main

import (
	"fmt"

	"github.com/timhugh/ticket-engine/common"
)

const mongodb_uri string = "mongodb://devbox/tickets"

func main() {
	adapter, err := common.NewMongoAdapter(mongodb_uri)
	check(err)
	defer adapter.Close()

	repo := common.LocationRepository{adapter}

	storeLoc := common.Location{
		ID:           "test_id",
		SignatureKey: "test_signature",
	}
	err = repo.Store(storeLoc)
	check(err)

	findLoc, findErr := repo.Find("test_id")
	if findErr != nil {
		panic(findErr)
	} else {
		fmt.Printf("Found location id '%s' with signature '%s'", findLoc.ID, findLoc.SignatureKey)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
