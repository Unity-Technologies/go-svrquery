package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/multiplay/go-svrquery/lib/svrquery"
	"github.com/multiplay/go-svrquery/lib/svrquery/protocol"
)

const (
	numWorkers = 100
)

var (
	errNoItem       = errors.New("no item")
	errEntryInvalid = errors.New("invalid entry")
)

// BulkResponseItem contains the information about the query being performed
// against a single server.
type BulkResponseItem struct {
	Address    string                      `json:"address"`
	ServerInfo *BulkResponseServerInfoItem `json:"serverInfo,omitempty"`
	Error      string                      `json:"error,omitempty"`
}

// BulkResponseServerInfoItem container the server information.
type BulkResponseServerInfoItem struct {
	CurrentPlayers int64  `json:"currentPlayers"`
	MaxPlayers     int64  `json:"maxPlayers"`
	Map            string `json:"map"`
}

// BulkResponseItemWork is an item returned by an worker containing the data item
// plus any terminal error it encountered.
type BulkResponseItemWork struct {
	Item *BulkResponseItem
	Err  error
}

// queryBulk queries a bulk set of servers using a query file.
func queryBulk(file string) error {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Make a jobs channel and a number of workers to processes
	// work off of the channel.
	jobChan := make(chan string)
	resultsChan := make(chan BulkResponseItemWork)
	for w := 1; w <= numWorkers; w++ {
		go worker(jobChan, resultsChan)
	}

	items := make([]BulkResponseItem, 0, 1000)

	// Queue work onto the channel
	scanner := bufio.NewScanner(f)
	jobCount := 0
	for scanner.Scan() {
		jobCount++
		jobChan <- scanner.Text()
	}

	// Receive results from workers.
	for i := 0; i < jobCount; i++ {
		v := <-resultsChan
		switch {
		case errors.Is(v.Err, errNoItem):
			// Not fatal, but no response for this entry was created.
			continue
		case v.Err != nil:
			// We had a major issue, force immediate stop.
			return fmt.Errorf("fatal error from worker: %w", v.Err)
		}

		// add the item to our list of items.
		items = append(items, *v.Item)
	}

	b, err := json.MarshalIndent(items, "", "\t")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", b)
	return nil
}

// worker is run in a goroutine to provide processing for the items.
func worker(jobChan <-chan string, results chan<- BulkResponseItemWork) {
	for entry := range jobChan {
		item, err := processBulkEntry(entry)
		results <- BulkResponseItemWork{
			Item: item,
			Err:  err,
		}
	}
}

// processBulkEntry processes an entry and returns an item containing the result or error.
func processBulkEntry(entry string) (*BulkResponseItem, error) {
	querySection, addressSection, err := parseEntry(entry)
	if err != nil {
		return nil, fmt.Errorf("parse file entry: %w", err)
	}

	item := &BulkResponseItem{
		Address: addressSection,
	}

	// If the query contains any options retrieve them and
	querySection, options, err := parseOptions(querySection)
	if err != nil {
		// These errors are non fatal, as we know which server it is for
		item.Error = err.Error()
		return item, nil
	}

	if !protocol.Supported(querySection) {
		item.Error = fmt.Sprintf("unsupported protocol: %s", querySection)
		return item, nil
	}

	client, err := svrquery.NewClient(querySection, addressSection, options...)
	if err != nil {
		item.Error = fmt.Sprintf("create client: %s", err.Error())
		return item, nil
	}

	resp, err := client.Query()
	if err != nil {
		item.Error = fmt.Sprintf("query client: %s", err.Error())
		return item, nil
	}

	item.ServerInfo = &BulkResponseServerInfoItem{
		CurrentPlayers: resp.NumClients(),
		MaxPlayers:     resp.MaxClients(),
		Map:            "UNKNOWN",
	}

	if currentMap, ok := resp.(protocol.Mapper); ok {
		item.ServerInfo.Map = currentMap.Map()
	}
	return item, nil
}

func parseEntry(entry string) (querySection, addressSection string, err error) {
	entry = strings.TrimSpace(entry)
	if entry == "" {
		return "", "", fmt.Errorf("process entry: %w", errNoItem)
	}
	sections := strings.Split(entry, " ")
	if len(sections) != 2 {
		return "", "", fmt.Errorf("%w: wrong number of sections", errEntryInvalid)
	}

	return sections[0], sections[1], nil
}

func parseOptions(querySection string) (baseQuery string, options []svrquery.Option, error error) {
	options = make([]svrquery.Option, 0)
	protocolSections := strings.Split(querySection, ",")
	for i := 1; i < len(protocolSections); i++ {
		keyVal := strings.SplitN(protocolSections[i], "=", 2)
		if len(keyVal) != 2 {
			return "", nil, fmt.Errorf("key value pair invalid: %v", keyVal)

		}
		// Only support key at the moment.
		switch strings.ToLower(keyVal[0]) {
		case "key":
			options = append(options, svrquery.WithKey(keyVal[1]))
		}
	}
	return protocolSections[0], options, nil
}
