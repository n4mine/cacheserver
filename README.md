# Cacheserver

* in memory tsdb

# Capability

* store 20M series with 2 hours data in memory (use less than 100GB memory)

# Architecture

<img src="https://n4mine.github.io/img/cc.jpg"/>



# How to Build && Run

```
go get -u -v github.com/n4mine/cacheserver
cd $GOPATH/github.com/n4mine/cacheserver
make

./cacheserver
```

# How to Push && Get data

## For debug/development

### Push data

use `series1` as name

```
curl -s -XPOST "http://127.0.0.1:7000/pushdata?name=series1&ts=$(date +%s)&value=$(shuf -i 0-100 -n1)"
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

* use rpc in [rpc/cacheserver.go](./rpc/cacheserver.go)
* need to unpack data in client side

```go
type DataResp struct {
	// code == 0, normal
	// code >  0, exception
	Code int    `msg:"code"`
	Msg  string `msg:"msg"`
	Key  string `msg:"key"`
	From int64  `msg:"from"`
	To   int64  `msg:"to"`
	Step int    `msg:"step"`
	RRA  int    `msg:"rra"`
	Data []Iter `msg:"data"`
}

type Iter struct {
	*tsz.Iter
}
```
