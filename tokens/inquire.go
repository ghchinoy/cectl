package tokens

import (
	"github.com/AlecAivazis/survey"
)

var usernameQuestion = survey.Question{
	Name:     "un",
	Prompt:   &survey.Input{Message: "Username"},
	Validate: survey.Required,
}

var passwordQuestion = survey.Question{
	Name:     "pwd",
	Prompt:   &survey.Password{Message: "Password"},
	Validate: survey.Required,
}
var environmentQuestion = survey.Question{
	Name: "env",
	Prompt: &survey.Select{
		Message: "Choose an environment:",
		Options: []string{"snapshot", "staging", "production", "uk"},
		Default: "snapshot",
	},
	Transform: environmentTransformer,
}
var outputQuestion = survey.Question{
	Name: "output",
	Prompt: &survey.Select{
		Message: "Choose an output method:",
		Options: []string{"toml", "json"},
		Default: "toml",
	},
}

// not used
var prompts = []*survey.Question{
	&usernameQuestion,
	&passwordQuestion,
	&environmentQuestion,
	&outputQuestion,
}

var u, p, e, o string

var environments map[string]string

func init() {
	environments = make(map[string]string)
	environments["snapshot"] = "https://snapshot.cloud-elements.com/elements/api-v2"
	environments["staging"] = "https://staging.cloud-elements.com/elements/api-v2"
	environments["production"] = "https://api.cloud-elements.com/elements/api-v2"
	environments["uk"] = "https://console.cloud-elements.co.uk/elements/api-v2"

}

func LoginInquiry() (string, string, string, error) {

	var org, user string
	var config Config

	survey.Ask([]*survey.Question{&usernameQuestion}, &config)
	survey.Ask([]*survey.Question{&passwordQuestion}, &config)
	survey.Ask([]*survey.Question{&environmentQuestion}, &config)
	config.Output = "toml"

	token, err := ObtainCEToken(config)
	if err != nil {
		return config.Environment, org, user, err
	}

	return config.Environment, token.Organization, token.User, nil

}

// environmentTransformer is used by the AlecAivazis.survey package
// to look up the Cloud Elements endpoint base URL from a string input
func environmentTransformer(answer interface{}) interface{} {
	s, ok := answer.(string)
	if !ok {
		return nil
	}
	return environments[s]
}

/*
// askFor prompts for user input, optionally a password field
func askFor(prompt string, isPassword bool) (string, error) {
	var result string
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s ", prompt)
	if isPassword {
		pwdBytes, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return result, err
		}
		result = string(pwdBytes)
		fmt.Println()
	} else {
		var err error
		result, err = reader.ReadString('\n')
		if err != nil {
			return result, err
		}
	}
	return strings.TrimSpace(result), nil
}
*/
