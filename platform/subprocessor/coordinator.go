package subprocessor

import (
	"fmt"
	"github.com/satori/go.uuid"
	"log"
	"sync"
)

const (
	defaultRoutineCount = 10
)

type Coordinator struct {
	coordinationPool chan chan Resource
	ResourceChannel  chan Resource
	routines         []*Routine
	quit             chan bool
	id               string
}

var instance *Coordinator
var once sync.Once

func NewCoordinator() *Coordinator {

	coordinator := &Coordinator{}

	uuid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	coordinator.id = fmt.Sprintf("%s", uuid)
	log.Printf("Starting coordinator: %s", coordinator.id)
	//Create channels
	//Two way coordination channel
	coordinator.coordinationPool = make(chan chan Resource, defaultRoutineCount)
	//Input channel
	coordinator.ResourceChannel = make(chan Resource)

	coordinator.quit = make(chan bool)

	//Create routines

	for i := 0; i < defaultRoutineCount; i++ {
		rt := NewRoutine(coordinator.coordinationPool)
		coordinator.routines = append(coordinator.routines, rt)
		rt.Start() //Start each routine
	}
	log.Printf("Finished building coordinator: %s", coordinator.id)
	return coordinator
}
func (coordinator *Coordinator) Destroy() {
	go func() {
		log.Printf("Coordinator %s is being destroyed", coordinator.id)
		coordinator.stop()

		for _, routine := range coordinator.routines {
			routine.quit <- true
		}
	}()
}
func (coordinator *Coordinator) stop() {

	coordinator.quit <- true
	log.Printf("Coordinator %s will stop receiving input", coordinator.id)
}
func (coordinator *Coordinator) Run() {

	for {
		select {
		case msg := <-coordinator.ResourceChannel: //External Coordinator input channel
			go func(msg Resource) {
				next := <-coordinator.coordinationPool
				next <- msg
				log.Println("Processed new message!")
			}(msg)
		case _ = <-coordinator.quit:
			return
		}
	}
}
