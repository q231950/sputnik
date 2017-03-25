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
	"fmt"

	"github.com/q231950/sputnik/keymanager"
	"github.com/spf13/cobra"
)

// identitydeleteCmd represents the identitydelete command
var identitydeleteCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes the signing identity",
	Long: `
	This command is destructive!

	'remove' removes the current signing identity. This makes the key ID in the Cloudkit Dashboard useless. After running this command you should also revoke the key ID from https://icloud.developer.apple.com/dashboard/.`,
	Run: func(cmd *cobra.Command, args []string) {
		keyManager := keymanager.New()
		if keyManager.SigningIdentityExists() {
			removeSigningIdentity(&keyManager)
		} else {
			fmt.Println("There is no signing identity to remove.")
		}
	},
}

func removeSigningIdentity(keyManager keymanager.KeyManager) {
	pub := keyManager.PublicKey()
	err := keyManager.RemoveSigningIdentity()
	if err != nil {
		fmt.Println("An error occurred while removing the signing identity", err)
	} else {
		fmt.Println("Your signing identity has been removed. Make sure to revoke the corresponding KeyID in the Cloudkit Dashboard.")
		fmt.Println("The identity with the following public key was removed:")
		fmt.Println(pub)
	}
}

func init() {
	eckeyCmd.AddCommand(identitydeleteCmd)
}
