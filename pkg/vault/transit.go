package vault

import (
	"encoding/base64"
	"path"
	"regexp"

	vaultapi "github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
)

var (
	transitEncryptedValue = regexp.MustCompile(`^vault:v\d+:.+$`)
)

// Transit is the transit engine vault client
type Transit struct {
	*vaultapi.Client
}

// IsEncrypted check with regexp that value encrypter by Vault transit secret engine
func (t *Transit) IsEncrypted(value string) bool {
	return transitEncryptedValue.MatchString(value)
}

// Decrypt decrypts the given encrypted value using the specified transit engine path and key
func (t *Transit) Decrypt(mountpath, key string, ciphertext []byte) ([]byte, error) {
	if !t.IsEncrypted(string(ciphertext)) {
		log.Debugf("[vault/transit] %s is not encrypted. Skipping.", ciphertext)
		return ciphertext, nil
	}
	out, err := t.Logical().Write(
		path.Join(mountpath, "decrypt", key),
		map[string]interface{}{
			"ciphertext": string(ciphertext),
		},
	)
	if err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(out.Data["plaintext"].(string))
}

// Encrypt encrypts the given value using the specified transit engine path and key
func (t *Transit) Encrypt(mountpath, key string, value []byte) ([]byte, error) {
	if t.IsEncrypted(string(value)) {
		log.Infof("[vault/transit] %s is already encrypted. Skipping.", value)
		return value, nil
	}
	out, err := t.Logical().Write(
		path.Join(mountpath, "encrypt", key),
		map[string]interface{}{
			"plaintext": base64.StdEncoding.EncodeToString(value),
		},
	)
	if err != nil {
		return nil, err
	}
	return []byte(out.Data["ciphertext"].(string)), nil
}
