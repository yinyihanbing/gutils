package qcloud_cdn

import (
	"testing"
	"fmt"
)

func TestAggregate(t *testing.T) {
	/**get SecretKey & SecretId from https://console.qcloud.com/capi**/
	var Requesturl string = "cdn.api.qcloud.com/v2/index.php"
	var SecretKey string = "kIelgPwWvAD3hARVPAj6e1CD2e4IJ8kA"
	var Method string = "POST"

	/**params to signature**/
	params := make(map[string]interface{})
	params["SecretId"] = "AKIDdbNxLX2LxNabP1LdeTq1zlLjuseGsssf"
	params["Action"] = "RefreshCdnUrl"
	params["urls.0"] = "http://pkgo.blbl666.com/index.html"

	/*use qcloudcdn_api.Signature to obtain signature and params with correct signature**/
	signature, request_params := Signature(SecretKey, params, Method, Requesturl)
	fmt.Println("signature : ", signature)

	/*use qcloudcdn_api.SendRequest to send request**/
	response, _ := SendRequest(Requesturl, request_params, Method)
	fmt.Println(response.Code)
}
