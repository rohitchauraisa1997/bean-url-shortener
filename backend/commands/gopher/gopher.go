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
package gopher

import (
	"sync"

	// msentry "backend/packages/sentry"
	// "backend/repositories"
	// "backend/services"

	"github.com/getsentry/sentry-go"
	bean "github.com/retail-ai-inc/bean/v2"
	"github.com/retail-ai-inc/bean/v2/helpers"
	"github.com/spf13/cobra"
)

type gopherRepositories struct {
	// TODO: You can add your repositories here

	// helloRepo        repositories.HelloRepository
	// worldRepo        repositories.WorldRepository
}

type gopherServices struct {
	// TODO: You can add your services here

	// helloSvc         services.HelloService
	// worldSvc         services.WorldService
}

var GopherCmd = &cobra.Command{
	Use:   "gopher [command]",
	Short: "This command requires a sub command parameter of your own.",
	Long:  "This command requires a sub command parameter. You can create a new sub command by creating a new go file under `gopher` directory. An example can be found here: https://github.com/retail-ai-inc/bean#make-your-own-commands",
}

var b *bean.Bean
var once sync.Once

func initBean(isInitDB ...bool) *bean.Bean {
	once.Do(func() {
		// Prepare sentry options before initialize bean.
		if bean.BeanConfig.Sentry.On {
			bean.BeanConfig.Sentry.ClientOptions = &sentry.ClientOptions{
				Debug:       bean.BeanConfig.Sentry.Debug,
				Dsn:         bean.BeanConfig.Sentry.Dsn,
				Environment: bean.BeanConfig.Environment,
				//BeforeSend:       msentry.CustomBeforeSend, // Default beforeSend function. You can initialize your own custom function.
				AttachStacktrace: true,
				TracesSampleRate: helpers.FloatInRange(bean.BeanConfig.Sentry.TracesSampleRate, 0.0, 1.0),
			}
		}

		// Create a bean object
		b = bean.New()
		b.Echo.AcquireContext().Reset(nil, nil)

		// IMPORTANT - This is very useful when you run some cloudfunction/command in GCP/AWS/Azure and you cannot connect
		// your memorystore/SQL/mongo server from cloudfunction/VM using usual `host` ip. Therfore, you can set
		// a different host ip by setting a different host parameter inside your `TenantConnections` table under
		// `Connections` JSON.
		bean.TenantAlterDbHostParam = ""

		// Init DB dependency.
		if len(isInitDB) == 0 || isInitDB[0] == true {
			b.InitDB()
		}
	})

	return b
}
