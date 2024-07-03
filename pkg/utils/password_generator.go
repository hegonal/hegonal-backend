package utils

import "golang.org/x/crypto/bcrypt"

func NormalizePassword(p string) []byte {
	return []byte(p)
}

func GeneratePassword(p string) string {
	bytePwd := NormalizePassword(p)

	hash, err := bcrypt.GenerateFromPassword(bytePwd, bcrypt.MinCost)
	if err != nil {
		return err.Error()
	}

	return string(hash)
}

func ComparePasswords(hashedPwd, inputPwd string) bool {
	byteHash := NormalizePassword(hashedPwd)
	byteInput := NormalizePassword(inputPwd)

	if err := bcrypt.CompareHashAndPassword(byteHash, byteInput); err != nil {
		return false
	}

	return true
}
