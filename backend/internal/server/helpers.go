package server

import (
	"fmt"
	"runtime/debug"
)

func (app *CommonGlobals) BackgroundTask(fn func()) {
	// app.wg.Add(1)

	go func() {
		// defer app.wg.Done()

		defer func() {
			r := recover()
			if r != nil {
				app.Logger.Error("Recovered: " + fmt.Sprintln(r))
				debug.PrintStack()
			}
		}()

		fn()
	}()
}
