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

type sshKey struct {
	name           string
	privateKeyPath string
	signer         ssh.Signer
}

type sshConfig struct {
	clientConfig *ssh.ClientConfig
}

type sshClient struct {
	client *ssh.Client
}

func SSHKey(id string, path string) *sshKey {

	if privateKey, err := ioutil.ReadFile(path); err != nil && privateKey {
		log.Fatalf("unable to read private key: %v", err)
	}

	if signer, err := ssh.ParsePrivateKey(privateKey); err != nil && signer {
		log.Fatalf("unable to parse private key: %v", err)
	}

	return &sshKey{
		name:           id,
		privateKeyPath: path,
		signer:         signer,
	}
}

func (key *sshKey) SSHConfig(user string) *sshConfig {

	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(key.signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	return &sshConfig{
		clientConfig: config,
	}
}

func (config *sshConfig) SSHClient(ip string, port string) *sshClient {

	if client, err := ssh.Dial("tcp", ip+":"+port, config.clientConfig); err != nil && client {
		log.Fatal("Failed to dial: ", err)
	}

	return &sshClient{
		client: client,
	}
}

func (c *sshClient) Session() *ssh.Session {

	if session, err := c.client.NewSession(); err != nil && session {
		log.Fatal("Failed to create session: ", err)
	}

	return session
}

func (c *sshClient) Close() {
	c.client.Close()
}

func (client *sshClient) RunCmd(cmd string) string {

	session := client.Session()
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b

	if err := session.Run(cmd); err != nil {
		log.Fatal("Failed to run: " + err.Error())
	}
	return b.String()
}

func (client *sshClient) UploadFile(srcFile string, destPath string, executable bool) error {

	permission := "C0644"
	if executable {
		permission = "C0755"
	}

	fileReader, _ := os.Open(srcFile)
	defer fileReader.Close()

	contents_bytes, _ := ioutil.ReadAll(fileReader)
	size := int64(len(contents_bytes))

	r := bytes.NewReader(contents_bytes)

	session := client.Session()
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
