package conf

type App struct {
	DefaultChannel string `env:"DEFAULT_CHANNEL" envDefault:"movie"`

	WebAddr string `env:"WEB_ADDR" envDefault:":8080"`
	HlsAddr string `env:"HLS_ADDR" envDefault:":7002"`

	EnableIngestor bool   `env:"ENABLE_INGESTOR" envDefault:"true"`
	IngestorUpload string `env:"INGESTOR_UPLOAD" envDefault:"rtmp://localhost:11935/live"`

	IngestorWatch    []string `env:"INGESTOR_WATCH" envSeparator:"," envDefault:"http://localhost:9090/"`
	AvailableQuality []string `env:"AVAILABLE_QUALITY" envSeparator:"," envDefault:"1080p,720p,480p,360p,240p"`
	OriginalPrefix   string   `env:"ORIGINAL_PREFIX" envDefault:"live/"`
	EncodedPrefix    string   `env:"ENCODED_PREFIX" envDefault:"shakaled/"`
}
