package vault

import (
	"encoding/json"
	"path"
	"regexp"
	"strings"

	vaultapi "github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
)

var (
	secretExp = regexp.MustCompile(`^(/.*)/(data)(.*)`)
)

// KV is the kv secret engine
type KV struct {
	*vaultapi.Client
}

// Put writes the secret to the kv path
func (kv *KV) Put(key string, data map[string]interface{}) error {
	d := make(map[string]interface{})
	d["data"] = data

	out, err := kv.Logical().Write(
		key,
		d,
	)
	if err != nil {
		return err
	}
	logdata, _ := json.Marshal(out.Data)
	log.Infof("[vault/kv] wrote secret %s %s", key, string(logdata))
	return nil
}

// Read reads the secret for the input key
func (kv *KV) ReadSecretData(key string) (map[string]interface{}, error) {
	log.Debugf("[vault/kv] reading secret data: %s", key)
	out, err := kv.Logical().Read(key)
	if err != nil {
		log.Warn("[vault/kv] ", err)
		return nil, err
	}
	if out == nil {
		log.Warnf("[vault/kv] secret not found %s", key)
	}
	logdata, _ := json.Marshal(out.Data)
	log.Debugf("[vault/kv] secret %s data: %s", key, string(logdata))
	return out.Data, nil
}

// GetKeys returns all keys for the mount and path recursively
func (kv *KV) GetKeys(key string) ([]string, error) {
	mount := secretExp.ReplaceAllString(key, "$1")
	key = secretExp.ReplaceAllString(key, "$3")

	var keys []string
	var f func(string)
	f = func(p string) {
		metadataPath := path.Join(mount, "metadata", p)
		log.Debugf("[vault/kv] reading metadatapath: %s, key: %s", metadataPath, p)

		// check if path is a secret. if it is, append it to the keys slice
		if !strings.HasSuffix(metadataPath, "/") && p != "/" && p != "" {
			secretResponse, err := kv.Logical().Read(metadataPath)
			if secretResponse != nil && err == nil {
				log.Debugf("[vault/kv] found secret %s is a secret", metadataPath)
				log.Debugf("[vault/kv] appending key: %s", p)
				//keys = append(keys, key)
				keys = append(keys, p)
			}
			if err != nil {
				log.Warn("[vault/kv] ", err)
			}
		}

		// check if the path is a folder
		secretFolderResponse, err := kv.Logical().List(metadataPath)
		if err != nil {
			log.Warn("[vault/kv]", err)
			return
		}
		if secretFolderResponse == nil && err == nil {
			log.Warnf("[vault/kv] path %s not found", metadataPath)
			return
		}

		// process the folder's keys
		currentKeys := secretFolderResponse.Data["keys"].([]interface{})
		log.Debugf("[vault/kv] processing path %s:  keys: %s", p, currentKeys)

		for _, k := range currentKeys {
			currentKey := path.Join(p, k.(string))
			log.Debugf("[vault/kv] current key: %s", currentKey)

			if !strings.HasSuffix(k.(string), "/") {
				log.Debugf("[vault/kv] appending key: %s", currentKey)
				keys = append(keys, currentKey)
			} else {
				log.Debugf("[vault/kv] %s is a folder. populating sub keys", currentKey)
				f(currentKey)
			}
		}
	}

	f(key)
	log.Debugf("[vault/kv] result - path: %s keys: %s", key, keys)
	for i, k := range keys {
		keys[i] = path.Join(mount, "data", k)
	}

	return keys, nil
}

// GetSecrets returns vault secrets on the given path
func (kv *KV) GetSecrets(key string) ([]map[string]interface{}, error) {
	var secrets []map[string]interface{}

	keys, err := kv.GetKeys(key)
	log.Debugf("[vault/kv] keys: %s", strings.Join(keys, ","))
	if err != nil {
		return nil, err
	}
	for _, i := range keys {
		secret, err := kv.ReadSecretData(i)
		if err != nil {
			return nil, err
		}
		secrets = append(secrets, secret)
	}
	return secrets, nil
}
