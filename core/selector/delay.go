package selector

import (
	"github.com/sipt/gocks/core"
	"sync/atomic"
	"time"
	"net/http"
	"net"
	"context"
)

func init() {
	core.RegisterSelector(DelaySelector, func(group *core.ServerGroup) (core.ISelector, error) {
		selector := &delaySelector{
			group:     group,
			timer:     time.NewTimer(10 * time.Minute),
			cancel:    make(chan bool, 1),
			resetChan: make(chan bool, 1),
			status:    1,
		}
		go func() {
			var (
				sg         *core.ServerGroup
				s          *core.Server
				min        *core.Server
				minTime, t time.Duration
				ok         bool
			)
			for {
				select {
				case <-selector.timer.C:
				case <-selector.resetChan:
				case <-selector.cancel:
					return
				}
				for _, v := range selector.group.Servers {
					if sg, ok = v.(*core.ServerGroup); ok {
						s, _ = sg.Selector.Get("")
					} else {
						s = v.(*core.Server)
					}
					t = urlTest(s)
					if minTime > t {
						min = s
						minTime = t
					}
				}
				selector.selected = min
				minTime = time.Hour
			}
		}()
		return selector, nil
	})
}

const DelaySelector = "delay"

type delaySelector struct {
	selected  interface{}
	group     *core.ServerGroup
	status    uint32
	timer     *time.Timer
	cancel    chan bool
	resetChan chan bool
}

func (d *delaySelector) Get(_ string) (*core.Server, error) {
	if d.selected == nil {
		d.Start()
		d.selected = d.group.Servers[0]
		if sg, ok := d.selected.(*core.ServerGroup); ok {
			return sg.Selector.Get("")
		}
		s := d.selected.(*core.Server)
		return s, nil
	}
	if sg, ok := d.selected.(*core.ServerGroup); ok {
		return sg.Selector.Get("")
	}
	s := d.selected.(*core.Server)
	return s, nil
}

func (d *delaySelector) Refresh() {
	d.timer.Reset(10 * time.Minute)
	d.resetChan <- true
}

func (d *delaySelector) Reset(group *core.ServerGroup) {
}

func (d *delaySelector) Start() {
	ok := atomic.CompareAndSwapUint32(&d.status, 0, 1)
	if ok {
		d.Refresh()
	}
}

func (d *delaySelector) Stop() {
	ok := atomic.CompareAndSwapUint32(&d.status, 1, 0)
	if ok {
		d.timer.Stop()
	}
}

const url = "http://www.gstatic.com/generate_204"

func urlTest(s *core.Server) time.Duration {
	start := time.Now()
	tr := &http.Transport{
		DialContext: func(_ context.Context, _, addr string) (net.Conn, error) {
			return s.Conn(addr)
		},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil || resp.Status != "204" {
		return time.Hour
	}
	return time.Now().Sub(start)
}
