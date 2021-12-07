package sensor

import (
	"context"
	"os"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/internalapi/central"
	sensorInternal "github.com/stackrox/rox/generated/internalapi/sensor"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/centralsensor"
	"github.com/stackrox/rox/pkg/clusterid"
	"github.com/stackrox/rox/pkg/env"
	"github.com/stackrox/rox/pkg/expiringcache"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/grpc"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/namespaces"
	"github.com/stackrox/rox/pkg/protoutils"
	"github.com/stackrox/rox/pkg/satoken"
	"github.com/stackrox/rox/sensor/common"
	"github.com/stackrox/rox/sensor/common/admissioncontroller"
	"github.com/stackrox/rox/sensor/common/centralclient"
	"github.com/stackrox/rox/sensor/common/certdistribution"
	"github.com/stackrox/rox/sensor/common/clusterentities"
	"github.com/stackrox/rox/sensor/common/compliance"
	"github.com/stackrox/rox/sensor/common/config"
	"github.com/stackrox/rox/sensor/common/deployment"
	"github.com/stackrox/rox/sensor/common/detector"
	"github.com/stackrox/rox/sensor/common/externalsrcs"
	"github.com/stackrox/rox/sensor/common/image"
	"github.com/stackrox/rox/sensor/common/networkflow/manager"
	"github.com/stackrox/rox/sensor/common/networkflow/service"
	"github.com/stackrox/rox/sensor/common/processfilter"
	"github.com/stackrox/rox/sensor/common/processsignal"
	"github.com/stackrox/rox/sensor/common/reprocessor"
	"github.com/stackrox/rox/sensor/common/sensor"
	"github.com/stackrox/rox/sensor/common/sensor/helmconfig"
	signalService "github.com/stackrox/rox/sensor/common/signal"
	k8sadmctrl "github.com/stackrox/rox/sensor/kubernetes/admissioncontroller"
	"github.com/stackrox/rox/sensor/kubernetes/client"
	"github.com/stackrox/rox/sensor/kubernetes/clusterhealth"
	"github.com/stackrox/rox/sensor/kubernetes/clusterstatus"
	"github.com/stackrox/rox/sensor/kubernetes/enforcer"
	"github.com/stackrox/rox/sensor/kubernetes/fake"
	"github.com/stackrox/rox/sensor/kubernetes/listener"
	"github.com/stackrox/rox/sensor/kubernetes/listener/resources"
	"github.com/stackrox/rox/sensor/kubernetes/networkpolicies"
	"github.com/stackrox/rox/sensor/kubernetes/operator"
	"github.com/stackrox/rox/sensor/kubernetes/orchestrator"
	"github.com/stackrox/rox/sensor/kubernetes/telemetry"
	"github.com/stackrox/rox/sensor/kubernetes/upgrade"
)

var (
	log = logging.LoggerForModule()
)

func loadHelmConfig(client client.Interface, deploymentIdentification *storage.SensorDeploymentIdentification) (*central.HelmManagedConfigInit, operator.Operator, error) {
	var helmManagedConfig *central.HelmManagedConfigInit
	if configFP := helmconfig.HelmConfigFingerprint.Setting(); configFP != "" {
		var err error
		helmManagedConfig, err = helmconfig.Load()
		if err != nil {
			return nil, nil, errors.Wrap(err, "loading Helm cluster config")
		}
		if helmManagedConfig.GetClusterConfig().GetConfigFingerprint() != configFP {
			return nil, nil, errors.Errorf("fingerprint %q of loaded config does not match expected fingerprint %q, config changes can only be applied via 'helm upgrade' or a similar chart-based mechanism", helmManagedConfig.GetClusterConfig().GetConfigFingerprint(), configFP)
		}
		log.Infof("Loaded Helm cluster configuration with fingerprint %q", configFP)

		if err := helmconfig.CheckEffectiveClusterName(helmManagedConfig); err != nil {
			return nil, nil, errors.Wrap(err, "validating cluster name")
		}
	}

	if helmManagedConfig.GetClusterName() == "" {
		certClusterID, err := clusterid.ParseClusterIDFromServiceCert(storage.ServiceType_SENSOR_SERVICE)
		if err != nil {
			return nil, nil, errors.Wrap(err, "parsing cluster ID from service certificate")
		}
		if certClusterID == centralsensor.InitCertClusterID {
			return nil, nil, errors.New("a sensor that uses certificates from an init bundle must have a cluster name specified")
		}
	}

	var sensorOperator operator.Operator
	if helmManagedConfig != nil {
		sensorOperator = operator.New(client.Kubernetes(), deploymentIdentification.GetAppNamespace())
		if err := sensorOperator.Initialize(context.Background()); err == nil {
			log.Error(err)
		} else {
			helmManagedConfig.GetClusterConfig().HelmReleaseRevision = sensorOperator.GetHelmReleaseRevision()
		}
	}

	return helmManagedConfig, sensorOperator, nil
}

