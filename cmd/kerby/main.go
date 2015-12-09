// Copyright 2015 Andrew E. Bruno
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This is an example implementation for client/server SPNEGO-based Kerberos
// HTTP authentication
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ubccr/kerby"
	"github.com/ubccr/kerby/khttp"
)

var (
	negotiateHeader       = "Negotiate"
	wwwAuthenticateHeader = "WWW-Authenticate"
	authorizationHeader   = "Authorization"
)

// Kerberos based HTTP authentication server
func handler(w http.ResponseWriter, req *http.Request) {
	authReq := strings.Split(req.Header.Get(authorizationHeader), " ")
	if len(authReq) != 2 || authReq[0] != negotiateHeader {
		w.Header().Set(wwwAuthenticateHeader, negotiateHeader)
		http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
		return
	}

	ks := new(kerby.KerbServer)
	err := ks.Init("")
	if err != nil {
		log.Printf("KerbServer Init Error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer ks.Clean()

	err = ks.Step(authReq[1])
	w.Header().Set(wwwAuthenticateHeader, negotiateHeader+" "+ks.Response())

	if err != nil {
		log.Printf("KerbServer Step Error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	user := ks.UserName()
	fmt.Fprintf(w, "Hello, %s", user)
}

// Kerberos based HTTP authentication client
func client(url, keytab string) {
	req, err := http.NewRequest("GET", url, nil)
	t := &khttp.Transport{KeyTab: keytab}
	client := &http.Client{Transport: t}

	res, err := client.Do(req)
	if err != nil {
		log.Printf("Request failed with error: %s\n", err.Error())
		return
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Failed reading request body: %s\n", err.Error())
	}

	fmt.Printf("%d - %s\n", res.StatusCode, data)
}

var mode string
var url string
var keytab string
var port int

func init() {
	flag.StringVar(&mode, "mode", "client", "client or server")
	flag.StringVar(&keytab, "keytab", "", "path to keytab file")
	flag.StringVar(&url, "url", "http://localhost:8080", "the url to request in client mode")
	flag.IntVar(&port, "port", 8080, "the port to bind to in server mode")
}

func main() {
	flag.Parse()
	if "server" == mode {
		if len(keytab) > 0 {
			os.Setenv("KRB5_KTNAME", keytab)
		}
		http.HandleFunc("/", handler)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	} else {
		client(url, keytab)
	}
}
