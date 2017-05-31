package collectors

import (
	"bytes"
	"sync"
	"time"

	"errors"

	log "github.com/Sirupsen/logrus"
	processTool "github.com/shirou/gopsutil/process"
)

// Process collects mem and cpu metric for each N process
type Process struct {
	processes []*processTool.Process

	mutex     sync.RWMutex
	sensision bytes.Buffer
	level     uint8
}

func newProcess(period uint, level uint8) (p *Process) {

	p = &Process{
		level: level,
	}

	if level == 0 {
		return
	}

	tick := time.Tick(time.Duration(period) * time.Millisecond)
	go func() {
		for range tick {
			if err := p.scrape(); err != nil {
				log.Error(err)
			}
		}
	}()

	return
}

// Metrics delivers metrics.
func (p *Process) Metrics() (res *bytes.Buffer) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	res.Write(p.sensision.Bytes())
	return
}

func (p *Process) scrape() (err error) {

	pids, err := processTool.Pids()
	if err != nil {
		return
	}

	p.processes = p.processes[0:0]

	for _, pid := range pids {
		prcs, err := processTool.NewProcess(pid)
		if err != nil {
			return errors.New("Failed to get process infos :" + err.Error())
		}
		p.processes = append(p.processes, prcs)
	}

	return nil
}
