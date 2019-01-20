package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

type questions []question

type question struct {
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

	questions := initQuestions()
	promptQuestions(&questions)
	save(&questions)

}

//initialize questions with defaults
func initQuestions() questions {
	var questions = make([]question, 10)

	questions[0] = question{
		q:   "[HP]Enter Relay Server IP: ",
		key: "RelayIP",
	}

	questions[1] = question{
		q:   "[HP]Enter Relay Server Port: ",
		key: "RelayPort",
	}

	questions[2] = question{
		q:   "[HP]Enter User ID: ",
		key: "UID",
	}

	questions[3] = question{
		q:   "Enter your host name: ",
		key: "Host",
	}

	questions[4] = question{
		q:   "Enter your service name: ",
		key: "ServiceName",
	}

	questions[5] = question{
		q:   "Enter your service Weight: ",
		key: "ServiceWeight",
	}

	questions[6] = question{
		q:   "Enter your TTL: ",
		key: "TTL",
	}

	questions[7] = question{
		q:   "Enter your service priority: ",
		key: "Priority",
	}

	questions[8] = question{
		q:   "Enter your local port: ",
		key: "LocalPort",
	}

	questions[9] = question{
		q:   "Enter your local IP: ",
		key: "LocalIP",
	}

	return questions

}

func promptQuestions(qs *questions) {
	buf := bufio.NewReader(os.Stdin)
	fmt.Println("Press Enter to accept the default")
	for i, q := range *qs {
		fmt.Print(q.q)
		a, _ := buf.ReadString('\n')
		a = strings.TrimSpace(a)

		if len(a) > 0 {
			(*qs)[i].answer = a
		} else {
			(*qs)[i].answer = viper.GetString(q.key)
		}
	}
}

func save(values *questions) {
	for _, q := range *values {
		viper.Set(q.key, q.answer)
	}

	viper.WriteConfig()
}

func init() {
	RootCmd.AddCommand(initCmd)
}
