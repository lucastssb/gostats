package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Data struct {
	TotalMemory     uint64  `json:"total_memory"`
	FreeMemory      uint64  `json:"free_memory"`
	UsedMemory      uint64  `json:"used_memory"`
	UsedPercent     float64 `json:"used_percent"`
	CpuUsagePercent float64 `json:"cpu_usage_percent"`
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	for {
		time.Sleep(1 * time.Second)

		info, err := getSystemInfo()
		if err != nil {
			log.Println("Error getting system info:", err)
		}

		data, err := json.Marshal(info)
		if err != nil {
			log.Println("Error marshaling data:", err)
		}

		err = c.WriteMessage(1, []byte(data))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func getSystemInfo() (Data, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Println("Error getting memory info:", err)
	}
	p, err := cpu.Percent(0, false)
	if err != nil {
		log.Println("Error getting cpu info:", err)
	}

	data := Data{
		TotalMemory:     v.Total,
		FreeMemory:      v.Free,
		UsedMemory:      v.Used,
		UsedPercent:     v.UsedPercent,
		CpuUsagePercent: p[0],
	}
	return data, nil
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	log.Fatal(http.ListenAndServe(":8888", nil))
}
