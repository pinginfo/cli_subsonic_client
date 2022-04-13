package server

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"ping/cli_subsonic_client/api"
	"ping/cli_subsonic_client/command"
	"ping/cli_subsonic_client/player"
	"strconv"
)

var (
	connection *api.SubsonicConnection
	myPlayer   *player.Player
)

func InitServer(username string, password string, host string) {
	var err error
	myPlayer, err = player.InitPlayer()
	if err != nil {
		fmt.Println("Player init error: ", err.Error())
		return
	}
	connection = &api.SubsonicConnection{
		Username:       username,
		Password:       password,
		Host:           host,
		DirectoryCache: make(map[string]api.SubsonicResponse),
	}

	socket, err := net.Listen("tcp", "0.0.0.0:9000")
	if err != nil {
		fmt.Println("Socket error: ", err.Error())
		os.Exit(0)
	}
	defer socket.Close()

	fmt.Println("Server listing 0.0.0.0:9000")

	for {
		conn, err := socket.Accept()
		if err != nil {
			fmt.Println("Socket accept error: ", err.Error())
			return
		}
		go handleCommand(conn)
	}
}

func handleCommand(conn net.Conn) {
	buffer := make([]byte, 2048)
	l, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Socket read error: ", err.Error())
		return
	}

	var cmd command.Command
	err = json.Unmarshal(buffer[:l], &cmd)
	if err != nil {
		fmt.Println("Json unmarshal error: ", err.Error())
		return
	}

	if cmd.Command == "test" {
		conn.Write([]byte("nil"))
	} else if cmd.Command == "volume" {
		if len(cmd.Values) > 0 {
			inc, err := strconv.ParseInt(cmd.Values[0], 10, 64)
			if err != nil {
				conn.Write([]byte("Bad volume value"))
				return
			}
			vol := strconv.FormatInt(myPlayer.AdjustVolume(inc), 10)
			conn.Write([]byte(vol))
		} else {
			vol := strconv.FormatInt(myPlayer.AdjustVolume(0), 10)
			conn.Write([]byte(vol))
		}
	} else if cmd.Command == "search" {
		if cmd.Values[0] == "album" {
			albumId, err := connection.GetAlbumIdByName(cmd.Values[1])
			if err != nil {
				conn.Write([]byte(err.Error()))
				return
			}
			musics, err := connection.GetMusicsFromAlbumId(*albumId)
			if err != nil {
				conn.Write([]byte(err.Error()))
				return
			}
			conn.Write([]byte(getNamesMusics(musics)))
		} else {
			conn.Write([]byte("not implemented"))
		}
	} else if cmd.Command == "add" {
		if cmd.Values[0] == "album" {
			albumId, err := connection.GetAlbumIdByName(cmd.Values[1])
			if err != nil {
				conn.Write([]byte(err.Error()))
				return
			}
			musics, err := connection.GetMusicsFromAlbumId(*albumId)
			if err != nil {
				conn.Write([]byte(err.Error()))
				return
			}
			for _, m := range musics {
				myPlayer.AddInQueue(connection.GetPlayUrl(&m), m.Title, m.Artist, m.Duration)
			}
			conn.Write([]byte(getNamesMusics(musics)))
		} else if cmd.Values[0] == "playlist" {
			musics, err := connection.GetMusicsPlaylistByName(cmd.Values[1])
			if err != nil {
				conn.Write([]byte(err.Error()))
				return
			}
			for _, m := range musics {
				myPlayer.AddInQueue(connection.GetPlayUrl(&m), m.Title, m.Artist, m.Duration)
			}
			conn.Write([]byte(getNamesMusics(musics)))
		} else {
			conn.Write([]byte("not implemented"))
		}
	} else if cmd.Command == "queued" {
		var str string
		for _, m := range myPlayer.Queue {
			str += m.Title + "\n"
		}
		conn.Write([]byte(str))
	} else if cmd.Command == "current" {
		conn.Write([]byte(myPlayer.Actual()))
	} else if cmd.Command == "next" {
		myPlayer.PlayNextTrack()
		conn.Write([]byte(myPlayer.Actual()))
	} else if cmd.Command == "pause" {
		myPlayer.Pause()
		conn.Write([]byte(myPlayer.Actual()))
	} else if cmd.Command == "play" {
		myPlayer.Play()
		conn.Write([]byte(myPlayer.Actual()))
	} else if cmd.Command == "prev" {
		myPlayer.PlayPrevTrack()
		conn.Write([]byte(myPlayer.Actual()))
	} else if cmd.Command == "random" {
		if cmd.Values[0] == "off" {
			myPlayer.SetRandom(false)
		} else {
			myPlayer.SetRandom(true)
		}
		conn.Write([]byte(myPlayer.Actual()))
	} else if cmd.Command == "repeat" {
		if cmd.Values[0] == "off" {
			myPlayer.SetRepeat(false)
		} else {
			myPlayer.SetRepeat(true)
		}
		conn.Write([]byte(myPlayer.Actual()))
	} else if cmd.Command == "stop" {
		myPlayer.Stop()
		conn.Write([]byte(myPlayer.Actual()))
	} else if cmd.Command == "status" {
		conn.Write([]byte(myPlayer.Actual()))
	} else if cmd.Command == "clean" {
		myPlayer.CleanQueue()
		conn.Write([]byte(myPlayer.Actual()))
	} else if cmd.Command == "favorite" {
		saveSongInFile(myPlayer.Actual())
		conn.Write([]byte("song added into favorite file"))
	} else {
		if myPlayer.IsPaused() {
			conn.Write([]byte("nil"))
		} else {
			str := myPlayer.Actual()
			if len(str) > 20 {
				str = str[:20]
			}
			conn.Write([]byte(str))
		}
	}
}

func getNamesMusics(musics []api.SubsonicEntity) string {
	str := ""
	for _, music := range musics {
		str += music.Title + "\n"
	}
	return str
}

func saveSongInFile(name string) {
	f, err := os.OpenFile(os.Getenv("HOME")+"/.favorites_song", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, err = fmt.Fprintln(f, name)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return
}
