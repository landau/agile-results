package main

import (
	"bufio"
	"flag"
	"landau/agile-results/src/app"
	"landau/agile-results/src/ollert"
	"landau/agile-results/src/prompt"
	"os"

	"github.com/sirupsen/logrus"
)

// Flags - Supported CLI flags
type Flags struct {
	verbose      bool
	hasChecklist bool
}

// parseFlags - Parses CLI flags
func parseFlags() Flags {
	verbose := flag.Bool("v", false, "Show debug logging")
	hasChecklist := flag.Bool("checklist", false, "Add a checklist to the card")

	flag.Parse()
	return Flags{verbose: *verbose, hasChecklist: *hasChecklist}
}

func main() {
	flags := parseFlags()
	logger := logrus.New()

	if flags.verbose {
		logger.SetLevel(logrus.DebugLevel)
	}

	// TODO: assert that these values are set
	apiKey := os.Getenv("TRELLO_API_KEY")
	token := os.Getenv("TRELLO_API_TOKEN")
	listID := os.Getenv("TRELLO_LIST_ID")
	boardID := os.Getenv("TRELLO_BOARD_ID")

	client := ollert.NewClient(apiKey, token, logger)

	_, err := app.RunApp(&app.Config{
		BoardID:      boardID,
		Client:       client,
		HasChecklist: flags.hasChecklist,
		// TODO: move to CardCreator
		ListID: listID,
		Logrus: logger,
		Prompter: prompt.New(
			bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout),
		),
	})

	if err != nil {
		logrus.Fatalf("Failed to create card: %v", err)
	}
}
