package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

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

// BulkResponseServerInfoItem containing basic server information.
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
	// To simplify the workerpool load all the entries we are going to work on
	lines := fileLines(file)

	if len(lines) > 10000 {
		return fmt.Errorf("too many servers requested %d (max 10000)", len(lines))
	}

	// Make a jobs channel and a number of workers to processes
	// work off of the channel.
	jobChan := make(chan string, len(lines))
	resultsChan := make(chan BulkResponseItemWork)
	wg := sync.WaitGroup{}
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(jobChan, resultsChan, &wg)
	}

	items := make([]BulkResponseItem, 0, len(lines))

	// Queue work onto the channel
	for _, line := range lines {
		jobChan <- line
	}

	// Receive results from workers.
	var err error
	for i := 0; i < len(lines); i++ {
		v := <-resultsChan
		switch {
		case errors.Is(v.Err, errNoItem):
			// Not fatal, but no response for this entry was created.
			continue
		case v.Err != nil:
			// We had a major issue processing the list
			if err == nil {
				err = fmt.Errorf("fatal error: %w", v.Err)
				continue
			}
			err = fmt.Errorf("additional error: %w", v.Err)
			continue
		}

		// add the item to our list of items.
		items = append(items, *v.Item)
	}
	close(jobChan)
	wg.Wait()

	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(items, "", "\t")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", b)
	return nil
}

func fileLines(file string) []string {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	result := make([]string, 0, 1000)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		result = append(result, line)
	}
	return result
}

// worker is run in a goroutine to provide processing for the items.
func worker(jobChan <-chan string, results chan<- BulkResponseItemWork, wg *sync.WaitGroup) {
	defer wg.Done()
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
