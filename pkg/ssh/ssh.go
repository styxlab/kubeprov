package ssh

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"golang.org/x/crypto/ssh"
)

// Key holds the key data for secure authentication
type Key struct {
	name           string
	privateKeyPath string
	signer         ssh.Signer
}

// Config holds the config data of the ssh client
type Config struct {
	clientConfig *ssh.ClientConfig
}

// Client holds the connection handle
type Client struct {
	address string
	port int
	client *ssh.Client
}

// AuthKey reads a private SSH key and creates a signiture
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

// Config holds the ssh connection configuration
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

// Client establishes a connection
func (config *Config) Client(ip string, port int) *Client {

	interval := 2 * time.Second
	timeout := 60 * time.Second

	log.Printf("Wait for open port...")
	if err := WaitForOpenPort(ip, port, interval, timeout); err != nil {
		log.Fatal("Port is closed. Check your firewall.")
	}

	endpoint := fmt.Sprintf("%s:%d", ip, port);
	client, err := ssh.Dial("tcp", endpoint, config.clientConfig)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}

	return &Client{
		address: ip,
		port: port,
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

// Close the connection
func (c *Client) Close() {
	c.client.Close()
}

// RunCmd executes a shell command on the remote server
func (c *Client) RunCmd(cmd string) string {

	session := c.session()
	defer session.Close()

	var b bytes.Buffer
	var e bytes.Buffer
	session.Stdout = &b
	session.Stderr = &e

	if err := session.Run(cmd); err != nil {
		log.Fatalf("Failed to run command:%s\nstdout:%s\nstderr:%v ", cmd, b.String(), err)
	}
	return b.String()
}

//UploadFile copies a file from a local to a remote machine
func (c *Client) UploadFile(localFile string, remotePath string, executable bool) error {

	permission := "C0644"
	if executable {
		permission = "C0755"
	}

	fileReader, _ := os.Open(localFile)
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
		fmt.Fprintln(w, permission, size, path.Base(localFile))
		io.Copy(w, r)
		fmt.Fprint(w, "\x00")
		w.Close()
	}()

	if err := session.Run("/usr/bin/scp -t " + remotePath); err != nil {
		fmt.Println(b.String())
		log.Fatalf("write failed:%v", err.Error())
	}

	return nil
}

//UploadFile copies a file from a local to a remote machine
func (c *Client) DownloadFile(remoteFile string, localPath string, useSudo bool) error {

	filename := path.Base(remoteFile)
 	localFile := filepath.Join(localPath, filename)

	fileReader, _ := os.Create(localFile)
	defer fileReader.Close()

	cmd := fmt.Sprintf("dd if=%s", remoteFile)
	if useSudo {
		cmd = fmt.Sprintf("sudo %s", cmd)
	}

	session := c.session()
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b

	if err := session.Run(cmd); err != nil {
		log.Fatal("Failed to run: " + err.Error())
	}

	_, err := fileReader.Write(b.Bytes())
	if err != nil {
		log.Fatalf("write failed:%v", err)
	}

	return nil
}

