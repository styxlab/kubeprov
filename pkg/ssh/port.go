package ssh

import (
    "fmt"
    "errors"
    "time"
    "net"
)

// WaitForOpenPort scans and waits for open port until timeout
func WaitForOpenPort(ip string, port int, interval time.Duration, timeout time.Duration) error {

    if timeout < interval {
        return errors.New("port closed")
    }

    target := fmt.Sprintf("%s:%d", ip, port)
    
    conn, err := net.DialTimeout("tcp", target, interval)    
    if err != nil {
        time.Sleep(interval)
        return WaitForOpenPort(ip, port, interval, timeout - interval)
    }

    conn.Close()
    return nil
}
