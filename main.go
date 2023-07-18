package main

import (
	"compress/gzip"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type config struct {
	list   bool
	output io.Writer
	dest   string
}

var (
	ErrNoRootProvided = errors.New("please provide a root")
	ErrGetArchiveInfo = errors.New("can't get the archive info")
	ErrInvalidPath    = errors.New("provided an invalid path")
)

func main() {
	root := flag.String("root", "", "Root directory from where to unarchive")
	destination := flag.String("dest", "", "Destination directory where to store unarchived files")
	list := flag.Bool("list", false, "List the archives info")
	logFile := flag.String("log", "", "The outfile for logging")

	flag.Parse()
	var out io.Writer = os.Stdout
	if *logFile != "" {
		f, err := os.OpenFile(*logFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		out = f
		defer f.Close()
	}

	cfg := config{
		list:   *list,
		output: out,
		dest:   *destination,
	}
	if *root != "" {
		if err := run(*root, cfg); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func run(root string, cfg config) error {
	return filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if cfg.list {
			if err := getArchiveInfo(path, cfg.output); err != nil {
				return fmt.Errorf("%w, %s", ErrGetArchiveInfo, err)
			}
		}

		// Filter directories because it's going to fail if you want to read a directory as a file
		if info.IsDir() {
			return nil
		}

		return unarchive(root, path, cfg.dest, info)
	})
}

func unarchive(root, path, dest string, file os.FileInfo) error {
	inf, err := os.Stat(dest)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrGetArchiveInfo, err)
	}
	// TODO create the directory if it doesn't exist
	if !inf.IsDir() {
		return fmt.Errorf("%w: %s", ErrInvalidPath, "the destination is not a directory")
	}
	input, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrGetArchiveInfo, err)
	}
	gr, err := gzip.NewReader(input)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrGetArchiveInfo, err)
	}
	// this will get the directory where the file exists relative to the root
	relDir, err := filepath.Rel(root, filepath.Dir(path))
	if err != nil {
		return err
	}

	destBase := gr.Name
	// this is creating the path to where the archive will be saved
	targetPath := filepath.Join(dest, relDir, destBase)
	// this will create all the directories at once; if they exist it will
	// do nothing
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	out, err := os.OpenFile(targetPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	if _, err := io.Copy(out, gr); err != nil {
		return fmt.Errorf("%w: %s", ErrGetArchiveInfo, err)
	}

	if err := input.Close(); err != nil {
		return fmt.Errorf("%w: %s", ErrGetArchiveInfo, err)
	}

	if err := gr.Close(); err != nil {
		return fmt.Errorf("%w: %s", ErrGetArchiveInfo, err)
	}

	return out.Close()
}

func getArchiveInfo(path string, w io.Writer) error {
	_, err := fmt.Fprintln(w, path)
	return err
}
