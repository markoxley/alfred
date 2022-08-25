package state

import "sync"

type State struct {
	allowLook  bool
	allowSpeak bool
	isTracking bool
}

var (
	safetyMutex sync.Mutex
	Ohbot       *State
)

func Init() *State {
	if Ohbot == nil {
		Ohbot = &State{
			allowLook:  true,
			allowSpeak: true,
			isTracking: true,
		}
	}
	return Ohbot
}

func (s *State) Look() bool {
	safetyMutex.Lock()
	defer safetyMutex.Unlock()
	return s.allowLook
}

func (s *State) Speak() bool {
	safetyMutex.Lock()
	defer safetyMutex.Unlock()
	return s.allowLook
}

func (s *State) setAllowLook(v bool) {
	safetyMutex.Lock()
	defer safetyMutex.Unlock()
	s.allowLook = v
}

func (s *State) setTracking(v bool) {
	safetyMutex.Lock()
	defer safetyMutex.Unlock()
	s.isTracking = v
}

func (s *State) setAllowSpeak(v bool) {
	safetyMutex.Lock()
	defer safetyMutex.Unlock()
	s.allowSpeak = v
}

func (s *State) DisableLook() {
	s.setAllowLook(false)
}

func (s *State) EnableLook() {
	s.setAllowLook(true)
}

func (s *State) DisableSpeak() {
	s.setAllowSpeak(false)
}

func (s *State) EnableSpeak() {
	s.setAllowSpeak(true)
}

func (s *State) StartTracking() {
	s.setTracking(true)
}

func (s *State) StopTracking() {
	s.setTracking(false)
}
