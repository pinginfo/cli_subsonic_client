package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (connection *SubsonicConnection) GetMusicsPlaylistByName(name string) ([]SubsonicEntity, error) {
	playlists, err := connection.getPlaylists()
	if err != nil {
		return nil, err
	}

	for _, playlist := range playlists.Playlists.Playlist {
		if playlist.Name == name {
			musics, err := connection.getMusicsPlaylistById(playlist.Id)
			if err != nil {
				return nil, err
			}
			return musics, nil
		}
	}
	return nil, errors.New("Playlist not found")
}

func (connection *SubsonicConnection) getPlaylists() (*SubsonicResponse, error) {
	query := defaultQuery(connection)
	requestUrl := connection.Host + "/rest/getPlaylists" + "?" + query.Encode()
	fmt.Println(requestUrl)
	res, err := http.Get(requestUrl)
	if err != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	responseBody, readErr := ioutil.ReadAll(res.Body)

	if readErr != nil {
		return nil, err
	}

	fmt.Println(responseBody)
	var decodedBody responseWrapper
	err = json.Unmarshal(responseBody, &decodedBody)

	if err != nil {
		return nil, err
	}

	return &decodedBody.Response, nil
}

func (connection *SubsonicConnection) getPlaylist(id string) (*SubsonicResponse, error) {
	query := defaultQuery(connection)
	query.Set("id", id)
	requestUrl := connection.Host + "/rest/getPlaylist" + "?" + query.Encode()
	fmt.Println(requestUrl)
	res, err := http.Get(requestUrl)
	if err != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	responseBody, readErr := ioutil.ReadAll(res.Body)

	if readErr != nil {
		return nil, err
	}

	var decodedBody responseWrapper
	err = json.Unmarshal(responseBody, &decodedBody)

	if err != nil {
		return nil, err
	}

	return &decodedBody.Response, nil
}

func (connection *SubsonicConnection) getMusicsPlaylistById(id string) ([]SubsonicEntity, error) {
	playlist, err := connection.getPlaylist(id)
	if err != nil {
		return nil, err
	}
	return playlist.Playlist.Entities, nil
}
