package ce

// Job represents an scheduled job on the platform
type Job struct {
	ID                 string     `json:"id"`
	DisallowConcurrent bool       `json:"disallowConcurrent"`
	Data               JobData    `json:"data"`
	Name               string     `json:"name"`
	Description        string     `json:"description"`
	Trigger            JobTrigger `json:"trigger"`
}

// JobData represents the data of the scheduled job
type JobData struct {
	ID            int         `json:"id"`
	ElementKey    string      `json:"elementKey"`
	Topic         string      `json:"topic"`
	Notifications interface{} `json:"notifications"`
}

// JobTrigger is the trigger that kicks off the job
type JobTrigger struct {
	ID           string `json:"ID"`
	CalendarName string `json:"calendarName"`
	MayFireAgain bool   `json:"mayFireAgain"`
	NextFireTime int    `json:"nextFireTime"`
	Description  string `json:"Description"`
	StartTime    int    `json:"startTime"`
	EndTime      int    `json:"endTime"`
	Priority     int    `json:"priority"`
	State        string `json:"state"`
}
