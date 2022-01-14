package bloat

import (
	"os"
	"path/filepath"
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
