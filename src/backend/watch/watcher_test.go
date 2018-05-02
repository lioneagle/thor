package watch

import (
	//"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/lioneagle/goutil/src/logger"
	"github.com/lioneagle/goutil/src/test"
)

func newWatcher(t *testing.T) *Watcher {
	watcher, err := NewWatcher()
	test.ASSERT_EQ(t, err, nil, "")
	return watcher
}

func watch(t *testing.T, watcher *Watcher, name string, callback interface{}) {
	err := watcher.Watch(name, callback)
	test.ASSERT_EQ(t, err, nil, "Couldn' Watch %s : %s", name, err)
}

func unwatch(t *testing.T, watcher *Watcher, name string, callback interface{}) {
	err := watcher.UnWatch(name, callback)
	test.ASSERT_EQ(t, err, nil, "Couldn' UnWatch %s : %s", name, err)
}

func testWatched(t *testing.T, watched map[string][]interface{}, expWatched []string) {
	test.EXPECT_EQ(t, len(watched), len(expWatched), "Expected watched %v keys equal to %v", watched, expWatched)

	for _, p := range expWatched {
		absp, err := filepath.Abs(p)
		test.EXPECT_EQ(t, err, nil, "Failed to Abs(%s): %s", p, err)

		_, exist := watched[absp]
		test.EXPECT_TRUE(t, exist, "Expected %s exist in watched", absp)
	}
}

func testWatchers(t *testing.T, watchers []string, expWatchers []string) {
	test.EXPECT_EQ(t, len(watchers), len(expWatchers), "Expected watchers %v keys equal to %v", watchers, expWatchers)

	for i, p := range expWatchers {
		absp, err := filepath.Abs(p)
		test.EXPECT_EQ(t, err, nil, "Failed to Abs(%s): %s", p, err)
		test.EXPECT_EQ(t, watchers[i], absp, "Expected watchers %s to be %s", watchers[i], absp)
	}
}

func TestNewWatcher(t *testing.T) {
	watcher := newWatcher(t)
	defer watcher.Close()

	test.EXPECT_EQ(t, len(watcher.dirs), 0, "Expected len(dirs) of new watcher %d, but got %d", 0, len(watcher.dirs))
	test.EXPECT_EQ(t, len(watcher.watchers), 0, "Expected len(watchers) of new watcher %d, but got %d", 0, len(watcher.watchers))
}

type dummy struct {
	name    string
	c       chan bool
	lock    sync.Mutex
	created bool
	changed bool
	renamed bool
	removed bool
}

func newDummy(name string) *dummy {
	return &dummy{name: name, c: make(chan bool, 5)}
}

func (d *dummy) reset() {
	d.created = false
	d.changed = false
	d.renamed = false
	d.removed = false
}

func (d *dummy) done(name string, got *bool) {
	// fmt.Println("Dummy: ", got, name == d.name, name, d.name)
	if name != d.name {
		return
	}
	d.lock.Lock()
	defer func() { d.c <- true }() // make sure Unlock() is called first
	defer d.lock.Unlock()          // in order to avoid deadlocks
	*got = true
}

func (d *dummy) FileChanged(name string) {
	d.done(name, &d.changed)
}

func (d *dummy) FileCreated(name string) {
	d.done(name, &d.created)
}

func (d *dummy) FileRemoved(name string) {
	d.done(name, &d.removed)
}

func (d *dummy) FileRenamed(name string) {
	d.done(name, &d.renamed)
}

func (d *dummy) Wait() {
	<-d.c
	for {
		select {
		case <-d.c:
			continue
		case <-time.After(10 * time.Millisecond):
			return
		}
	}
}

func TestWatch(t *testing.T) {
	path := filepath.FromSlash(os.Args[len(os.Args)-1] + "/src/")
	/*testdata := []struct {
		paths       []string
		expWatched  []string
		expWatchers []string
	}{
		{
			[]string{"testdata/dummy.txt", "testdata/test.txt"},
			[]string{"testdata/dummy.txt", "testdata/test.txt"},
			[]string{"testdata/dummy.txt", "testdata/test.txt"},
		},
		{
			[]string{"testdata", "testdata/dummy.txt", "testdata/test.txt"},
			[]string{"testdata", "testdata/dummy.txt", "testdata/test.txt"},
			[]string{"testdata"},
		},
		{
			[]string{"testdata/dummy.txt", "testdata/test.txt", "testdata"},
			[]string{"testdata", "testdata/dummy.txt", "testdata/test.txt"},
			[]string{"testdata"},
		},
	}

	for i, v := range testdata {
		v := v

		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			//t.Parallel()

			watcher := newWatcher(t)
			for _, name := range v.paths {
				watch(t, watcher, name, newDummy(name))
			}
			testWatched(t, watcher.watched, v.expWatched)
			testWatchers(t, watcher.watchers, v.expWatchers)
			defer watcher.Close()
		})
	}*/
	tests := []struct {
		paths       []string
		expWatched  []string
		expWatchers []string
	}{
		{
			[]string{path + "testdata\\dummy.txt", path + "testdata\\test.txt"},
			[]string{path + "testdata\\dummy.txt", path + "testdata\\test.txt"},
			[]string{path + "testdata\\dummy.txt", path + "testdata\\test.txt"},
		},
		/*{
			[]string{path + "testdata", path + "testdata/dummy.txt", path + "testdata/test.txt"},
			[]string{path + "testdata", path + "testdata/dummy.txt", path + "testdata/test.txt"},
			[]string{path + "testdata"},
		},
		{
			[]string{path + "testdata/dummy.txt", path + "testdata/test.txt", path + "testdata"},
			[]string{path + "testdata", path + "testdata/dummy.txt", path + "testdata/test.txt"},
			[]string{path + "testdata"},
		},*/
	}

	logger.SetLevel(logger.DEBUG)
	logger.ShowShortFileName()

	for _, test := range tests {
		watcher := newWatcher(t)
		for _, name := range test.paths {
			watch(t, watcher, name, newDummy(name))
		}
		testWatched(t, watcher.watched, test.expWatched)
		testWatchers(t, watcher.watchers, test.expWatchers)
		defer watcher.Close()
	}
}

