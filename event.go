package nylas

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/go-querystring/query"
)

// EventParticipant represents an event participant in the Nylas system.
type EventParticipant struct {
	// (Optional) The participant's full name.
	Name string `json:"name"`
	// The participant's email address.
	Email string `json:"email"`
	// The participant's attendance status. Allowed values are yes, maybe, no and noreply.
	// The default value is noreply.
	Status string `json:"status"`
	// 	(Optional) A comment by the participant.
	Comment string `json:"comment"`
}

// EventTimeSubobject represents an event time subobject.
type EventTimeSubobject interface {
	isEventTimeSubobject()
}

// EventTime subobject corresponds a single moment in time, which has no duration.
// Reminders or alarms are represented as time subobjects.
type EventTime struct {
	// Event time in UTC.
	Time time.Time `json:"time"`
	// If timezone is present, then the value for time will be read with timezone. Timezone using IANA formatted string.
	Timezone *TimeZone `json:"timezone"`
}

func (t *EventTime) isEventTimeSubobject() {}

// EventTimespan represents a span of time with a specific beginning and end time.
// An hour lunch meeting would be represented as timespan subobjects.
type EventTimespan struct {
	// The start time of the event.
	StartTime time.Time `json:"start_time"`
	// The end time of the event.
	EndTime time.Time `json:"end_time"`
	// start_timezone and end_timezone must be submitted together. Timezone using IANA formatted string.
	StartTimezone *TimeZone `json:"start_timezone"`
	// start_timezone and end_timezone must be submitted together. Timezone using IANA formatted string.
	EndTimezone *TimeZone `json:"end_timezone"`
}

func (t *EventTimespan) isEventTimeSubobject() {}

// EventDate represents a specific date for an event, without a clock-based start or end time.
// Your birthday and holidays would be represented as date subobjects.
type EventDate struct {
	Date time.Time `json:"date"`
}

func (t *EventDate) isEventTimeSubobject() {}

// EventDatespan a span of entire days without specific times.
// A business quarter or academic semester would be represented as datespan subobjects.
type EventDatespan struct {
	// The start date of the event.
	StartDate time.Time `json:"start_date"`
	// The end date of the event.
	EndDate time.Time `json:"end_date"`
}

func (t *EventDatespan) isEventTimeSubobject() {}

// EventRecurrence represents an event recurrence in the Nylas system.
type EventRecurrence struct {
	// An array of RRULE and EXDATE strings. See RFC-5545 for more details.
	// Please note that EXRULE and RDATE strings are not supported for POST or PUT requests at this time.
	// This tool is helpful in understanding the RRULE spec.
	RRule []string `json:"rrule"`
	// The participant's email address.
	Timezone *TimeZone `json:"timezone"`
}

// Event represents an event in the Nylas system.
type Event struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	AccountID string `json:"account_id"`

	// A reference to the parent calendar object.
	CalendarID string `json:"calendar_id"`
	// The title of the event, usually short (maximum string length of 1024 characters).
	Title string `json:"title"`
	// The description of the event, which may contain more details or an agenda
	// (maximum string length of 8192 characters).
	Description string `json:"description"`
	// Unique identifier as defined in RFC5545. It is used to uniquely identify events across
	// calendaring systems. Can be null.
	ICalUID string             `json:"ical_uid"`
	When    EventTimeSubobject `json:"-"`
	// A location, such as a physical address or meeting room name.
	Location string `json:"location"`
	// The owner of the event, usually specified with their email or name and email.
	Owner string `json:"owner"`
	// An array of other participants invited to the event. Keys are email, name, status.
	// Participants may also be rooms or resources.
	Participants []EventParticipant `json:"participants"`
	// One of the following values: confirmed, tentative, or cancelled.
	Status string `json:"status"`
	// Indicates whether the event can be modified.
	ReadOnly bool `json:"read_only"`
	// On shared or public calendars, indicates whether to show this event's time block as available.
	// (Also called transparency in some systems.)
	Busy bool `json:"busy"`
	// Included if the event is a master recurring event.
	Recurrence EventRecurrence `json:"recurrence"`
	// Only included in exceptions (overrides) to recurring events, the ID of the recurring event.
	MasterEventID string `json:"master_event_id"`
	// Only included in exceptions (overrides) to recurring events, the start time of the recurring event.
	OriginalStartTime time.Time `json:"original_start_time"`
	// A key-value pair added to an event object to store data
	Metadata json.RawMessage `json:"metadata"`
}

