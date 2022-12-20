package splitlog

import (
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var lvls = []zerolog.Level{
	zerolog.TraceLevel,
	zerolog.DebugLevel,
	zerolog.InfoLevel,
	zerolog.WarnLevel,
	zerolog.ErrorLevel,
	zerolog.FatalLevel,
	zerolog.PanicLevel,
	zerolog.NoLevel,
}

func init() {
	rand.Seed(time.Now().Unix())
}

func logWriter(log zerolog.Logger, split bool, done *sync.WaitGroup, when time.Duration, stop <-chan struct{}) {
	defer done.Done()

	ticker := time.NewTicker(when)
	defer ticker.Stop()

	lvl := zerolog.DebugLevel

	for {
		select {
		case <-stop:
			return
		case tl, ok := <-ticker.C:
			if !ok {
				return
			}
			if split {
				lvl = lvls[rand.Intn(len(lvls))]
			}
			log.WithLevel(lvl).Time("when", tl).Int("val", rand.Int()).Msg("just random msg")
		}
	}
}

var dir1, dir2, dir3 string

func TestMain(m *testing.M) {
	var err error
	dir1, err = os.MkdirTemp("", "splitlog")
	if err != nil {
		panic(err)
	}

	dir2, err = os.MkdirTemp("", "rotatelog")
	if err != nil {
		panic(err)
	}

	dir3, err = os.MkdirTemp("", "splitlvllog")
	if err != nil {
		panic(err)
	}

	defer os.RemoveAll(dir1)
	defer os.RemoveAll(dir2)
	defer os.RemoveAll(dir3)
	m.Run()
}

func TestSplitRotateLog(t *testing.T) {
	zrm := ZeroLogRotateManager{
		SplitLogLevel:   true,
		DefaultLogLevel: zerolog.InfoLevel,
		Dir:             dir1,
		FormatFilename:  "2006_01_02_15_04_05_0000",
		FirstRotation:   5 * time.Second,
		RotateDuration:  5 * time.Second,
	}
	defer zrm.Stop()

	require.NoError(t, zrm.Init())

	wg := &sync.WaitGroup{}
	stop := make(chan struct{}, 10)
	when := 200 * time.Millisecond

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go logWriter(zrm.Logger, true, wg, when, stop)
	}

	// wait 30 second
	<-time.After(31 * time.Second)

	for i := zerolog.TraceLevel; i < zerolog.Disabled; i++ {
		lvldir := filepath.Join(dir1, i.String())
		assert.DirExists(t, lvldir)

		entries, err := filepath.Glob(filepath.Join(lvldir, "*.log"))
		require.NoError(t, err)
		if i == zerolog.NoLevel {
			assert.Equal(t, 14, len(entries), "dir: %s", lvldir)
		} else {
			assert.Equal(t, 7, len(entries), "dir: %s", lvldir)
		}
	}

	entries, err := filepath.Glob(filepath.Join(dir1, "*.log"))
	require.NoError(t, err)
	assert.Equal(t, 14, len(entries))
}

func TestSplitRotateLogLevel(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zrm := ZeroLogRotateManager{
		SplitLogLevel:   true,
		DefaultLogLevel: zerolog.InfoLevel,
		Dir:             dir3,
		FormatFilename:  "2006_01_02_15_04_05_0000",
		FirstRotation:   5 * time.Second,
		RotateDuration:  5 * time.Second,
	}
	defer zrm.Stop()

	require.NoError(t, zrm.Init())

	wg := &sync.WaitGroup{}
	stop := make(chan struct{}, 10)
	when := 200 * time.Millisecond

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go logWriter(zrm.Logger, true, wg, when, stop)
	}

	// wait 30 second
	<-time.After(31 * time.Second)

	for i := zerolog.TraceLevel; i < zerolog.Disabled; i++ {
		lvldir := filepath.Join(dir3, i.String())
		if i > zerolog.DebugLevel {
			assert.DirExists(t, lvldir)

			entries, err := filepath.Glob(filepath.Join(lvldir, "*.log"))
			require.NoError(t, err)
			if i == zerolog.NoLevel {
				assert.Equal(t, 14, len(entries), "dir: %s", lvldir)
			} else {
				assert.Equal(t, 7, len(entries), "dir: %s", lvldir)
			}
		} else {
			assert.NoDirExists(t, lvldir)
		}
	}

	entries, err := filepath.Glob(filepath.Join(dir3, "*.log"))
	require.NoError(t, err)
	assert.Equal(t, 14, len(entries))
}

func TestRotateLog(t *testing.T) {
	zrm := ZeroLogRotateManager{
		SplitLogLevel:   false,
		DefaultLogLevel: zerolog.InfoLevel,
		Dir:             dir2,
		FormatFilename:  "2006_01_02_15_04_05_0000",
		FirstRotation:   5 * time.Second,
		RotateDuration:  5 * time.Second,
	}
	defer zrm.Stop()

	require.NoError(t, zrm.Init())

	wg := &sync.WaitGroup{}
	stop := make(chan struct{}, 10)
	when := 200 * time.Millisecond

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go logWriter(zrm.Logger, true, wg, when, stop)
	}

	// wait 30 second
	<-time.After(31 * time.Second)

	for i := zerolog.TraceLevel; i < zerolog.Disabled; i++ {
		lvldir := filepath.Join(dir2, zerolog.TraceLevel.String())
		assert.NoDirExists(t, lvldir)
	}

	entries, err := filepath.Glob(filepath.Join(dir2, "*.log"))
	require.NoError(t, err)
	assert.Equal(t, 14, len(entries))
}
