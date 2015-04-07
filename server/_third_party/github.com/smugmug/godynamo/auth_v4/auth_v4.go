// Manages AWS Auth v4 requests to DynamoDB.
// See http://docs.aws.amazon.com/general/latest/gr/signature-version-4.html
// for more information on v4 signed requests. For examples, see any of
// the package in the `endpoints` directory.
package auth_v4

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/auth_v4/tasks"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/aws_const"
	"github.com/FiloSottile/Heartbleed/server/_third_party/github.com/smugmug/godynamo/conf"
	"hash"
	"hash/crc32"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	IAM_WARN_MESSAGE = "check roles sources and make sure you have run one of the roles " +
		"management functions in package conf_iam, such as GoIAM"
)

// Client for executing requests.
var Client *http.Client

// Initialize package-scoped client.
func init() {
	// The timeout seems too-long, but it accomodates the exponential decay retry loop.
	// Programs using this can either change this directly or use goroutine timeouts
	// to impose a local minimum.
	tr := &http.Transport{MaxIdleConnsPerHost: 250,
		ResponseHeaderTimeout: time.Duration(20) * time.Second}
	Client = &http.Client{Transport: tr}
}

// GetRespReqID retrieves the unique identifier from the AWS Response
func GetRespReqID(response http.Response) (string, error) {
	if amz_reqid_list, reqid_ok := response.Header["X-Amzn-Requestid"]; reqid_ok {
		if len(amz_reqid_list) == 1 {
			return amz_reqid_list[0], nil
		}
	}
	return "", errors.New("auth_v4.GetRespReqID: no X-Amzn-Requestid found")
}

// MatchCheckSum will perform a local crc32 on the response body and match it against the aws crc32
// *** WARNING ***
// There seems to be a mismatch between what Go calculates and what AWS (java?) calculates here,
// I believe related to utf8 (go) vs utf16 (java), but I don't know enough about encodings to
// solve it. So until that issue is solved, don't use this.
func MatchCheckSum(response http.Response, respbody []byte) (bool, error) {
	if amz_crc_list, crc_ok := response.Header["X-Amz-Crc32"]; crc_ok {
		if len(amz_crc_list) == 1 {
			amz_crc_int32, amz_crc32_err := strconv.Atoi(amz_crc_list[0])
			if amz_crc32_err == nil {
				client_crc_int32 := int(crc32.ChecksumIEEE(respbody))
				if amz_crc_int32 != client_crc_int32 {
					_ = fmt.Sprintf("auth_v4.MatchCheckSum: resp crc mismatch: amz %d client %d",
						amz_crc_int32, client_crc_int32)
					return false, nil
				}
			}
		} else {
			return false, errors.New("auth_v4.MatchCheckSum: X-Amz-Crc32 malformed")
		}
	} else {
		return false, errors.New("auth_v4.MatchCheckSum: no X-Amz-Crc32 found")
	}
	return true, nil
}

