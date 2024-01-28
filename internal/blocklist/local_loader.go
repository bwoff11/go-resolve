package blocklist

import (
	"os"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

// LoadLocalBlockLists loads blocklists from provided file paths.
func (bl *BlockList) loadLocalBlockLists(paths []string) {
	for _, path := range paths {
		bl.loadBlockListFromFile(path)
	}
}

func (bl *BlockList) loadBlockListFromFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to read local blocklist file: %s", path)
		return
	}

	var blocks []Block
	if err := yaml.Unmarshal(data, &blocks); err != nil {
		log.Error().Err(err).Msgf("Failed to unmarshal YAML from local blocklist file: %s", path)
		return
	}

	for _, block := range blocks {
		bl.AddBlock(block)
	}

	log.Info().Int("count", len(blocks)).Msgf("Loaded local blocklist from: %s", path)
}
