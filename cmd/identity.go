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
	"github.com/spf13/cobra"
)

// identityCmd represents the eckey command
var eckeyCmd = &cobra.Command{
	Use:   "identity",
	Short: "Show the signing identity",
	Long:  `Show the signing identity that is used for signing the CloudKit requests.`,
	Run: func(cmd *cobra.Command, args []string) {
		keyManager := keymanager.New()
		keyExists := keyManager.SigningIdentityExists()
		if keyExists {
			identity := keyManager.ECKey()
			fmt.Println(identity)
		} else {
			fmt.Println("A signing identity could not be found. You can create one by `./sputnik identity create`")
		}
	},
}

func init() {
	RootCmd.AddCommand(eckeyCmd)
}