// UnmarshalJSON defines an Event unmarshaller that infers the `when` subobject.
func (e *Event) UnmarshalJSON(data []byte) error {
	type EventAlias Event
	ea := &struct {
		*EventAlias
		When map[string]interface{} `json:"when"`
	}{
		EventAlias: (*EventAlias)(e),
	}

	if err := json.Unmarshal(data, &ea); err != nil {
		return err
	}

	switch w := ea.When; {
	case w["time"] != nil:
		t, ok := w["time"].(float64)
		if !ok {
			return errors.New("invalid time for event time")
		}

		loc, err := time.LoadLocation(w["timezone"].(string))
		if err != nil {
			return errors.New("invalid timezone for event time")
		}

		ea.EventAlias.When = &EventTime{
			Time:     time.Unix(int64(t), 0),
			Timezone: &TimeZone{Location: loc},
		}
	case w["start_time"] != nil:
		st, ok := w["start_time"].(float64)
		if !ok {
			return errors.New("invalid start time for event timespan")
		}
		et, ok := w["end_time"].(float64)
		if !ok {
			return errors.New("invalid end time for event timespan")
		}

		sLoc, err := time.LoadLocation(w["start_timezone"].(string))
		if err != nil {
			return errors.New("invalid start timezone for event timespan")
		}

		eLoc, err := time.LoadLocation(w["end_timezone"].(string))
		if err != nil {
			return errors.New("invalid end timezone for event timespan")
		}

		ea.EventAlias.When = &EventTimespan{
			StartTime:     time.Unix(int64(st), 0),
			EndTime:       time.Unix(int64(et), 0),
			StartTimezone: &TimeZone{Location: sLoc},
			EndTimezone:   &TimeZone{Location: eLoc},
		}
	case w["date"] != nil:
		t, err := time.Parse("2006-01-02", w["date"].(string))
		if err != nil {
			return errors.New("invalid date for event date")
		}
		ea.EventAlias.When = &EventDate{
			Date: t,
		}
	case w["start_date"] != nil:
		st, err := time.Parse("2006-01-02", w["start_date"].(string))
		if err != nil {
			return errors.New("invalid start date for event date")
		}

		et, err := time.Parse("2006-01-02", w["end_date"].(string))
		if err != nil {
			return errors.New("invalid end date for event date")
		}

		ea.EventAlias.When = &EventDatespan{
			StartDate: st,
			EndDate:   et,
		}
	}

	return nil
}

// EventsOptions represents request options.
type EventsOptions struct {
	ShowCancelled   string `url:"show_cancelled,omitempty"`
	Limit           int    `url:"limit,omitempty"`
	Offset          int    `url:"offset,omitempty"`
	EventID         string `url:"event_id,omitempty"`
	CalendarID      string `url:"calendar_id,omitempty"`
	Title           string `url:"title,omitempty"`
	Description     string `url:"description,omitempty"`
	Location        string `url:"location,omitempty"`
	StartsBefore    string `url:"starts_before,omitempty"`
	StartsAfter     string `url:"starts_after,omitempty"`
	EndsBefore      string `url:"ends_before,omitempty"`
	EndsAfter       string `url:"ends_after,omitempty"`
	MetadataKey     string `url:"metadata_key,omitempty"`
	MetadataValue   string `url:"metadata_value,omitempty"`
	ExpandRecurring bool   `url:"expand_recurring,omitempty"`
	MetadataPair    string `url:"metadata_pair,omitempty"`
	Busy            bool   `url:"busy,omitempty"`
}

// Events returns all events.
// See: https://developer.nylas.com/docs/api/#get/events
func (c *Client) Events(ctx context.Context, opts *EventsOptions) ([]Event, error) {
	req, err := c.newUserRequest(ctx, http.MethodGet, "/events", nil)
	if err != nil {
		return nil, err
	}

	if opts != nil {
		vs, err := query.Values(opts)
		if err != nil {
			return nil, err
		}
		appendQueryValues(req, vs)
	}

	var resp []Event
	return resp, c.do(req, &resp)
}

// Event returns an event by id.
// See: https://developer.nylas.com/docs/api/#get/events/id
func (c *Client) Event(ctx context.Context, id string) (Event, error) {
	req, err := c.newUserRequest(ctx, http.MethodGet, "/events/"+id, nil)
	if err != nil {
		return Event{}, err
	}

	var resp Event
	return resp, c.do(req, &resp)
}
