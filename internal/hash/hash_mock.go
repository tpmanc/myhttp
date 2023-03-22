package hash

type Mock struct{}

func (m Mock) GenerateHash(_ []byte) string {
	return "hash string"
}
