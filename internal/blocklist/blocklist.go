package blocklist

import (
	"io"
	"net/http"
	"sync"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

// BlockList is a slice of Block, representing a list of blocked domains and their details.
type BlockList []Block

// Block represents a single blocked domain with its category and reason for being blocked.
type Block struct {
	Domain   string `yaml:"domain"`   // Domain name to be blocked.
	Category string `yaml:"category"` // Category of the reason for blocking (e.g., advertising).
	Reason   string `yaml:"reason"`   // Long reason for the domain being blocked.
}

// New initializes a BlockList from a list of URLs pointing to blocklists.
// It downloads and parses the blocklists concurrently for efficiency.
func New(blockListURLs []string) *BlockList {
	var combinedBlockList BlockList
	var mutex sync.Mutex
	var wg sync.WaitGroup

	// Iterate over the provided URLs to download and parse their blocklists concurrently.
	for _, url := range blockListURLs {
		wg.Add(1)
		go addBlocklist(url, &combinedBlockList, &mutex, &wg)
	}

	wg.Wait()
	log.Info().Int("count", len(combinedBlockList)).Msg("Block list compiled")
	return &combinedBlockList
}

// addBlocklist handles the asynchronous download and parsing of a blocklist from a URL.
// It appends the parsed blocklist to the combinedBlockList, ensuring thread-safe access.
func addBlocklist(url string, combinedBlockList *BlockList, mutex *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	rawYAML, err := downloadBlockList(url)
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("Failed to download block list")
		return
	}

	blockList, err := parseBlockList(rawYAML)
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("Failed to parse block list")
		return
	}

	mutex.Lock()
	*combinedBlockList = append(*combinedBlockList, *blockList...)
	mutex.Unlock()

	log.Info().Str("url", url).Msg("Block list downloaded and parsed")
}

// downloadBlockList downloads a block list from a given URL and returns its content as a byte slice.
func downloadBlockList(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// parseBlockList parses raw YAML content into a BlockList.
func parseBlockList(rawYAML []byte) (*BlockList, error) {
	var blockList BlockList
	err := yaml.Unmarshal(rawYAML, &blockList)
	return &blockList, err
}

// Query checks if a domain is present in the block list and returns the corresponding Block if found.
func (bl *BlockList) Query(domain string) *Block {
	for _, block := range *bl {
		if block.Domain == domain {
			return &block
		}
	}
	return nil
}

// Check checks if a domain is present in the block list and returns true if found.
func (bl *BlockList) Check(domain string) bool {
	return bl.Query(domain) != nil
}
