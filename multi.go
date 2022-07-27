package progressbar

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

var Out = os.Stdout
var RefreshInterval = time.Millisecond * 50

type MultiProgress struct {
	Out       io.Writer
	barsNames []string
	mapBars   map[string]*ProgressBar

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
		barsNames:       make([]string, 0),
		mapBars:         make(map[string]*ProgressBar),
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

func (m *MultiProgress) AddBar(name string, bar *ProgressBar) error {
	if bar == nil {
		return errors.New("bar is nil")
	}
	if _, ok := m.mapBars[name]; ok {
		return errors.New("bar name already exists")
	}
	m.rwMtx.Lock()
	defer m.rwMtx.Unlock()
	m.mapBars[name] = bar
	m.barsNames = append(m.barsNames, name)
	return nil
}

func (m *MultiProgress) Add64Bar(name string, max int64) (*ProgressBar, error) {
	if _, ok := m.mapBars[name]; ok {
		return nil, errors.New("bar name already exists")
	}
	b := New64(max)
	return b, m.AddBar(name, b)
}

func (m *MultiProgress) AddDefaultBar(name string, max int64, description ...string) (*ProgressBar, error) {
	if _, ok := m.mapBars[name]; ok {
		return nil, errors.New("bar name already exists")
	}
	b := Default(max, description...)
	return b, m.AddBar(name, b)
}

func (m *MultiProgress) BarAdd(name string, n int) {
	m.BarAdd64(name, int64(n))
}

func (m *MultiProgress) BarAdd64(name string, n int64) {
	m.rwMtx.Lock()
	defer m.rwMtx.Unlock()
	bar, ok := m.mapBars[name]
	if !ok {
		return
	}
	bar.Add64(n)
}

func (m *MultiProgress) BarSet64(name string, n int64) {
	m.rwMtx.Lock()
	defer m.rwMtx.Unlock()
	bar, ok := m.mapBars[name]
	if !ok {
		return
	}
	bar.Set64(n)
}

func (m *MultiProgress) BarSet(name string, n int) {
	m.BarSet64(name, int64(n))
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
	finishCount := 0
	//把已经完成的bar放在list最前面
	for i := 0; i < len(m.barsNames); i++ {
		bar := m.mapBars[m.barsNames[i]]
		bar.render()
		if bar.IsFinished() {
			tmp := make([]string, 0, len(m.barsNames))
			tmp = append(tmp, m.barsNames[i])
			tmp = append(tmp, m.barsNames[:i]...)
			m.barsNames = append(tmp, m.barsNames[i+1:]...)
		}
	}
	//输出bar
	for i := 0; i < len(m.barsNames); i++ {
		bar, ok := m.mapBars[m.barsNames[i]]
		if !ok {
			continue
		}
		//when bar is done, pop
		s := bar.String()
		if s == "" {
			continue
		}
		_, _ = fmt.Fprintln(m.lw, s)
		if bar.IsFinished() {
			m.rwMtx.Lock()
			delete(m.mapBars, m.barsNames[i])
			m.barsNames = append(m.barsNames[:i], m.barsNames[i+1:]...)
			finishCount++
			m.rwMtx.Unlock()
			i--
		}
	}
	_ = m.lw.Flush()
	m.lw.lineCount -= finishCount
}

func (m *MultiProgress) Start() {
	go m.Listen()
}

func (m *MultiProgress) Stop() {
	m.tDone <- true
	<-m.tDone
}
