package file

import (
	"errors"
	"fmt"
	"github.com/chendeke/config/config/source"
	"github.com/chendeke/config/config/utils"
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
)

type watcher struct {
	f *file

	fw   *fsnotify.Watcher
	exit chan bool
}

func newWatcher(f *file) (source.Watcher, error) {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	fw.Add(f.path)

	return &watcher{
		f:    f,
		fw:   fw,
		exit: make(chan bool),
	}, nil
}

func (w *watcher) Next() (*source.ChangeSet, error) {
	// is it closed?
	select {
	case <-w.exit:
		return nil, errors.New("watcher stopped")
	default:
	}
	// configFile := filepath.Clean(w.f.path)
	configDir, _ := filepath.Split(filepath.Clean(w.f.path))
	configFile, err := path_utils.FindConfigFile("config", configDir)
	if err != nil {
		return nil, err
	}
	fmt.Println("configFile:", configFile)
	realConfigFile, _ := filepath.EvalSymlinks(w.f.path)
	fmt.Println("realConfigFile:", realConfigFile)
	// try get the event
	select {
	case event, _ := <-w.fw.Events:

		fmt.Println("recv fw events:", event.String())

		if event.Op == fsnotify.Rename {
			// check existence of file, and add watch again
			_, err := os.Stat(event.Name)
			if err == nil || os.IsExist(err) {
				w.fw.Add(event.Name)
			}
		}

		if event.Op == fsnotify.Remove {
			// Since the symlink was removed, we must
			// re-register the file to be watched
			w.fw.Remove(event.Name)
			w.fw.Add(event.Name)
		}

		// viper 方案
		// currentConfigFile, _ := filepath.EvalSymlinks(configFile)
		// fmt.Println("currentConfigFile:", currentConfigFile)
		// fmt.Println("event.Name:", event.Name)
		// const writeOrCreateMask = fsnotify.Write | fsnotify.Create
		// if (filepath.Clean(event.Name) == configFile && event.Op&writeOrCreateMask != 0) ||
		// 	(currentConfigFile != "" && currentConfigFile != realConfigFile) {
		// 	realConfigFile = currentConfigFile
		// 	c, err := w.f.Read()
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	return c, nil
		// }

		// 原生方案
		c, err := w.f.Read()
		if err != nil {
			return nil, err
		}
		return c, nil

		// return nil, nil
	case err := <-w.fw.Errors:
		return nil, err
	case <-w.exit:
		return nil, errors.New("watcher stopped")
	}
}

func (w *watcher) Stop() error {
	return w.fw.Close()
}
