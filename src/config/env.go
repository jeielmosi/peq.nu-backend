package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/joho/godotenv"
)

func getEnvsPath() string {
	_, b, _, _ := runtime.Caller(0)
	configPath := filepath.Dir(b)
	srcPath := filepath.Dir(configPath)
	rootPath := filepath.Dir(srcPath)
	envsPath := filepath.Join(rootPath, "envs")

	return envsPath
}

func LoadEnv(envName string) {

	envsPath := getEnvsPath()
	envPath := filepath.Join(envsPath, envName)

	firebasePath := filepath.Join(envPath, "firebase.json")
	err := os.Setenv(envName+"_FIREBASE_PATH", firebasePath)
	if err != nil {
		log.Fatalf("Error setting firebase path file: %s", err.Error())
		os.Exit(1)
	}

	dotEnvPath := filepath.Join(envPath, ".env")
	err = godotenv.Overload(dotEnvPath)
	if err != nil {
		log.Fatalf("Error loading " + dotEnvPath + " file")
		os.Exit(1)
	}
}

var once sync.Once

func loadDefaulfEnv() {
	var err error = nil
	once.Do(
		func() {
			envsPath := getEnvsPath()
			currentEnvPath := filepath.Join(envsPath, sharedFile)
			err = godotenv.Load(currentEnvPath)
		},
	)

	if err != nil {
		log.Fatalf("Error loading " + sharedFile + " file")
		os.Exit(1)
	}
}

func Load() {
	loadDefaulfEnv()
	envName := os.Getenv(CURRENT_ENV)
	LoadEnv(envName)
}
