package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net"
	"os"

	"ping/cli_subsonic_client/command"
	"ping/cli_subsonic_client/server"
)

type Config struct {
	username string `yaml:"username"`
	password string `yaml:"password"`
	host     string `yaml:"host"`
}

func main() {
	if len(os.Args[1:]) > 0 && os.Args[1] == "server" {
		var config Config
		var erru bool
		var errp bool
		var errh bool

		config.username, erru = os.LookupEnv("CLI_SUBSONIC_USERNAME")
		config.password, errp = os.LookupEnv("CLI_SUBSONIC_PASSWORD")
		config.host, errh = os.LookupEnv("CLI_SUBSONIC_HOST")
		if !erru || !errp || !errh {
			data, err := ioutil.ReadFile("~/.config/cli_subsonic/config.yaml")
			if err != nil {
				fmt.Println(err)
				return
			}

			err = yaml.Unmarshal(data, &config)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		for {
			server.InitServer(config.username, config.password, config.host)
		}
	} else {
		var values []string
		commandArg := ""

		if len(os.Args) > 1 {
			commandArg = os.Args[1]
		}

		if len(os.Args) > 2 {
			for _, a := range os.Args[2:] {
				values = append(values, a)
			}
		}

		cmd := command.Command{
			commandArg,
			values,
		}
		send(cmd, "127.0.0.1:9000")
	}
}

func send(cmd command.Command, address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Socket dial error: ", err.Error())
		return
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		fmt.Println("Json marshal error: ", err.Error())
		return
	}
	conn.Write(data)
	buffer := make([]byte, 2048)
	l, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Socket read error: ", err.Error())
		return
	}
	if string(buffer[:l]) == "nil" {
		return
	}
	if string(buffer[:l]) != "void" {
		fmt.Println(string(buffer[:l]))
	}
	conn.Close()
}
