// MIT License

// Copyright (c) The RAI Authors

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// ./backend gopher helloworld

package gopher

import (
	"errors"

	bean "github.com/retail-ai-inc/bean/v2"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "helloworld",
		Short: "Hello world!",
		Long:  `This command will just print hello world otherwise hello mars`,
		RunE: func(cmd *cobra.Command, args []string) error {
			NewHellowWorld()
			err := helloWorld("hello")
			if err != nil {
				// If you turn on `sentry` via env.json then the error will be captured by sentry otherwise ignore.
				bean.SentryCaptureException(nil, err)
			}

			return err
		},
	}

	GopherCmd.AddCommand(cmd)
}

func NewHellowWorld() {
	// IMPORTANT: If you pass `false` then database connection will not be initialized.
	_ = initBean(false)
}

func helloWorld(h string) error {
	if h == "hello" {
		return nil
	}

	return errors.New("hello mars")
}
