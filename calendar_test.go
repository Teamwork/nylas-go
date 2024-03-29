package nylas

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestCalendar(t *testing.T) {
	accessToken := "accessToken"
	id := "br57kcekhf1hsjq04y8aonkit"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, accessToken, "")
		assertMethodPath(t, r, http.MethodGet, "/calendars/"+id)

		_, _ = w.Write(getCalendarJSON)
	}))
	defer ts.Close()

	client := NewClient("", "", withTestServer(ts), WithAccessToken(accessToken))
	got, err := client.Calendar(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("loading timezone: %v", err)
	}

	want := Calendar{
		ID:          "8cid1lhd0m7x9k5wjrkpufs1a",
		Object:      "calendar",
		AccountID:   "43jf3n4e***",
		Name:        "name",
		Description: "description",
		IsPrimary:   true,
		JobStatusID: "status01",
		ReadOnly:    false,
		TimeZone:    &TimeZone{Location: loc},
	}

	if diff := cmp.Diff(got, want, cmp.Comparer(compareTimeZones)); diff != "" {
		t.Errorf("Calendar: (-got +want):\n%s", diff)
	}
}

func TestCalendars(t *testing.T) {
	accessToken := "accessToken"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, accessToken, "")
		assertMethodPath(t, r, http.MethodGet, "/calendars")
		assertQueryParams(t, r, url.Values{
			"offset": {"10"},
			"limit":  {"1"},
		})
		_, _ = w.Write([]byte(fmt.Sprintf("[%s]", getCalendarJSON)))
	}))
	defer ts.Close()

	client := NewClient("", "", withTestServer(ts), WithAccessToken(accessToken))
	got, err := client.Calendars(context.Background(), &CalendarsOptions{
		Limit:  1,
		Offset: 10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("loading timezone: %v", err)
	}

	want := []Calendar{{
		ID:          "8cid1lhd0m7x9k5wjrkpufs1a",
		Object:      "calendar",
		AccountID:   "43jf3n4e***",
		Name:        "name",
		Description: "description",
		IsPrimary:   true,
		JobStatusID: "status01",
		ReadOnly:    false,
		TimeZone:    &TimeZone{Location: loc},
	}}

	if diff := cmp.Diff(got, want, cmp.Comparer(compareTimeZones)); diff != "" {
		t.Errorf("Calendar: (-got +want):\n%s", diff)
	}
}

func compareTimeZones(x, y *TimeZone) bool {
	if x == nil && y == nil {
		return true
	} else if x == nil || y == nil {
		return false
	}
	return x.Location.String() == y.Location.String()
}

var getCalendarJSON = []byte(`{
    "id": "8cid1lhd0m7x9k5wjrkpufs1a",
    "account_id": "43jf3n4e***",
    "object": "calendar",
    "name": "name",
    "description": "description",
    "is_primary": true,
    "job_status_id": "status01",
    "read_only": false,
	"timezone": "America/New_York"
}`)
