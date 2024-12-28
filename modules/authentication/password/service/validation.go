package service

import "unicode"

func (s *PasswordService) isPasswordStrong(password string) bool {
	if len(password) < 8 {
		return false
	}

	var hasSpecial bool
	for _, char := range password {
		if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			hasSpecial = true
			break
		}
	}

	return hasSpecial
}

func (s *PasswordService) isPasswordCommon(password string) bool {
	return s.commonPasswords[password]
}
