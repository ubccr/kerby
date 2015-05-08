===============================================================================
Kerby - Go wrapper for Kerberos GSSAPI 
===============================================================================

This is a port of the PyKerberos library in Go. The main motivation for this
library was to provide HTTP client authentication using Kerberos. The khttp
package provides a transport that authenticates all outgoing requests using
SPNEGO (negotiate authentication) http://tools.ietf.org/html/rfc4559. 

The C code is adapted from PyKerberos http://calendarserver.org/wiki/PyKerberos.

------------------------------------------------------------------------
Usage
------------------------------------------------------------------------

Note: You need the have the krb5-libs/GSSAPI packages installed for your OS.

Install using go tools::

    $ go get github.com/ubccr/kerby

To run the unit tests you must have a valid Kerberos setup on the test machine
and you should ensure that you have valid Kerberos tickets (run 'klist' on the
command line). If you're authentication using a client keytab file you can
optionally export the env variable KRB5_CLIENT_KTNAME::

    $ export KRB5_CLIENT_KTNAME=/path/to/client.keytab
    $ export KERBY_TEST_SERVICE="service@REALM"
    $ export KERBY_TEST_PRINC="princ@REALM"
    $ go test

Example HTTP Kerberos client authentication using a client keytab file::

    package main

    import (
        "fmt"
        "io/ioutil"
        "bytes"
        "net/http"

        "github.com/ubccr/kerby/khttp"
    )

    func main() {
        payload := []byte(`{"method":"hello_world"}`)
        req, err := http.NewRequest(
            "POST",
            "https://server.example.com/json",
            bytes.NewBuffer(payload))

        req.Header.Set("Content-Type", "application/json")

        t := &khttp.Transport{
            KeyTab: "/path/to/client.keytab",
            Principal: "principal@REALM"}

        client := &http.Client{Transport: t}

        res, err := client.Do(req)
        if err != nil {
            panic(err)
        }
        defer res.Body.Close()

        data, err := ioutil.ReadAll(res.Body)
        if err != nil {
            panic(err)
        }

        fmt.Printf("%d\n", res.StatusCode)
        fmt.Printf("%s", data)
    }

------------------------------------------------------------------------
License
------------------------------------------------------------------------

Kerby is released under the Apache 2.0 License. See the LICENSE file.


