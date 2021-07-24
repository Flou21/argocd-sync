package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	log.Println("checking env variables")
	server, err := getEnvVar("ARGOCD_SERVER")
	if err != nil {
		log.Fatalln("ARGOCD_SERVER variable is not set")
	}

	token, err := getEnvVar("ARGOCD_AUTH_TOKEN")
	if err != nil {
		log.Fatalln("ARGOCD_TOKEN variable is not")
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

	if err != nil {
		log.Fatalln("problem trying to sync the application "+appName+"\n", err)
	}
	log.Println(output)

	log.Println("going to wait for the application " + appName)
	output, err = runCommand("argocd app wait --sync "+appName, server, token)
	if err != nil {
		log.Fatalln("problem trying to wait the application "+appName+"\n", err)
	}
	log.Println(output)

	// wait additinal 10 seconds

	log.Println("application is synced going to wait another 10 seconds")
	time.Sleep(10 * time.Second)
}

func runCommand(cmd string, server string, token string) (string, error) {

	cmdString := "export ARGOCD_SERVER=\"" + server + "\" && ARGOCD_AUTH_TOKEN=\"" + token + "\" && " + cmd
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
