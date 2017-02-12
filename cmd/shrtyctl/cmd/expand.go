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

package cmd

import (
	"context"
	"fmt"

	"google.golang.org/grpc/grpclog"

	"git.kono.sh/bkono/shrty"
	"github.com/spf13/cobra"
)

// expandCmd represents the expand command
var expandCmd = &cobra.Command{
	Use:   "expand",
	Short: "Expands a token",
	Long: `Expands a token back to the original url. Counts as a view for the
	servers metrics.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("serverAddr = %+v\n", serverAddr)
		if len(args) == 0 {
			fmt.Println("token required")
			return
		}

		fmt.Printf("expand called for %v\n", args[0])
		c, conn := NewClient(serverAddr)
		defer conn.Close()

		ctx := context.Background()
		er := &shrty.ExpandRequest{Token: args[0]}

		resp, err := c.Expand(ctx, er)
		if err != nil {
			grpclog.Fatalf("failed to expand: %v", err)
		}

		fmt.Printf("url = %+v\n", resp.URL)
	},
}

func init() {
	RootCmd.AddCommand(expandCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// expandCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// expandCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
