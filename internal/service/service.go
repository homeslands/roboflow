package service

type Service interface {
	QRLocation() QRLocationService
	Raybot() RaybotService
	RaybotCommand() RaybotCommandService
	Workflow() WorkflowService
	WorkflowExecution() WorkflowExecutionService
	StepExecution() StepExecutionService
}
