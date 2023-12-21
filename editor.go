package texteditor

import (
	"fmt"
	"github.com/coconutLatte/texteditor/revert"
	"os"
	"os/exec"
	"strings"
)

func EditorStatic(inContent []byte) ([]byte, error) {
	var f *os.File
	var err error
	var path string

	// Detect the text editor to use
	editor := os.Getenv("VISUAL")
	if editor == "" {
		editor = os.Getenv("EDITOR")
		if editor == "" {
			for _, p := range []string{"vim", "vi", "editor", "emacs", "nano"} {
				_, err := exec.LookPath(p)
				if err == nil {
					editor = p
					break
				}
			}
			if editor == "" {
				return []byte{}, fmt.Errorf("No text editor found, please set the EDITOR environment variable")
			}
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("get wd failed, %w", err)
	}

	f, err = os.CreateTemp(wd, "coconut_Latte_editor_")
	if err != nil {
		return []byte{}, err
	}

	r := revert.New()
	defer r.Fail()
	r.Add(func() {
		_ = f.Close()
		_ = os.Remove(f.Name())
	})

	err = os.Chmod(f.Name(), 0600)
	if err != nil {
		return []byte{}, err
	}

	_, err = f.Write(inContent)
	if err != nil {
		return []byte{}, err
	}

	err = f.Close()
	if err != nil {
		return []byte{}, err
	}

	path = fmt.Sprintf("%s.tmp", f.Name())
	err = os.Rename(f.Name(), path)
	if err != nil {
		return []byte{}, err
	}

	r.Success()
	r.Add(func() { _ = os.Remove(path) })

	cmdParts := strings.Fields(editor)
	cmd := exec.Command(cmdParts[0], append(cmdParts[1:], path)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return []byte{}, err
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return []byte{}, err
	}

	return content, nil
}
