package event

import "github.com/veandco/go-sdl2/sdl"

//TODO proper docs
//TODO do something with mouse coordinates, dunno

//MouseStateObserver is used to ensure that upon clicking the mouse a buttons OnClick function is only called once
//instead of continuously
type MouseStateObserver struct {
	mx, my     int32
	ms         uint32
	previousMs uint32
}

//NewMouseObserver function
func NewMouseObserver() *MouseStateObserver {
	return &MouseStateObserver{
		0, 0, 0, 0,
	}
}

//MousePosition function
func (observer *MouseStateObserver) MousePosition() (int32, int32) {
	return observer.mx, observer.my
}

//Update function
func (observer *MouseStateObserver) Update() (uint32, bool) {
	observer.previousMs = observer.ms
	observer.mx, observer.my, observer.ms = sdl.GetMouseState()

	if observer.ms == 0 && observer.previousMs != 0 {
		return observer.previousMs, true
	}

	return 0, false
}