// CreateSensor takes in a client interface and returns a sensor instantiation
func CreateSensor(client client.Interface, workloadHandler *fake.WorkloadManager) (*sensor.Sensor, error) {
	admCtrlSettingsMgr := admissioncontroller.NewSettingsManager(resources.DeploymentStoreSingleton(), resources.PodStoreSingleton())

	deploymentIdentification := fetchDeploymentIdentification(context.Background(), client.Kubernetes())
	log.Infof("Determined deployment identification: %s", protoutils.NewWrapper(deploymentIdentification))

	helmManagedConfig, sensorOperator, err := loadHelmConfig(client, deploymentIdentification)
	if err != nil {
		return nil, err
	}

	auditLogEventsInput := make(chan *sensorInternal.AuditEvents)
	auditLogCollectionManager := compliance.NewAuditLogCollectionManager()

	o := orchestrator.New(client.Kubernetes())
	complianceService := compliance.NewService(o, auditLogEventsInput, auditLogCollectionManager)

	configHandler := config.NewCommandHandler(admCtrlSettingsMgr, deploymentIdentification, helmManagedConfig, auditLogCollectionManager)
	enforcer, err := enforcer.New(client)
	if err != nil {
		return nil, errors.Wrap(err, "creating enforcer")
	}

	imageCache := expiringcache.NewExpiringCache(env.ReprocessInterval.DurationSetting())

	policyDetector := detector.New(enforcer, admCtrlSettingsMgr, resources.DeploymentStoreSingleton(), imageCache, auditLogEventsInput, auditLogCollectionManager)

	admCtrlMsgForwarder := admissioncontroller.NewAdmCtrlMsgForwarder(admCtrlSettingsMgr, listener.New(client, configHandler, policyDetector))

	upgradeCmdHandler, err := upgrade.NewCommandHandler(configHandler)
	if err != nil {
		return nil, errors.Wrap(err, "creating upgrade command handler")
	}

	imageService := image.NewService(imageCache)
	complianceCommandHandler := compliance.NewCommandHandler(complianceService)

	// Create Process Pipeline
	indicators := make(chan *central.MsgFromSensor)
	processPipeline := processsignal.NewProcessPipeline(indicators, clusterentities.StoreInstance(), processfilter.Singleton(), policyDetector)
	processSignals := signalService.New(processPipeline, indicators)
	networkFlowManager :=
		manager.NewManager(clusterentities.StoreInstance(), externalsrcs.StoreInstance(), policyDetector)
	components := []common.SensorComponent{
		admCtrlMsgForwarder,
		enforcer,
		networkFlowManager,
		networkpolicies.NewCommandHandler(client.Kubernetes()),
		clusterstatus.NewUpdater(client),
		clusterhealth.NewUpdater(client.Kubernetes(), 0),
		complianceCommandHandler,
		processSignals,
		telemetry.NewCommandHandler(client.Kubernetes()),
		upgradeCmdHandler,
		externalsrcs.Singleton(),
		admissioncontroller.AlertHandlerSingleton(),
		auditLogCollectionManager,
	}

	if features.VulnRiskManagement.Enabled() {
		components = append(components, reprocessor.NewHandler(resources.DeploymentStoreSingleton()))
	}

	sensorNamespace, err := satoken.LoadNamespaceFromFile()
	if err != nil {
		log.Errorf("Failed to determine namespace from service account token file: %s", err)
	}
	if sensorNamespace == "" {
		sensorNamespace = os.Getenv("POD_NAMESPACE")
	}
	if sensorNamespace == "" {
		sensorNamespace = namespaces.StackRox
		log.Warnf("Unable to determine Sensor namespace, defaulting to %s", sensorNamespace)
	}

	if admCtrlSettingsMgr != nil {
		components = append(components, k8sadmctrl.NewConfigMapSettingsPersister(client.Kubernetes(), admCtrlSettingsMgr, sensorNamespace))
	}

	centralClient, err := centralclient.NewClient(env.CentralEndpoint.Setting())
	if err != nil {
		return nil, errors.Wrap(err, "creating central client")
	}

	s := sensor.NewSensor(
		configHandler,
		policyDetector,
		imageService,
		centralClient,
		sensorOperator,
		components...,
	)

	if workloadHandler != nil {
		workloadHandler.SetSignalHandlers(processPipeline, networkFlowManager)
	}

	networkFlowService := service.NewService(networkFlowManager)
	apiServices := []grpc.APIService{
		networkFlowService,
		processSignals,
		complianceService,
		imageService,
		deployment.NewService(resources.DeploymentStoreSingleton(), resources.PodStoreSingleton()),
	}

	if admCtrlSettingsMgr != nil {
		apiServices = append(apiServices, admissioncontroller.NewManagementService(admCtrlSettingsMgr, admissioncontroller.AlertHandlerSingleton()))
	}

	apiServices = append(apiServices, certdistribution.NewService(client.Kubernetes(), sensorNamespace))

	s.AddAPIServices(apiServices...)
	return s, nil
}
