
## gollog

> few helper/wrapper functions for simplistic level based logging
>
> extracted from an older pandora box of such packages at [abhishekkr/gol](https://github.com/abhishekkr/gol)

### Public Functions

* `LogOnce(logfile, level, msg string)`, can be used without calling `Start()`

* `Start() error`, to be called at init/main before following Log Level functions to be used.
* `Debug(msg string)`
* `Debugf(msgTmpl string, param ...interface{})`
* `Info(msg string)`
* `Infof(msgTmpl string, param ...interface{})`
* `Warn(msg string)`
* `Warnf(msgTmpl string, param ...interface{})`
* `Error(msg string)`
* `Errorf(msgTmpl string, param ...interface{})`
* `Panic(msg string)`
* `Panicf(msgTmpl string, param ...interface{})`


#### Public Variables available to be tuned

* `Level` ("Debug", "Info", "Warn", "Error", "Panic")
* `SaveAt` (if to be persisted, log file path)
* `Display` (if to be printed on display, set to true)
* `Persist` (if to be persisted, set to true)

### Wrappers

* `SetGinLog(*gin.Engine)` to be used for [Gin](https://github.com/gin-gonic/gin) to log json details for requests

---
