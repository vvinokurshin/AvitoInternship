package pkg

import (
	"math/rand"
	"time"
)

func PercentageIDs(IDs []uint64, percent int) []uint64 {
	newCount := len(IDs) * percent / 100
	newIDs := make([]uint64, newCount)

	rand.Seed(time.Now().UnixNano())
	randIndexes := rand.Perm(len(IDs))[:newCount]

	for i, index := range randIndexes {
		newIDs[i] = IDs[index]
	}

	return newIDs
}
