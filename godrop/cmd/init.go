package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// var questions = map[string]string{
// 	"Enter Relay Server IP: ":   "",
// 	"Enter Relay Server Port: ": "",
// 	"Enter User ID: ":           "",
// }

type Questions struct {
	qs []Question
}

type Question struct {
	q      string
	key    string
	answer string
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize godrop by creating a config file",
	Long:  "Initialize some godrop values. You may skip through each field to accept the default.\nQuestions prefixed with [HP] are only used when using godrop in a TCP holepunch context",
	Run:   execInitCommand,
}

func execInitCommand(command *cobra.Command, args []string) {
	// questions := initQuestions()
	// promptQuestions(&questions)

	// if err := writeConfig(&questions); err != nil {
	// 	fmt.Println("Error writing config file.")
	// 	command.Usage()
	// 	os.Exit(1)
	// }

	// home, _ := homedir.Dir()

	// fmt.Printf("Godrop is initialized. Config file written in %s\n", path.Join(home, config))

}

//initialize questions with defaults
func initQuestions() Questions {
	var questions = Questions{
		qs: make([]Question, 3),
	}

	questions.qs[0] = Question{
		q:      "[HP]Enter Relay Server IP: ",
		key:    "RelayIP",
		answer: "127.0.0.1",
	}

	questions.qs[1] = Question{
		q:      "[HP]Enter Relay Server Port: ",
		key:    "RelayPort",
		answer: "8080",
	}

	questions.qs[2] = Question{
		q:      "[HP]Enter User ID: ",
		key:    "UID",
		answer: "godrop",
	}

	return questions

}

func promptQuestions(qs *Questions) {
	buf := bufio.NewReader(os.Stdin)
	for i, q := range (*qs).qs {
		fmt.Print(q.q)
		a, _ := buf.ReadString('\n')
		(*qs).qs[i].answer = a
	}
}

func writeConfig(values *Questions) (err error) {
	home, err := homedir.Dir()

	pathToConf := path.Join(home, config)

	var f *os.File

	if _, err := os.Stat(pathToConf); err == nil {
		//file does exist delete it first
		fmt.Println("Deleting old file")
		os.Remove(pathToConf)
	}

	if f, err = os.Create(pathToConf); f != nil {
		write(f, values)
	}

	return
}

func write(file *os.File, values *Questions) {
	for _, v := range (*values).qs {
		value := v.key + " - " + v.answer
		file.WriteString(value)
	}
}

func init() {
	RootCmd.AddCommand(initCmd)
}
