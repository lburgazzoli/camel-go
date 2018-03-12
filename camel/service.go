package camel

// ServiceStatus --
type ServiceStatus int

const (
	// ServiceStatusSTOPPING --
	ServiceStatusSTOPPING ServiceStatus = 10

	// ServiceStatusSTOPPED --
	ServiceStatusSTOPPED ServiceStatus = 11

	// ServiceStatusSUSPENDING --
	ServiceStatusSUSPENDING ServiceStatus = 20

	// ServiceStatusSUSPENDED --
	ServiceStatusSUSPENDED ServiceStatus = 21

	// ServiceStatusSTARTING --
	ServiceStatusSTARTING ServiceStatus = 30

	// ServiceStatusSTARTED --
	ServiceStatusSTARTED ServiceStatus = 31
)

// Service --
type Service interface {
	Start()
	Stop()
}
