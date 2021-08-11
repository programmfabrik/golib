package lib

type IOWriter func([]byte) (int, error)

func (iow IOWriter) Write(bts []byte) (n int, err error) {
	return iow(bts)
}
