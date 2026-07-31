package hash

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Параметры Argon2id, подобранные под интерактивный вход (см. рекомендации OWASP)
const (
	argonTime    = 1
	argonMemory  = 64 * 1024
	argonThreads = 4
	argonKeyLen  = 32
	saltLen      = 16
)

var ErrInvalidHash = errors.New("hash: invalid encoded hash format")

// Hash хеширует пароль с помощью Argon2id и возвращает самодостаточную строку
// в формате $argon2id$v=19$m=...,t=...,p=...$<salt>$<hash>
func Hash(password string) (string, error) {
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	digest := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	encoded := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		argonMemory,
		argonTime,
		argonThreads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(digest),
	)

	return encoded, nil
}

// Verify сверяет пароль с ранее сохранённым хешем, сформированным Hash
func Verify(password, encoded string) (bool, error) {
	var (
		version                       int
		memory, timeCost, parallelism uint32
	)

	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false, ErrInvalidHash
	}

	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return false, ErrInvalidHash
	}
	if version != argon2.Version {
		return false, ErrInvalidHash
	}

	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &timeCost, &parallelism); err != nil {
		return false, ErrInvalidHash
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, ErrInvalidHash
	}

	want, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, ErrInvalidHash
	}

	if parallelism > math.MaxUint8 {
		return false, ErrInvalidHash
	}

	got := argon2.IDKey([]byte(password), salt, timeCost, memory, uint8(parallelism), uint32(len(want))) //nolint:gosec

	return subtle.ConstantTimeCompare(got, want) == 1, nil
}
