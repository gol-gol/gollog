package gollog

import (
	"fmt"
	"io"
	"log"
	"os"
)

const DefaultLogFile = "default-gollog.log"

var (
	Level          string
	Thread         chan GollogMsg
	SaveAt         string
	LogFileHandler *os.File
	Display        bool
	Persist        bool
)

// Open, Log, Close
func LogOnce(logfile, level, msg string) {
	var err error
	LogFileHandler, err = openLogFile(logfile)
	if err != nil {
		log.Println(err)
		return
	}
	defer closeLogFile()
	golMsg := GollogMsg{Level: GollogLevel(level), Msg: msg}
	processMsg(golMsg)
}

// start Log Action
func Start(proc func()) {
	Thread = make(chan GollogMsg)
	if Level == "" {
		Level = "Info"
		Display = true
	}
	if Persist && SaveAt == "" {
		SaveAt = DefaultLogFile
	}

	Persist = false
	go proc()
}

func Stop() {
	Thread <- GollogMsg{Level: "STOP"}
}

func LogIt() {
	var err error
	if SaveAt != "" {
		if LogFileHandler, err = openLogFile(SaveAt); err == nil {
			Persist = true
		} else {
			Display = true
			log.Println(err)
		}
	}
	abovePriority := getPriority(Level)
	for {
		msg := <-(Thread)
		if msg.Level == "STOP" {
			closeLogFile()
			break
		} else if abovePriority <= getPriority(msg.Level) {
			processMsg(msg)
		}
	}
}

func processMsg(m GollogMsg) {
	if Display {
		fmt.Printf("[%v] %v\n", m.Level, m.Msg)
	}
	if Persist {
		lyn := fmt.Sprintf("%v\n", m.Msg)
		n, err := io.WriteString(LogFileHandler, lyn)
		if err != nil {
			fmt.Println(n, err)
		}
	}
}

// Logfile sends back handle of opened logfile, remember to defer F.Close() at usage.
func openLogFile(logfile string) (*os.File, error) {
	return os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
}

// close LogFile
func closeLogFile() {
	if LogFileHandler != nil {
		LogFileHandler.Close()
	}
}

func Debug(msg string) {
	golMsg := GollogMsg{Level: LevelDebug, Msg: msg}
	Thread <- golMsg
}

func Debugf(msgTmpl string, param ...interface{}) {
	msg := fmt.Sprintf(msgTmpl, param...)
	golMsg := GollogMsg{Level: LevelDebug, Msg: msg}
	Thread <- golMsg
}

func Info(msg string) {
	golMsg := GollogMsg{Level: LevelInfo, Msg: msg}
	Thread <- golMsg
}

func Infof(msgTmpl string, param ...interface{}) {
	msg := fmt.Sprintf(msgTmpl, param...)
	golMsg := GollogMsg{Level: LevelInfo, Msg: msg}
	Thread <- golMsg
}

func Warn(msg string) {
	golMsg := GollogMsg{Level: LevelWarning, Msg: msg}
	Thread <- golMsg
}

func Warnf(msgTmpl string, param ...interface{}) {
	msg := fmt.Sprintf(msgTmpl, param...)
	golMsg := GollogMsg{Level: LevelWarning, Msg: msg}
	Thread <- golMsg
}

func Error(msg string) {
	golMsg := GollogMsg{Level: LevelError, Msg: msg}
	Thread <- golMsg
}

func Errorf(msgTmpl string, param ...interface{}) {
	msg := fmt.Sprintf(msgTmpl, param...)
	golMsg := GollogMsg{Level: LevelError, Msg: msg}
	Thread <- golMsg
}

func Panic(msg string) {
	golMsg := GollogMsg{Level: LevelPanic, Msg: msg}
	Thread <- golMsg
}

func Panicf(msgTmpl string, param ...interface{}) {
	msg := fmt.Sprintf(msgTmpl, param...)
	golMsg := GollogMsg{Level: LevelPanic, Msg: msg}
	Thread <- golMsg
}
