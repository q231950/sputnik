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
	log "github.com/apex/log"

	"github.com/q231950/sputnik/keymanager"
	"github.com/spf13/cobra"
)

// The identity command shows the current identity.
var eckeyCmd = &cobra.Command{
	Use:   "identity",
	Short: "Show the signing identity",
	Long:  `Show the signing identity that is used for signing the iCloud requests.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("Attempting to retrieve the current identity...")
		keyManager := keymanager.New()
		keyExists, err := keyManager.SigningIdentityExists()
		if err != nil {
			log.Errorf("Error in SigningIdentityExists: %s", err)
		}

		if keyExists {
			identity := keyManager.ECKey()
			log.Debug("The current identity you can create a new server-to-server key with in the iCloud Dashboard:")
			log.Infof("\n%s", identity)

			keyID := keyManager.KeyID()
			if len(keyID) == 0 {
				log.Error("No iCloud KeyID specified. Please either provide one by `sputnik keyid store <your KeyID>` or set the environment variable `SPUTNIK_CLOUDKIT_KEYID`.")
			}
		} else {
			log.Error("A signing identity could not be found. A signing identity can be created by `./sputnik identity create`")
		}
	},
}

func init() {
	RootCmd.AddCommand(eckeyCmd)
}