func Testwatch(t *testing.T) {
	watcher := newWatcher(t)
	defer watcher.Close()

	err := watcher.watch("testdata/dummy.txt", false)
	test.ASSERT_EQ(t, err, nil, "Couldn't watch %s", "testdata/dummy.txt")

	err = watcher.watch("testdata/test.txt", false)
	test.ASSERT_EQ(t, err, nil, "Couldn't watch %s", "testdata/test.txt")

	testWatched(t, watcher.watched, []string{"testdata/dummy.txt", "testdata/test.txt"})
	testWatchers(t, watcher.watchers, []string{"testdata/dummy.txt", "testdata/test.txt"})
}

func TestAdd(t *testing.T) {
	watcher := newWatcher(t)
	defer watcher.Close()
	d := newDummy("test")
	watcher.addCallback("test", d)
	callback := watcher.watched["test"][0]
	test.EXPECT_EQ(t, callback, d, "Expected watcher['test'][0] callback equal to %v, but got %v", d, callback)
}

func TestFlushDir(t *testing.T) {
	name := "testdata/dummy.txt"
	dir, _ := filepath.Abs("testdata")
	watcher := newWatcher(t)
	defer watcher.Close()
	d := newDummy(name)
	watch(t, watcher, name, d)
	testWatchers(t, watcher.dirs, []string{})
	testWatchers(t, watcher.watchers, []string{name})
	watcher.flushDir(dir)
	testWatchers(t, watcher.dirs, []string{dir})
	testWatchers(t, watcher.watchers, []string{})
}

func TestUnWatch(t *testing.T) {
	name, _ := filepath.Abs("testdata/dummy.txt")
	watcher := newWatcher(t)
	defer watcher.Close()
	d := newDummy(name)
	watch(t, watcher, name, d)
	unwatch(t, watcher, name, d)
	test.EXPECT_EQ(t, len(watcher.watched), 0, "Expected watcheds be empty, but got %v", watcher.watched)
}

func TestUnWatchAll(t *testing.T) {
	name, _ := filepath.Abs("testdata/dummy.txt")
	watcher := newWatcher(t)
	defer watcher.Close()
	d1 := new(dummy)
	d2 := new(dummy)
	watch(t, watcher, name, d1)
	watch(t, watcher, name, d2)

	l := len(watcher.watched[name])
	test.EXPECT_EQ(t, l, 2, "Expected len of watched['%s'] be %d, but got %d", name, 2, l)

	unwatch(t, watcher, name, nil)
	_, exist := watcher.watched[name]
	test.EXPECT_FALSE(t, exist, "Expected all %s watched be removed", name)

	testWatchers(t, watcher.watchers, []string{})
}

func TestUnWatchDirectory(t *testing.T) {
	name := "testdata/dummy.txt"
	absname, _ := filepath.Abs(name)
	dir, _ := filepath.Abs("testdata")
	watcher := newWatcher(t)
	defer watcher.Close()
	d := newDummy(absname)
	watch(t, watcher, name, d)
	watch(t, watcher, dir, nil)
	testWatchers(t, watcher.watchers, []string{dir})
	unwatch(t, watcher, dir, nil)
	testWatchers(t, watcher.watchers, []string{name})
}

func TestUnWatchOneOfSubscribers(t *testing.T) {
	name, _ := filepath.Abs("testdata/dummy.txt")
	watcher := newWatcher(t)
	defer watcher.Close()
	d1 := new(dummy)
	d2 := new(dummy)
	watch(t, watcher, name, d1)
	watch(t, watcher, name, d2)
	if len(watcher.watched[name]) != 2 {
		t.Fatalf("Expected watched[%s] length be %d, but got %d", name, 2, len(watcher.watched[name]))
	}
	unwatch(t, watcher, name, d1)
	testWatchers(t, watcher.watchers, []string{name})
	if len(watcher.watched[name]) != 1 {
		t.Errorf("Expected watched[%s] length be %d, but got %d", name, 1, len(watcher.watched[name]))
	}
}

