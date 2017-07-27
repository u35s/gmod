package gnet

type Processor interface {
	Marshal(i interface{}) ([]byte, error)
	UnMarshal(p []byte, q chan interface{}) (n int, err error)
}
