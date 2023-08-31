package pkg

import (
	"github.com/robfig/cron"
)

func CronInit(spec string, task func()) error {
	c := cron.New()
	c.Start()

	err := c.AddFunc(spec, task)
	if err != nil {
		return err
	}

	go func() {
		select {}
	}()

	return nil
}