func TestunWatch(t *testing.T) {
	name := "testdata/dummy.txt"
	watcher := newWatcher(t)
	defer watcher.Close()
	d1 := new(dummy)
	d2 := new(dummy)
	watch(t, watcher, name, d1)
	watch(t, watcher, name, d2)
	if err := watcher.unWatch(name); err != nil {
		t.Fatalf("Couldn't unWatch %s: %s", name, err)
	}
	if _, exist := watcher.watched[name]; exist {
		t.Errorf("Expected all %s watched be removed", name)
	}
	// if !reflect.DeepEqual(watcher.watchers, []string{}) {
	// 	t.Errorf("Expected watchers be empty but got %v", watcher.watchers)
	// }
	testWatchers(t, watcher.watchers, []string{})
}

func TestRemoveWatch(t *testing.T) {
	name, _ := filepath.Abs("testdata/dummy.txt")
	watcher := newWatcher(t)
	defer watcher.Close()
	d := newDummy(name)
	watch(t, watcher, name, d)
	watcher.removeWatch(name)
	// if !reflect.DeepEqual(watcher.watchers, []string{}) {
	// 	t.Errorf("Expected watchers be empty but got %v", watcher.watchers)
	// }
	testWatchers(t, watcher.watchers, []string{})
}

func TestRemoveDir(t *testing.T) {
	name, _ := filepath.Abs("testdata/dummy.txt")
	dir, _ := filepath.Abs("testdata")
	watcher := newWatcher(t)
	defer watcher.Close()
	d := newDummy(name)
	watch(t, watcher, dir, d)
	watch(t, watcher, name, d)
	testWatchers(t, watcher.watchers, []string{dir})
	testWatchers(t, watcher.dirs, []string{dir})
	watcher.removeDir(dir)
	testWatchers(t, watcher.dirs, []string{})
	testWatchers(t, watcher.watchers, []string{dir, name})
}

func TestObserve(t *testing.T) {
	name := "testdata/test.txt"
	absname, _ := filepath.Abs(name)
	watcher := newWatcher(t)
	defer ioutil.WriteFile(name, []byte(""), 0644)
	defer watcher.Close()
	d := newDummy(absname)
	watch(t, watcher, name, d)

	if err := ioutil.WriteFile(name, []byte("test"), 0644); err != nil {
		t.Fatalf("WriteFile error: %s", err)
	}

	d.Wait()
	if !d.changed {
		t.Errorf("Expected dummy Text %s, but got %#v", "Changed", d)
	}
}

func TestCreateEvent(t *testing.T) {
	name := "testdata/new.txt"
	absname, _ := filepath.Abs(name)
	os.Remove(name)
	defer os.Remove(name)
	watcher := newWatcher(t)
	defer watcher.Close()
	d := newDummy(absname)
	watch(t, watcher, name, d)

	testWatchers(t, watcher.watchers, []string{"testdata"})

	if f, err := os.Create(name); err != nil {
		t.Fatalf("File creation error: %s", err)
	} else {
		f.Close()
	}
	d.Wait()
	if !d.created {
		t.Errorf("Expected dummy Text %s, but got %#v", "Created", d)
	}
}

func TestDeleteEvent(t *testing.T) {
	if os.ExpandEnv("$TRAVIS") != "" {
		// This test just times out on travis (ie the callback is never called).
		// See https://github.com/limetext/lime/issues/438
		t.Skip("Skipping test as it doesn't work with travis")
		return
	}
	name := "testdata/dummy.txt"
	absname, _ := filepath.Abs(name)
	watcher := newWatcher(t)
	defer watcher.Close()
	d := newDummy(absname)
	watch(t, watcher, name, d)

	if err := os.Remove(name); err != nil {
		t.Fatalf("Couldn't remove file %s: %s", name, err)
	}
	d.Wait()
	if !d.removed {
		t.Errorf("Expected dummy Text %s, but got %#v", "Removed", d)
	}
	if f, err := os.Create(name); err != nil {
		t.Errorf("Couldn't create file: %s", err)
	} else {
		f.Close()
	}
	d.Wait()
	if !d.created {
		t.Errorf("Expected dummy Text %s, but got %#v", "Created", d)
	}
}

func TestRenameEvent(t *testing.T) {
	name := "testdata/test.txt"
	absname, _ := filepath.Abs(name)
	defer os.Rename("testdata/rename.txt", name)
	watcher := newWatcher(t)
	defer watcher.Close()
	d := newDummy(absname)
	watch(t, watcher, name, d)

	os.Rename(name, "testdata/rename.txt")
	d.Wait()
	if !d.renamed {
		t.Errorf("Expected dummy Text %s, but got %#v", "Renamed", d)
	}
}
