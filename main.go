package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("checking env variables")
	server, err := getEnvVar("ARGOCD_SERVER")
	if err != nil {
		log.Panic().Msg("ARGOCD_SERVER variable is not set")
	}

	token, err := getEnvVar("ARGOCD_AUTH_TOKEN")
	if err != nil {
		log.Panic().Msg("ARGOCD_AUTH_TOKEN variable is not")
	}

	log.Info().Msg("going to check command arguments")

	appName, err := getEnvVar("APP_NAME")
	if err != nil {
		log.Panic().Err(err).Msg("app name is not set")
	}

	log.Info().Msgf("going to sync the application %s", appName)
	_, err = runCommand("argocd app sync "+appName, server, token)

	syncFailed := false
	if err != nil {
		log.Error().Err(err).Msgf("problem trying to sync the application %s", appName)
		syncFailed = true
	}

	log.Info().Msgf("going to wait for the application %s", appName)
	output, err := runCommand("argocd app wait --sync "+appName, server, token)
	fmt.Print(output)
	if err != nil {
		log.Error().Err(err).Msgf("problem trying to wait the application %s", appName)
		syncFailed = true
	}

	if syncFailed {
		log.Info().Msg("waiting for 3 minutes, sync failed so wait until automatic poll")
		time.Sleep(time.Minute * 3)
		return
	}

	log.Info().Msg("finished")
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
