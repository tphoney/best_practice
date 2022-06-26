package outputter

import "os"

func WriteToFile(path, content string) error {
	// open file for writing
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	// write content to file
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	// close file
	err = file.Close()
	if err != nil {
		return err
	}
	return nil
}
