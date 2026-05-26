package pkg

import (
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func Init() {
	InfoLogger = log.New(os.Stdout, "INFO:", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stdout, "ERROR:", log.Ldate|log.Ltime|log.Lshortfile)
}

func LogInfo(message string) {
	InfoLogger.Println(message)
}

func LogError(err string) {
	ErrorLogger.Println(err)
}
func LogFattal(message string) {
	ErrorLogger.Fatal(message)
}
