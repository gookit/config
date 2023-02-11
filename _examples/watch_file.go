package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"github.com/gookit/goutil"
	"github.com/gookit/goutil/cliutil"
)

func main() {
	config.AddDriver(yaml.Driver)
	config.WithOptions(
		config.ParseEnv,
		config.WithHookFunc(func(event string, c *config.Config) {
			if event == config.OnReloadData {
				cliutil.Cyanln("config reloaded, you can do something ....")
			}
		}),
	)

	// load app config files
	err := config.LoadFiles(
		"testdata/json_base.json",
		"testdata/yml_base.yml",
		"testdata/yml_other.yml",
	)
	if err != nil {
		panic(err)
	}

	// mock server running
	done := make(chan bool)

	// watch loaded config files
	err = watchConfigFiles(config.Default())
	goutil.PanicErr(err)

	cliutil.Infoln("loaded config files is watching ...")
	<-done
}

func watchConfigFiles(cfg *config.Config) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	//noinspection GoUnhandledErrorResult
	defer watcher.Close()

	// get loaded files
	files := cfg.LoadedFiles()
	if len(files) == 0 {
		return nil
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok { // 'Events' channel is closed
					cliutil.Infoln("'Events' channel is closed ...", event)
					return
				}

				// if event.Op > 0 {
				cliutil.Infof("file event: %s\n", event)

				if event.Op&fsnotify.Write == fsnotify.Write {
					cliutil.Infof("modified file: %s\n", event.Name)

					err := cfg.ReloadFiles()
					if err != nil {
						cliutil.Errorf("reload config error: %s\n", err.Error())
					}
				}
				// }

			case err, ok := <-watcher.Errors:
				if ok { // 'Errors' channel is not closed
					cliutil.Errorf("watch file error: %s\n", err.Error())
				}
				if err != nil {
					cliutil.Errorf("watch file error2: %s\n", err.Error())
				}
				return
			}
		}
	}()

	for _, path := range files {
		cliutil.Infof("add watch file: %s\n", path)
		if err := watcher.Add(path); err != nil {
			return err
		}
	}
	return nil
}
