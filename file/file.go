package file

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func GetCurrentDir() string {
	// Получаем текущую рабочую директорию (проекта)
	currentDir, err := os.Getwd()

	if err != nil {
		fmt.Println("GetCurrentDir:", err)
	}

	return currentDir
}

// GeneratePath работает с директорией assets в корне модуля
func GeneratePath(p string) string {
	return filepath.Join(GetCurrentDir(), "/assets/", p)
}

// OpenFile - Open / Create file
// p - path to file
func OpenFile(p string) *os.File {
	path := GeneratePath(p)

	if !Exists(path) {
		// Открываем файл для записи (или создаем новый, если его нет) и добавляем данные в конец
		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatalf("Cannot open file %s: %v", path, err)
		}

		fmt.Println("Create file to path:", path)

		return f
	} else {
		fmt.Println("File exist to path:", path)
		return nil
	}

}
func CloseFile(f *os.File) {
	err := f.Close()

	if err != nil {
		log.Fatalf("Cannot close file %s: %v", f.Name(), err)
	}
}
func WriteString(t string, f *os.File) {
	if f != nil {
		_, err := f.WriteString(t)
		if err != nil {
			log.Fatalf("Cannot write t in file %s: %v", f.Name(), err)
		}

		CloseFile(f)
	}
}

// Exists returns whether the given file or directory exists
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
func CreateDirToPath(p string) {
	path := GeneratePath(p)

	if !Exists(path) {
		err := os.Mkdir(path, 0777)
		if err != nil {
			log.Fatalf("Error create dir to path %s: %v", path, err)
		}

		fmt.Println("Create dir to path:", path)
	} else {
		fmt.Println("Dir exist to path:", path)
	}
}
