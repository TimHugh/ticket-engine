package main

import (
	"fmt"

	"github.com/timhugh/ticket-engine/common"
)

const mongodb_uri string = "http://devbox:3000/tickets"

func main() {
	adapter := common.NewMongoAdapter(mongdb_uri, mongodb_db)
	repo := common.LocationRepository{adapter}

	storeLoc := Location{
		ID:           "test_id",
		SignatureKey: "test_signature",
	}
	repo.Store(storeLoc)

	findLoc := repo.Find("test_id")
	fmt.Printf("Found location id '%s' with signature '%s'", findLoc.ID, findLoc.SignatureKey)
}
