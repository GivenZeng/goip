package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ipipdotnet/ipdb-go"
	geoip2 "github.com/oschwald/geoip2-golang"
)

const usage = "curl http://127.0.0.1:18010/ipip?ip=123.123.123.123 or \ncurl http://127.0.0.1:18010/mmdb?ip=123.123.123.123"

type Response struct {
	IP      string
	Country string
	City    string
	Region  string

	Msg string
}

type IPSrv struct {
	ipip *ipdb.City
	mmdb *geoip2.Reader
}

func NewSrv() (srv *IPSrv, err error) {
	binPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return nil, err
	}
	// ipip database
	ipdbFilePath := binPath + "/etc/ipip.ipdb"
	ipip, err := ipdb.NewCity(ipdbFilePath)
	if err != nil {
		return nil, err
	}

	// maxmind database
	mmdbPath := binPath + "/etc/GeoLite2-Country.mmdb"
	mmdb, err := geoip2.Open(mmdbPath)
	if err != nil {
		log.Fatal(err)
	}

	srv = &IPSrv{
		ipip: ipip,
		mmdb: mmdb,
	}

	go func() {
		// update ipdb
		ticker := time.NewTicker(time.Second * 60)
		for {
			<-ticker.C
			fmt.Println("reload ipdb: " + ipdbFilePath)
			err := srv.ipip.Reload(ipdbFilePath) // 更新 ipdb 文件后可调用 Reload 方法重新加载内容
			if err != nil {
				fmt.Println(err)
			}
		}
	}()
	return srv, nil
}

func (srv *IPSrv) ipHandler(w http.ResponseWriter, req *http.Request) {
	res := &Response{Msg: "success"}
	// get ip and validate
	ip := getIP(req)
	if ip == nil || ip.String() == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// find ip msg
	rst, err := srv.ipip.FindMap(ip.String(), "CN")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	res.IP = ip.String()
	res.Country = rst["country_name"]
	res.City = rst["city_name"]
	res.Region = rst["region_name"]

	resBytes, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(w, string(resBytes))
}

func (srv *IPSrv) mmdbHandler(w http.ResponseWriter, req *http.Request) {
	res := &Response{Msg: "success"}
	ip := getIP(req)
	if ip == nil || ip.String() == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	record, err := srv.mmdb.City(ip)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.IP = ip.String()
	res.Country = record.Country.Names["en"]
	res.Region = res.Country

	resBytes, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(w, string(resBytes))
}

func getIP(req *http.Request) (ip net.IP) {
	ipStr := req.URL.Query().Get("ip")
	if ipStr == "" {
		parts := strings.Split(req.RemoteAddr, ":")
		if len(parts) > 1 {
			ipStr = parts[0]
		}
	}
	ip = net.ParseIP(ipStr)
	return ip
}

func main() {
	port := flag.String("port", "127.0.0.1:18010", "api port, default: 127.0.0.1:18010")
	flag.Parse()
	srv, err := NewSrv()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", srv.ipHandler)
	http.HandleFunc("/ip", srv.ipHandler)
	http.HandleFunc("/ipip", srv.ipHandler)
	http.HandleFunc("/maxmind", srv.mmdbHandler)
	http.HandleFunc("/mmdb", srv.mmdbHandler)

	fmt.Println("start serving at: " + *port + "........\n" + usage + "\n")
	log.Fatal(http.ListenAndServe(*port, nil))
}
