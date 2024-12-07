package request

import "SchoolManagement/model"

type CourseScheduleRequest struct {
	CourseId  string `json:"course_id" validate:"required"`
	Room      string `json:"room" validate:"required"`
	StartTime int64  `json:"start_time" validate:"required"`
	EndTime   int64  `json:"end_time" validate:"required"`
}

func (req *CourseScheduleRequest) ToCourseSchedule() model.CourseSchedule {
	return model.CourseSchedule{
		CourseId:  req.CourseId,
		Room:      req.Room,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}
}
