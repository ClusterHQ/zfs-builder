package main

import (
	"bufio"
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"os/exec"
	"strings"
)

type Settings map[string]string

func main() {
	settings := getSettings()
	err, lines := runBuild()
	kernel, channel := getBuildEnv()
	if err != nil {
		// This means the build command outputted a valid artifact.
		// Upload it to github.
	}
	sendReport(settings, err, lines, kernel, channel)
}

func getBuildEnv() (string, string) {
	kernelVersion, err := exec.Command("uname", "-r").Output()
	if err != nil {
		log.Fatal(err)
	}
	updateFile, err := ioutil.ReadFile("/etc/coreos/update.conf")
	if err != nil {
		log.Fatal(err)
	}
	coreOsChannel := strings.Split(strings.Split(string(updateFile), "\n")[0], "=")[1]
	return string(kernelVersion), coreOsChannel
}

func runBuild() (error, []byte) {
	var buffer bytes.Buffer
	var result []byte

	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	in := bufio.NewScanner(io.MultiReader(stdout, stderr))

	for in.Scan() {
		line := in.Text()
		buffer.WriteString(line + "\n")
		log.Printf(line)
	}
	err = cmd.Wait()
	result = []byte(buffer.String())
	return err, result
}

func getSettings() Settings {
	var settings Settings
	settingsFile, err := ioutil.ReadFile("settings.yml")
	if err != nil {
		log.Fatal(err)
	}
	yaml.Unmarshal(settingsFile, &settings)
	return settings
}

func sendReport(settings Settings, reportErr error,
	buffer []byte, kernel string, channel string) {
	var stringResult string
	if reportErr != nil {
		stringResult = fmt.Sprintf("failure: %v", reportErr)
	} else {
		stringResult = "success"
	}
	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("",
			settings["gmail_smtp_username"],
			settings["gmail_smtp_password"],
			"smtp.gmail.com"),
		settings["email_from"],
		[]string{settings["email_to"]},
		[]byte(fmt.Sprintf(`From: %s
To: %s
Subject: coreos result: %s on %s (CoreOS %s)

Build results:

%s`,
			settings["email_from"],
			settings["email_to"],
			stringResult, kernel, channel,
			buffer)))
	if err != nil {
		log.Fatal(err)
	}
}
