package edge

type Edge struct {
	ID           string  `json:"id" validate:"required,uuid"`
	Type         string  `json:"type" validate:"required"`
	Source       string  `json:"source" validate:"required,uuid"`
	Target       string  `json:"target" validate:"required,uuid"`
	SourceHandle string  `json:"source_handle" validate:"required"`
	TargetHandle string  `json:"target_handle" validate:"required"`
	Label        string  `json:"label" validate:"required"`
	Animated     bool    `json:"animated" validate:"required"`
	SourceX      float32 `json:"source_x" validate:"required"`
	SourceY      float32 `json:"source_y" validate:"required"`
	TargetX      float32 `json:"target_x" validate:"required"`
	TargetY      float32 `json:"target_y" validate:"required"`
}
