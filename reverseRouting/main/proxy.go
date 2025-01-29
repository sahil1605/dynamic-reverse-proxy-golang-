// Proxy runs the actual proxy and will look at the
// hostnames requested from the received request. It will
// then translate that to the inside hostname and forward the
// request
func Proxy(c *gin.Context) {

	// Get if HTTP or HTTPS
	scheme := GetScheme(c)

	log.Println(scheme, c.Request.Host, c.Request.URL.String())

	// If this is a registration request, save it and
	// then stop processing this request
	if IsRegistrationRequest(c) {

		err := SaveRegistrationRequest(c)

		if err != nil {
			log.Println(err)
			c.String(400, "Couldnt Register Host")
			return
		}

		c.String(201, "Host Registered")
		return
	}

	// Translate the outside hostname to the inside hostname
	forwardTo, ok := KnownAddresses[c.Request.Host]

	if !ok {
		log.Printf("Unkown Host: %v", c.Request.Host)
		c.String(400, "Unkown Host")
		return
	}

	rUrl := fmt.Sprintf("%v://%v%v", scheme, forwardTo, c.Request.URL)

	remote, err := url.Parse(rUrl)

	if err != nil {
		log.Println(err)
		c.String(500, "Error Proxying Host")
		return
	}

	log.Println("Forwarding request to", remote)

	// Forward the request to the inside remote server
	// https://pkg.go.dev/net/http/httputil#NewSingleHostReverseProxy
	proxy := httputil.NewSingleHostReverseProxy(remote)

	// Director is a function which modifies
	// the request into a new request to be sent
	// https://pkg.go.dev/net/http/httputil#ReverseProxy
	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = c.Param("path")
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}