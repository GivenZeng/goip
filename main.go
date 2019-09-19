package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ipipdotnet/ipdb-go"
)

type Response struct {
	IP      string
	Country string
	City    string
	Region  string

	Msg string
}

type IPSrv struct {
	db *ipdb.City
}

func NewSrv() (srv *IPSrv, err error) {
	binPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return nil, err
	}
	ipdbFilePath := binPath + "/etc/ipip.ipdb"
	db, err := ipdb.NewCity(ipdbFilePath)
	if err != nil {
		return nil, err
	}
	srv = &IPSrv{db: db}

	go func() {
		// update ipdb
		ticker := time.NewTicker(time.Second * 60)
		for {
			<-ticker.C
			fmt.Println("reload ipdb: " + ipdbFilePath)
			err := srv.db.Reload(ipdbFilePath) // 更新 ipdb 文件后可调用 Reload 方法重新加载内容
			if err != nil {
				fmt.Println(err)
			}
		}
	}()
	return srv, nil
}

func (srv *IPSrv) ipHandler(w http.ResponseWriter, req *http.Request) {
	// fmt.Println(db.IsIPv4())    // check database support ip type
	// fmt.Println(db.IsIPv6())    // check database support ip type
	// fmt.Println(db.BuildTime()) // database build time
	// fmt.Println(db.Languages()) // database support language
	// fmt.Println(srv.db.Fields()) // database support fields

	// fmt.Println(srv.db.FindInfo("2001:250:200::", "CN")) // return CityInfo
	// fmt.Println(srv.db.Find("1.1.1.1", "CN"))       // return []string
	// fmt.Println(srv.db.FindInfo("127.0.0.1", "CN")) // return CityInfo
	res := &Response{Msg: "success"}
	defer func() {
		resBytes, err := json.Marshal(res)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprintf(w, string(resBytes))
	}()

	// get ip and validate
	ipStr := req.URL.Query().Get("ip")
	if ipStr == "" {
		parts := strings.Split(req.RemoteAddr, ":")
		if len(parts) > 1 {
			ipStr = parts[0]
		}
		res.Msg = "ip should not be empty, example = http://domain/ip?ip=123.123.123.123"
	}
	// find ip msg
	rst, err := srv.db.FindMap(ipStr, "CN")
	if err != nil {
		res = &Response{Msg: "error occured: " + err.Error()}
		return
	}
	res.IP = ipStr
	res.Country = rst["country_name"]
	res.City = rst["city_name"]
	res.Region = rst["region_name"]
}

func main() {
	srv, err := NewSrv()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", srv.ipHandler)
	port := "127.0.0.1:18010"
	fmt.Println("start serve at: " + port + "........")
	log.Fatal(http.ListenAndServe(port, nil))
}
