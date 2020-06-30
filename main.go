package main

import (
	"bufio"
	"flag"
	"landau/agile-results/src/app"
	"landau/agile-results/src/prompt"
	"os"

	"github.com/adlio/trello"
	"github.com/sirupsen/logrus"
)

// Flags - Supported CLI flags
type Flags struct {
	verbose bool
}

// parseFlags - Parses CLI flags
func parseFlags() Flags {
	verbose := flag.Bool("v", false, "Verbose")
	flag.Parse()
	return Flags{verbose: *verbose}
}

func main() {
	flags := parseFlags()
	logger := logrus.New()

	if flags.verbose {
		logger.SetLevel(logrus.DebugLevel)
	}

	// TODO: assert that these values are set
	apiKey := os.Getenv("TRELLO_API_KEY")
	token := os.Getenv("TRELLO_TOKEN")
	listID := os.Getenv("TRELLO_LIST_ID")
	boardID := os.Getenv("TRELLO_BOARD_ID")

	client := trello.NewClient(apiKey, token)
	client.Logger = logger

	// TODO: refactor this into RunApp
	board, err := client.GetBoard(boardID, trello.Defaults())
	if err != nil {
		logrus.Fatalf("Failed: %v", err)
	}

	labels, err := board.GetLabels(trello.Defaults())

	if err != nil {
		logrus.Fatalf("Failed: %v", err)
	}

	err = app.RunApp(&app.Config{
		CardCreator: client,
		Labels:      labels,
		ListID:      listID,
		Logrus:      logger,
		Prompter: prompt.New(
			bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout),
		),
	})

	if err != nil {
		logrus.Fatalf("Failed: %v", err)
	}
}
