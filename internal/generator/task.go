package generator

import (
	"github.com/zenthangplus/goccm"
	"log"
	"promogen/internal/utils"
)

func StartTask(in []Instance, Task func(c Instance)) {
	cfg, err := utils.LoadConfig("config.json")
	if err != nil {
		log.Println(err)
		return
	}
	routines := len(in)
	if cfg.Cfg.Limit > 0 {
		routines = cfg.Cfg.Limit
	}
	wg := goccm.New(routines)

	for i := 0; i < len(in); i++ {
		c := in[i]
		wg.Wait()
		go func(i int, c Instance) {
			defer wg.Done()
			Task(c)
		}(i, c)
	}
	wg.WaitAllDone()
}
