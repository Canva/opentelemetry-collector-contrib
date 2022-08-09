# Reporter Name Processor

This is a custom processor produced at Canva that decodes a two-character reporter name encoding sent from the frontend into its actual name. For example, this processor could receive from the frontend a span with `ReporterName = "-a"` and process it into `ReporterName = "home_page_main"`.

For more information, see the associated [Design Doc](https://docs.google.com/document/d/1GM0FB5PGM4tn9R_JCbJdNxXaV1FEq4HEPyPRjJNRWNY/edit).

## Processor Configuration

Please refer to [config.go](./config.go) for the config spec.

Examples:

```yaml
processors:
  reportername:
    s3_bucket_url: "URL to S3 bucket containing the decoding map for reporter names in JSON for each Canva release"
```

## S3 Bucket

The S3 bucket should follow the schema of `s3:///reporter-names.canva.com/[RELEASE_NAME]/reporter_names.json` where `reporter_names.json` is an object with encoded names as keys and decoded names as values. For example, `reporter_names.json` ought to look like:

```json
{
  "r0": "telemetry_test",
  "-a": "home_page_main"
}
```
