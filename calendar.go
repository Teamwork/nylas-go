package nylas

import (
	"context"
	"net/http"
	"time"
)

type TimeZone struct {
	*time.Location
}

// Calendar represents a file in the Nylas system.
type Calendar struct {
	// Globally unique object identifier
	ID string `json:"id"`
	// A string describing the type of object (value is "calendar")
	Object string `json:"object"`
	// string	Reference to parent account object
	AccountID string `json:"account_id"`

	// Name of the Calendar
	Name string `json:"name"`
	// Description of the Calendar
	Description string `json:"description"`
	// A boolean denoting whether this is the primary calendar associated with a account
	IsPrimary bool `json:"is_primary"`
	// Job status ID for the calendar modification.
	JobStatusID string `json:"job_status_id"`
	// True if the Calendar is read only
	ReadOnly bool `json:"read_only"`
	// IANA time zone database formatted string (e.g. America/New_York).
	TimeZone *TimeZone `json:"timezone"`
}

func (tz *TimeZone) UnmarshalJSON(b []byte) (err error) {
	loc, err := time.LoadLocation(string(b))
	if err != nil {
		return err
	}
	tz.Location = loc
	return nil
}

// Returns all calendars.
// See: https://developer.nylas.com/docs/api/#get/calendars
func (c *Client) Calendars(ctx context.Context) ([]Calendar, error) {
	req, err := c.newUserRequest(ctx, http.MethodGet, "/calendars", nil)
	if err != nil {
		return []Calendar{}, err
	}

	var resp []Calendar
	return resp, c.do(req, &resp)
}

// Returns a calendar by ID.
// See: https://developer.nylas.com/docs/api/#get/calendars/id
func (c *Client) Calendar(ctx context.Context, id string) (Calendar, error) {
	req, err := c.newUserRequest(ctx, http.MethodGet, "/calendars/"+id, nil)
	if err != nil {
		return Calendar{}, err
	}

	var resp Calendar
	return resp, c.do(req, &resp)
}
