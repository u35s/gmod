package gnet

type Processor interface {
	Unmarshal(bts []byte) (interface{}, error)
	Marshal(interface{}) ([]byte, error)
}
