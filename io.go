package dockerlib

import (
	"bytes"
	"io"
)

func ReadLinesAsBytes(reader io.Reader) <-chan []byte {
	lines := make(chan []byte)
	buffer := make([]byte, 1024)
	leftover := make([]byte, 0)
	sep := []byte("\n")
	go func() {
		for {
			n, err := reader.Read(buffer)

			if n > 0 {
				parts := bytes.Split(buffer[:n], sep)

				// send first part plus any continuation
				logger.Debugf("Before: %s %s", leftover, parts[0])
				leftover = append(leftover, parts[0]...)
				logger.Debugf("After : %s", leftover)

				leftover = bytes.Trim(leftover, "\r")
				lines <- leftover

				// send middle parts
				for i := 1; i < len(parts)-1; i++ {
					lines <- bytes.Trim(parts[i], "\r")
				}

				// save continuation
				leftover = bytes.Trim(parts[len(parts)-1], "\r")
			}

			switch {
			case err == io.EOF:
				if len(leftover) > 0 {
					lines <- leftover
				}
				close(lines)
				return
			case err != nil:
				logger.Errorf("Unexpected error from Reader: %v", err)
				lines <- leftover
				close(lines)
				break
			}
		}
	}()

	return lines
}
