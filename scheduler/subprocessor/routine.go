package subprocessor

import (
	"fmt"
	"github.com/satori/go.uuid"
	logger "github.com/sirupsen/logrus"
)

type Routine struct {
	coordinationPool chan chan Resource
	ResourceChannel  chan Resource
	quit             chan bool
	id               string
}

func NewRoutine(cp chan chan Resource) *Routine {
	return &Routine{
		coordinationPool: cp,
		ResourceChannel:  make(chan Resource),
		quit:             make(chan bool),
	}
}

func (r *Routine) Start() {

	uuid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	r.id = fmt.Sprintf("%s", uuid)
	logger.Infof("Starting routine %s", r.id)

	go func() {
		defer func() {
			logger.Infof("Shutting down routine %s", r.id)
		}()
		for {
			// Add my channel into the pool
			r.coordinationPool <- r.ResourceChannel
			select {
			case msg := <-r.ResourceChannel: //Poll my channel for bound queue msg

				msg.Process()
			case <-r.quit:
				return
			}
		}
	}()

}

func (r *Routine) Stop() {
	go func() {
		r.quit <- true
	}()
}
