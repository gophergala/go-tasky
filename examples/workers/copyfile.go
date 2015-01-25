package workers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/gophergala/go-tasky/tasky"
)

type CopyFile struct {
}

type cp struct {
	Src string
	Dst string
}

type cpUsage struct {
	Usage cp
}

type cpError struct {
	Error string
}

func (d *CopyFile) Name() string {
	return "CopyFile"
}

func (d *CopyFile) Usage() string {
	c := cp{"<source file>", "destination file>"}
	u := cpUsage{c}

	jsonStr, err := json.Marshal(u)
	if err != nil {
		e := cpError{err.Error()}
		estr, _ := json.Marshal(e)
		return string(estr)
	}

	return string(jsonStr)
}

// cpFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func cpFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = cpFileContents(src, dst)
	return
}

// cpFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func cpFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func (d *CopyFile) Perform(job []byte, dataCh chan []byte, errCh chan error, quitCh chan bool) {
	c := cp{}

	err := json.Unmarshal(job, &c)
	if err != nil {
		errCh <- err
		return
	}

	done := make(chan bool)
	go func() {
		err := cpFile(c.Src, c.Dst)
		if err != nil {
			errCh <- err
			done <- true
			return
		}

		dataCh <- []byte("Filed copied successfully.")
		done <- true
	}()

	select {
	case <-done:
		return

	case <-quitCh:
		return
	}
}

func (d *CopyFile) Status() string {
	return tasky.Enabled
}

func (d *CopyFile) Signal(act tasky.Action) bool {
	return true
}

func (d *CopyFile) MaxNumTasks() uint64 {
	return 10
}
