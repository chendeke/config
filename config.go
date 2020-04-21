package config

import (
	"fmt"
	conf "github.com/chendeke/config/config"
	"github.com/chendeke/config/config/err_code"
	"github.com/chendeke/config/config/reader"
	"github.com/chendeke/config/config/source"
	"github.com/chendeke/config/config/source/env"
	"github.com/chendeke/config/config/source/file"
	"github.com/chendeke/config/config/utils"
	"io"
	"os"
	"strings"
)

type Config = conf.Config
type Watcher = conf.Watcher
type Options = conf.Options
type Option = conf.Option
type Value = reader.Value

func init() {
	runMode := os.Getenv("runmode")
	workPath, _ := os.Getwd()

	sources := make([]source.Source, 0)
	base := workPath + "/conf/config.yaml"
	dev := workPath + "/conf/dev.yaml"
	test := workPath + "/conf/test.yaml"
	prod := workPath + "/conf/prod.yaml"

	sources = append(sources, file.NewSource(file.WithPath(base)))
	if runMode == "dev" {
		sources = append(sources, file.NewSource(file.WithPath(dev)))
	} else if runMode == "test" {
		sources = append(sources, file.NewSource(file.WithPath(test)))
	} else if runMode == "prod" {
		sources = append(sources, file.NewSource(file.WithPath(prod)))
	} else if len(runMode) == 0 {
		if path_utils.IsExist(dev) {
			sources = append(sources, file.NewSource(file.WithPath(dev)))
		}

		if path_utils.IsExist(test) {
			sources = append(sources, file.NewSource(file.WithPath(test)))
		}

		if path_utils.IsExist(prod) {
			sources = append(sources, file.NewSource(file.WithPath(prod)))
		}
	}

	sources = append(sources, env.NewSource())
	// flag.NewSource(),
	err := Load(sources...)
	if err != nil {
		// panic(err)
	}
}

// NewConfig returns new config
func NewConfig(opts ...Option) Config {
	return conf.NewConfig(opts...)
}

// Return config as raw json
func Bytes() []byte {
	return conf.Bytes()
}

// Return config as a map
func Map() map[string]interface{} {
	return conf.Map()
}

// Scan values to a go type
func Scan(v interface{}) error {
	return conf.Scan(v)
}

// Scan values to a go type
func ScanKey(key string, v interface{}) error {
	return conf.Get(key).Scan(v)
}

// Force a source changeset sync
func Sync() error {
	return conf.Sync()
}

// Get a value from the config
func Get(path ...string) Value {
	if len(path) == 1 {
		segments := strings.Split(path[0], ".")
		return conf.Get(segments...)
	}

	return conf.Get(path...)
}

// Load config sources
func Load(source ...source.Source) error {
	return conf.Load(source...)
}

// Watch a value for changes
func Watch(path ...string) (Watcher, error) {
	return conf.Watch(path...)
}

type watchCloser struct {
	exit chan struct{}
}

func (w watchCloser) Close() error {
	fmt.Println("close")
	w.exit <- struct{}{}
	return nil
}

func WatchFunc(handle func(reader.Value), paths ...string) (io.Closer, error) {
	path := make([]string, 0, len(paths))
	for _, v := range paths {
		path = append(path, strings.Split(v, ".")...)
	}

	exit := make(chan struct{})
	w, err := Watch(path...)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			v, err := w.Next()
			if err == err_code.WatchStoppedError {
				return
			}
			if err != nil {
				continue
			}

			if v.Empty() {
				continue
			}

			if handle != nil {
				handle(v)
			}
		}
	}()

	go func() {
		select {
		case <-exit:
			_ = w.Stop()
		}
	}()

	return watchCloser{exit: exit}, nil
}
