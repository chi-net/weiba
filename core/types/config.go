package types

// for configutation

type YmlConfigurationData struct {
	AdminUID        int64                      `yaml:"admin_uid"`
	BotToken        string                     `yaml:"bot_token"`
	Mode            string                     `yaml:"mode"`
	TransformDigits int                        `yaml:"transform_digits"`
	Features        ConfigurationFeaturesData  `yaml:"features"`
	Whitelists      ConfigurationWhitelistData `yaml:"whitelists"`
	AI              ConfigurationAIData        `yaml:"ai"`
	Ranking         ConfigurationRankingData   `yaml:"ranking"`
}

type ConfigurationFeaturesData struct {
	GicAuth           bool `yaml:"gic_auth"`
	MonitorMembers    bool `yaml:"monitor_members"`
	UnpinChannelPosts bool `yaml:"unpin_channel_posts"`
	AnonymousChat     bool `yaml:"anonymous_chat"`
	Debug             bool `yaml:"debug"`
	Tietie            bool `yaml:"tietie"`
	Waifu             bool `yaml:"waifu"`
	Ranking           bool `yaml:"ranking"`
	SafetyMonitor     bool `yaml:"safety_monitor"`
}

type ConfigurationWhitelistData struct {
	GicAuth           []int64 `yaml:"gic_auth"`
	MonitorMembers    []int64 `yaml:"monitor_members"`
	UnpinChannelPosts []int64 `yaml:"unpin_channel_posts"`
	Ranking           []int64 `yaml:"ranking"`
	Tietie            []int64 `yaml:"tietie"`
	MonitorStatus     []int64 `yaml:"monitor_status"`
}

type ConfigurationWhiteListMonitorMembersData struct {
	Id   int64    `yaml:"id"`
	Tags []string `yaml:"tags"`
}

type ConfigurationAIData struct {
	GeminiAPIKey string `yaml:"gemini_api_key"`
}

type ConfigurationRankingData struct {
	Global   bool                               `yaml:"global"`
	Features ConfigurationRankingFeaturesData   `yaml:"features"`
	Category ConfigurationRankingCategoriesData `yaml:"categories"`
}

type ConfigurationRankingFeaturesData struct {
	Cai bool `yaml:"cai"`
	Xm  bool `yaml:"xm"`
}

type ConfigurationRankingCategoriesData struct {
	Name     string   `yaml:"name"`
	Triggers []string `yaml:"triggers"`
}
