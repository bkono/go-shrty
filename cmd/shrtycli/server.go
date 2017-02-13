// Copyright © 2017 Bryan Konowitz <bryan@kono.sh>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package shrtycli

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"git.kono.sh/bkono/shrty"
	"github.com/spf13/cobra"
)

var baseURL string
var httpPort, grpcPort int

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Launches the shrty server process",
	Long: `This command starts up an instance of the shrty server listening for HTTP and gRPC requests.

Shrty server listens for requests at <baseURL>/<token> and auto redirects to the expanded address.
It also exposes a gRPC API for expanding and shortening URLs programmatically or via the client 
functions of this binary.

Shrty server requires a baseURL to provide back when shortening. 
Example would be http://<yourdomain>/ `,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("server called with baseURL=%v httpPort=%v grpcPort=%v\n", baseURL, httpPort, grpcPort)

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
		ts = shrty.NewTokenService("some secret salt") // extract to param or env

		// Setup ShortenedURLService
		s := shrty.NewShortenedURLService(baseURL, db, ts)

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			expandHandler(w, r, s)
		})
		go http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil)
		shrty.RunGRPCServer(s, grpcPort)
	},
}

func init() {
	RootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVar(&baseURL, "baseURL", "localhost", "URL to append tokens to when shortening")
	serverCmd.Flags().IntVar(&httpPort, "httpPort", 3000, "HTTP listening port")
	serverCmd.Flags().IntVar(&grpcPort, "grpcPort", 3001, "gRPC listening port")
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
