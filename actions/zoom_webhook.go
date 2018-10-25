package actions

import (
	"errors"
	"fmt"

	"github.com/gobuffalo/buffalo"
)

const recordingCompleted = "RECORDING_MEETING_COMPLETED"

type wh struct {
	MeetingID   string `json:"id"`
	MeetingUUID string `json:"uuid"`
	Status      string `json:"status"`
	HostID      string `json:"host_id"`
}

func zoomWebhook(c buffalo.Context) error {
	logger := c.Logger()
	// TODO: zoom -> webhook auth
	hook := new(wh)
	if err := c.Bind(hook); err != nil {
		logger.Errorf("webhook body was wrong: %s", err)
		return c.Error(400, errors.New("webhook body was wrong"))
	}
	if hook.Status != recordingCompleted {
		logger.Errorf("hook status was not %s", recordingCompleted)
		// TODO: is 400 the right thing to send back to zoom?
		return c.Error(400, fmt.Errorf("hook status was not %s", recordingCompleted))
	}

	// TODO: use the zoom client that I built in ./pkg/zoom
	// something to get all the meeting's recordings
	// see https://zoom.github.io/api/#retrieve-a-meeting-s-all-recordings for that

	// TODO: meeting ID or UUID?
	recordingList := client.GetRecordingList(hook.MeetingID)
	//
	// step 2: choose the recording for video + speaker view
	//
	// step 3: download the actual recording to somewhere. IDK how to DL it yet (there's a URL for it somewhere maybe?)
	//
	// step 4: upload it to YouTube. we plan to do that in a live coding session:
	// https://github.com/go-study-group/agendas/issues/3

	// I don't think contexts have a way to return a bare status code, error
	// or otherwise
	return c.Error(200, errors.New("looks good :)"))
}
