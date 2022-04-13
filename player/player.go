package player

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/yourok/go-mpv/mpv"
)

const (
	PlayerStopped = 0
	PlayerPlaying = 1
	PlayerPaused  = 2
)

type QueueItem struct {
	Uri      string
	Title    string
	Artist   string
	Duration int
}

type Player struct {
	Instance          *mpv.Mpv
	Queue             []QueueItem
	RandomQueue       []int
	ReplaceInProgress bool
	IndexQueue        int
	Repeat            bool
	Random            bool
}

func eventListener(m *mpv.Mpv, p *Player) {
	go func() {
		for {
			e := m.WaitEvent(1)
			if e == nil {
				return
			}
			if e.Event_Id == mpv.EVENT_END_FILE && e.Data.(mpv.EventEndFile).Reason == mpv.END_FILE_REASON_EOF {
				p.PlayNextTrack()
			}
		}
	}()
}

func InitPlayer() (*Player, error) {
	mpvInstance := mpv.Create()
	rand.Seed(time.Now().UnixNano())

	// TODO figure out what other mpv options we need
	mpvInstance.SetOptionString("audio-display", "no")
	mpvInstance.SetOptionString("video", "no")

	err := mpvInstance.Initialize()
	if err != nil {
		mpvInstance.TerminateDestroy()
		return nil, err
	}
	player := Player{mpvInstance, nil, nil, false, 0, false, false}
	eventListener(mpvInstance, &player)
	return &player, nil
}

func (p *Player) Play() {
	song, err := p.ActualSong()
	if err != nil {
		fmt.Println("Player play error: ", err.Error())
		return
	}
	p.Instance.Command([]string{"loadfile", song.Uri})
}

func (p *Player) PlayNextTrack() {
	p.next()
	p.Play()
}

func (p *Player) PlayPrevTrack() {
	p.prev()
	p.Play()
}

func (p *Player) genRandomQueue() {
	p.RandomQueue = nil
	for i := 0; i < len(p.Queue); i++ {
		random := rand.Intn(len(p.Queue))
		for contains(p.RandomQueue, random) {
			random = rand.Intn(len(p.Queue))
		}

		p.RandomQueue = append(p.RandomQueue, random)
	}
}

func contains(array []int, value int) bool {
	for _, i := range array {
		if i == value {
			return true
		}
	}
	return false
}

func (p *Player) next() {
	if p.IndexQueue+1 == len(p.Queue) && p.Repeat {
		p.IndexQueue = 0
	} else {
		p.IndexQueue += 1
	}
}

func (p *Player) prev() {
	if p.IndexQueue > 0 {
		p.IndexQueue -= 1
	}
}

func (p *Player) Stop() {
	p.Instance.Command([]string{"stop"})
}

func (p *Player) IsSongLoaded() bool {
	idle, _ := p.Instance.GetProperty("idle-active", mpv.FORMAT_FLAG)
	return !idle.(bool)
}

func (p *Player) IsPaused() bool {
	pause, _ := p.Instance.GetProperty("pause", mpv.FORMAT_FLAG)
	return pause.(bool)
}

func (p *Player) Pause() int {
	loaded := p.IsSongLoaded()
	pause := p.IsPaused()

	if loaded {
		if pause {
			p.Instance.SetProperty("pause", mpv.FORMAT_FLAG, false)
			return PlayerPlaying
		} else {
			p.Instance.SetProperty("pause", mpv.FORMAT_FLAG, true)
			return PlayerPaused
		}
	} else {
		if len(p.Queue) != 0 {
			p.Play()
			return PlayerPlaying
		} else {
			return PlayerStopped
		}
	}
}

func (p *Player) AdjustVolume(increment int64) int64 {
	volume, _ := p.Instance.GetProperty("volume", mpv.FORMAT_INT64)

	if volume == nil {
		return 0
	}

	newVolume := volume.(int64) + increment

	if newVolume > 100 {
		newVolume = 100
	} else if newVolume < 0 {
		newVolume = 0
	}

	p.Instance.SetProperty("volume", mpv.FORMAT_INT64, newVolume)

	return newVolume
}

func (p *Player) AddInQueue(uri string, title string, artist string, duration int) {
	p.Queue = append(p.Queue, QueueItem{uri, title, artist, duration})
	p.genRandomQueue()
}

func (p *Player) CleanQueue() {
	p.Queue = nil
}

func (p *Player) ActualSong() (*QueueItem, error) {
	if p.IndexQueue < len(p.Queue) {
		if p.Random {
			return &p.Queue[p.RandomQueue[p.IndexQueue]], nil
		} else {
			return &p.Queue[p.IndexQueue], nil
		}
	}
	return nil, errors.New("no actual song")
}

func (p *Player) Actual() string {
	var str string
	if len(p.Queue) > 0 && p.IndexQueue < len(p.Queue) {
		song, err := p.ActualSong()
		if err != nil {
			return err.Error()
		}
		str += song.Artist + " - " + song.Title
	} else {
		str = "no queue"
	}
	return str
}

func (p *Player) SetRepeat(b bool) {
	p.Repeat = b
}

func (p *Player) SetRandom(b bool) {
	p.Random = b
}
