package err_code

import (
	"errors"
)

var (
	WatchStoppedError = errors.New("watcher stopped")
	NotFoundError = errors.New("not found")
)

