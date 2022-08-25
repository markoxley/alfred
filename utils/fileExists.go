package utils

import "os"

// TestFile checks a file exists. If it does not, a default file is created with the contents of cnt

func TestFile(fp, cnt string) error {
	_, err := os.Stat(fp)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil {
		return nil
	}

	file, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString(cnt)
	return nil
}
