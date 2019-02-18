package cmd

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	. "github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type questions []question

type question struct {
	q      string
	key    string
	answer string
	canAsk bool
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
	questions[1].answer = getGodropDomain()
	save(&questions)
	initText()
}

//initialize questions with defaults
func initQuestions() questions {
	var questions = make([]question, 3)

	questions[0] = question{
		q:      "[HP]Enter User ID: ",
		key:    "UID",
		canAsk: true,
	}

	questions[1] = question{
		q:      "",
		key:    "Host",
		canAsk: false,
	}

	questions[2] = question{
		q:      "Enter your local port: ",
		key:    "LocalPort",
		canAsk: true,
	}

	return questions

}

func promptQuestions(qs *questions) {
	buf := bufio.NewReader(os.Stdin)
	fmt.Println("Press Enter to accept the default")
	for i, q := range *qs {

		if !q.canAsk {
			continue
		}

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

func getGodropDomain() string {
	response, err := http.Get("http://104.248.183.179:80/domain")

	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(response.Body)
	var domain = ""
	for scanner.Scan() {
		domain += scanner.Text()
	}
	return strings.TrimSpace(domain)
}

func init() {
	RootCmd.AddCommand(initCmd)
}

func initText() {
	fmt.Println(Bold("Godrop has been initialized."))
	fmt.Println()
	fmt.Printf("Host: %s\nPort: %d\nUID: %s\n", Bold(viper.GetString("Host")), Bold(viper.GetInt("LocalPort")), Bold(viper.GetString("UID")))
	fmt.Println()
	fmt.Printf("You may run %s to generate your private key and TLS certificate\n", Bold("godrop gencert"))
	fmt.Println()
}
