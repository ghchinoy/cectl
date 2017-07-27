package ce

const (
	InstancesURI                         = "/instances"
	InstancesFormatURI                   = "/instances/%s"
	InstanceConfigurationURI             = "/instances/configuration"
	InstanceConfigurationFormatURI       = "/instances/configuration/%s"
	InstanceDocsURI                      = "/instances/docs"
	InstanceOperationDocsFormatURI       = "/instances/docs/%s"
	InstancesEventsURI                   = "/instances/events"
	InstancesEventsAnalyticsAccountsURI  = "/instances/events/analytics/accounts"
	InstancesEventsAnalyticsInstancesURI = "/instances/events/analytics/instances"
	InstancesEventsFormatURI             = "/instances/events/%s"
	InstancesObjectsDefinitionsURI       = "/instances/objects/definitions"
	InstancesTransformationsURI          = "/instances/transformations"
	InstanceTransformationsFormatURI     = "/instances/%s/transformations"
	InstanceDocFormatURI                 = "/instances/%s/docs"
)

func GetAllInstances() {}

func GetInstanceInfo(id string) {}

func GetInstanceDocs(id string) {}

func GetInstanceTransformations(id string) {}
