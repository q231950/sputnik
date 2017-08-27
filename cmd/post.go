// Copyright © 2017 Martin Kim Dung-Pham <kim@elbedev.com>
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
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/apex/log"
	"github.com/q231950/sputnik/keymanager"
	"github.com/q231950/sputnik/requesthandling"
	"github.com/spf13/cobra"
)

// postCmd represents the post command
var postCmd = &cobra.Command{
	Use:   "post",
	Short: "post allows you to send post requests with a given payload",
	Long: `post allows you to send post requests with a given payload:

	Here, a payload from file is post'ed with the operation modify
	./sputnik requests post --operation "modify" --payload-file-path 'path/to/file.json'
	./sputnik requests post -o "modify" --pf 'path/to/file.json'

	./sputnik requests post --operation "modify" --payload '<json payload>'
	./sputnik requests post -o "modify" -p '<json payload>'

`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithFields(log.Fields{
			"Operation": operation}).Info("Attempting to POST...")

		if operation == "" {
			log.Error("Missing operation, please provide one. See `sputnik help requests post`")
		}

		var payloadToUse string
		if payloadFilePath != "" {
			payloadToUse = payloadFromFile(payloadFilePath)
			log.WithFields(log.Fields{
				"Payload": payloadToUse}).Info("Payload from file")

		} else {
			payloadToUse = payload
		}

		keyManager := keymanager.New()

		if container != "" {
			config := requesthandling.RequestConfig{Version: "1", Database: "public", ContainerID: container}
			requestManager := requesthandling.New(config, &keyManager)

			request, err := requestManager.PostRequest(operation, payloadToUse)
			if err != nil {
				log.Error(err.Error())
			} else {
				client := http.Client{}
				resp, err := client.Do(request)
				if err != nil {
					panic(err)
				} else {
					body, _ := ioutil.ReadAll(resp.Body)
					s, _ := json.MarshalIndent(string(body), "", "    ")
					log.Info(string(s))
				}
			}
		} else {
			log.Error("Missing container, please provide one. See `sputnik help requests post`")
		}
	},
}

func init() {
	requestsCmd.AddCommand(postCmd)
	postCmd.Flags().StringVarP(&payloadFilePath, "json-file-path", "j", "", "A path to a file that contains the json payload")
	postCmd.Flags().StringVarP(&payload, "payload", "p", "", "A json payload as string")
	postCmd.Flags().StringVarP(&operation, "operation", "o", "", "The operation to execute: Depending on your intention, operation-specific subpaths may be of [modify, query, lookup, changes, resolve, accept]")
	postCmd.Flags().StringVarP(&container, "container", "c", "", "The CloudKit container to access. (normally `iCloud.your.bundle.identifier`)")
}
