package edenutil

import "os"

func CreateFileIfNotExists(path string) error {
	var _, err = os.Stat(path)
	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return err
		}
		err = file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func CreatePathIfNotExists(path string) error {
	var _, err = os.Stat(path)
	// create file if not exists
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
