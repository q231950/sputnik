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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		keyManager := keymanager.New()
		exists := keyManager.SigningIdentityExists()
		if exists {
			log.Debugf("Printing the public/private keys:\n%s", keyManager.PrivatePublicKeyWriter())
		} else {
			log.Debug("The ec key does not exist, need to create, one moment, please")
			keyManager.CreateSigningIdentity()

			log.Debugf("Ok done. This is it: \n%s", keyManager.PrivatePublicKeyWriter())
		}
	},
}

func init() {
	eckeyCmd.AddCommand(printCmd)
}
