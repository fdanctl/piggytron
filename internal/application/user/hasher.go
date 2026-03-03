package user

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

type PasswordHasher struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
	saltLen uint32
}

// type PasswordHasher interface {
// 	Hash(password string) (string, error)
// 	Verify(hash, password string) bool
// }

func NewPasswordHasher(
	time uint32,
	memory uint32,
	threads uint8,
	keyLen uint32,
	saltLen uint32,
) *PasswordHasher {
	return &PasswordHasher{
		time:    time,
		memory:  memory,
		threads: threads,
		keyLen:  keyLen,
		saltLen: saltLen,
	}
}

// hashes string as
// $argon2id$v=<version>$m=<memory>,t=<time>,p=<threads>$<salt>$<hash>
func (h *PasswordHasher) Hash(password string) (string, error) {
	salt := make([]byte, h.saltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		h.time,
		h.memory,
		h.threads,
		h.keyLen,
	)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		h.memory, h.time, h.threads,
		b64Salt,
		b64Hash,
	)

	return encoded, nil
}

func (h *PasswordHasher) Verify(hash, password string) (bool, error) {
	// $argon2id$v=<version>$m=<memory>,t=<time>,p=<threads>$<salt>$<hash>
	arr := strings.Split(hash, "$")

	mtp := make([]int, 3)
	for i, p := range strings.Split(arr[3], ",") {
		kv := strings.Split(p, "=")
		v, err := strconv.Atoi(kv[1])
		if err != nil {
			return false, err
		}
		mtp[i] = v
	}

	m := uint32(mtp[0])
	t := uint32(mtp[1])
	p := uint8(mtp[2])

	salt, err := base64.RawStdEncoding.DecodeString(arr[4])
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	hash1, err := base64.RawStdEncoding.DecodeString(arr[5])
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	hash2 := argon2.IDKey(
		[]byte(password),
		salt,
		t,
		m,
		p,
		uint32(len(hash1)),
	)

	return subtle.ConstantTimeCompare(hash1, hash2) == 1, nil
}
