package ssh

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"golang.org/x/crypto/ssh"
)

//Key hold the key data for secure authentication
type Key struct {
	name           string
	privateKeyPath string
	signer         ssh.Signer
}

//Config holds the config data of the ssh client
type Config struct {
	clientConfig *ssh.ClientConfig
}

//Client holds the connection handle
type Client struct {
	client *ssh.Client
}

//AuthKey reads a private SSH key and creates a signiture
func AuthKey(id string, filePath string) *Key {

	privateKey, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	return &Key{
		name:           id,
		privateKeyPath: filePath,
		signer:         signer,
	}
}

//Config holds the ssh connection configuration
func (key *Key) Config(user string) *Config {

	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(key.signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         60 * time.Second,
	}

	return &Config{
		clientConfig: config,
	}
}

//Client establishes a connection
func (config *Config) Client(ip string, port string) *Client {

	client, err := ssh.Dial("tcp", ip+":"+port, config.clientConfig)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}

	return &Client{
		client: client,
	}
}

func (c *Client) session() *ssh.Session {

	session, err := c.client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}

	return session
}

//Close the connection
func (c *Client) Close() {
	c.client.Close()
}

//RunCmd executes a shell command
func (c *Client) RunCmd(cmd string) string {

	session := c.session()
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b

	if err := session.Run(cmd); err != nil {
		log.Fatal("Failed to run: " + err.Error())
	}
	return b.String()
}

//UploadFile copies a file from a local to a remote machine
func (c *Client) UploadFile(srcFile string, destPath string, executable bool) error {

	permission := "C0644"
	if executable {
		permission = "C0755"
	}

	fileReader, _ := os.Open(srcFile)
	defer fileReader.Close()

	contentsBytes, _ := ioutil.ReadAll(fileReader)
	size := int64(len(contentsBytes))

	r := bytes.NewReader(contentsBytes)

	session := c.session()
	defer session.Close()

	var b bytes.Buffer
	session.Stderr = &b

	go func() {
		w, _ := session.StdinPipe()
		fmt.Fprintln(w, permission, size, path.Base(srcFile))
		io.Copy(w, r)
		fmt.Fprint(w, "\x00")
		w.Close()
	}()

	if err := session.Run("/usr/bin/scp -t " + destPath); err != nil {
		fmt.Println(b.String())
		log.Fatalf("write failed:%v", err.Error())
	}

	return nil
}

//Scans for open port
func ScanPort(ip string, port int, timeout time.Duration) bool {
    target := fmt.Sprintf("%s:%d", ip, port)
    conn, err := net.DialTimeout("tcp", target, timeout)
    
    if err != nil {
        if strings.Contains(err.Error(), "too many open files") {
            time.Sleep(timeout)
            ScanPort(ip, port, timeout)
        } else {
            fmt.Println(port, "closed")
        }
        return false
    }
    
    conn.Close()
    fmt.Println(port, "open")
    return true
}
