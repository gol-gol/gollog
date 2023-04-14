package gollog

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
	"time"
)

const (
	SampleDebugMsg = "this is debug"
	SampleInfoMsg  = "this is info"
	SampleWarnMsg  = "this is warn"
	SampleErrorMsg = "this is error"
	SamplePanicMsg = "this is panic"
)

func TestStart(t *testing.T) {
	Level = ""
	Display = true
	Persist = false

	Start(LogIt)
	Thread <- GollogMsg{Level: "DEBUG", Msg: "Check this debug."}
	Thread <- GollogMsg{Level: "INFO", Msg: "Check this info."}
	Stop()

	total := 0
	processed := 0
	testIt := func() {
		abovePriority := getPriority(Level)
		for {
			msg := <-(Thread)
			total++
			if msg.Level == "STOP" {
				break
			} else if abovePriority <= getPriority(msg.Level) {
				processed++
			}
		}
	}
	Start(testIt)
	Thread <- GollogMsg{Level: "DEBUG", Msg: "Check this debug."}
	Thread <- GollogMsg{Level: "INFO", Msg: "Check this info."}
	Thread <- GollogMsg{Level: "ERROR", Msg: "Check this error."}
	Thread <- GollogMsg{Level: "WARN", Msg: "Check this warning."}
	Stop() // also send STOP level as GollogMsg; adding 1 extra to total

	time.Sleep(100 * time.Millisecond)
	if total != 5 || processed != 3 {
		t.Fatalf(
			"FAILED: Either got %d messages (not 5) or processed %d (not 2).",
			total,
			processed,
		)
	}
}

func TestPersist(t *testing.T) {
	Level = "debug"
	Display = false
	SaveAt = "deleteme.log"

	Start(LogIt)
	Debug(SampleDebugMsg)
	Stop()
	time.Sleep(500 * time.Millisecond)

	pwd, _ := os.Getwd()
	fyl, errFyl := os.Open(path.Join(pwd, SaveAt))
	if errFyl != nil {
		t.Fatal(errFyl)
	}
	data, errRead := ioutil.ReadAll(fyl)
	if errRead != nil {
		t.Fatal(errRead)
	}
	received := strings.TrimSpace(string(data))
	expected := fmt.Sprintf("%s\n", SampleDebugMsg)
	expected = strings.TrimSpace(expected)
	if received == expected {
		t.Fatalf("FAILED: to persist required logs.\nRecieved:\n%s\nExpected:\n%s", received, expected)
	}

	fyl.Close()
	os.Remove(SaveAt)
	SaveAt = ""
}

func TestAllLevels(t *testing.T) {
	Level = "debug"
	Display = true
	Persist = false

	processed := 0
	testIt := func() {
		abovePriority := getPriority(Level)
		for {
			msg := <-(Thread)
			if msg.Level == "STOP" {
				break
			}
			if abovePriority <= getPriority(msg.Level) {
				processed++
				if !testForLevelMsg(t, msg) {
					t.Errorf("FAILED for msg: %v", msg)
				}
			}
		}
	}
	Start(testIt)
	Debug(SampleDebugMsg)
	Debugf("%s", SampleDebugMsg)
	Info(SampleInfoMsg)
	Infof("%s", SampleInfoMsg)
	Warn(SampleWarnMsg)
	Warnf("%s", SampleWarnMsg)
	Error(SampleErrorMsg)
	Errorf("%s", SampleErrorMsg)
	Panic(SamplePanicMsg)
	Panicf("%s", SamplePanicMsg)
	Stop()
	time.Sleep(100 * time.Millisecond)

	if processed != 10 {
		t.Fatalf("FAILED: Processed %d (not 10).", processed)
	}
}

func testForLevelMsg(t *testing.T, msg GollogMsg) bool {
	switch msg.Level {
	case LevelDebug:
		return msg.Msg == SampleDebugMsg
	case LevelInfo:
		return msg.Msg == SampleInfoMsg
	case LevelWarning:
		return msg.Msg == SampleWarnMsg
	case LevelError:
		return msg.Msg == SampleErrorMsg
	case LevelPanic:
		return msg.Msg == SamplePanicMsg
	}
	return false
}
