package hash

import "errors"

// hashAlgo interface for different algorithm realizations
type hashAlgo interface {
	generate(msg []byte) string
}

// Hash service generates hash with defined algorithm
type Hash struct {
	hashAlgo hashAlgo
}

// New returns a new instance of Hash
func New(hashAlgo hashAlgo) (Hash, error) {
	if hashAlgo == nil {
		return Hash{}, errors.New("hash algorithm mustn't be nil")
	}

	return Hash{
		hashAlgo: hashAlgo,
	}, nil
}

// GenerateHash returns hash string made by defined algorithm
func (h Hash) GenerateHash(msg []byte) string {
	return h.hashAlgo.generate(msg)
}
