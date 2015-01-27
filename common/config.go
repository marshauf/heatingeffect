package common

import (
	"encoding/json"
	"os"
	"time"
)

// MongoDBConfig contains dialing information and collection names.
type MongoDBConfig struct {
	// Addrs holds the addresses for the seed servers.
	Addrs []string `json:"addrs"`

	// Timeout is the amount of time to wait for a server to respond when
	// first connecting and on follow up operations in the session. If
	// timeout is zero, the call may block forever waiting for a connection
	// to be established.
	Timeout time.Duration `json:"timeout"`

	// Database is the default database name used when the Session.DB method
	// is called with an empty name, and is also used during the intial
	// authenticatoin if Source is unset.
	Database string `json:"database"`

	// Username and Password inform the credentials for the initial authentication
	// done on the database defined by the Source field.
	Username string `json:"username"`
	Password string `json:"password"`

	// collection name for chillingeffects.Notices
	NoticesCollectionName string `json:"notices_collection_name"`

	NoticesSendToStatCollectionName string `json:"notices_sendto_stat_collection_name"`
}

// Config represents an decoded json config file.
type Config struct {
	MongoDB           *MongoDBConfig `json:"mongodb"`
	RunMode           string         `json:"runmode"`
	IDRange           *IDRange       `json:"id_range"`
	RequestsPerWorker int            `json:"requests_per_worker"`
	IronIO            *IronIO        `json:"iron_io"`
}

// IDRange holds range information for Notice IDs.
type IDRange struct {
	Low  int `json:"low"`
	High int `json:"high"`
}

// IronIO holds config values specific for iron.io worker.
type IronIO struct {
	ScheduleTasksLimit int `json:"schedule_tasks_limit"`

	// Schedule
	CodeName string `json:"code_name"`
	Name     string `json:"name"`
	Label    string `json:"label"`
	Cluster  string `json:"cluster"`
}

// NewConfig returns a decoded *Config from data.
func NewConfig(data string) (*Config, error) {
	config := &Config{}
	err := json.Unmarshal([]byte(data), config)
	return config, err
}

// LoadConfig loads the file fileName and returns a decoded *Config.
func LoadConfig(fileName string) (*Config, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(file)
	config := &Config{}
	err = decoder.Decode(config)
	return config, err
}
