package api

import "errors"

func (connection *SubsonicConnection) GetAlbumList() ([]SubsonicEntity, error) {
	var albumList []SubsonicEntity
	entity, err := connection.GetIndexes()
	if err != nil {
		return nil, err
	}
	for _, index := range entity.Indexes.Index {
		for _, artist := range index.Artists {
			dir, err := connection.GetMusicDirectory(artist.Id)
			if err != nil {
				return nil, err
			}
			for _, album := range dir.Directory.Entities {
				albumList = append(albumList, album)
			}
		}
	}
	return albumList, nil
}

func (connection *SubsonicConnection) GetAlbumIdByName(name string) (*string, error) {
	albumList, err := connection.GetAlbumList()
	if err != nil {
		return nil, err
	}
	for _, album := range albumList {
		if album.Title == name {
			return &album.Id, nil
		}
	}
	return nil, errors.New("No album found")
}

func (connection *SubsonicConnection) GetMusicsFromAlbumId(id string) ([]SubsonicEntity, error) {
	dir, err := connection.GetMusicDirectory(id)
	if err != nil {
		return nil, err
	}
	return dir.Directory.Entities, nil
}
