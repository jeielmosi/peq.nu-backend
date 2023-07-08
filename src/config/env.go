package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/joho/godotenv"
)

func getEnvPath() string {
	_, b, _, _ := runtime.Caller(0)
	configPath := filepath.Dir(b)
	srcPath := filepath.Dir(configPath)
	rootPath := filepath.Dir(srcPath)
	envPath := filepath.Join(rootPath, "env")

	err := filepath.Walk(envPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
			return err
		}
		log.Printf("dir: %v: name: %s\n", info.IsDir(), path)
		return nil
	})
	if err != nil {
		log.Println(err)
	}

	return envPath
}

var once sync.Once

func loadDefaulfEnv() {
	var err error = nil
	once.Do(
		func() {
			envPath := getEnvPath()
			currentEnvPath := filepath.Join(envPath, envFile)
			err = godotenv.Load(currentEnvPath)
			if err != nil {
				log.Println("Error loading " + envFile + " file: " + err.Error())
				return
			}

			firebasePath := filepath.Join(envPath, firebaseFile)
			err = os.Setenv(FIREBASE_PATH, firebasePath)
			if err != nil {
				log.Println("Error setting firebase path file: " + err.Error())
				return
			}

		},
	)

	if err != nil {
		os.Exit(1)
	}
}

func Load() {
	loadDefaulfEnv()
}
