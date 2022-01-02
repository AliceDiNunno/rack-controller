package crypto

import "golang.org/x/crypto/bcrypt"

func HashAndSalt(token string) (string, error) {
	println("Hashing token:", token)
	//TODO: cost should be a variable
	hash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ComparePasswords(password string, comparison string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(comparison))

	if err != nil {
		return false, err
	}

	return true, err
}
