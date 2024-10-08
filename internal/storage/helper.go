package storage

import "github.com/alexedwards/argon2id"

// HashUserPasswords hashes the username and password using argon2id
// return the hashed string of 98 characters
func HashUserPasswords(password, username string) (string, error) {

	// Parameters based on OWASP recommendation for argon2id.
	params := &argon2id.Params{
		Memory:      32 * 1024,
		Iterations:  2,
		Parallelism: 1,
		SaltLength:  16,
		KeyLength:   32,
	}
	hash, err := argon2id.CreateHash(password+username, params)
	if err != nil {
		return "", err
	}

	return hash, nil
}

// CompareProvidedPassword Compares the provided password with the hash.
func CompareProvidedPassword(password, username, hash string) (bool, error) {
	check, err := argon2id.ComparePasswordAndHash(password+username, hash)
	if err != nil {
		return false, err
	}

	return check, nil
}
