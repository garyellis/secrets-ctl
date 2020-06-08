package vault

import (
	"encoding/json"

	vaultapi "github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
)

// KV is the kv secret engine
type KV struct {
	*vaultapi.Client
}

// Put writes the secret to the kv path
func (kv *KV) Put(path string, data map[string]interface{}) error {
	d := make(map[string]interface{})
	d["data"] = data

	out, err := kv.Logical().Write(
		path,
		d,
	)
	if err != nil {
		return err
	}
	logdata, _ := json.Marshal(out.Data)
	log.Infof("[vault/kv] wrote secret %s %s", path, string(logdata))
	return nil
}

// Read reads the secret for the input secret
func (kv *KV) Read(path string, data map[string]interface{}) (map[string]interface{}, error) {
	out, err := kv.Logical().Read(
		path,
	)
	if err != nil {
		return nil, err
	}
	logdata, _ := json.Marshal(out.Data)
	log.Infof("[vault/kv] read secret %s", string(logdata))
	return out.Data, nil
}

// List returns a slice of vault kv paths recursively
func (kv *KV) List(path string) ([]string, error) {
	return nil, nil
}
