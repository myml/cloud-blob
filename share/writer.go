package share

import "io"

type Writer struct {
	err error
	io.WriteCloser
	ErrChan chan error
}

func (w *Writer) Write(data []byte) (int, error) {
	select {
	case err := <-w.ErrChan:
		w.err = err
		return 0, w.err
	default:
		if w.err != nil {
			return 0, w.err
		}
		return w.WriteCloser.Write(data)
	}
}
func (w *Writer) Close() error {
	err := w.WriteCloser.Close()
	if err != nil {
		return err
	}
	err, ok := <-w.ErrChan
	if ok {
		w.err = err
	}
	return w.err
}
