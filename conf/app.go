package conf

type App struct {
	DefaultChannel string `env:"DEFAULT_CHANNEL" envDefault:"movie"`

	WebAddr string `env:"WEB_ADDR" envDefault:":8080"`
	FlvAddr string `env:"FLV_ADDR" envDefault:":7001"`
	HlsAddr string `env:"HLS_ADDR" envDefault:":7002"`
}
