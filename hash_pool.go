package merkle-dag

import "hash"

type HashPool interface {
	Get() hash.Hash
}
