// A collection of consts reused in various packages.
package aws_const

const (
	DYNAMODB                 = "DynamoDB"
	ISO8601FMT               = "2006-01-02T15:04:05Z"
	ISO8601FMT_CONDENSED     = "20060102T150405Z"
	ISODATEFMT               = "20060102"
	PORT                     = "80"
	METHOD                   = "POST"
	CTYPE                    = "application/x-amz-json-1.0"
	AMZ_TARGET_HDR           = "X-Amz-Target"
	CONTENT_MD5_HDR          = "Content-MD5"
	CONTENT_TYPE_HDR         = "Content-Type"
	DATE_HDR                 = "Date"
	CURRENT_API_VERSION      = "DynamoDB_20120810"
	ENDPOINT_PREFIX          = CURRENT_API_VERSION + "."
	X_AMZ_DATE_HDR           = "X-Amz-Date"
	X_AMZ_SECURITY_TOKEN_HDR = "X-Amz-Security-Token"
	X_AMZN_AUTHORIZATION_HDR = "X-Amzn-Authorization"
	RETRIES                  = 7
	EXCEEDED_MSG             = "ProvisionedThroughputExceededException"
	UNRECOGNIZED_CLIENT_MSG  = "UnrecognizedClientException"
	THROTTLING_MSG           = "ThrottlingException"
)
