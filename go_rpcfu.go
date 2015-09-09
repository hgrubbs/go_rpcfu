/*
   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
)

type Response struct {
	Body   string `json:"body"`
	Arg1   string `json:"arg1"`
	Method string `json:"method"`
}

func aHandler(w http.ResponseWriter, r *http.Request) {
	querystring_args := r.URL.Query()
	arg1_array, exists := querystring_args["arg1"]
	if exists == false {
		http.Error(w, "Missing 'arg1' variable from in URL", 406)
		return
	}

	arg1 := arg1_array[0]
	body, _ := ioutil.ReadAll(r.Body)

	// encode JSON
	var m Response
	m.Method = r.Method
	m.Body = string(body)
	m.Arg1 = string(arg1)
	msg, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, string(msg)) // return JSON
}

func main() {
	// Parse command-line args
	concurrency := flag.Int("cpus", 1, "Concurrency factor for multiple CPUs")
	bind_ip := flag.String("ip", "0.0.0.0", "Network address to listen on")
	bind_port := flag.String("port", "8080", "Network port to listen on")
	flag.Parse()

	runtime.GOMAXPROCS(*concurrency)

	// URL mapping
	http.HandleFunc("/ahandler/", aHandler)

	fmt.Printf("Listening on %s:%s\n", *bind_ip, *bind_port)
	http.ListenAndServe((*bind_ip + ":" + *bind_port), nil)
}
