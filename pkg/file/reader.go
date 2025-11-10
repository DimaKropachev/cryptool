package file

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func ReadDecryptedFile(path string, blockSize int) (<-chan []byte, <-chan error, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil, fmt.Errorf("file not found")
		}
		return nil, nil, err
	}

	outCh := make(chan []byte)
	errCh := make(chan error)

	go func() {
		defer f.Close()
		defer close(outCh)
		defer close(errCh)

		buf := make([]byte, blockSize)

		for {
			n, err := f.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				errCh <- err
			}

			outCh <- buf[:n]
		}
	}()

	return outCh, errCh, nil
}

type Content struct {
	Nonce []byte
	Buf   []byte
}

func ReadEncryptedFile(f *os.File, nonceSize, blockSize, tagSize int) (<-chan Content, <-chan error, error) {
	outCh := make(chan Content)
	errCh := make(chan error)

	go func() {
		defer close(outCh)
		defer close(errCh)

		nonce := make([]byte, nonceSize)
		buf := make([]byte, blockSize+tagSize)
		for {
			n, err := f.Read(nonce)
			if err != nil {
				if err == io.EOF {
					break
				}
				errCh <- err
			}

			m, err := f.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				errCh <- err
			}

			outCh <- Content{
				Nonce: nonce[:n],
				Buf:   buf[:m],
			}
		}
	}()

	return outCh, errCh, nil
}
