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

// printCmd represents the print command
var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Print the private key",
	Long:  `Use this command to get the private key to paste into the CloudKit Dashboard when granting API access.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Attempting to print the private key.")

		keyManager := keymanager.New()
		exists, err := keyManager.SigningIdentityExists()
		if err != nil {
			log.Errorf("Error in SigningIdentityExists: %s", err)
		}
		if exists {
			log.Infof("Printing the public/private keys:\n%s", keyManager.PublicKeyString())
		} else {
			log.Info("The ec key does not exist, need to create, one moment, please")
			keyManager.CreateSigningIdentity()

			log.Infof("Ok done. This is it: \n%s", keyManager.PublicKeyString())
		}
	},
}

func init() {
	eckeyCmd.AddCommand(printCmd)
}
