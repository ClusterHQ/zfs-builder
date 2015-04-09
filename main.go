package main

import (
	"bufio"
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"os/exec"
	"strings"
)

type Settings map[string]string

func main() {
	settings := getSettings()
	kernel, channel := getBuildEnv()
	operatingSystem := "coreos"
	exists, err := checkReleaseExists(operatingSystem, channel, kernel)
	if err != nil {
		sendReport(settings, err, []byte("Error from checkReleaseExists"), kernel, channel)
		return
	}
	if exists {
		sendReport(settings, nil, []byte("Build already exists, skipping..."), kernel, channel)
		return
	}
	// We got here so this is a new kernel version never seen before. Build it!
	err, lines := runBuild()
	if err == nil {
		// This means the build command outputted a valid artifact.
		// Upload it to github.
		pushToGit(operatingSystem, channel, kernel)
	}
	sendReport(settings, err, lines, kernel, channel)
}

func runCommand(cmds ...string) []byte {
	log.Printf("Running command %s", strings.Join(cmds, " "))
	out, cmdErr := exec.Command("rm", "-rf", "zfs-binaries").CombinedOutput()
	if cmdErr != nil {
		log.Fatal(cmdErr, "\n\n", string(out))
	}
	return out
}

func pushToGit(operatingSystem string, channel string, kernel string) {
	gentooDir := "/home/core/gentoo"
	gitDir := "/home/core/zfs-binaries"
	releaseFile := fmt.Sprintf("zfs-%s.tar.gz", kernel)
	runCommand("git", "clone", "git@github.com:clusterhq/zfs-binaries")
	runCommand("mkdir", "-p", fmt.Sprintf("%s/%s", gitDir, operatingSystem))
	runCommand("cp", fmt.Sprintf("%s/%s", gentooDir, releaseFile),
		fmt.Sprintf("zfs-binaries/%s/", operatingSystem))
	cmdErr := os.Chdir(gitDir)
	if cmdErr != nil {
		log.Fatal(cmdErr)
	}
	runCommand("git", "add", releaseFile)
	runCommand("git", "commit", "-m",
		fmt.Sprintf("Automated build for kernel %s on %s %s.",
			kernel, operatingSystem, channel))
	runCommand("git", "push")
}

func checkReleaseExists(operatingSystem string, channel string, kernel string) (bool, error) {
	// Check whether a release exists on GitHub... returns true/false, or an
	// error.
	url := "https://raw.githubusercontent.com/ClusterHQ/zfs-binaries/master"
	resp, err := http.Head(fmt.Sprintf("%s/%s/zfs-%s.tar.gz", url, operatingSystem, kernel))
	if err != nil {
		return false, err
	}
	if resp.StatusCode == 200 {
		return true, nil
	} else if resp.StatusCode == 404 {
		return false, nil
	} else {
		return false, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}
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

func sendReport(settings Settings, reportErr error, buffer []byte, kernel string, channel string) {
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
