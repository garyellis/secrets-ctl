# secrets-ctl [![Build Status](https://travis-ci.org/garyellis/secrets-ctl.svg?branch=master)](https://travis-ci.org/garyellis/secrets-ctl)
Secrets file storage utility for IAC pipelines.


## Features
* vault transit secret engine encryption/decryption
* decrypt/encrypt secrets to yaml configuration files
* write secret yaml configuration files to vault kv
* export vault kv secrets to secret yaml configuration files


## Use Cases
* secrets storage in git
* cd pipeline for vault kv secrets (solves how do we maintain a pipeline that write secrets into to vault)
* backup vault kv secrets engine secrets to yaml


## Usage


write a file and encrypt it
```bash
$ cat > secret.yaml <<EOF
secretconfig:
  encryption:
    transit_mount: /transit
    transit_key: demo-key
  secrets:
    - vault_kv_path: secret/data/secret1
      data:
        foo1: '{ "foo": "abcd", "bar": "efgh" }'
        foo2: "abcd"
    - vault_kv_path: secret/data/secret2
      data:
        APIKEY: "pretendkey"
EOF


$ export VAULT_ADDR=<my-vault-server>
$ export VAULT_TOKEN=<set-as-needed>


$ secrets-ctl encrypt --path . --filter secret.yaml --backup=false


$ cat secret.yaml
secretconfig:
  encryption:
    transit_mount: /transit
    transit_key: demo-key
  secrets:
  - vault_kv_path: secret/data/example/secret1
    data:
      foo1: vault:v1:HGhX1kSZswsvfDON4ZIWHcZJYEBwDkng7wb//pBliIka6VEY1+jKOmIP8B/DLHiCNNZMVmacJuRcghWr
      foo2: vault:v1:bTPRokegzYryLsoNO1NSO3eA+qUVzAE5lcSGOk3kelI=
  - vault_kv_path: secret/data/example/secret2
    data:
      APIKEY: vault:v1:1kXDbR4WAobetY1GHlIPDb4pHaZYp5iwI6jxKDhHktKr6DtTKv4=
```

decrypt the file
```bash
$ secrets-ctl decrypt --path . --filter secret.yaml --backup=false


$ cat secret.yaml
secretconfig:
  encryption:
    transit_mount: /transit
    transit_key: demo-key
  secrets:
  - vault_kv_path: secret/data/example/secret1
    data:
      foo1: '{ "foo": "abcd", "bar": "efgh" }'
      foo2: abcd
  - vault_kv_path: secret/data/example/secret2
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
$ vault kv get /secret/example/secret1
====== Metadata ======
Key              Value
---              -----
created_time     2020-06-12T04:30:42.200527167Z
deletion_time    n/a
destroyed        false
version          2

==== Data ====
Key     Value
---     -----
foo1    { "foo": "abcd", "bar": "efgh" }
foo2    abcd


$ vault kv get /secret/example/secret2
====== Metadata ======
Key              Value
---              -----
created_time     2020-06-12T04:30:42.350926653Z
deletion_time    n/a
destroyed        false
version          2

===== Data =====
Key       Value
---       -----
APIKEY    pretendkey
```
