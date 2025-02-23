# stuttgart-things/homerun-chaos-catcher

## PRE-TASKS

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
PATH_TO_KUBECONFIG="/home/sthings/.kube/config"
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
