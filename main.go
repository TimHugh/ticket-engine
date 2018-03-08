package main

import (
	"fmt"

	"github.com/timhugh/ticket-engine/common"
)

const mongodb_uri string = "mongodb://devbox/tickets"

func main() {
	adapter, err := common.NewMongoAdapter(mongodb_uri)
	if err != nil {
		panic(err)
	}
	repo := common.LocationRepository{adapter}

	storeLoc := common.Location{
		ID:           "test_id",
		SignatureKey: "test_signature",
	}
	repo.Store(storeLoc)

	findLoc, _ := repo.Find("test_id")
	fmt.Printf("Found location id '%s' with signature '%s'", findLoc.ID, findLoc.SignatureKey)
}
