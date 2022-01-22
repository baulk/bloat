package bloat

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func findSingleChildPath(src string) (string, error) {
	ds, err := os.ReadDir(src)
	if err != nil {
		return "", err
	}
	if len(ds) != 1 || !ds[0].IsDir() || strings.EqualFold(ds[0].Name(), "bin") {
		return "", nil
	}
	return filepath.Join(src, ds[0].Name()), nil
}

func MakeFlattened(src, dest string) error {
	subfirst, err := findSingleChildPath(src)
	if err != nil {
		return err
	}
	if len(subfirst) != 0 {
		return nil
	}
	current := subfirst
	for i := 0; i < 10; i++ {
		next, err := findSingleChildPath(current)
		if err != nil {
			return err
		}
		if len(next) == 0 {
			break
		}
		current = next
	}
	items, err := os.ReadDir(current)
	if err != nil {
		return err
	}
	for _, item := range items {
		_ = os.Rename(
			filepath.Join(current, item.Name()),
			filepath.Join(dest, item.Name()))
	}
	return os.RemoveAll(subfirst)
}

var (
	ErrDangerousPathAccessDenied = errors.New("dangerous path access denied")
)

func JoinSanitizePath(parent string, elem ...string) (string, error) {
	var buf strings.Builder
	_, _ = buf.WriteString(parent)
	for _, e := range elem {
		_ = buf.WriteByte(os.PathSeparator)
		_, _ = buf.WriteString(e)
	}
	out := filepath.Clean(buf.String())
	if len(out) <= len(parent) {
		return "", ErrDangerousPathAccessDenied
	}
	if strings.HasPrefix(out, parent) && os.IsPathSeparator(out[len(parent)]) {
		return out, nil
	}
	if runtime.GOOS != "windows" && parent == "/" {
		return out, nil
	}
	return "", ErrDangerousPathAccessDenied
}

func JoinSanitizePathSlow(parent string, elem ...string) (string, error) {
	parent = filepath.Clean(parent)
	return JoinSanitizePath(parent, elem...)
}

func Symlink(oldname string, newname string) error {
	if err := os.MkdirAll(filepath.Dir(newname), 0755); err != nil {
		return fmt.Errorf("%s: making directory for file: %v", newname, err)
	}

	if _, err := os.Lstat(newname); err == nil {
		if err = os.Remove(newname); err != nil {
			return fmt.Errorf("%s: failed to unlink: %+v", newname, err)
		}
	}

	if err := os.Symlink(oldname, newname); err != nil {
		return fmt.Errorf("%s: making symbolic link for: %v", newname, err)
	}
	return nil
}
