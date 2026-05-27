package pkg

import (
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	RepoInfo    *log.Logger
	RepoError   *log.Logger
)

func Init() {
	InfoLogger = log.New(os.Stdout, "[INFO:]", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stdout, "[ERROR:]", log.Ldate|log.Ltime|log.Lshortfile)
	RepoInfo = log.New(os.Stdout, "[INFO::Repo]", log.Ldate|log.Ltime|log.Lshortfile)
	RepoError = log.New(os.Stdout, "[ERROR:Repo]", log.Ldate|log.Ltime|log.Lshortfile)

}

func LogInfo(message string) {
	InfoLogger.Println(message)
}
func LogRepo(mesasge string) {
	RepoInfo.Println(mesasge)
}
func LogRepoError(err error) {
	RepoError.Println(err)
}

func LogError(err error) {
	ErrorLogger.Println(err)
}
func LogFattal(message string) {
	ErrorLogger.Fatal(message)
}
