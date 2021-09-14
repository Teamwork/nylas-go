package nylas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestEvents(t *testing.T) {
	accessToken := "accessToken"
	wantQuery := url.Values{
		"limit":  {"3"},
		"offset": {"1"},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, accessToken, "")
		assertMethodPath(t, r, http.MethodGet, "/events")
		assertQueryParams(t, r, wantQuery)
		_, _ = w.Write(eventsJSON)
	}))
	defer ts.Close()

	client := NewClient("", "", withTestServer(ts), WithAccessToken(accessToken))
	got, err := client.Events(context.Background(), &EventsOptions{
		Offset: 1,
		Limit:  3,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("loading timezone: %v", err)
	}

	meta := json.RawMessage(`{
		"hello": "goodbye"
	}`)

	want := []Event{
		{
			ID:          "{event_id}",
			Object:      "event",
			AccountID:   "{account_id}",
			Busy:        true,
			CalendarID:  "{calendar_id}",
			Description: "Coffee meeting",
			ICalUID:     "{ical_uid}",
			Location:    "string",
			Owner:       "<some_email@email.com>",
			Participants: []EventParticipant{{
				Name:    "Dorothy Vaughan",
				Email:   "dorothy@spacetech.com",
				Status:  "noreply",
				Comment: "string",
			}},
			ReadOnly: true,
			Title:    "Remote Event: Group Yoga Class",
			When: &EventTimespan{
				StartTime:     time.Unix(1409594400, 0),
				EndTime:       time.Unix(1409598000, 0),
				StartTimezone: &TimeZone{Location: loc},
				EndTimezone:   &TimeZone{Location: loc},
			},
			Status: "confirmed",
			Recurrence: EventRecurrence{
				RRule:    []string{"RRULE:FREQ=WEEKLY;BYDAY=MO"},
				Timezone: &TimeZone{Location: loc},
			},
			Metadata: meta,
		},
	}

	if diff := cmp.Diff(got, want, cmp.Comparer(compareTimeZones)); diff != "" {
		t.Errorf("Events: (-got +want):\n%s", diff)
	}
}

var eventsJSON = []byte(fmt.Sprintf("[%s]", eventJSON))

func TestEvent(t *testing.T) {
	accessToken := "accessToken"

	id := "{event_id}"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, accessToken, "")
		assertMethodPath(t, r, http.MethodGet, "/events/"+id)
		_, _ = w.Write(eventJSON)
	}))
	defer ts.Close()

	client := NewClient("", "", withTestServer(ts), WithAccessToken(accessToken))
	got, err := client.Event(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("loading timezone: %v", err)
	}

	meta := json.RawMessage(`{
		"hello": "goodbye"
	}`)

	want := Event{
		ID:          "{event_id}",
		Object:      "event",
		AccountID:   "{account_id}",
		Busy:        true,
		CalendarID:  "{calendar_id}",
		Description: "Coffee meeting",
		ICalUID:     "{ical_uid}",
		Location:    "string",
		Owner:       "<some_email@email.com>",
		Participants: []EventParticipant{{
			Name:    "Dorothy Vaughan",
			Email:   "dorothy@spacetech.com",
			Status:  "noreply",
			Comment: "string",
		}},
		ReadOnly: true,
		Title:    "Remote Event: Group Yoga Class",
		When: &EventTimespan{
			StartTime:     time.Unix(1409594400, 0),
			EndTime:       time.Unix(1409598000, 0),
			StartTimezone: &TimeZone{Location: loc},
			EndTimezone:   &TimeZone{Location: loc},
		},
		Status: "confirmed",
		Recurrence: EventRecurrence{
			RRule:    []string{"RRULE:FREQ=WEEKLY;BYDAY=MO"},
			Timezone: &TimeZone{Location: loc},
		},
		Metadata: meta,
	}

	if diff := cmp.Diff(got, want, cmp.Comparer(compareTimeZones)); diff != "" {
		t.Errorf("Events: (-got +want):\n%s", diff)
	}
}

var eventJSON = []byte(`{
	"account_id": "{account_id}",
	"busy": true,
	"calendar_id": "{calendar_id}",
	"description": "Coffee meeting",
	"ical_uid": "{ical_uid}",
	"id": "{event_id}",
	"location": "string",
	"object": "event",
	"owner": "<some_email@email.com>",
	"participants": [
		{
			"name": "Dorothy Vaughan",
			"email": "dorothy@spacetech.com",
			"status": "noreply",
			"comment": "string"
		}
	],
	"read_only": true,
	"title": "Remote Event: Group Yoga Class",
	"when": {
		"start_time": 1409594400,
		"end_time": 1409598000,
		"start_timezone": "America/New_York",
		"end_timezone": "America/New_York"
	},
	"status": "confirmed",
	"conferencing": {
		"provider": "WebEx",
		"details": {
			"password": "string",
			"pin": "string",
			"url": "string"
		}
	},
	"recurrence": {
		"rrule": [
			"RRULE:FREQ=WEEKLY;BYDAY=MO"
		],
		"timezone": "America/New_York"
	},
	"metadata": {
		"hello": "goodbye"
	}
}`)
