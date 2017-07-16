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
	log "github.com/apex/log"

	"github.com/q231950/sputnik/keymanager"
	"github.com/spf13/cobra"
)

// keyidCmd represents the keyid command
var keyidCmd = &cobra.Command{
	Use:   "keyid",
	Short: "Show the key ID that is currently in use",
	Long: `You can provide a Cloudkit key ID by 2 methods
	#1 use the Sputnik command 'keyid store <your key id>'
	#2 by setting an environment variable 'SPUTNIK_CLOUDKIT_KEYID'`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Attempting to retrieve the current key id...")
		keyManager := keymanager.New()
		keyID := keyManager.KeyID()
		if len(keyID) > 0 {
			log.Infof("The following key id is stored: `%s`", keyID)
		} else {
			log.Error("No iCloud key id specified. Please either provide one by `sputnik keyid store <your KeyID>` or set the environment variable `SPUTNIK_CLOUDKIT_KEYID`.")
		}
	},
}

func init() {
	RootCmd.AddCommand(keyidCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// keyidCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// keyidCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