// RawReq will sign and transmit the request to the AWS DynanoDB endpoint.
// This method is DynamoDB-specific.
// returns []byte respBody, string aws reqID, int http code, error
func RawReq(reqJSON []byte, amzTarget string) ([]byte, string, int, error) {

	// shadow conf vars in a read lock to minimize contention
	conf.Vals.ConfLock.RLock()
	conf_url_str := conf.Vals.Network.DynamoDB.URL
	conf_host := conf.Vals.Network.DynamoDB.Host
	conf_port_str := conf.Vals.Network.DynamoDB.Port
	conf_zone := conf.Vals.Network.DynamoDB.Zone
	conf_useIAM := conf.Vals.UseIAM
	conf_IAMSecret := conf.Vals.IAM.Credentials.Secret
	conf_IAMAccessKey := conf.Vals.IAM.Credentials.AccessKey
	conf_IAMToken := conf.Vals.IAM.Credentials.Token
	conf_AuthSecret := conf.Vals.Auth.Secret
	conf_AuthAccessKey := conf.Vals.Auth.AccessKey
	conf.Vals.ConfLock.RUnlock()

	// initialize req with body reader
	body := strings.NewReader(string(reqJSON))
	request, req_err := http.NewRequest(aws_const.METHOD, conf_url_str, body)
	if req_err != nil {
		e := fmt.Sprintf("auth_v4.RawReq:failed init conn %s", req_err.Error())
		return nil, "", 0, errors.New(e)
	}

	// add headers
	// content type
	request.Header.Add(aws_const.CONTENT_TYPE_HDR, aws_const.CTYPE)
	// amz target
	request.Header.Add(aws_const.AMZ_TARGET_HDR, amzTarget)
	// dates
	now := time.Now()
	request.Header.Add(aws_const.X_AMZ_DATE_HDR,
		now.UTC().Format(aws_const.ISO8601FMT_CONDENSED))

	// encode request json payload
	var h256 hash.Hash = sha256.New()
	h256.Write(reqJSON)
	hexPayload := string(hex.EncodeToString([]byte(h256.Sum(nil))))

	// create the various signed formats aws uses for v4 signed reqs
	service := strings.ToLower(aws_const.DYNAMODB)
	canonical_request := tasks.CanonicalRequest(
		conf_host,
		conf_port_str,
		request.Header.Get(aws_const.X_AMZ_DATE_HDR),
		request.Header.Get(aws_const.AMZ_TARGET_HDR),
		hexPayload)
	str2sign := tasks.String2Sign(now, canonical_request,
		conf_zone,
		service)

	// obtain the aws secret credential from the global Auth or from IAM
	var secret string
	if conf_useIAM == true {
		secret = conf_IAMSecret
	} else {
		secret = conf_AuthSecret
	}
	if secret == "" {
		panic("auth_v4.cacheable_hmacs: no Secret defined; " + IAM_WARN_MESSAGE)
	}

	signature := tasks.MakeSignature(str2sign, conf_zone, service, secret)

	// obtain the aws accessKey credential from the global Auth or from IAM
	// if using IAM, read the token while we have the lock
	var accessKey, token string
	if conf_useIAM == true {
		accessKey = conf_IAMAccessKey
		token = conf_IAMToken
	} else {
		accessKey = conf_AuthAccessKey
	}
	if accessKey == "" {
		panic("auth_v4.RawReq: no Access Key defined; " + IAM_WARN_MESSAGE)
	}

	v4auth := "AWS4-HMAC-SHA256 Credential=" + accessKey +
		"/" + now.UTC().Format(aws_const.ISODATEFMT) + "/" +
		conf_zone + "/" + service + "/aws4_request," +
		"SignedHeaders=content-type;host;x-amz-date;x-amz-target," +
		"Signature=" + signature

	request.Header.Add("Authorization", v4auth)
	if conf_useIAM == true {
		if token == "" {
			panic("auth_v4.RawReq: no Token defined;" + IAM_WARN_MESSAGE)
		}
		request.Header.Add(aws_const.X_AMZ_SECURITY_TOKEN_HDR, token)
	}

	// where we finally send req to aws
	response, rsp_err := Client.Do(request)

	if rsp_err != nil {
		return nil, "", 0, rsp_err
	}

	respbody, read_err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if read_err != nil && read_err != io.EOF {
		e := fmt.Sprintf("auth_v4.RawReq:err reading resp body: %s", read_err.Error())
		return nil, "", 0, errors.New(e)
	}

	amz_requestid, amz_requestid_err := GetRespReqID(*response)
	if amz_requestid_err != nil {
		return nil, "", 0, amz_requestid_err
	}

	return respbody, amz_requestid, response.StatusCode, nil
}

// Req is just a wrapper for RawReq if we need to massage data
// before dispatch.
func Req(reqJSON []byte, amzTarget string) ([]byte, string, int, error) {
	return RawReq(reqJSON, amzTarget)
}
