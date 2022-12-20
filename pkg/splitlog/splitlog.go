package splitlog

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	ErrLogNotAvailable = errors.New("no log available")
)

// RotatableLog a log that provide functionality to rotate log file
type RotatableLog interface {
	zerolog.LevelWriter
	io.WriteCloser
	Rotate(filename string) error
}

const (
	DefaultZipStructureLog = 7 * 24 * time.Hour
	DefaultConsoleLog      = 3 * 24 * time.Hour
)

type ZeroLogRotateManager struct {
	zerolog.Logger

	// SplitLogLevel indicate whether the log should split into different file
	// base on the log level
	SplitLogLevel bool

	// DefaultLogLevel use for write a log that has no log level
	DefaultLogLevel zerolog.Level

	// RotateDuration interval duration when next file log created
	RotateDuration time.Duration

	// FirstRotation interval until first rotate log file.
	FirstRotation time.Duration

	// Dir directory of where the log write to
	Dir string

	// FormatFilename a time format use to generate log filename
	FormatFilename string

	// ZipStructureLog an interval to compress log structure into a zip file.
	// By default the structure log will be compress within 7 days.
	ZipStructureLog time.Duration

	// PurgableConsoleLog an internal indicate when to purge last log file.
	// By default the console log file will be purge within 3 days.
	// Note: the console log file does not split into multiple file event if `SplitLogLevel` is set to `true`
	PurgableConsoleLog time.Duration

	// Console indicate whether os.Stdout should be include into multi level log.
	Console bool

	splitLog   RotatableLog
	consoleLog RotatableLog
	stop       chan struct{}
	ticker     *time.Ticker
	first      atomic.Bool
}

func (zlm *ZeroLogRotateManager) Init() error {
	if err := zlm.createLog(); err != nil {
		return err
	}
	zlm.first.Store(true)
	zlm.ticker = time.NewTicker(zlm.FirstRotation)
	zlm.stop = make(chan struct{}, 1)
	go zlm.rotate()
	return nil
}

func (zlm *ZeroLogRotateManager) Stop() {
	if zlm.stop != nil {
		close(zlm.stop)
	}

	if zlm.ticker != nil {
		zlm.ticker.Stop()
	}
}

func (zlm *ZeroLogRotateManager) rotate() {
	for {
		select {
		case <-zlm.stop:
			return
		case <-zlm.ticker.C:
			if zlm.first.Load() {
				zlm.first.Store(false)
				zlm.ticker.Stop()
				zlm.ticker = time.NewTicker(zlm.RotateDuration)
			}

			filename := time.Now().Format(zlm.FormatFilename)
			if err := zlm.consoleLog.Rotate(filename + ".console"); err != nil {
				log.Err(err).Msg("unable to rotate console log file")
			}
			if err := zlm.splitLog.Rotate(filename + ".struct"); err != nil {
				log.Err(err).Msg("unable to rotate structure log file")
			}
		}
	}
}

// createLog create multi level writer log
func (zlm *ZeroLogRotateManager) createLog() error {
	filename := time.Now().Format(zlm.FormatFilename)

	var structW RotatableLog
	var err error
	if zlm.SplitLogLevel {
		structW, err = newSplitLogLevel(zlm.DefaultLogLevel, zlm.Dir, filename+".struct")
	} else {
		structW, err = newMergeLog(zlm.Dir, filename+".struct")
	}

	if err != nil {
		log.Err(err).Msg("unable to renew rotated structure log file")
		return err
	}

	consoleFileW, err := newMergeLog(zlm.Dir, filename+".console")
	if err != nil {
		defer structW.Close()
		log.Err(err).Msg("unable to renew rotated console log file")
		return err
	}

	zlm.splitLog = structW
	zlm.consoleLog = consoleFileW

	ws := make([]io.Writer, 2, 3)
	ws[0] = structW
	ws[1] = zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.NoColor = true
		w.Out = consoleFileW
	})

	// output into stdout
	if zlm.Console {
		ws = append(ws, zerolog.NewConsoleWriter())
	}

	zlm.Logger = log.Output(zerolog.MultiLevelWriter(ws...))
	return nil
}

// MergeLog write all log into a single file
type MergeLog struct {
	lck      sync.RWMutex
	wg       sync.WaitGroup
	it       io.WriteCloser
	dir      string
	filename string
	once     sync.Once
}

