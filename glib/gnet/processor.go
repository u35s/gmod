package gnet

type Processor interface {
	Pack(i interface{}) ([]byte, error)
	Marshal(v interface{}) ([]byte, error)
	UnPack(p []byte, q chan interface{}) (n int, err error)
	Unmarshal(data []byte, v interface{}) error
}
