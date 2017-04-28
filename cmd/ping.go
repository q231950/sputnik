// Copyright © 2016 Martin Kim Dung-Pham <kim@elbedev.com>
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

package cmd

import (
	"fmt"
	"github.com/q231950/sputnik/keymanager"
	"github.com/q231950/sputnik/requesthandling"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/http"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Send a request to cloudkit if at least something works",
	Long: `Ping creates a GET request and sends it off`,
	Run: func(cmd *cobra.Command, args []string) {
		keyManager := keymanager.New()
		// subpath := "records/query"
		// subpath := "users/lookup/email"
		// subpath := "users/caller"
		subpath := "records/modify"
		database := "public"
		requestManager := requesthandling.New(&keyManager, database, subpath)
		request, err := requestManager.PingRequest()
		if err == nil {
			fmt.Println(request)
		} else {
			log.Fatal("Failed to create ping request")
		}

		client := &http.Client{}
    resp, err := client.Do(request)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))
	},
}

func init() {
	RootCmd.AddCommand(pingCmd)
}
