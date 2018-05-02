package watch

import (
	"errors"
	//"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/lioneagle/goutil/src/chars"
	"github.com/lioneagle/goutil/src/logger"

	"github.com/rjeczalik/notify"
)

type FileCreatedCallback interface {
	FileCreated(string)
}

type FileChangedCallback interface {
	FileChanged(string)
}

type FileRemovedCallback interface {
	FileRemoved(string)
}

type FileRenamedCallback interface {
	FileRenamed(string)
}

type Watcher struct {
	sync.Mutex
	fsEvent  chan notify.EventInfo
	watched  map[string][]interface{}
	watchers []string
	dirs     []string
}

func NewWatcher() (*Watcher, error) {
	watcher := &Watcher{fsEvent: make(chan notify.EventInfo, 5)}
	watcher.watched = make(map[string][]interface{})
	watcher.watchers = make([]string, 0)
	watcher.dirs = make([]string, 0)
	go watcher.observe()

	return watcher, nil
}

func (this *Watcher) Close() {
	notify.Stop(this.fsEvent)
}

func (this *Watcher) observe() {
	for {
		select {
		case event, ok := <-this.fsEvent:
			if !ok {
				// We get here only when this.fsEvent is stopped when closing the watcher
				this.watched = nil
				this.watchers = nil
				this.dirs = nil
				close(this.fsEvent)
				this.fsEvent = nil
				return
			}
			this.parseEvent(event)
		}
	}
}

func (this *Watcher) parseEvent(event notify.EventInfo) {
	logger.Info("watcher event %s", event)

	this.Lock()
	defer this.Unlock()

	path := event.Path()
	this.apply(path, event.Event())
	// currently fsnotify pushs remove event for files
	// inside directory when a directory is removed but
	// when the directory is renamed there is no event for
	// files inside directory
	if event.Event()&notify.Rename != 0 && chars.Exists(this.dirs, path) {
		for v, _ := range this.watched {
			if filepath.Dir(v) == path {
				this.apply(v, event.Event())
			}
		}
	}
	dir := filepath.Dir(path)
	// The watcher will be removed if the file is deleted
	// so we need to watch the parent directory for when the
	// file is created again
	if event.Event()&notify.Remove != 0 {
		this.watchers = chars.Remove(this.watchers, path)
		this.Unlock()
		this.Watch(dir, nil)
		this.Lock()
	}
	// If the event is create we will apply FileCreated callback
	// for the parent directory to because when new file is created
	// inside directory we won't get any event for the watched directory.
	// we need this feature to detect new packages(plugins, settings, etc)
	if callbacks, exist := this.watched[dir]; event.Event()&notify.Create != 0 && exist {
		for _, v := range callbacks {
			if callback, ok := v.(FileCreatedCallback); ok {
				this.Unlock()
				callback.FileCreated(path)
				this.Lock()
			}
		}
	}
}

func (watcher *Watcher) Watch(name string, callback interface{}) error {
	if !filepath.IsAbs(name) {
		var err error
		name, err = filepath.Abs(name)
		if err != nil {
			return err
		}
	}

	logger.Info("Watch(\"%s\")", name)

	fileInfo, err := os.Stat(name)
	isDir := err == nil && fileInfo.IsDir()

	// If the file doesn't exist currently we will add watcher for file
	// directory and look for create event inside the directory
	if os.IsNotExist(err) {
		logger.Info("\"%s\" doesn't exist, Watching parent directory", name)
		if err := watcher.Watch(filepath.Dir(name), nil); err != nil {
			return err
		}
	}

	watcher.Lock()
	defer watcher.Unlock()

	if err := watcher.addCallback(name, callback); err != nil {
		if !isDir {
			logger.Error("\"%s\" is not directory", name)
			return err
		}
		if chars.Exists(watcher.dirs, name) {
			logger.Info("\"%s\" is watched already", name)
			return nil
		}
	}
	// If exists in watchers we are already watching the path
	// Or
	// If the file is under one of watched dirs
	//
	// no need to create watcher
	if chars.Exists(watcher.watchers, name) || (!isDir && chars.Exists(watcher.dirs, filepath.Dir(name))) {
		return nil
	}

	if err := watcher.watch(name, isDir); err != nil {
		return err
	}

	if isDir {
		watcher.flushDir(name)
	}
	return nil
}

