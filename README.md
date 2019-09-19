A service used to get an IP's country based on ipip.

## Build & Run
```
make goip
cp -r etc_sample bin/etc
make run


## Usage
```
curl "http://127.0.0.1:18010/ip?ip=123.123.123.123"

// if you want to specify port: bin/goip -port 127.0.0.1:your_port
```

response:
```
{
    "IP": "123.123.123.123",
    "Country": "中国",
    "City": "北京",
    "Region": "北京",
    "Msg": "success"
}
```


## My server
```
curl https://common.givenzeng.cn/ip?ip=123.123.123.123
```
