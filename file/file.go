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
	// Открываем файл для записи (или создаем новый, если его нет) и добавляем данные в конец
	f, err := os.OpenFile(GeneratePath(p), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Не удалось открыть файл: %v", err)
	}

	fmt.Println("Create file to path:", GeneratePath(p))

	return f
}
func ClosFile(f *os.File) {
	err := f.Close()

	if err != nil {
		log.Fatalf("Cannot close file %s: %v", f.Name(), err)
	}
}
func WriteString(t string, f *os.File) {
	_, err := f.WriteString(t)
	if err != nil {
		log.Fatalf("Cannot write t in file %s: %v", f.Name(), err)
	}

	ClosFile(f)
}

// Exists returns whether the given file or directory exists
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
func CreateDirToPath(p string) {
	if !Exists(p) {
		err := os.Mkdir(GeneratePath(p), 0777)
		if err != nil {
			panic(err)
		}

		fmt.Println("Create dir to path:", GeneratePath(p))
	} else {
		fmt.Println("Dir exist to path:", GeneratePath(p))
	}
}
