package pkg

import (
	"io"
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	RepoInfo    *log.Logger
	RepoError   *log.Logger
	HttpInfo    *log.Logger
	HttpError   *log.Logger
)

func Init() {
	out := io.MultiWriter(os.Stdout, recentLogs)
	InfoLogger = log.New(out, "[INFO:]", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(out, "[ERROR:]", log.Ldate|log.Ltime|log.Lshortfile)
	RepoInfo = log.New(out, "[INFO::Repo]", log.Ldate|log.Ltime|log.Lshortfile)
	RepoError = log.New(out, "[ERROR:Repo]", log.Ldate|log.Ltime|log.Lshortfile)
	HttpInfo = log.New(out, "[INFO::HTTP]", log.Ldate|log.Ltime|log.Lshortfile)
	HttpError = log.New(out, "[ERROR:HTTP]", log.Ldate|log.Ltime|log.Lshortfile)
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
func LogHttp(mesasge string) {
	HttpInfo.Println(mesasge)
}
func LogHttpError(err error) {
	HttpError.Println(err)
}

func LogError(err error) {
	ErrorLogger.Println(err)
}
func LogFattal(message string) {
	ErrorLogger.Fatal(message)
}
