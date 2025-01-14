package tasks

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	ErrCantAddEmptyTask = errors.New("Can't add an empty task")
	ErrCouldNotOpenFile = errors.New("Could not open file")
)

type Task struct {
	Description string
}

func New(text string) (*Task, error) {
	if text == "" {
		return nil, ErrCantAddEmptyTask
	}

	return &Task{Description: text}, nil
}

func ReadOpenTasks(filename string) ([]string, error) {
	tasks, err := ReadTasks(filename)
	if err != nil {
		return nil, err
	}

	result := []string{}
	for _, task := range tasks {
		if strings.HasPrefix(task, "- [ ]") {
			result = append(result, task)
		}
	}

	return result, nil
}

func ReadClosedTasks(filename string) ([]string, error) {
	tasks, err := ReadTasks(filename)
	if err != nil {
		return nil, err
	}

	result := []string{}
	for _, task := range tasks {
		if strings.HasPrefix(task, "- [x]") {
			result = append(result, task)
		}
	}

	return result, nil
}

func ReadTasks(filename string) ([]string, error) {
	result := []string{}

	file, err := os.Open(filename)
	if err != nil {
		return nil, ErrCouldNotOpenFile
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "- [ ] ") || strings.HasPrefix(line, "- [x] ") {
			result = append(result, line)
		}
	}

	return result, nil
}

func CompleteTask(filename string, taskNum int) error {
	taskNum = taskNum - 1
	lines := []string{}

	file, err := os.Open(filename)
	if err != nil {
		return ErrCouldNotOpenFile
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	counter := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "- [ ] ") {
			if counter == taskNum {
				line = strings.TrimPrefix(line, "- [ ] ")
				line = "- [x] " + line
			}
			counter++
		}
		lines = append(lines, line)
	}

	file, err = os.Create(filename)
	if err != nil {
		return ErrCouldNotOpenFile
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	// Ensure everything is written to the file
	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

func AddNewTask(filename, task string) error {
	task = "- [ ] " + task
	newContents := []string{}

	// If file doesn't exist yet
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("File doesn't exist, creating new file with name %s in dailies directory\n", filename)
		file, err := os.Create(filename)
		if err != nil {
			return ErrCouldNotOpenFile
		}
		defer file.Close()

		writer := bufio.NewWriter(file)
		_, err = writer.WriteString(task + "\n")
		if err != nil {
			return err
		}

		// Ensure everything is written
		return writer.Flush()
	}

	// Read existing file contents
	file, err := os.Open(filename)
	if err != nil {
		return ErrCouldNotOpenFile
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		newContents = append(newContents, line)
	}
	newContents = append(newContents, task) // Append new task

	// Reopen file in write mode
	file, err = os.Create(filename) // Overwrites file with new contents
	if err != nil {
		return ErrCouldNotOpenFile
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range newContents {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	// Ensure everything is written
	return writer.Flush()
}
