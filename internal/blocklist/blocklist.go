package blocklist

import "github.com/bwoff11/go-resolve/internal/config"

// BlockList holds a list of blocked domains with reasons and categories.
type BlockList struct {
	Blocks []Block
}

// Block holds information about a blocked domain.
type Block struct {
	Domain   string `yaml:"domain"`
	Category string `yaml:"category"`
	Reason   string `yaml:"reason"`
}

// New creates a new BlockList.
func New(cfg config.BlockListConfig) *BlockList {
	bl := &BlockList{}

	bl.loadLocalBlockLists(cfg.Local)
	bl.loadRemoteBlockLists(cfg.Remote)

	return bl
}

// AddBlock adds a new block to the list.
func (bl *BlockList) AddBlock(block Block) {
	bl.Blocks = append(bl.Blocks, block)
}

// Query finds if a domain is in the blocklist and returns the corresponding block entry.
func (bl *BlockList) Query(domain string) *Block {
	for _, block := range bl.Blocks {
		if block.Domain+"." == domain {
			return &block
		}
	}
	return nil
}