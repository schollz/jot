package sdees

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

type Config struct {
	Remote, Editor, CurrentDocument string
}

func SetupConfig() {
	var configParameters Config

	var yesno string
	err := errors.New("Incorrect remote")
	for {
		fmt.Print("Enter remote (e.g.: git@github.com:USER/REPO.git): ")
		fmt.Scanln(&yesno)
		cwd, _ := os.Getwd()
		os.Chdir(CachePath)
		os.RemoveAll(HashString(yesno))
		// cmd := exec.Command("git", "clone", yesno, HashString(yesno))
		// _, err := cmd.Output()

		cmd := exec.Command("git", "clone", yesno, HashString(yesno))
		// out1, _ := cmd.StdoutPipe()
		out2, _ := cmd.StderrPipe()
		cmd.Start()
		// _, _ := ioutil.ReadAll(out1)
		out2b, _ := ioutil.ReadAll(out2)
		cmd.Wait()
		os.Chdir(cwd)
		fmt.Println(strings.TrimSpace(string(out2b)))
		if !strings.Contains(string(out2b), "Cloning into ") {
			// logger.Debug("Tried command '%s' in path %s", "git clone "+yesno+" "+HashString(yesno), CachePath)
			fmt.Println("Could not clone, please re-enter")
		} else {
			break
		}
	}
	configParameters.Remote = yesno

	fmt.Printf("Which editor do you want to use: vim (default), nano, or emacs? ")
	fmt.Scanln(&yesno)
	if strings.TrimSpace(strings.ToLower(yesno)) == "nano" {
		configParameters.Editor = "nano"
	} else if strings.TrimSpace(strings.ToLower(yesno)) == "emacs" {
		configParameters.Editor = "emacs"
	} else {
		configParameters.Editor = "vim"
	}
	configParameters.CurrentDocument = ""

	b, err := json.Marshal(configParameters)
	if err != nil {
		log.Println(err)
	}
	ioutil.WriteFile(path.Join(ConfigPath, "config.json"), b, 0644)
}

func LoadConfiguration() {
	defer timeTrack(time.Now(), "Loaded and saved configuration")
	var c Config
	data, err := ioutil.ReadFile(path.Join(ConfigPath, "config.json"))
	if err != nil {
		logger.Error("Could not load config.json")
		return
	}
	json.Unmarshal(data, &c)
	if len(CurrentDocument) == 0 {
		CurrentDocument = c.CurrentDocument
	}
	Editor = c.Editor
	Remote = c.Remote
	RemoteFolder = path.Join(CachePath, HashString(Remote))
	if len(Remote) == 0 {
		SetupConfig()
	}
}

func SaveConfiguration(editor string, remote string, currentdoc string) {
	defer timeTrack(time.Now(), "Saved configuration")
	c := Config{Editor: editor, Remote: remote, CurrentDocument: currentdoc}
	b, err := json.Marshal(c)
	if err != nil {
		log.Println(err)
	}
	ioutil.WriteFile(path.Join(ConfigPath, "config.json"), b, 0644)
}
