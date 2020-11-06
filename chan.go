package plog

func Chan(c chan *Msg) Opener {
	return func() (Messenger, error) {
		return chanMes(c), nil
	}
}

type chanMes chan *Msg

func (c chanMes) Recv() (*Msg, error) {
	return <-c, nil
}

func (c chanMes) Send(msg *Msg) error {
	c <- msg
	return nil
}

func (c chanMes) Close() error {
	close(c)
	return nil
}
