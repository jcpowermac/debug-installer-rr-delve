package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/go-delve/delve/service/api"
	"github.com/sirupsen/logrus"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-delve/delve/service/rpc2"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/wait"
)

var (
	wg           sync.WaitGroup
	functionName string
	address      string
)

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

func detach(client *rpc2.RPCClient) {
	err := client.Detach(true)
	if err != nil {
		exitWithError(err)
	}
}

func initCliFlags() {
	var rootCmd = &cobra.Command{
		Use:     "fterraform",
		Short:   "fterraform finds terraform provider functions within a rr recording",
		GroupID: "",
		Long:    ``,
		Example: "",
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}
	rootCmd.PersistentFlags().StringVar(&functionName, "function", "", "name of function in terraform provider")
	rootCmd.PersistentFlags().StringVar(&address, "address", "127.0.0.1:2345", "address of the delve remote instance")

	err := rootCmd.MarkPersistentFlagRequired("function")

	rootCmd.Flag("function").Shorthand = "f"
	rootCmd.Flag("address").Shorthand = "a"

	if err != nil {
		exitWithError(rootCmd.Help())
	}

	if err := rootCmd.Execute(); err != nil {
		exitWithError(err)
	}
}

func waitForAddress(address string) (net.Conn, error) {
	var conn net.Conn
	timeout := time.Minute * 2
	waitContext, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()

	fmt.Print("Connecting")

	wait.Until(func() {
		var err error
		conn, err = net.DialTimeout("tcp", address, timeout)
		if err == nil {
			cancel()
		}
		if err != nil {
			fmt.Print(".")
		}
	}, time.Second*2, waitContext.Done())
	fmt.Print("\n")

	return conn, waitContext.Err()
}

func createBreakpoint(client *rpc2.RPCClient, functionName string) {
	breakpoint := api.Breakpoint{
		Name:         "findTerraformFunction",
		FunctionName: functionName,
	}

	_, err := client.CreateBreakpoint(&breakpoint)
	if err != nil {
		exitWithError(err)
	}
}

func findFunctionByName(client *rpc2.RPCClient, functionName string) (string, error) {
	functionNameRegex := fmt.Sprintf(".*%s", functionName)
	functions, err := client.ListFunctions(functionNameRegex)
	if err != nil {
		return "", err
	}

	if len(functions) == 0 {
		return "", errors.New("no functions found")
	}

	if len(functions) > 1 {
		logrus.Warnf("more than one function available %v", functions)
	}
	return functions[0], nil
}

func main() {
	initCliFlags()

	wg.Add(1)

	conn, err := waitForAddress(address)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			exitWithError(err)
		}
	}

	client := rpc2.NewClientFromConn(conn)

	defer detach(client)

	if !client.Recorded() {
		exitWithError(errors.New("expecting this to be used with rr recording"))
	}

	foundFunctionName, err := findFunctionByName(client, functionName)
	if err != nil {
		exitWithError(err)
	}

	createBreakpoint(client, foundFunctionName)

	found := false

	go func() {
		once := true
		defer wg.Done()
		ch := client.Continue()
		for {
			state := <-ch
			if state.Running {
				time.Sleep(time.Second * 1)
				continue
			} else if once {
				for _, t := range state.Threads {
					if t != nil {
						if t.Function != nil {
							if t.Breakpoint != nil {
								if t.Function.Name() == foundFunctionName {
									found = true
									spew.Dump(t)
									return
								}
							}
						}
					}
				}
				once = false
			}
			if state.Exited {
				found = false
				return
			}
		}
	}()

	wg.Wait()

	if found {
		fmt.Printf("found thread with function %s", foundFunctionName)
	}
}
