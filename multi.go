package progressbar

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

var Out = os.Stdout
var RefreshInterval = time.Millisecond * 50

type MultiProgress struct {
	Out  io.Writer
	Bars []*ProgressBar

	RefreshInterval time.Duration
	ticker          *time.Ticker
	tDone           chan bool
	rwMtx           *sync.RWMutex
	lw              *Writer
}

func NewMultiProgress() *MultiProgress {
	mp := &MultiProgress{
		Out:             Out,
		RefreshInterval: RefreshInterval,
		tDone:           make(chan bool),
		rwMtx:           &sync.RWMutex{},
	}
	mp.lw = NewWriter()
	mp.lw.Out = Out
	return mp
}

func (m *MultiProgress) SetRefreshInterval(t time.Duration) {
	m.rwMtx.Lock()
	defer m.rwMtx.Unlock()
	m.RefreshInterval = t
}

func (m *MultiProgress) AddBar(bar *ProgressBar) {
	m.rwMtx.Lock()
	defer m.rwMtx.Unlock()
	m.Bars = append(m.Bars, bar)
}

func (m *MultiProgress) Add64Bar(max int64) *ProgressBar {
	b := New64(max)
	m.AddBar(b)
	return b
}

func (m *MultiProgress) AddDefaultBar(max int64, description ...string) *ProgressBar {
	b := Default(max, description...)
	m.AddBar(b)
	return b
}

func (m *MultiProgress) Listen() {
	for {
		m.rwMtx.RLock()
		interval := m.RefreshInterval
		m.rwMtx.RUnlock()

		select {
		case <-time.After(interval):
			m.print()
		case <-m.tDone:
			close(m.tDone)
			return
		}
	}
}

func (m *MultiProgress) print() {
	//输出bar
	for _, bar := range m.Bars {
		//when bar is done, pop
		renderProgressBar(bar.config, &bar.state)
		_, _ = fmt.Fprintln(m.lw, bar.String())
		if bar.IsFinished() {
			fmt.Println("finished")
		}
	}
	_ = m.lw.Flush()
}

func (m *MultiProgress) Start() {
	go m.Listen()
}

func (m *MultiProgress) Stop() {
	m.tDone <- true
	<-m.tDone
}
