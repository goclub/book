# rolling-update

```shell
docker service create --name goclub_book_rolling_update --replicas 2 name:version
docker service ps name
docker service udpate --image name/version
```