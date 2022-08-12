package moodle

import (
	"bkel-fetching/model"
	"encoding/json"
	"io"
	"log"
	"net/http"
	url2 "net/url"
	"strconv"
)

var baseURL = url2.URL{
	Host:   "e-learning.hcmut.edu.vn",
	Scheme: "http",
	Path:   "webservice/rest/server.php",
}

func cloneURLWithParams(url url2.URL, params url2.Values) url2.URL {
	return url2.URL{
		Host:     url.Host,
		Scheme:   url.Scheme,
		Path:     url.Path,
		RawQuery: params.Encode(),
	}
}

// Send a GET request using the url (i.e., the school Moodle URL).
// Return response as []byte
func sendRequest(url url2.URL, returnResponse bool) (response []byte) {
	resp, err := http.Get(url.String())
	if err != nil {
		log.Panicf("Error sending request, err: %v", err)
	}

	if resp.StatusCode != 200 {
		log.Panic("Request isn't successful.")
	}

	if !returnResponse {
		return nil
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panicf("Error reading response, err: %v", err)
	}

	if err = resp.Body.Close(); err != nil {
		log.Printf("Can't close GetMoodleID Body, err: %v", err)
	}

	return respBody
}

func GetMoodleID(moodleToken string) string {
	params := url2.Values{}
	params.Add("moodlewsrestformat", "json")
	params.Add("wstoken", moodleToken)
	params.Add("wsfunction", "core_webservice_get_site_info")

	url := cloneURLWithParams(baseURL, params)
	resp := sendRequest(url, true)

	var respBody model.UserInfo
	if err := json.Unmarshal(resp, &respBody); err != nil {
		log.Panicf("Can't parsed response, err: %v", err)
	}

	return strconv.Itoa(respBody.UserId)
}

func FetchMessages(moodleToken string, moodleId string) []model.Message {
	fetchMsgParams := url2.Values{}
	fetchMsgParams.Add("moodlewsrestformat", "json")
	fetchMsgParams.Add("wstoken", moodleToken)
	fetchMsgParams.Add("useridto", moodleId)
	fetchMsgParams.Add("wsfunction", "core_message_get_messages")
	fetchMsgParams.Add("type", "both")
	fetchMsgParams.Add("limitnum", "50")
	fetchMsgParams.Add("read", "0")

	url := cloneURLWithParams(baseURL, fetchMsgParams)
	resp := sendRequest(url, true)

	var respBody model.MsgResponse
	if err := json.Unmarshal(resp, &respBody); err != nil {
		log.Panicf("Can't parsed response, err: %v", err)
	}

	return respBody.Messages
}

func MarkChatAsRead(moodleToken string, msgId int) {
	params := url2.Values{}
	params.Add("moodlewsrestformat", "json")
	params.Add("wstoken", moodleToken)
	params.Add("wsfunction", "core_message_mark_message_read")
	params.Add("messageid", strconv.Itoa(msgId))

	url := cloneURLWithParams(baseURL, params)
	_ = sendRequest(url, false)
}

func MarkAllNotificationsAsRead(moodleToken string) {
	params := url2.Values{}
	params.Add("moodlewsrestformat", "json")
	params.Add("wstoken", moodleToken)
	params.Add("wsfunction", "core_message_mark_all_notifications_as_read")
	params.Add("useridto", "0")

	url := cloneURLWithParams(baseURL, params)
	_ = sendRequest(url, false)
}
