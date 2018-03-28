package healer

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/childe/glog"
)

type ProducerConfig struct {
	BootstrapServers         string
	ClientID                 string
	Acks                     int16
	CompressionType          string
	BatchSize                int
	MessageMaxCount          int
	FlushIntervalMS          int
	MetadataMaxAgeMS         int
	FetchTopicMetaDataRetrys int

	// TODO
	ConnectionsMaxIdleMS int
	Retries              int
	RequestTimeoutMS     int32
}

var DefaultProducerConfig *ProducerConfig = &ProducerConfig{
	ClientID:                 "healer",
	Acks:                     1,
	CompressionType:          "none",
	BatchSize:                16384,
	MessageMaxCount:          1024,
	FlushIntervalMS:          200,
	MetadataMaxAgeMS:         300000,
	FetchTopicMetaDataRetrys: 3,
	ConnectionsMaxIdleMS:     540000,

	Retries:          0,
	RequestTimeoutMS: 30000,
}

var (
	messageMaxCountError   = errors.New("message.max.count must > 0")
	flushIntervalMSError   = errors.New("flush.interval.ms must > 0")
	unknownCompressionType = errors.New("unknown compression type")
	bootstrapServersNotSet = errors.New("bootstrap servers not set")
)

func (config *ProducerConfig) checkValid() error {
	if config.BootstrapServers == "" {
		return bootstrapServersNotSet
	}
	if config.MessageMaxCount <= 0 {
		return messageMaxCountError
	}
	if config.FlushIntervalMS <= 0 {
		return flushIntervalMSError
	}

	switch config.CompressionType {
	case "none":
	case "gzip":
	case "snappy":
	case "lz4":
	default:
		return unknownCompressionType
	}
	return nil
}

type BrokerConfig struct {
	ConnectTimeoutMS    int
	TimeoutMS           int
	TimeoutMSForEachAPI []int
}

func DefaultBrokerConfig() *BrokerConfig {
	return &BrokerConfig{
		ConnectTimeoutMS:    60000,
		TimeoutMS:           30000,
		TimeoutMSForEachAPI: make([]int, 0),
	}
}

var (
	brokerAddressNotSet = errors.New("broker address not set in broker config")
)

func (c *BrokerConfig) checkValid() error {
	return nil
}

type ConsumerConfig struct {
	BootstrapServers     string `json:"bootstrap.servers"`
	ClientID             string `json:"client.id"`
	GroupID              string `json:"group.id"`
	SessionTimeoutMS     int32  `json:"session.timeout.ms"`
	FetchMaxWaitMS       int32  `json:"fetch.max.wait.ms"`
	FetchMaxBytes        int32  `json:"fetch.max.bytes"`
	FetchMinBytes        int32  `json:"fetch.min.bytes"`
	FromBeginning        bool   `json:"frombeginning"`
	AutoCommit           bool   `json:"auto.commit"`
	CommitAfterFetch     bool   `json:"commit.after.fetch"`
	AutoCommitIntervalMS int    `json:"auto.commit.interval.ms"`
	OffsetsStorage       int    `json:"offset.storage"`
	ConnectTimeoutMS     int    `json:"connect.timeout.ms"`
	TimeoutMS            int    `json:"timeout.ms"`
}

func DefaultConsumerConfig() *ConsumerConfig {
	return &ConsumerConfig{
		ClientID:             "healer",
		GroupID:              "",
		SessionTimeoutMS:     30000,
		FetchMaxWaitMS:       100,
		FetchMaxBytes:        10 * 1024 * 1024,
		FetchMinBytes:        1,
		FromBeginning:        false,
		AutoCommit:           true,
		CommitAfterFetch:     false,
		AutoCommitIntervalMS: 5000,
		OffsetsStorage:       1,
		ConnectTimeoutMS:     30000,
		TimeoutMS:            60000,
	}
}

//  fetch.min.bytes -> FetchMinBytes
//  session.timeout.ms -> SessionTimeoutMS
//  client.id -> ClientID
func convertKey(s string) string {
	parts := strings.Split(s, ".")
	switch parts[len(parts)-1] {
	case "id":
		parts[len(parts)-1] = "ID"
	case "ms":
		parts[len(parts)-1] = "MS"
	}
	convertedParts := make([]string, len(parts))
	for i, part := range parts {
		convertedParts[i] = strings.Title(part)
	}
	return strings.Join(convertedParts, "")
}

func GetConsumerConfig(config map[string]interface{}) (*ConsumerConfig, error) {
	b, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	glog.Info(string(b))

	consumerConfig := DefaultConsumerConfig()
	err = json.Unmarshal(b, consumerConfig)
	if err != nil {
		return nil, err
	}
	//configValue := reflect.ValueOf(consumerConfig).Elem()
	//for k, v := range config {
	//convertedKey := convertKey(k)
	//field := configValue.FieldByName(convertedKey)
	//if field.CanSet() {
	//field.Set(reflect.ValueOf(v))
	//} else {
	//}
	//}
	return consumerConfig, nil
}

var (
	emptyGroupID = errors.New("group.id is empty")
)

func (config *ConsumerConfig) checkValid() error {
	if config.BootstrapServers == "" {
		return bootstrapServersNotSet
	}
	if config.GroupID == "" {
		return emptyGroupID
	}
	return nil
}