func (this *Watcher) addCallback(name string, callback interface{}) error {
	logger.Info("Adding \"%s\" callback", name)
	numok := 0
	if _, ok := callback.(FileChangedCallback); ok {
		numok++
	}
	if _, ok := callback.(FileCreatedCallback); ok {
		numok++
	}
	if _, ok := callback.(FileRemovedCallback); ok {
		numok++
	}
	if _, ok := callback.(FileRenamedCallback); ok {
		numok++
	}
	if numok == 0 {
		return errors.New("The callback argument does satisfy any File*Callback interfaces")
	}
	this.watched[name] = append(this.watched[name], callback)
	return nil
}

func (this *Watcher) watch(name string, isDir bool) error {
	watchPath := name
	if isDir {
		watchPath = filepath.Join(watchPath, "...")
	}
	logger.Info("Creating watcher on \"%s\"", name)
	if err := notify.Watch(watchPath, this.fsEvent, notify.All); err != nil {
		return err
	}
	this.watchers = append(this.watchers, name)
	return nil
}

// Remove watchers created on files under this directory because
// one watcher on the parent directory is enough for all of them
func (this *Watcher) flushDir(name string) {
	logger.Info("Flusing watched directory %s", name)
	this.dirs = append(this.dirs, name)
	for _, p := range this.watchers {
		if filepath.Dir(p) == name && !chars.Exists(this.dirs, p) {
			if err := this.removeWatch(p); err != nil {
				logger.Error("Couldn't unwatch file %s: %s", p, err)
			}
		}
	}
}

func (this *Watcher) UnWatch(name string, callback interface{}) error {
	if !filepath.IsAbs(name) {
		var err error
		name, err = filepath.Abs(name)
		if err != nil {
			return err
		}
	}

	logger.Info("UnWatch(%s)", name)
	this.Lock()
	defer this.Unlock()

	if callback == nil {
		return this.unWatch(name)
	}
	for i, c := range this.watched[name] {
		if c == callback {
			this.watched[name][i] = this.watched[name][len(this.watched[name])-1]
			this.watched[name][len(this.watched[name])-1] = nil
			this.watched[name] = this.watched[name][:len(this.watched[name])-1]
			break
		}
	}
	if len(this.watched[name]) == 0 {
		this.unWatch(name)
	}
	return nil
}

func (this *Watcher) unWatch(name string) error {
	delete(this.watched, name)
	if err := this.removeWatch(name); err != nil {
		return err
	}
	return nil
}

func (this *Watcher) removeWatch(name string) error {
	logger.Info("removing watcher from \"%s\"", name)
	// TODO
	// notify.Stop(w.fsEvent) would stop ALL watchers, and the only way for
	// unwatching is notify.Stop(). So we won't unwatch from the notify for
	// now, just we will remove it from Watcher struct untill there is a
	// better solution
	this.watchers = chars.Remove(this.watchers, name)
	if chars.Exists(this.dirs, name) {
		this.removeDir(name)
	}
	return nil
}

// Put back watchers on watching files under the directory
func (this *Watcher) removeDir(name string) {
	for p, _ := range this.watched {
		if filepath.Dir(p) == name {
			stat, err := os.Stat(p)
			if err != nil {
				logger.Error("Stat error: %s", err)
			}
			if err := this.watch(p, stat.IsDir()); err != nil {
				logger.Error("Could not watch: %s", err)
				continue
			}
		}
	}
	this.dirs = chars.Remove(this.dirs, name)
}

func (this *Watcher) apply(path string, flags notify.Event) {
	for _, v := range this.watched[path] {
		if flags&notify.Create != 0 {
			if callback, ok := v.(FileCreatedCallback); ok {
				callback.FileCreated(path)
			}
		}
		if flags&notify.Write != 0 {
			if c, ok := v.(FileChangedCallback); ok {
				c.FileChanged(path)
			}
		}
		if flags&notify.Remove != 0 {
			if c, ok := v.(FileRemovedCallback); ok {
				c.FileRemoved(path)
			}
		}
		if flags&notify.Rename != 0 {
			if c, ok := v.(FileRenamedCallback); ok {
				c.FileRenamed(path)
			}
		}
	}
}
