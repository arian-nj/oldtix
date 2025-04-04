package server

import (
	"fmt"
)

func (app *CommonGlobals) BackgroundTask(fn func()) {
	// app.wg.Add(1)

	go func() {
		// defer app.wg.Done()

		defer func() {
			err := recover()
			if err != nil {
				app.ReportError(fmt.Errorf("%s", err))
			}
		}()

		fn()
	}()
}
