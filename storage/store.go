package storage

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
)

type TransformKeyFunc func(string) string

type Store struct {
	Root        string
	transformFn TransformKeyFunc
}

func NewStore(root string, transformFn TransformKeyFunc) *Store {
	return &Store{
		Root:        root,
		transformFn: transformFn,
	}
}

func (s *Store) Store(key string, r io.Reader) error {
	// store key to local folder
	keyPath := s.transformFn(key)
	keyPathWithRoot := path.Join(s.Root, filepath.Dir(keyPath))
	err := os.MkdirAll(keyPathWithRoot, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create key path: %v", err)
	}

	keyFileWithRoot := path.Join(s.Root, keyPath)
	f, err := os.Create(keyFileWithRoot)
	if err != nil {
		return fmt.Errorf("failed to create key file: %v", err)
	}
	defer f.Close()

	// store data to the file
	if r == nil {
		log.Printf("r is nil")
		return nil
	}

	_, err = io.Copy(f, r)
	if err != nil {
		return fmt.Errorf("failed to copy data to file: %v", err)
	}
	return nil
}

func (s *Store) HasKey(key string) bool {
	keyPath := s.transformFn(key)
	keyFileWithRoot := path.Join(s.Root, keyPath)
	_, err := os.Stat(keyFileWithRoot)
	return err == nil
}

func (s *Store) Get(key string) (io.ReadCloser, error) {
	keyPath := s.transformFn(key)
	keyFileWithRoot := path.Join(s.Root, keyPath)
	f, err := os.Open(keyFileWithRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}

	return f, nil
}
