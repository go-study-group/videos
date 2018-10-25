package zoom

// RecordingList is the struct returned by the API:
// GET https://api.zoom.us/v2/meetings/{meetingId}/recordings
type RecordingList struct {
	UUID           string          `json:"uuid"`
	RecordingFiles []RecordingFile `json:"recording_files"`
}

// RecordingFile is one element in the list of 'recording_files' in the API:
// GET https://api.zoom.us/v2/meetings/{meetingId}/recordings
type RecordingFile struct {
	RecordingType string `json:"recording_type"`
	DownloadURL   string `json:"download_url"`
}

// GetRecordingList does the Zoom API call for this:
//
func (c *Client) GetRecordingList(meetingID string) (*RecordingList, error) {
	// GET https://api.zoom.us/v2/meetings/{meetingId}/recordings
	list := new(RecordingList)
	err := c.get("GET", "/meetings/"+meetingID+"/recordings", list)
	if err != nil {
		return nil, err
	}
	return list, nil
}
