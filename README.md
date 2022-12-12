A service used to get an IP's country based on ipip and maxmind.

## Build & Run
```
make goip # make a excutable binanry
mkdir bin/etc 
cp -r etc_sample/* bin/etc # copy the sample config to ect dir
make run # run
# or: bin/goip -port 127.0.0.1:your_port
```

## Usage
```
curl "http://127.0.0.1:18010/ip?ip=123.123.123.123"

curl "http://127.0.0.1:18010/mmdb?ip=123.123.123.123"

// if you want a specified port: 
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


## other
you can download the newest ip database from [ipip](https://www.ipip.net/product/ip.html#ipv4city)，and replace the database file in directory etc
