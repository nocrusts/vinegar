//go:build nogui || nosplash

package splash

import (
	"errors"
	"log"
)

var ErrClosed = errors.New("window closed")

type Splash struct {
	LogPath string
}

func (ui *Splash) SetMessage(msg string) {
}

func (ui *Splash) SetDesc(desc string) {
}

func (ui *Splash) SetProgress(progress float32) {
}

func (ui *Splash) Close() {
}

func (ui *Splash) Invalidate() {
}

func (ui *Splash) IsClosed() bool {
	return true
}

func (ui *Splash) Dialog(txt string, user bool) bool {
	log.Printf("splash: dialog(%d, false): %s", user, txt)
	return false
}

func New(cfg *Config) *Splash {
	return &Splash{}
}

func (ui *Splash) Run() error {
	return nil
}
