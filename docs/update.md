# To update the AV run the following:

```bash
$ docker run --name=fprot malice/fprot update
```

## Then to use the updated fprot container:

```bash
$ docker commit fprot malice/fprot:updated
$ docker rm fprot # clean up updated container
$ docker run --rm malice/fprot:updated EICAR
```
