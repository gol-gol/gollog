package gollog

import (
	"log"
	"strings"
)

type GollogLevel string

const (
	LevelDebug   GollogLevel = "DEBUG"
	LevelInfo    GollogLevel = "INFO"
	LevelWarning GollogLevel = "WARN"
	LevelError   GollogLevel = "ERROR"
	LevelPanic   GollogLevel = "PANIC"
)

var (
	LevelPriority = map[GollogLevel]uint8{
		LevelDebug:   0,
		LevelInfo:    1,
		LevelWarning: 2,
		LevelError:   3,
		LevelPanic:   4,
	}
)

func getPriority(level interface{}) uint8 {
	switch l := level.(type) {
	case GollogLevel:
		return LevelPriority[l]
	case string:
		l2 := GollogLevel(strings.ToUpper(l))
		return LevelPriority[l2]
	default:
		log.Printf("Providing unrecognised type of level %v\n", l)
	}
	return 0
}
