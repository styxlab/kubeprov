package ssh

import (
    "fmt"
    "time"
    "net"
)

//Scans for open port
func ScanPort(ip string, port int, interval time.Duration, timeout time.Duration) bool {

    if timeout < interval {
        fmt.Println(port, "closed")
        return false
    }

    target := fmt.Sprintf("%s:%d", ip, port)
    
    conn, err := net.DialTimeout("tcp", target, interval)
    
    if err != nil {
        time.Sleep(interval)
        return ScanPort(ip, port, interval, timeout - interval)
    }

    conn.Close()
    fmt.Println(port, "open")
    return true
}
