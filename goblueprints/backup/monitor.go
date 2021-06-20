package backup

import (
	"fmt"
	"path/filepath"
	"time"
)

type Monitor struct {
	Paths    map[string]string
	Archiver Archiver
	Dest     string
}

func (m Monitor) Now() (int, error) {
	counter := 0
	for path, lastHash := range m.Paths {
		newHash, err := DirHash(path)
		if err != nil {
			return counter, err
		}
		if newHash != lastHash {
			err := m.act(path)
			if err != nil {
				return counter, err
			}
			counter++
			m.Paths[path] = newHash
		}
	}
	return counter, nil
}

func (m Monitor) act(path string) error {
	dirname := filepath.Base(path)
	filename := fmt.Sprintf(m.Archiver.DestFmt(), time.Now().UnixNano())
	return m.Archiver.Archive(path, filepath.Join(m.Dest, dirname, filename))
}
