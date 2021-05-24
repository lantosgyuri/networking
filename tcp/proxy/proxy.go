package proxy

import "io"

func proxy(from io.Reader, to io.Writer) error {
	fromWriter, okWriter := from.(io.Writer)
	toReader, okReader := to.(io.Reader)

	if okWriter && okReader {
		go func() {
			_, _ = io.Copy(fromWriter, toReader)
		}()
	}

	_, err := io.Copy(to, from)

	return err
}
