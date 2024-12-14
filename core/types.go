package core

type ImportedGICAuthData struct {
	LastUpdate int            `json:"lastupdate"`
	Data       []gicGroupData `json:"data"`
}

type gicGroupData struct {
	ID   int64    `json:"i"`
	Data []string `json:"d"`
}

type YmlConfigurationData struct {
	AdminUID                      int64   `yaml:"admin_uid"`
	BotToken                      string  `yaml:"bot_token"`
	Mode                          string  `yaml:"mode"`
	TransformDigits               int     `yaml:"transform_digits"`
	EnhancedMonitorChannelMembers bool    `yaml:"enhanced_monitor_channel_members"`
	ChatWhitelist                 []int64 `yaml:"chat_whitelist"`
	UnpinChannelPosts             bool    `yaml:"unpin_channel_posts"`
	UnpinChannelPostsGroups       []int64 `yaml:"unpin_channel_posts_groups"`
}

type AuthMaps struct {
	// user <-> auth channel
	Data map[int64]int64
	// user <-> 10 group ids
	GroupIds map[int64][]int64
	// user <-> 10 group message ids
	GroupMessages map[int64][]int64
	// user <-> challenge steps
	Steps map[int64]int64
}
