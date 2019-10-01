A service used to get an IP's country based on ipip and maxmind.

## Build & Run
```
make goip
mkdir bin/etc
cp -r etc_sample/* bin/etc
make run
```

## Usage
```
curl "http://127.0.0.1:18010/ip?ip=123.123.123.123"

curl "http://127.0.0.1:18010/mmdb?ip=123.123.123.123"

// if you want a specified port: bin/goip -port 127.0.0.1:your_port
```

response:
```
{
    "IP": "123.123.123.123",
    "Country": "China",
    "City": "Beijing",
    "Region": "Beijing",
    "Msg": "success"
}
```
---

## My server
```
curl https://common.givenzeng.cn/ipip?ip=123.123.123.123
curl https://common.givenzeng.cn/mmdb?ip=123.123.123.123
```
