package consumergroup

import (
	"errors"
	"time"

	"github.com/Shopify/sarama"
	"github.com/wvanbergen/kazoo-go"
)

type Config struct {
	*sarama.Config

	Zookeeper *kazoo.Config

	Offsets struct {
		Initial           int64         // The initial offset method to use if the consumer has no previously stored offset. Must be either sarama.OffsetOldest (default) or sarama.OffsetNewest.
		ProcessingTimeout time.Duration // Time to wait for all the offsets for a partition to be processed after stopping to consume from it. Defaults to 1 minute.
		CommitInterval    time.Duration // The interval between which the processed offsets are commited.
		ResetOffsets      bool          // Resets the offsets for the consumergroup so that it won't resume from where it left off previously.
	}
}

func NewConfig() *Config {
	config := &Config{}
	config.Config = sarama.NewConfig()
	config.Zookeeper = kazoo.NewConfig()
	config.Offsets.Initial = sarama.OffsetOldest
	config.Offsets.ProcessingTimeout = 30 * time.Second
	config.Offsets.CommitInterval = time.Minute

	return config
}

func (cgc *Config) Validate() error {
	if cgc.Zookeeper.Timeout <= 0 {
		return sarama.ConfigurationError("ZookeeperTimeout should have a duration > 0")
	}

	if cgc.Offsets.CommitInterval <= 0 {
		return sarama.ConfigurationError("CommitInterval should have a duration > 0")
	}

	if cgc.Offsets.Initial != sarama.OffsetOldest && cgc.Offsets.Initial != sarama.OffsetNewest {
		return errors.New("Offsets.Initial should be sarama.OffsetOldest or sarama.OffsetNewest.")
	}

	if cgc.Config != nil {
		if err := cgc.Config.Validate(); err != nil {
			return err
		}
	}

	return nil
}
