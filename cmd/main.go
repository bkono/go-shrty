package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	"git.kono.sh/shrty"
)

var (
	base     = flag.String("baseURL", "http://localhost:3000/", "The base URL when shortening")
	httpPort = flag.Int("httpPort", 3000, "The HTTP server port")
	grpcPort = flag.Int("grpcPort", 3001, "The gRPC server port")
)

func main() {
	flag.Parse()

	fmt.Printf("base = %+v\n", *base)

	// Setup db
	db := shrty.NewDBClient()
	db.Path = "shrty.db"
	err := db.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Setup TokenService
	var ts shrty.TokenService
	ts = shrty.NewTokenService("some secret salt")

	// Setup ShortenedURLService
	s := shrty.NewShortenedURLService(*base, db, ts)

	shrtygRPC := shrty.NewgRPCServer(s)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPort))
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		expandHandler(w, r, s)
	})
	go http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), nil)

	fmt.Println("past ListenAndServe")
	grpcServer := grpc.NewServer()
	shrty.RegisterShrtyServer(grpcServer, shrtygRPC)
	grpcServer.Serve(lis)
}

func expandHandler(w http.ResponseWriter, r *http.Request, s shrty.ShortenedURLService) {
	tk := strings.TrimLeft(r.URL.Path, "/")
	url, err := s.Expand(tk)
	if err != nil {
		log.Printf("Error while expanding token = %+v\n", tk)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("redirecting token %v to %v", tk, url)
	http.Redirect(w, r, url, http.StatusFound)
}

// TODOs:
// - take flags for the url base
// - cache urls, don't recreate them
// - grpc endpoint for create and expand
// - web endpoint for metrics
// - move metrics to async task
// - swap to chi, for the fun of it
