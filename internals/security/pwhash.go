package internal_security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

func PasswordHashArgon2(password string) (string, error) {
	// Generate a random salt (16 bytes is common)
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Argon2id parameters (see recommendations below)
	time := uint32(1)           // iterations
	memory := uint32(64 * 1024) // 64 MiB
	threads := uint8(4)         // parallelism
	keyLen := uint32(32)        // 32-byte key

	// Derive the hash
	hash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)

	// Encode salt and hash to base64 for storage
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Format: $argon2id$v=19$m=65536,t=1,p=4$<salt>$<hash>
	encoded := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, memory, time, threads, b64Salt, b64Hash)

	return encoded, nil
}

func PasswordVerifyArgon2(password, encodedHash string) (bool, error) {
	// Parse the encoded hash
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, fmt.Errorf("invalid hash format")
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return false, err
	}
	if version != argon2.Version {
		return false, fmt.Errorf("incompatible version")
	}

	var memory, time uint32
	var threads uint8
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads); err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	// Re-derive the hash with the same parameters
	computedHash := argon2.IDKey([]byte(password), salt, time, memory, threads, uint32(len(hash)))

	// Constant-time comparison
	return subtle.ConstantTimeCompare(hash, computedHash) == 1, nil
}
