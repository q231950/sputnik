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

	"github.com/spf13/cobra"
	"github.com/q231950/sputnik/eckeyhandling"
)

// eckeyCmd represents the eckey command
var eckeyCmd = &cobra.Command{
	Use:   "eckey",
	Short: "Read the ec key",
	Long: `The ec key is used for server to server communication between CloudKit and everyone else.`,
	Run: func(cmd *cobra.Command, args []string) {
		keyExists := eckeyhandling.ECKeyExists()
		if keyExists {
			_ = eckeyhandling.ECKey()
		} else {
			fmt.Println("The ec key does not exist, need to create one... I'll do this for you...\n")
			createErr := eckeyhandling.CreateECKey()
			if createErr != nil {
				fmt.Println("Sorry, failed to create the ec key\n")
			} else {
				path, _ := eckeyhandling.SecretsFolder()
				fmt.Println("Ok, this is your key. It's named eckey.pem and located under", path, "\n")
				_ = eckeyhandling.ECKey()
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(eckeyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// eckeyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// eckeyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
