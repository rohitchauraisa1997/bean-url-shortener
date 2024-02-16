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
package commands

import (
	"os"
	"sort"
	"strings"

	"backend/routers"

	"github.com/getsentry/sentry-go"
	"github.com/olekukonko/tablewriter"
	bean "github.com/retail-ai-inc/bean/v2"
	"github.com/spf13/cobra"
)

var (
	// routeCmd represents the `route` command.
	routeCmd = &cobra.Command{
		Use:   "route [command]",
		Short: "This command requires a sub command parameter.",
		Long:  "",
	}
)

var (
	// listCmd represents the route list command.
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "Display the route list.",
		Long:  `This command will display all the URI listed in route.go file.`,
		Run:   routeList,
	}
)

func init() {
	routeCmd.AddCommand(listCmd)
	rootCmd.AddCommand(routeCmd)
}

func routeList(cmd *cobra.Command, args []string) {
	// Just initialize a plain sentry client option structure if sentry is on.
	if bean.BeanConfig.Sentry.On {
		bean.BeanConfig.Sentry.ClientOptions = &sentry.ClientOptions{
			Debug:       false,
			Dsn:         bean.BeanConfig.Sentry.Dsn,
			Environment: bean.BeanConfig.Environment,
		}
	}

	// Create a bean object
	b := bean.New()

	// Create an empty database dependency.
	b.DBConn = &bean.DBDeps{}

	// Init different routes.
	routers.Init(b)

	// Consider the allowed methods to display only URI path that's support it.
	allowedMethod := bean.BeanConfig.HTTP.AllowedMethod

	table := tablewriter.NewWriter(os.Stdout)
	header := []string{"Path", "Method", "Handler"}
	table.SetHeader(header)

	for _, r := range b.Echo.Routes() {

		if strings.Contains(r.Name, "glob..func1") {
			continue
		}

		// XXX: IMPORTANT - `allowedMethod` has to be a sorted slice.
		i := sort.SearchStrings(allowedMethod, r.Method)

		if i >= len(allowedMethod) || allowedMethod[i] != r.Method {
			continue
		}

		row := []string{r.Path, r.Method, strings.TrimRight(r.Name, "-fm")}
		table.Append(row)
	}

	table.Render()
}
