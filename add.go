package shell

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	files "github.com/ipfs/go-ipfs-files"
)

type object struct {
	Hash string
	Name string
	Size string
}

type AddOpts = func(*RequestBuilder) error

func OnlyHash(enabled bool) AddOpts {
	return func(rb *RequestBuilder) error {
		rb.Option("only-hash", enabled)
		return nil
	}
}

func Pin(enabled bool) AddOpts {
	return func(rb *RequestBuilder) error {
		rb.Option("pin", enabled)
		return nil
	}
}

func Progress(enabled bool) AddOpts {
	return func(rb *RequestBuilder) error {
		rb.Option("progress", enabled)
		return nil
	}
}

func RawLeaves(enabled bool) AddOpts {
	return func(rb *RequestBuilder) error {
		rb.Option("raw-leaves", enabled)
		return nil
	}
}

// Hash allows for selecting the multihash type
func Hash(hash string) AddOpts {
	return func(rb *RequestBuilder) error {
		rb.Option("hash", hash)
		return nil
	}
}

// CidVersion allows for selecting the CID version that ipfs should use.
func CidVersion(version int) AddOpts {
	return func(rb *RequestBuilder) error {
		rb.Option("cid-version", version)
		return nil
	}
}

// Add adds a file to ipfs pinning it with the given options
func (s *Shell) Add(r io.Reader, options ...AddOpts) (string, error) {
	fr := files.NewReaderFile(r)
	slf := files.NewSliceDirectory([]files.DirEntry{files.FileEntry("", fr)})
	fileReader := files.NewMultiFileReader(slf, true)

	var out object
	rb := s.Request("add")
	for _, option := range options {
		option(rb)
	}
	return out.Hash, rb.Body(fileReader).Exec(context.Background(), &out)
}

// AddNoPin adds a file to ipfs without pinning it
// Deprecated: Use Add() with option functions instead
func (s *Shell) AddNoPin(r io.Reader) (string, error) {
	return s.Add(r, Pin(false))
}

// AddWithOpts adds a file to ipfs with some additional options
// Deprecated: Use Add() with option functions instead
func (s *Shell) AddWithOpts(r io.Reader, pin bool, rawLeaves bool) (string, error) {
	return s.Add(r, Pin(pin), RawLeaves(rawLeaves))
}

func (s *Shell) AddLink(target string) (string, error) {
	link := files.NewLinkFile(target, nil)
	slf := files.NewSliceDirectory([]files.DirEntry{files.FileEntry("", link)})
	reader := files.NewMultiFileReader(slf, true)

	var out object
	return out.Hash, s.Request("add").Body(reader).Exec(context.Background(), &out)
}

// AddDir adds a directory recursively with all of the files under it
func (s *Shell) AddDir(dir string) (object, error) {
	stat, err := os.Lstat(dir)
	if err != nil {
		return object{}, err
	}

	sf, err := files.NewSerialFile(dir, false, stat)
	if err != nil {
		return object{}, err
	}
	slf := files.NewSliceDirectory([]files.DirEntry{files.FileEntry(filepath.Base(dir), sf)})
	reader := files.NewMultiFileReader(slf, true)

	resp, err := s.Request("add").
		Option("recursive", true).
		Body(reader).
		Send(context.Background())
	//fmt.Println(resp.Body)
	if err != nil {
		return object{}, err
	}

	defer resp.Close()

	if resp.Error != nil {
		return object{}, resp.Error
	}

	dec := json.NewDecoder(resp.Output)

	var final object
	for {
		var out object
		err = dec.Decode(&out)
		if err != nil {
			if err == io.EOF {
				break
			}
			return object{}, err
		}
		final = out
		//fmt.Println(out.Name, out.Size)
	}
	return final, nil
}

// AddDir adds a directory recursively with all of the files under it
func (s *Shell) SwanAddDir(dir string) (*http.Response, error) {
	stat, err := os.Lstat(dir)
	if err != nil {
		return nil, err
	}

	sf, err := files.NewSerialFile(dir, false, stat)
	if err != nil {
		return nil, err
	}
	slf := files.NewSliceDirectory([]files.DirEntry{files.FileEntry(filepath.Base(dir), sf)})
	reader := files.NewMultiFileReader(slf, true)

	resp, err := s.Request("add").
		Option("recursive", true).
		Body(reader).
		SwanSend(context.Background())
	fmt.Println(resp.Body)
	return resp, nil
}
