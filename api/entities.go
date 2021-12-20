package api

type SubsonicConnection struct {
	Username       string
	Password       string
	Host           string
	DirectoryCache map[string]SubsonicResponse
}

type SubsonicError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SubsonicArtist struct {
	Id         string
	Name       string
	AlbumCount int
}

type SubsonicDirectory struct {
	Id       string           `json:"id"`
	Parent   string           `json:"parent"`
	Name     string           `json:"name"`
	Entities []SubsonicEntity `json:"child"`
}

type SubsonicEntity struct {
	Id          string `json:"id"`
	IsDirectory bool   `json:"isDir"`
	Parent      string `json:"parent"`
	Title       string `json:"title"`
	Artist      string `json:"artist"`
	Duration    int    `json:"duration"`
	Track       int    `json:"track"`
	DiskNumber  int    `json:"diskNumber"`
	Path        string `json:"path"`
}

type SubsonicIndexes struct {
	Index []SubsonicIndex
}

type SubsonicIndex struct {
	Name    string           `json:"name"`
	Artists []SubsonicArtist `json:"artist"`
}

type SubsonicResponse struct {
	Status    string            `json:"status"`
	Version   string            `json:"version"`
	Indexes   SubsonicIndexes   `json:"indexes"`
	Directory SubsonicDirectory `json:"directory"`
	Error     SubsonicError     `json:"error"`
	Playlists SubsonicPlaylists `json:"playlists"`
	Playlist  SubsonicPlaylist  `json:"playlist"`
}

type responseWrapper struct {
	Response SubsonicResponse `json:"subsonic-response"`
}

type SubsonicPlaylists struct {
	Playlist []SubsonicPlaylist
}

type SubsonicPlaylist struct {
	Id        string           `json:"id"`
	Name      string           `json:"name"`
	Owner     string           `json:"owner"`
	Public    bool             `json:"public"`
	SongCount int              `json:"songCount"`
	Duration  int              `json:"duration"`
	Created   string           `json:"created"`
	Changed   string           `json:"changed"`
	Entities  []SubsonicEntity `json:"entry"`
}
