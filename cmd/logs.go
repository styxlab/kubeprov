package cmd

import (
	"fmt"
	"net"
	"io"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	loggingCmd = &cobra.Command	{
		Use:     "logs",
		Aliases: []string{"l"},
		Short:   "prints detailed log messages from kubeprov cluster commands",
		Run: printLogMessages,
	}

	portFlag string
)

func init() {
	loggingCmd.Flags().StringVarP(&portFlag, "port", "p", "9090", "Port for logging")
}

func printLogMessages(cmd *cobra.Command, args []string) {
	
	port, err := strconv.Atoi(portFlag)
	if err != nil {
		fmt.Println("Port number must be an integer.")
		return
	}

	connectLoop(port, true)
}

func connectLoop(port int, firstrun bool){

	endpoint := fmt.Sprintf("127.0.0.1:%d", port)

	if firstrun {
		fmt.Printf("Waiting for connection on endpoint %s ... ", endpoint)
	}

	conn, err := net.Dial("tcp", endpoint)
	for err != nil {
		conn, err = net.Dial("tcp", endpoint)
	}

	if firstrun {
		fmt.Println("connected. Press Ctrl+C to stop.")
	}

	notify := make(chan error)

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				notify <- err
				if io.EOF == err {
					return
				}
			}
			if n > 0 {
				fmt.Printf("%s", buf[:n])
			}
		}
	}()

	select {
	case <-notify:
		break
	}

	connectLoop(port, false)
}