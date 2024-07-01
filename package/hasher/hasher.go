package hasher

import "crypto/sha256"

type HasherSHA struct {
	salt string
}

func NewHasher(salt string) *HasherSHA { return &HasherSHA{salt: salt} }

func (h *HasherSHA) Hash(str string) (string, error) {
	hash := sha256.New()

	if _, err := hash.Write([]byte(str)); err != nil {
		return "", err
	}

	return string(hash.Sum([]byte(h.salt))), nil
}
