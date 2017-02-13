// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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

package shrtycli

import (
	"context"
	"fmt"

	"git.kono.sh/bkono/shrty"
	"google.golang.org/grpc/grpclog"

	"github.com/spf13/cobra"
)

// shortenCmd represents the shorten command
var shortenCmd = &cobra.Command{
	Use:   "shorten",
	Short: "Shortens a url",
	Long: `Shortens the provided url, and returns a token as well as a fully qualifed short url reference.

Example: shrtyctl shorten https://google.com # output: Token: abc ShrtUrl: https://localhost/abc`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("url to shorten required")
			return
		}
		fmt.Printf("shorten called for %v\n", args[0])
		c, conn := NewClient(serverAddr)
		defer conn.Close()

		ctx := context.Background()
		sr := &shrty.ShortenRequest{URL: args[0]}
		resp, err := c.Shorten(ctx, sr)
		if err != nil {
			grpclog.Fatalf("failed to shorten: %v", err)
		}

		fmt.Printf("Token: %+v ShrtUrl: %v\n", resp.Token, resp.ShrtURL)
	},
}

func init() {
	RootCmd.AddCommand(shortenCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// shortenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// shortenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
