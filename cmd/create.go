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

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new signing identity",
	Long:  `For now, a file named eckey.pem will be put into the secrets folder.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("Attempting to create a new identity...")
		createECKey()
	},
}

func init() {
	eckeyCmd.AddCommand(createCmd)
}

func createECKey() {
	keyManager := keymanager.New()
	exists, err := keyManager.SigningIdentityExists()
	if err != nil {
		log.Errorf("Error in SigningIdentityExists: %s", err)
	}

	if exists {
		log.Error("There is an existing identity. You need to remove it with `./sputnik identity remove` before you can create a new one.")
		if len(keyManager.KeyID()) != 0 {
			// a key ID has been stored
			log.Infof("The current identity is linked with the following iCloud key ID:\n%s", keyManager.KeyID())
		} else {
			log.Infof("This is the current identity:\n%s", keyManager.ECKey())
		}
	} else {
		log.Info("Creating an identity")
		keyManager.CreateSigningIdentity()
		log.Info("Done")
	}
}
