package hetzner

import (
	"context"
	"os/user"
	"path/filepath"
	"log"
	"os"
	"strings"
	"io/ioutil"

	"github.com/spf13/viper"
)

const EnvForToken = "HCLOUD_TOKEN"
const ConfigType  = "json"
const SSHProtocol = "ed25519"
const DefaultKeyName = "default"

var DefaultConfigName string
var DefaultConfigRoot string
var DefaultConfigPath string
var DefaultConfigFile string

var publicKeyPath string
var publicKeyFile string

func init(){
	usr, err := user.Current()
	if err != nil {
		return
	}
	if usr.HomeDir != "" {
		DefaultConfigName = "hcloud"
		DefaultConfigRoot = filepath.Join(usr.HomeDir, ".config", "kubeprov")
		DefaultConfigPath = filepath.Join(DefaultConfigRoot, DefaultConfigName)
		DefaultConfigFile = DefaultConfigPath + "." + ConfigType

		//create empty file, if not existing
		createDirIfNotExist(DefaultConfigRoot)
		createFileIfNotExist(DefaultConfigFile)

		publicKeyPath = filepath.Join(usr.HomeDir, ".ssh")
		publicKeyFile = filepath.Join(publicKeyPath, "id_" + SSHProtocol + ".pub")
	}
}

type Config struct {
	context context.Context
	token string
	publicKeyName string
	publicKey string
	configPath string
}

func GetOrCreateConfig() *Config {

	c := &Config { 
		context: context.Background(),
		configPath: DefaultConfigPath,
	}

	//First search in Config File
	viper.SetConfigName(DefaultConfigName) 
	viper.AddConfigPath(DefaultConfigRoot)
	viper.SetConfigType(ConfigType)

	err := viper.ReadInConfig()
	if err == nil {
		c.token = viper.GetString("token")
		c.publicKey = viper.GetString("sshkey")
		if len(c.publicKey) == 0 {
			log.Fatalf("No public key found in config. Please delete and recreate the file at %s.\n", DefaultConfigFile)
		}
		c.publicKeyName = viper.GetString("keyname")
		if len(c.publicKeyName) == 0 {
			log.Fatalf("No public key name found in config. Please delete and recreate the file at %s.\n", DefaultConfigFile)
		}
		return c
	}
	log.Printf("unable to read config file %q: %s\n", DefaultConfigFile, err)

	c.publicKey = readPublicKey(publicKeyFile)
	c.publicKeyName = deriveNameFromPublicKey(c.publicKey)
	log.Printf("keyname %s\n", c.publicKeyName)
	
	//Second search in environment
	if token, ok := os.LookupEnv(EnvForToken); ok {
		token = strings.TrimSpace(token)
		viper.Set("token", token)
		viper.Set("sshkey", c.publicKey)
		viper.Set("keyname", c.publicKeyName)
		viper.WriteConfig()
		c.token = token
		return c
	}

	log.Fatalf("no %s environment variable set.\n", EnvForToken)
	return nil
}

func createDirIfNotExist(dir string) {
    if _, err := os.Stat(dir); os.IsNotExist(err) {
    	if err = os.MkdirAll(dir, 0755); err != nil {
            log.Fatalf("unable to create config directory. please check your directory permissions.\n")
        }
    }
}

func createFileIfNotExist(file string) {
	if _, err := os.OpenFile(file, os.O_RDONLY|os.O_CREATE, 0666); err != nil{
		log.Fatalf("unable to create file. please check your directory permissions.\n")
	}
}

func readPublicKey(file string) string {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Could not read public key file: %s. Need to generate a new key-pair with ssh-keygen?\n", file)
	}
	return strings.TrimSpace(string(data))
}

func deriveNameFromPublicKey(publicKey string) string {
	stringSlice := strings.Split(publicKey, " ")
	name := lastString(stringSlice)
	if len(name) < 1 {
		name = DefaultKeyName
	}
	return name
}

func lastString(ss []string) string {
    return ss[len(ss)-1]
}