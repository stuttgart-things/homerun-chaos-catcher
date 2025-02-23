# stuttgart-things/homerun-chaos-catcher

## USAGE

<details><summary>CONFIG EXMAPLE</summary>

```bash
cat <<EOF > chaos.yaml
---
chaosEvents:
  chaos1:
    systems:
      - tetris # all systems
    severity:
      - CHAOS1
    operation: delete
    count: 9
    resource: pod
    namespace: "*" # random namespace
  chaos2:
    systems:
      - tetris # all systems
    severity:
      - CHAOS2
    operation: add
    count: 1
    resource: deployment
    namespace: "*" # random namespace
EOF
```

</details>

<details><summary>RUN</summary>

```bash
export REDIS_SERVER=localhost
export REDIS_PORT=5000
export REDIS_PASSWORD=""
export REDIS_STREAM="homerun"
export REDIS_CONSUMER_GROUP="homerun-chaos-catcher"
export PROFILE_PATH="chaos.yaml"
export KUBECONFIG="/home/sthings/.kube/config"

homerun-chaos-catcher
```

</details>



## DEPLOYMENT

<details><summary>GITHUB RELEASE</summary>

```bash
VERSION=v1.3.0
BIN_DIR=/usr/bin
cd /tmp && wget https://github.com/stuttgart-things/homerun-chaos-catcher/releases/download/${VERSION}}/homerun-chaos-catcher_Linux_x86_64.tar.gz
tar xvfz homerun-chaos-catcher_Linux_x86_64.tar.gz
sudo mv homerun-chaos-catcher ${BIN_DIR}/homerun-chaos-catcher
sudo chmod +x ${BIN_DIR}/homerun-chaos-catcher
rm -rf CHANGELOG.md README.md LICENSE
cd -
```

</details>

## DEV

<details><summary>CREATE ENV FILE</summary>

.env file needed for Taskfile

```bash
cat <<EOF > .env
REDIS_SERVER=localhost
REDIS_PORT=5000
REDIS_PASSWORD=""
REDIS_STREAM="homerun"
REDIS_CONSUMER_GROUP="homerun-chaos-catcher"
PROFILE_PATH="tests/config.yaml"
KUBECONFIG="/home/sthings/.kube/config"
EOF
```

</details>

<details><summary>PORT-FORWARD REDIS</summary>

```bash
kubectl -n homerun port-forward services/redis-stack-headless 5000:6379
```

</details>


<details><summary>SEND TEST MESSAGE TO HOMERUN (GENERIC PITCHER)</summary>

```bash
ADDRESS=https://homerun.homerun-dev.sthings-vsphere.labul.sva.de/generic
curl -k -X POST "${ADDRESS}" \
    -H "Content-Type: application/json" \
    -H "X-Auth-Token: IhrGeheimerToken" \
    -d '{
           "title": "2 lines cleared",
           "message": "2 lines cleared at tetris",
           "severity": "CHAOS2",
           "author": "andreu",
           "timestamp": "2024-5-01T12:00:00Z",
           "system": "tetris",
           "tags": "tetris,lines,score",
           "assigneeaddress": "",
           "assigneename": "",
           "artifacts": "",
           "url": ""
    }'
```

</details>


## LICENSE

<details><summary><b>APACHE 2.0</b></summary>

Copyright 2023 patrick hermann.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

</details>

Author Information
------------------
Patrick Hermann, stuttgart-things 03/2025