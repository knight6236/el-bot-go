package eltype

type SingalType int

const (
	SingalTypeStop SingalType = iota
)

var DataRoot = "../../data"
var RssDataFileName = "rss.yml"
var ConfigRoot string = "../../config"
var SettingFullPath string = "../../plugins/MiraiAPIHTTP/setting.yml"
var FaceMapFullPath string = "../../config/face-map.yml"
var DefaultConfigFileName string = "default.yml"
var ImageFolder string = "../../plugins/MiraiAPIHTTP/images"
