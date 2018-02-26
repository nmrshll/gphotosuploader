package auth

import (
	"fmt"
	"log"
	"strings"

	"github.com/palantir/stacktrace"
)

type AuthenticationOptions struct {
	LaunchWebdriverConsent bool
	Silent                 bool
	AuthFilePath           string
}

func Authenticate(authOpts AuthenticationOptions) CookieCredentials {
	// Load authentication parameters
	credentials, err := NewCookieCredentialsFromFile(authOpts.AuthFilePath)
	if err != nil {
		credentials = nil
		log.Printf("Can't use auth file at %s : %s", authOpts.AuthFilePath, err.Error())
	}

	if credentials == nil {
		if authOpts.LaunchWebdriverConsent && !authOpts.Silent {
			fmt.Println("Auth cookies needed ...")
			fmt.Println("Would you like to run the WebDriver CookieCredentials Wizard ? (yes/no)")
			fmt.Println("Check README for more info")

			var answer string
			fmt.Scanln(&answer)
			startWizard := len(answer) > 0 && strings.ToLower(answer)[0] == 'y'

			if !startWizard {
				log.Fatalln("It's not possible to continue, sorry!")
			}
		}

		credentials, err = StartWebDriverCookieCredentialsWizard()
		if err != nil {
			log.Fatalf("Can't complete the login wizard, got: %v\n", err)
		} else {
			credentials.SerializeToFile(authOpts.AuthFilePath)
			if err != nil {
				log.Fatal(stacktrace.Propagate(err, "failed serializing to file"))
			}
		}

	}

	// // Get a new At token
	// if !authOpts.Silent {
	// 	log.Println("Getting a new At token ...")
	// }
	// token, err := NewAtTokenScraper(*credentials).ScrapeNewAtToken()
	// if err != nil {
	// 	log.Fatalf("Can't scrape a new At token (%v)\n", err)
	// }
	// credentials.RuntimeParameters.AtToken = token
	// if !authOpts.Silent {
	// 	log.Println("At token taken")
	// }
	return *credentials
}