func newMergeLog(dir, filename string) (*MergeLog, error) {
	return &MergeLog{dir: dir, filename: filename}, nil
}

func (nml *MergeLog) Close() (err error) {
	nml.lck.Lock()
	defer nml.lck.Unlock()
	nml.wg.Wait()
	if err := nml.it.Close(); err != nil {
		fmt.Fprintln(os.Stderr, "unable to close writer", err)
	}
	nml.it = nil
	return nil
}

func (nml *MergeLog) Write(p []byte) (n int, werr error) {
	return nml.WriteLevel(zerolog.NoLevel, p)
}

func (nml *MergeLog) WriteLevel(level zerolog.Level, p []byte) (n int, werr error) {
	nml.lck.RLock()
	defer nml.lck.RUnlock()
	defer nml.wg.Done()
	nml.wg.Add(1)

	// we need to create new file if sync.Once is refresh.
	// once will be update when trigger new rotation. It is prevent new file
	// being create if there is no log has occurred.
	nml.once.Do(func() {
		file := filepath.Join(nml.dir, nml.filename+".log")
		if err := os.Mkdir(nml.dir, 0744); err != nil && !os.IsExist(err) {
			werr = err
			return
		}

		w, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0744)
		if err != nil {
			werr = err
			return
		}
		nml.it = w
	})
	if werr != nil {
		return 0, werr
	}

	if nml.it != nil {
		return nml.it.Write(p)
	}
	return 0, ErrLogNotAvailable
}

func (nml *MergeLog) Rotate(filename string) error {
	nml.lck.Lock()
	defer nml.lck.Unlock()

	// wait until all write is finish
	nml.wg.Wait()
	// the file may never created if log of that specific level has not occurred
	if nml.it != nil {
		if err := nml.it.Close(); err != nil {
			fmt.Fprintln(os.Stderr, "unable to close writer during rotate log", err)
		}
	}
	nml.filename = filename
	nml.once = sync.Once{}
	return nil
}

// SplitLog write log into different file base on log level
type SplitLog struct {
	lvlDefault zerolog.Level
	logs       []*atomic.Pointer[MergeLog]
}

func newSplitLogLevel(def zerolog.Level, dir, filename string) (*SplitLog, error) {
	var llvl = zerolog.GlobalLevel()
	var totalLogs = int(zerolog.NoLevel) + 2
	var logs = make([]*atomic.Pointer[MergeLog], totalLogs)
	for i := 0; i < totalLogs; i++ {
		lvl := zerolog.Level(i - 1)
		// lower level than  global will be ignore by zerolog thus
		// the `Write([]byte)` or `WriteLevel(Level, []byte)` is never occurred
		if lvl >= llvl {
			mlog, err := newMergeLog(filepath.Join(dir, lvl.String()), filename)
			if err != nil {
				return nil, err
			}
			logs[i] = &atomic.Pointer[MergeLog]{}
			logs[i].Store(mlog)
		}
	}
	return &SplitLog{
		lvlDefault: def,
		logs:       logs,
	}, nil
}

func (sll *SplitLog) Close() (err error) {
	for _, ll := range sll.logs {
		if ll == nil {
			continue
		}
		if ml := ll.Swap(nil); ml != nil {
			ml.Close()
		}
	}
	return nil
}

func (sll *SplitLog) Write(p []byte) (n int, err error) {
	return sll.WriteLevel(sll.lvlDefault, p)
}

func (sll *SplitLog) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	lg := sll.logs[level+1]
	if lg == nil {
		return 0, nil
	} else if w := lg.Load(); w == nil {
		if zerolog.GlobalLevel() > level {
			return 0, ErrLogNotAvailable
		}
		return 0, nil
	} else {
		return w.Write(p)
	}
}

func (sll *SplitLog) Rotate(filename string) error {
	var llvl = zerolog.GlobalLevel()
	var totalLogs = int(zerolog.NoLevel) + 2
	for i := 0; i < totalLogs; i++ {
		lvl := zerolog.Level(i - 1)
		if lvl >= llvl {
			if it := sll.logs[i].Load(); it != nil {
				it.Rotate(filename)
			}
		}
	}
	return nil
}
