package services

import (
	"bufio"
	"fmt"
	"github.com/karrick/godirwalk"
	"github.com/relex/aini"
	"os"
	"path/filepath"
	"strings"
)

type Inventory struct {
	FilePath string
	PathTags []string
	Data     *aini.InventoryData
}

// Gather inventory files from a Parent directory
// Using a recursive scan. All non inventory files are ignored ( not .ini file )
// All sub parent directory added like label in the inventory
func readInventories(rootPath string, pathTags ...string) (result []*Inventory, err error) {
	absRoot, err := filepath.Abs(rootPath)

	if err != nil {
		return
	}

	err = godirwalk.Walk(absRoot, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if strings.Contains(osPathname, "vars") || strings.Contains(osPathname, "template") {
				return godirwalk.SkipThis
			}

			if !strings.Contains(filepath.Base(osPathname), ".ini") {
				return nil
			}
			pathMetas := strings.Split(strings.TrimSuffix(strings.TrimPrefix(osPathname, absRoot), fmt.Sprintf("/%s", de.Name())), "/")

			result = append(result, &Inventory{
				FilePath: osPathname,
				PathTags: pathMetas,
			})
			return nil
		},
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			return godirwalk.SkipNode
		},
		Unsorted: true,
	})
	return
}

func (path *Inventory) make() {
	if path == nil {
		return
	}

	if strings.Contains(filepath.Base(path.FilePath), ".ini") {
		if file, err := os.Open(path.FilePath); err == nil {
			reader := bufio.NewReader(file)
			if data, err := aini.Parse(reader); err == nil {
				path.Data = data
			}
		}
	}
}