// Package inputvalidation centralizes the normalization rules applied before
// user-provided values reach the domain or the database.
package inputvalidation

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

var (
	namePattern  = regexp.MustCompile(`^[\p{L}][\p{L}\p{M}'’.-]*(?: [\p{L}][\p{L}\p{M}'’.-]*)*$`)
	emailPattern = regexp.MustCompile(`^[A-Za-z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[A-Za-z0-9](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])?(?:\.[A-Za-z0-9](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])?)+$`)
)

// SingleLine removes control characters and collapses incidental whitespace.
// It does not HTML-escape text: API values are data, and output contexts must
// escape it correctly (React already does this for text nodes).
func SingleLine(value string) string {
	value = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, value)
	return strings.Join(strings.Fields(value), " ")
}

func Name(value string) (string, error) {
	value = SingleLine(value)
	if len([]rune(value)) < 2 || len([]rune(value)) > 80 || !namePattern.MatchString(value) {
		return "", errors.New("el nombre debe tener entre 2 y 80 caracteres y solo usar letras, espacios, apóstrofes, puntos o guiones")
	}
	return value, nil
}

func Email(value string) (string, error) {
	value = strings.ToLower(SingleLine(value))
	if len(value) > 254 || !emailPattern.MatchString(value) {
		return "", errors.New("correo electrónico inválido")
	}
	return value, nil
}

// Phone accepts common visual separators, stores a canonical digits-only form
// (optionally prefixed with +), and enforces the E.164 practical size range.
func Phone(value string) (string, error) {
	value = strings.TrimSpace(value)
	hasPlus := strings.HasPrefix(value, "+")
	var digits strings.Builder
	for _, r := range value {
		if unicode.IsDigit(r) {
			digits.WriteRune(r)
			continue
		}
		if !strings.ContainsRune(" +-().", r) {
			return "", errors.New("teléfono inválido")
		}
	}
	normalized := digits.String()
	if len(normalized) < 7 || len(normalized) > 15 {
		return "", errors.New("el teléfono debe contener entre 7 y 15 dígitos")
	}
	if hasPlus {
		normalized = "+" + normalized
	}
	return normalized, nil
}

func Password(value string) error {
	if len(value) < 8 || len(value) > 72 {
		return errors.New("la contraseña debe tener entre 8 y 72 caracteres")
	}
	var hasLower, hasUpper, hasDigit bool
	for _, r := range value {
		hasLower = hasLower || unicode.IsLower(r)
		hasUpper = hasUpper || unicode.IsUpper(r)
		hasDigit = hasDigit || unicode.IsDigit(r)
	}
	if !hasLower || !hasUpper || !hasDigit {
		return errors.New("la contraseña debe incluir una mayúscula, una minúscula y un número")
	}
	return nil
}

// LoginPassword only enforces a safe transport size. Existing accounts created
// before the stronger registration policy must still be able to authenticate.
func LoginPassword(value string) error {
	if len(value) == 0 || len(value) > 72 {
		return errors.New("credenciales inválidas")
	}
	return nil
}

func Service(value string) (string, error) {
	value = SingleLine(value)
	if len([]rune(value)) < 2 || len([]rune(value)) > 80 {
		return "", errors.New("tipo de servicio inválido")
	}
	return value, nil
}

// Notes retains intentional line breaks but removes null/control characters.
func Notes(value string, max int) (string, error) {
	value = strings.ReplaceAll(value, "\r\n", "\n")
	value = strings.ReplaceAll(value, "\r", "\n")
	value = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && r != '\n' && r != '\t' {
			return -1
		}
		return r
	}, value)
	value = strings.TrimSpace(value)
	if len([]rune(value)) > max {
		return "", errors.New("las notas no pueden superar el límite permitido")
	}
	return value, nil
}
