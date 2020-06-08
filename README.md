# secrets-ctl [![Build Status](https://travis-ci.org/garyellis/secrets-ctl.svg?branch=master)](https://travis-ci.org/garyellis/secrets-ctl)
Secrets file storage utility for IAC pipelines.


## Features
* vault transit secret engine encryption/decryption
* decrypt/encrypt secrets to yaml configuration files
* write yaml configuration files to vault kv


## Use Cases
* secrets storage in git
* cd pipeline for vault kv secrets (solves how do we maintain a pipeline that write secrets into to vault)


## Usage

write a file and encrypt it
```bash
$ cat > secret.yaml <<EOF
secret:
- vault_mount: transit
  vault_key: demo-key
  vault_kv_path: secret/data/secret1
  data:
    foo1: '{ "foo": "abcd", "bar": "efgh" }'
    foo2: "abcd"
- vault_mount: transit
  vault_key: demo-key
  vault_kv_path: secret/data/secret2
  data:
    APIKEY: "pretendkey"
EOF


$ export VAULT_ADDR=<my-vault-server>
$ export VAULT_TOKEN=<set-as-needed>


$ secrets-ctl encrypt --path . --filter secret.yaml --backup=false


$ cat secret.yaml
secret:
- vault_mount: transit
  vault_key: demo-key
  vault_kv_path: secret/data/secret1
  data:
    foo1: vault:v1:fi2AsBUa1cH77HI7TC1VwO0CczUVQwtbQYrYdaGzKCvv0fgZXSy3UuKLZrGK97+QXxwOSzpBrtdAtNA8
    foo2: vault:v1:0dxQvi20T0Aj3wnh80FEjwPQgY9WOuq0tQDndiovedg=
- vault_mount: transit
  vault_key: demo-key
  vault_kv_path: secret/data/secret2
  data:
    APIKEY: vault:v1:P6y6pBzl6xfqbO19Gvh8T2/3ibTINKMRDf4ua13rWnAw8PhTssQ=
```

decrypt the file
```bash
$ secrets-ctl decrypt --path . --filter secret.yaml --backup=false


$ cat secret.yaml
secret:
- vault_mount: transit
  vault_key: demo-key
  vault_kv_path: secret/data/secret1
  data:
    foo1: '{ "foo": "abcd", "bar": "efgh" }'
    foo2: abcd
- vault_mount: transit
  vault_key: demo-key
  vault_kv_path: secret/data/secret2
  data:
    APIKEY: pretendkey
```

encrypt a file and write it to to vault kv store
```bash
$ secrets-ctl encrypt --path . --filter secret.yaml --backup=false


$ secrets-ctl vault-kv write --path . --filter secret.yaml
```

read the secrets with the vault kv command.
```bash
$ vault kv get secret/secret1
====== Metadata ======
Key              Value
---              -----
created_time     2020-06-08T03:58:39.930912228Z
deletion_time    n/a
destroyed        false
version          1

==== Data ====
Key     Value
---     -----
foo1    { "foo": "abcd", "bar": "efgh" }
foo2    abcd


$ vault kv get secret/secret2
====== Metadata ======
Key              Value
---              -----
created_time     2020-06-08T03:58:40.072491126Z
deletion_time    n/a
destroyed        false
version          1

===== Data =====
Key       Value
---       -----
APIKEY    pretendkey
```
