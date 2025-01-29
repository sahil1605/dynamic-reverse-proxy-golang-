package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
	RunPort              = 2002                                           // The server port to run on
	ReverseServerAddr    = fmt.Sprint("0.0.0.0:", RunPort)                // this is our reverse server ip address
	InsideProxyHostname  = fmt.Sprint("proxy:", RunPort)                  // Requests from private network
	OutsideProxyHostname = fmt.Sprint("registration.localhost:", RunPort) // Requests from public network
	KnownAddresses       = map[string]string{}                            // Known Addresses
)