package blocklist

import (
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

// LoadRemoteBlockLists loads blocklists from provided URLs.
func (bl *BlockList) loadRemoteBlockLists(urls []string) {
	for _, url := range urls {
		bl.loadBlockListFromURL(url)
	}
}

func (bl *BlockList) loadBlockListFromURL(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to fetch remote blocklist: %s", url)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warn().Msgf("Remote blocklist fetch failed with status code %d: %s", resp.StatusCode, url)
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to read remote blocklist data: %s", url)
		return
	}

	var blocks []Block
	if err := yaml.Unmarshal(data, &blocks); err != nil {
		log.Error().Err(err).Msgf("Failed to unmarshal YAML from remote blocklist: %s", url)
		return
	}

	for _, block := range blocks {
		bl.AddBlock(block)
	}

	log.Info().Int("count", len(blocks)).Msgf("Loaded remote blocklist from: %s", url)
}
