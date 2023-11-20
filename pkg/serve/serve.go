package serve

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"tailscale.com/tsnet"
)

//go:embed form/*
var content embed.FS

func Serve(listenPort string, cogPort string, tailscale string) error {

	handler, err := buildHandler(cogPort)
	if err != nil {
		return err
	}
	fmt.Println("Proxy to listen on port", listenPort)

	if tailscale != "" {
		s := &tsnet.Server{Hostname: tailscale}
		defer s.Close()

		ln, err := s.ListenFunnel("tcp", ":443") // does TLS
		if err != nil {
			log.Fatal(err)
		}
		defer ln.Close()

		log.Fatal(http.Serve(ln, handler))
	}

	err = proxy(listenPort, handler)
	if err != nil {
		return err
	}

	return nil
}

func proxy(listenPort string, handler http.Handler) error {
	fmt.Println("Proxy is listening on port", listenPort)
	http.Handle("/", handler)
	return http.ListenAndServe(":"+(listenPort), nil)
}

func buildHandler(cogPort string) (http.Handler, error) {
	targetURL, err := url.Parse("http://localhost:" + cogPort)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
		if strings.HasPrefix(r.URL.Path, "/form") {
			http.FileServer(http.FS(content)).ServeHTTP(w, r)
		} else {
			proxy.ServeHTTP(w, r)
		}
	}), nil
}
