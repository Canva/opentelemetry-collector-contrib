package reporternamedecoderprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/redactionprocessor"

import (
	"go.opentelemetry.io/collector/config"
)

type Config struct {
	config.ProcessorSettings `mapstructure:",squash"`

	// URL to S3 bucket containing reporter name JSON.
	// The S3 bucket should follow the schema of
	// `s3:///reporter-names.canva.com/[RELEASE_NAME]/reporter_names.json`
	// where `reporter_names.json` is an object with encoded names as keys and
	// decoded names as values.
	S3BucketUrl string `mapstructure:"s3_bucket_url"`
}
