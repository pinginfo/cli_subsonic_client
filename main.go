package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"ping/cli_subsonic_client/command"
	"ping/cli_subsonic_client/server"
)

func main() {
	if len(os.Args[1:]) > 0 && os.Args[1] == "server" {
		username := os.Getenv("CLI_SUBSONIC_USERNAME")
		password := os.Getenv("CLI_SUBSONIC_PASSWORD")
		host := os.Getenv("CLI_SUBSONIC_HOST")
		server.InitServer(username, password, host)
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
