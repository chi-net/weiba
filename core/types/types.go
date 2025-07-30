package types

// for data.json

type ImportedGICAuthData struct {
	LastUpdate int            `json:"lastupdate"`
	Data       []gicGroupData `json:"data"`
}

type gicGroupData struct {
	ID   int64    `json:"i"`
	Data []string `json:"d"`
}

// for configutation

type YmlConfigurationData struct {
	AdminUID        int64                      `yaml:"admin_uid"`
	BotToken        string                     `yaml:"bot_token"`
	Mode            string                     `yaml:"mode"`
	TransformDigits int                        `yaml:"transform_digits"`
	Features        ConfigurationFeaturesData  `yaml:"features"`
	Whitelists      ConfigurationWhitelistData `yaml:"whitelists"`
}

type ConfigurationFeaturesData struct {
	GicAuth           bool `yaml:"gic_auth"`
	MonitorMembers    bool `yaml:"monitor_members"`
	UnpinChannelPosts bool `yaml:"unpin_channel_posts"`
	AnonymousChat     bool `yaml:"anonymous_chat"`
	Debug             bool `yaml:"debug"`
	Tietie            bool `yaml:"tietie"`
	Waifu             bool `yaml:"waifu"`
}

type ConfigurationWhitelistData struct {
	GicAuth           []int64 `yaml:"gic_auth"`
	MonitorMembers    []int64 `yaml:"monitor_members"`
	UnpinChannelPosts []int64 `yaml:"unpin_channel_posts"`
}

type ConfigurationWhiteListMonitorMembersData struct {
	Id   int64    `yaml:"id"`
	Tags []string `yaml:"tags"`
}

// for cache maps.

type AuthMaps struct {
	// user <-> auth channel
	Data map[int64]int64
	// user <-> 10 group ids
	GroupIds map[int64][]int64
	// user <-> 10 group message ids
	GroupMessages map[int64][]int64
	// user <-> challenge steps
	Steps map[int64]int64
	// user <-> if open chat tunnel
	ChatOpened map[int64]bool
}
