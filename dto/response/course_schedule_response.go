package response

type CourseScheduleResponse struct {
	Id        int    `json:"id"`
	CourseId  string `json:"course_id"`
	Room      string `json:"room"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
}
