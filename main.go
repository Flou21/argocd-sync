package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"os/exec"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	log.Println("checking env variables")
	server, err := getEnvVar("ARGOCD_SERVER")
	if err != nil {
		log.Fatalln("ARGOCD_SERVER variable is not set")
	}

	token, err := getEnvVar("ARGOCD_AUTH_TOKEN")
	if err != nil {
		log.Fatalln("ARGOCD_AUTH_TOKEN variable is not")
	}

	log.Println("going to check command arguments")
	var appName string
	flag.StringVar(&appName, "app", "", "the name of the argocd app")

	flag.Parse()

	if appName == "" {
		log.Fatalln("app flag is not given")
	}

	log.Println("going to sync the application " + appName)
	output, err := runCommand("argocd app sync "+appName, server, token)

	syncFailed := false

	if err != nil {
		log.Println("problem trying to sync the application " + appName)
		log.Println(err)
		syncFailed = true
	}
	log.Println(output)

	log.Println("going to wait for the application " + appName)
	output, err = runCommand("argocd app wait --sync "+appName, server, token)
	if err != nil {
		log.Println("problem trying to wait the application " + appName)
		log.Println(err)
		syncFailed = true
	}
	log.Println(output)

	if syncFailed {
		log.Println("waiting for 3 minutes, sync failed so wait until automatic poll")
		time.Sleep(time.Minute * 3)
		return
	}

	// wait additinal 10 seconds
	log.Println("application is synced going to wait another 10 seconds")
	time.Sleep(10 * time.Second)
}

func runCommand(cmd string, server string, token string) (string, error) {

	cmdString := "ARGOCD_SERVER=\"" + server + "\" && export ARGOCD_AUTH_TOKEN=\"" + token + "\" && " + cmd
	output, err := exec.Command("sh", "-c", cmdString).CombinedOutput()

	if err != nil {
		return "", err
	}
	return string(output), nil
}

func getEnvVar(env string) (string, error) {
	server := os.Getenv(env)

	if server == "" {
		return "", errors.New(env + " is not given")
	}

	return server, nil
}
