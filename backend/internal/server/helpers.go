package server

import (
	"log"
)

func (app *CommonGlobals) BackgroundTask(fn func() error) {
	// app.wg.Add(1)

	go func() {
		// defer app.wg.Done()

		defer func() {
			err := recover()
			if err != nil {
				log.Println(err)
			}
		}()

		err := fn()
		if err != nil {
			log.Println(err)
		}
	}()
}
