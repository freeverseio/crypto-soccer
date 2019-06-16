package fifo

type Msg struct {

}

type Fifo struct {
	fifo chan Msg
}

func FifoNew(bufferSize int) *Fifo {
	return &Fifo{
		fifo : make(chan Msg, bufferSize),
	}
}

func (m *Fifo) Push(msg Msg) {
	m.fifo <- msg
}

func (m* Fifo) Pop() Msg {
	return <-m.fifo
}