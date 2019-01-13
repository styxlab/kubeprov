package cmd

import (
	"fmt"
	"net"
	"io"

	"github.com/spf13/cobra"
)

var (
	loggingCmd = &cobra.Command	{
		Use:     "logs",
		Aliases: []string{"l"},
		Short:   "connects to the log server and prints log messages",
		Run: printLogMessages,
	}
)

func printLogMessages(cmd *cobra.Command, args []string) {
	
	connectLoop(true)
}

func connectLoop(firstrun bool){

	if firstrun {
		fmt.Print("Waiting for Connection... ")
	}

	conn, err := net.Dial("tcp", "127.0.0.1:9090")
	for err != nil {
		conn, err = net.Dial("tcp", "127.0.0.1:9090")
	}

	if firstrun {
		fmt.Print("connected.")
		fmt.Println(" Press Ctrl+C to stop.")
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

	connectLoop(false)

}