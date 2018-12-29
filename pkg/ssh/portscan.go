package ssh

import (
    "fmt"
    "time"
    "net"
)

//Scans for open port
func ScanPort(ip string, port int, interval time.Duration, timeout time.Duration) bool {

    if timeout < interval {
        return false
    }

    target := fmt.Sprintf("%s:%d", ip, port)
    
    conn, err := net.DialTimeout("tcp", target, interval)
    
    if err != nil {
        fmt.Println(port, "closed")
        time.Sleep(interval)
        ScanPort(ip, port, timeout - intervall)
        return false
    }

    conn.Close()
    fmt.Println(port, "open")
    return true
}
