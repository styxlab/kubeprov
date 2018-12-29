package ssh

import (
    "fmt"

    "time"
    "net"
    "strings"
)

//Scans for open port
func ScanPort(ip string, port int, timeout time.Duration) bool {

    interval := 2 * time.Seconds

    if timeout < interval {
        return false
    }
    
    target := fmt.Sprintf("%s:%d", ip, port)
    
    conn, err := net.DialTimeout("tcp", target, interval)
    
    if err != nil {
        if strings.Contains(err.Error(), "connection refused") {
            time.Sleep(interval)
            ScanPort(ip, port, timeout - intervall)
        } else {
            fmt.Println(port, "closed")
        }
        return false
    }
    
    conn.Close()
    fmt.Println(port, "open")
    return true
}
