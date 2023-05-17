package authutil

import (
	"golang.org/x/crypto/bcrypt"
)

// EncodePassword encode a given password with bcrypt and cost.
// Use bcrypt.DefaultConfig if provide zero cost value.
func EncodePassword(password string, cost int) (string, error) {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// MatchesPassword returns a nil if matched hashedPassword and raw password, otherwise returns a error
func MatchesPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
