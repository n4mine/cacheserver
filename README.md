# Cacheserver

* in memory tsdb

# Architecture

<img width="250" height="312" src="https://n4mine.github.io/img/cacheserver.png"/>



# BUILD

```
go get -u -v github.com/n4mine/cacheserver
cd $GOPATH/github.com/n4mine/cacheserver
make
```

# RUN

```
./cacheserver
```

# USAGE

## For debug/development

### Push data

use `series1` as name

```
curl -s -XPOST "http://127.0.0.1:7000/push?name=series1&ts=$(date +%s)&value=$(shuf -i 0-100 -n1)"
```

### Get data info

```
curl -s 'http://127.0.0.1:7000/getinfo?name=series1'
```

resp:

```
{
  "duration": 28,
  "newest": 1567001322,
  "oldest": 1567001294
}
```

### Get data

```
curl -s 'http://127.0.0.1:7000/getdata?name=series1&from=1567001294&to=1567001322'
```

resp:

```
{
  "1567001294": 10,
  "1567001295": 23,
  "1567001296": 14,
  "1567001319": 25,
  "1567001320": 69,
  "1567001321": 66,
  "1567001322": 71
}

```

## For production

* use [RPC](./rpc/cacheserver.go)
