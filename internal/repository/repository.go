package repository

type Repository interface {
	QRLocation() QRLocationRepository
	Raybot() RaybotRepository
	RaybotCommand() RaybotCommandRepository
	Workflow() WorkflowRepository
	WorkflowExecution() WorkflowExecutionRepository
	StepExecution() StepExecutionRepository
}
