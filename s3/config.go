package s3

// NewConfig is used to generate a config with defaults
func NewConfig(accessKey, secretKey string) Config {
	return Config{
		AccessKey:     accessKey,
		SecretKey:     secretKey,
		Bucket:        defaultBucket,
		Region:        defaultRegion,
		Endpoint:      "http://127.0.0.1:9000",
		RootDirectory: "",
		Secure:        false,
		Workers:       defaultWorkers,
	}
}
