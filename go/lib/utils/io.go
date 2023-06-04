package utils

import (
	"bufio"
	"os"
)

func ReadLines(name string, callback func(int64, string) bool) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	var i int64 = 0
	for fileScanner.Scan() {
		line := fileScanner.Text()
		ok := callback(i, line)
		if !ok {
			break
		}
		i++
	}
	return nil
}
