package api

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
)

// used for generating salt
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func authToken(password string) (string, string) {
	salt := randSeq(8)
	token := fmt.Sprintf("%x", md5.Sum([]byte(password+salt)))

	return token, salt
}

func defaultQuery(connection *SubsonicConnection) url.Values {
	token, salt := authToken(connection.Password)
	query := url.Values{}
	query.Set("u", connection.Username)
	query.Set("t", token)
	query.Set("s", salt)
	query.Set("v", "1.15.1")
	query.Set("c", "stmp")
	query.Set("f", "json")

	return query
}

// requests
func (connection *SubsonicConnection) GetServerInfo() (*SubsonicResponse, error) {
	query := defaultQuery(connection)
	requestUrl := connection.Host + "/rest/ping" + "?" + query.Encode()
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

func (connection *SubsonicConnection) GetIndexes() (*SubsonicResponse, error) {
	query := defaultQuery(connection)
	requestUrl := connection.Host + "/rest/getIndexes" + "?" + query.Encode()
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

func (connection *SubsonicConnection) GetMusicDirectory(id string) (*SubsonicResponse, error) {
	if cachedResponse, present := connection.DirectoryCache[id]; present {
		return &cachedResponse, nil
	}

	query := defaultQuery(connection)
	query.Set("id", id)
	requestUrl := connection.Host + "/rest/getMusicDirectory" + "?" + query.Encode()
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

	// on a sucessful request, cache the response
	if decodedBody.Response.Status == "ok" {
		connection.DirectoryCache[id] = decodedBody.Response
	}

	return &decodedBody.Response, nil
}

// note that this function does not make a request, it just formats the play url
// to pass to mpv
func (connection *SubsonicConnection) GetPlayUrl(entity *SubsonicEntity) string {
	// we don't want to call stream on a directory
	if entity.IsDirectory {
		return ""
	}

	query := defaultQuery(connection)
	query.Set("id", entity.Id)
	return connection.Host + "/rest/stream" + "?" + query.Encode()
}
