package moodle

import (
	"bkel-fetching/moodle/model"
	"encoding/json"
	"io"
	"log"
	"net/http"
	url2 "net/url"
	"strconv"
)

type Client struct {
	baseURL     url2.URL
	moodleToken string
	moodleId    string
}

func NewClient(hostname string, token string) *Client {
	baseURL := url2.URL{
		Host:   hostname,
		Scheme: "http",
		Path:   "webservice/rest/server.php",
	}

	client := &Client{
		baseURL:     baseURL,
		moodleToken: token,
	}

	client.moodleId = client.getMoodleID()

	log.Printf("Create MoodleClient for hostname: %v", hostname)

	return client
}

func (client *Client) sendRequest(params url2.Values, returnResponse bool) (response []byte) {
	url := client.baseURL
	url.RawQuery = params.Encode()

	resp, err := http.Get(url.String())
	if err != nil {
		log.Panicf("Error sending request, err: %v", err)
	}

	if resp.StatusCode != 200 {
		log.Panicf("Response status not 200: %v", resp.StatusCode)
	}

	if !returnResponse {
		return nil
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panicf("Error reading response, err: %v", err)
	}

	if err = resp.Body.Close(); err != nil {
		log.Printf("Failed to close GetMoodleID Body, err: %v", err)
	}

	return respBody
}

func (client *Client) getMoodleID() string {
	params := url2.Values{}
	params.Add("moodlewsrestformat", "json")
	params.Add("wstoken", client.moodleToken)
	params.Add("wsfunction", "core_webservice_get_site_info")

	resp := client.sendRequest(params, true)

	var respBody model.UserInfo
	if err := json.Unmarshal(resp, &respBody); err != nil {
		log.Panicf("Can't parsed response, err: %v", err)
	}

	return strconv.Itoa(respBody.UserId)
}

func (client *Client) FetchMessages() []model.Message {
	fetchMsgParams := url2.Values{}
	fetchMsgParams.Add("moodlewsrestformat", "json")
	fetchMsgParams.Add("wstoken", client.moodleToken)
	fetchMsgParams.Add("useridto", client.moodleId)
	fetchMsgParams.Add("wsfunction", "core_message_get_messages")
	fetchMsgParams.Add("type", "both")
	fetchMsgParams.Add("limitnum", "50")
	fetchMsgParams.Add("read", "0")

	resp := client.sendRequest(fetchMsgParams, true)

	var respBody model.MsgResponse
	if err := json.Unmarshal(resp, &respBody); err != nil {
		log.Panicf("Can't parsed response, err: %v", err)
	}

	return respBody.Messages
}

func (client *Client) MarkChatAsRead(msgId int) {
	params := url2.Values{}
	params.Add("moodlewsrestformat", "json")
	params.Add("wstoken", client.moodleToken)
	params.Add("wsfunction", "core_message_mark_message_read")
	params.Add("messageid", strconv.Itoa(msgId))

	_ = client.sendRequest(params, false)
}

func (client *Client) MarkNotificationAsRead(notificationId int) {
	params := url2.Values{}
	params.Add("moodlewsrestformat", "json")
	params.Add("wstoken", client.moodleToken)
	params.Add("wsfunction", "core_message_mark_notification_read")
	params.Add("notificationid", strconv.Itoa(notificationId))

	_ = client.sendRequest(params, false)
}

func (client *Client) MarkAllNotificationsAsRead() {
	// Deprecated: MarkAllNotificationsAsRead is depricated
	params := url2.Values{}
	params.Add("moodlewsrestformat", "json")
	params.Add("wstoken", client.moodleToken)
	params.Add("wsfunction", "core_message_mark_all_notifications_as_read")
	params.Add("useridto", "0")

	_ = client.sendRequest(params, false)
}

//func (client *Client) GetUpcomingEvents() []model.CalendarEvent {
//	params := url2.Values{}
//	params.Add("moodlewsrestformat", "json")
//	params.Add("wstoken", client.moodleToken)
//	params.Add("wsfunction", "core_calendar_get_calendar_upcoming_view")
//
//	resp := client.sendRequest(params, true)
//	var upcomingEvents model.UpcomingEventResponse
//
//	if err := json.Unmarshal(resp, &upcomingEvents); err != nil {
//		log.Panicf("Cannot get upcoming events, %v", err)
//	}
//
//	return upcomingEvents.Events
//}
