/*

Author:	Rakesh.N

HOW-TO:

Generate Rack file 			PATH:rakesh_script/checkall/Monitoring_Web/Rack_Gateway_Creation.sh
Generate variable file			PATH:rakesh_script/checkall/Monitoring_Web/Variable_File.sh

1.start script from 192.168.20.179	PATH:rakesh_script/checkall/Monitoring_Web/Web_Monitor_DataGenerator.py
2.start client from 192.168.20.112	PATH:WebMonitor/Web_Monitor_JSON_Client.py
3.start WebServer GoScript		PATH:WebMonitor/WebMonitor_Server.go  { go run FILE }

To View Monitoring:
http://192.168.CLIENT.IP:8080/CR	FOR CR
http://192.168.CLIENT.IP:8080/STADD	FOR Straddler

Use Button (TOP) to SWAP between both

To View Server LOGS:
/tmp/Monitor_Web.log

*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

//ProcessStatus Structure for data storage
type ProcessStatus struct {
	Monitor []struct {
		Name    string `json:"name"`
		Servers []struct {
			Status  string `json:"status"`
			IP      string `json:"ip"`
			Order   int    `json:"order"`
			Trades  int    `json:"trades"`
			Gateway string `json:"gate"`
			Rack    string `json:"rack"`
			CTO     int    `json:"cto"`
		} `json:"server"`
	} `json:"monitor"`
}

var GlobalData ProcessStatus

//StringModify To process button
func StringModify(Strategy string) string {
	var strat string
	if (strings.Compare(Strategy, "STADD")) == 0 {
		strat = "CR"
		return strat
	} else if (strings.Compare(Strategy, "CR")) == 0 {
		strat = "STADD"
		return strat
	}
	return Strategy
}

//WebMonitor for handling request
func WebMonitor(w http.ResponseWriter, r *http.Request) {
	var flag = true
	var range1, range2, range145std_1, range145std_2 net.IP
	var FILE = "/tmp/status.json"
	var counter = 11
	var OrderCount = 0
	var TradeCount = 0
	var index int
	message := r.URL.Path
	messages := strings.TrimPrefix(message, "/")
	log.Println("Request from", r.RemoteAddr, "For", message)
	range145std_1 = net.ParseIP("192.168.145.70")
	range145std_2 = net.ParseIP("192.168.145.90")
	if strings.Compare(messages, "CR") == 0 {
		index = 1
		range1 = net.ParseIP("192.168.145.1")
		range2 = net.ParseIP("192.168.145.254")

	} else if strings.Compare(messages, "STADD") == 0 {
		index = 0
		range1 = net.ParseIP("192.168.120.1")
		range2 = net.ParseIP("192.168.120.254")

	} else {
		message = "WRONG URL :" + message + "\n\n"
		panic(message)
	}
	plan, err1 := ioutil.ReadFile(FILE)
	if err1 != nil {
		panic(err1)
	}
	var data ProcessStatus
	err := json.Unmarshal(plan, &data)
	if err != nil {
		fmt.Println(err)
		data = GlobalData
		//panic(err)
	} else {
		GlobalData = data
	}
	fmt.Fprintf(w, "<html>")
	fmt.Fprintf(w, "<head>")
	fmt.Fprintf(w, "<meta http-equiv=\"refresh\" content=\"3\">")
	fmt.Fprintf(w, "<meta http-equiv=\"Cache-control\" content=\"no-cache, no-store, must-revalidate\">")
	fmt.Fprintf(w, "<meta http-equiv=\"Pragma\" content=\"no-cache\">")
	fmt.Fprintf(w, "<title>Monitoring App</title>")
	fmt.Fprintf(w, "<style>")
	fmt.Fprintf(w, ".hoverable{cursor:default;color:#000;text-decoration:none;} .hoverable .hover{display:none;} .hoverable:hover .normal{display:none;} .hoverable:hover .hover{display:inline;}")
	fmt.Fprintf(w, ".button {display: inline-block;border-radius: 8px; background-color: #555555; border:none; color: #FFFFFF;text-align:center;font-size:16px;padding:20px;width:150px;transition:all 0.5s;cursor:pointer;margin:5px;}")
	fmt.Fprintf(w, ".button span {cursor:pointer;display:inline-block;position:relative;transition:0.5s;}")
	fmt.Fprintf(w, ".button span:after{content: \"\\00bb\"; position: absolute; opacity:0; top:0; right: -20px;transition:0.5s;}")
	fmt.Fprintf(w, ".button:hover span{padding-right:20px;}")
	fmt.Fprintf(w, ".button:hover span:after{opacity:1;right:0;}")
	fmt.Fprintf(w, "table, td, th{ border: 1px solid black; padding: 5px; font-size: 11.3px; }")
	fmt.Fprintf(w, "tr:hover{background-color:#E6E6E6;}")
	fmt.Fprintf(w, "</style>")
	fmt.Fprintf(w, "</head>")
	fmt.Fprintf(w, "<body>")
	fmt.Fprintf(w, "<h1><form action=\"%[1]v\" method=\"POST\"><a href=\"javascript://\" class=\"hoverable\"><button type=\"submit\" class=\"button\" style=\"vertical-align:middle\"><span class=\"normal\">%v</span><span class=\"hover\">%[1]v</span></button></form></h1></a>", StringModify(strings.Trim(message, "/")), data.Monitor[index].Name)
	for i := range data.Monitor[index].Servers {
		var NetIP = net.ParseIP("192.168." + data.Monitor[index].Servers[i].IP)
		if flag {
			fmt.Fprintf(w, "<table style=\"float: left; padding: 0.5px;\">")
			fmt.Fprintf(w, "<col width=\"70\"> <col width=\"100\"> <col width=\"30\"> <col width=\"20\"> <col width=\"20\"> <col width=\"20\"> ")
			fmt.Fprintf(w, "<tr>")
			fmt.Fprintf(w, "<th>IPADDR</th>")
			fmt.Fprintf(w, "<th>STATUS</th>")
			fmt.Fprintf(w, "<th>ORD</th>")
			fmt.Fprintf(w, "<th>TRD</th>")
			fmt.Fprintf(w, "<th>RCK</th>")
			fmt.Fprintf(w, "<th font size=\"2\">GTY</th>")
			fmt.Fprintf(w, "</tr>")
			flag = false
			counter = 24
		}
		if ((bytes.Compare(NetIP, range1) >= 0 && bytes.Compare(NetIP, range2) <= 0) || (bytes.Compare(NetIP, range145std_1) >= 0 && bytes.Compare(NetIP, range145std_2) <= 0)) && data.Monitor[index].Servers[i].CTO == 1 {
			OrderCount = OrderCount + data.Monitor[index].Servers[i].Order
			TradeCount = TradeCount + data.Monitor[index].Servers[i].Trades
		}
		fmt.Fprintf(w, "<tr>")
		fmt.Fprintf(w, "<td>%v</td>", data.Monitor[index].Servers[i].IP)
		if strings.Compare(data.Monitor[index].Servers[i].Status, "running") == 0 {
			fmt.Fprintf(w, "<td bgcolor=#00E600>%v</td>", data.Monitor[index].Servers[i].Status)
		} else {
			fmt.Fprintf(w, "<td bgcolor=#E60000>%v</td>", data.Monitor[index].Servers[i].Status)
		}
		fmt.Fprintf(w, "<td>%v</td>", data.Monitor[index].Servers[i].Order)
		fmt.Fprintf(w, "<td>%v</td>", data.Monitor[index].Servers[i].Trades)
		fmt.Fprintf(w, "<td>%v</td>", data.Monitor[index].Servers[i].Rack)
		fmt.Fprintf(w, "<td>%v</td>", data.Monitor[index].Servers[i].Gateway)
		fmt.Fprintf(w, "</tr>")
		counter--
		if counter == 0 {
			flag = true
		}
	}
	fmt.Fprintf(w, "</table>")
	fmt.Fprintf(w, "<div style=\"clear:left;\"></div>")
	fmt.Fprintf(w, "<p></p>")
	fmt.Fprintf(w, "<table style=\"padding:0.01px;\">")
	fmt.Fprintf(w, "<col width=\"95\"> <col width=\"95\">")
	fmt.Fprintf(w, "<tr>")
	fmt.Fprintf(w, "<th>Orders</th><th>Trades</th>")
	fmt.Fprintf(w, "</tr><tr>")
	fmt.Fprintf(w, "<td>%v</td><td>%v</td>", OrderCount, TradeCount)
	fmt.Fprintf(w, "</tr></table>")
	fmt.Fprintf(w, "</body>")
	fmt.Fprintf(w, "</html>")
	log.Print("Client Metadata: \n", r, "\n\n")
}

func main() {
	f, err1 := os.OpenFile("/tmp/Monitor_Web.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err1 != nil {
		log.Fatalf("error opening file: %v", err1)
	}
	log.SetOutput(f)
	log.Println("Starting a server ...")
	http.HandleFunc("/", WebMonitor)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err1)
	}
}
