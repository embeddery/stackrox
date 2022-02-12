// Code generated by pg-bindings generator. DO NOT EDIT.

package postgres

import (
	"context"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stackrox/rox/central/globaldb"
	"github.com/stackrox/rox/central/metrics"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/logging"
	ops "github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/postgres/pgutils"
)

const (
	countStmt  = "SELECT COUNT(*) FROM deployments"
	existsStmt = "SELECT EXISTS(SELECT 1 FROM deployments WHERE Id = $1)"

	getStmt     = "SELECT serialized FROM deployments WHERE Id = $1"
	deleteStmt  = "DELETE FROM deployments WHERE Id = $1"
	walkStmt    = "SELECT serialized FROM deployments"
	getIDsStmt  = "SELECT Id FROM deployments"
	getManyStmt = "SELECT serialized FROM deployments WHERE Id = ANY($1::text[])"

	deleteManyStmt = "DELETE FROM deployments WHERE Id = ANY($1::text[])"
)

var (
	log = logging.LoggerForModule()

	table = "deployments"
)

type Store interface {
	Count() (int, error)
	Exists(id string) (bool, error)
	Get(id string) (*storage.Deployment, bool, error)
	Upsert(obj *storage.Deployment) error
	UpsertMany(objs []*storage.Deployment) error
	Delete(id string) error
	GetIDs() ([]string, error)
	GetMany(ids []string) ([]*storage.Deployment, []int, error)
	DeleteMany(ids []string) error

	Walk(fn func(obj *storage.Deployment) error) error
	AckKeysIndexed(keys ...string) error
	GetKeysToIndex() ([]string, error)
}

type storeImpl struct {
	db *pgxpool.Pool
}

func createTableDeployments(db *pgxpool.Pool) {
	// hack for testing, remove
	db.Exec(context.Background(), "DROP TABLE deployments CASCADE")

	table := `
create table if not exists deployments (
    Id varchar,
    Name varchar,
    Hash numeric,
    Type varchar,
    Namespace varchar,
    NamespaceId varchar,
    OrchestratorComponent bool,
    Replicas numeric,
    Labels jsonb,
    PodLabels jsonb,
    LabelSelector_MatchLabels jsonb,
    Created timestamp,
    ClusterId varchar,
    ClusterName varchar,
    Annotations jsonb,
    Priority numeric,
    Inactive bool,
    ImagePullSecrets text[],
    ServiceAccount varchar,
    ServiceAccountPermissionLevel integer,
    AutomountServiceAccountToken bool,
    HostNetwork bool,
    HostPid bool,
    HostIpc bool,
    RuntimeClass varchar,
    StateTimestamp numeric,
    RiskScore numeric,
    ProcessTags text[],
    serialized bytea,
    PRIMARY KEY(Id)
)
`

	_, err := db.Exec(context.Background(), table)
	if err != nil {
		panic("error creating table: " + table)
	}

	createTableDeploymentsRequirements(db)
	createTableDeploymentsContainers(db)
	createTableDeploymentsTolerations(db)
	createTableDeploymentsPorts(db)
}

func createTableDeploymentsRequirements(db *pgxpool.Pool) {
	// hack for testing, remove
	db.Exec(context.Background(), "DROP TABLE deployments_Requirements CASCADE")

	table := `
create table if not exists deployments_Requirements (
    parent_Id varchar,
    idx numeric,
    Key varchar,
    Op integer,
    Values text[],
    PRIMARY KEY(parent_Id, idx),
    CONSTRAINT fk_parent_table FOREIGN KEY (parent_Id) REFERENCES deployments(Id) ON DELETE CASCADE
)
`

	_, err := db.Exec(context.Background(), table)
	if err != nil {
		panic("error creating table: " + table)
	}

}

func createTableDeploymentsContainers(db *pgxpool.Pool) {
	// hack for testing, remove
	db.Exec(context.Background(), "DROP TABLE deployments_Containers CASCADE")

	table := `
create table if not exists deployments_Containers (
    parent_Id varchar,
    idx numeric,
    Id varchar,
    Config_Command text[],
    Config_Args text[],
    Config_Directory varchar,
    Config_User varchar,
    Config_Uid numeric,
    Config_AppArmorProfile varchar,
    Image_Id varchar,
    Image_Name_Registry varchar,
    Image_Name_Remote varchar,
    Image_Name_Tag varchar,
    Image_Name_FullName varchar,
    Image_NotPullable bool,
    SecurityContext_Privileged bool,
    SecurityContext_Selinux_User varchar,
    SecurityContext_Selinux_Role varchar,
    SecurityContext_Selinux_Type varchar,
    SecurityContext_Selinux_Level varchar,
    SecurityContext_DropCapabilities text[],
    SecurityContext_AddCapabilities text[],
    SecurityContext_ReadOnlyRootFilesystem bool,
    SecurityContext_SeccompProfile_Type integer,
    SecurityContext_SeccompProfile_LocalhostProfile varchar,
    Resources_CpuCoresRequest numeric,
    Resources_CpuCoresLimit numeric,
    Resources_MemoryMbRequest numeric,
    Resources_MemoryMbLimit numeric,
    Name varchar,
    LivenessProbe_Defined bool,
    ReadinessProbe_Defined bool,
    PRIMARY KEY(parent_Id, idx),
    CONSTRAINT fk_parent_table FOREIGN KEY (parent_Id) REFERENCES deployments(Id) ON DELETE CASCADE
)
`

	_, err := db.Exec(context.Background(), table)
	if err != nil {
		panic("error creating table: " + table)
	}

	createTableDeploymentsContainersEnv(db)
	createTableDeploymentsContainersVolumes(db)
	createTableDeploymentsContainersPorts(db)
	createTableDeploymentsContainersSecrets(db)
}

func createTableDeploymentsContainersEnv(db *pgxpool.Pool) {
	// hack for testing, remove
	db.Exec(context.Background(), "DROP TABLE deployments_Containers_Env CASCADE")

	table := `
create table if not exists deployments_Containers_Env (
    parent_parent_Id varchar,
    parent_idx numeric,
    idx numeric,
    Key varchar,
    Value varchar,
    EnvVarSource integer,
    PRIMARY KEY(parent_parent_Id, parent_idx, idx),
    CONSTRAINT fk_parent_table FOREIGN KEY (parent_parent_Id, parent_idx) REFERENCES deployments_Containers(parent_Id, idx) ON DELETE CASCADE
)
`

	_, err := db.Exec(context.Background(), table)
	if err != nil {
		panic("error creating table: " + table)
	}

}

func createTableDeploymentsContainersVolumes(db *pgxpool.Pool) {
	// hack for testing, remove
	db.Exec(context.Background(), "DROP TABLE deployments_Containers_Volumes CASCADE")

	table := `
create table if not exists deployments_Containers_Volumes (
    parent_parent_Id varchar,
    parent_idx numeric,
    idx numeric,
    Name varchar,
    Source varchar,
    Destination varchar,
    ReadOnly bool,
    Type varchar,
    MountPropagation integer,
    PRIMARY KEY(parent_parent_Id, parent_idx, idx),
    CONSTRAINT fk_parent_table FOREIGN KEY (parent_parent_Id, parent_idx) REFERENCES deployments_Containers(parent_Id, idx) ON DELETE CASCADE
)
`

	_, err := db.Exec(context.Background(), table)
	if err != nil {
		panic("error creating table: " + table)
	}

}

func createTableDeploymentsContainersPorts(db *pgxpool.Pool) {
	// hack for testing, remove
	db.Exec(context.Background(), "DROP TABLE deployments_Containers_Ports CASCADE")

	table := `
create table if not exists deployments_Containers_Ports (
    parent_parent_Id varchar,
    parent_idx numeric,
    idx numeric,
    Name varchar,
    ContainerPort numeric,
    Protocol varchar,
    Exposure integer,
    ExposedPort numeric,
    PRIMARY KEY(parent_parent_Id, parent_idx, idx),
    CONSTRAINT fk_parent_table FOREIGN KEY (parent_parent_Id, parent_idx) REFERENCES deployments_Containers(parent_Id, idx) ON DELETE CASCADE
)
`

	_, err := db.Exec(context.Background(), table)
	if err != nil {
		panic("error creating table: " + table)
	}

	createTableDeploymentsContainersPortsExposureInfos(db)
}

func createTableDeploymentsContainersPortsExposureInfos(db *pgxpool.Pool) {
	// hack for testing, remove
	db.Exec(context.Background(), "DROP TABLE deployments_Containers_Ports_ExposureInfos CASCADE")

	table := `
create table if not exists deployments_Containers_Ports_ExposureInfos (
    parent_parent_parent_Id varchar,
    parent_parent_idx numeric,
    parent_idx numeric,
    idx numeric,
    Level integer,
    ServiceName varchar,
    ServiceId varchar,
    ServiceClusterIp varchar,
    ServicePort numeric,
    NodePort numeric,
    ExternalIps text[],
    ExternalHostnames text[],
    PRIMARY KEY(parent_parent_parent_Id, parent_parent_idx, parent_idx, idx),
    CONSTRAINT fk_parent_table FOREIGN KEY (parent_parent_parent_Id, parent_parent_idx, parent_idx) REFERENCES deployments_Containers_Ports(parent_parent_Id, parent_idx, idx) ON DELETE CASCADE
)
`

	_, err := db.Exec(context.Background(), table)
	if err != nil {
		panic("error creating table: " + table)
	}

}

func createTableDeploymentsContainersSecrets(db *pgxpool.Pool) {
	// hack for testing, remove
	db.Exec(context.Background(), "DROP TABLE deployments_Containers_Secrets CASCADE")

	table := `
create table if not exists deployments_Containers_Secrets (
    parent_parent_Id varchar,
    parent_idx numeric,
    idx numeric,
    Name varchar,
    Path varchar,
    PRIMARY KEY(parent_parent_Id, parent_idx, idx),
    CONSTRAINT fk_parent_table FOREIGN KEY (parent_parent_Id, parent_idx) REFERENCES deployments_Containers(parent_Id, idx) ON DELETE CASCADE
)
`

	_, err := db.Exec(context.Background(), table)
	if err != nil {
		panic("error creating table: " + table)
	}

}

func createTableDeploymentsTolerations(db *pgxpool.Pool) {
	// hack for testing, remove
	db.Exec(context.Background(), "DROP TABLE deployments_Tolerations CASCADE")

	table := `
create table if not exists deployments_Tolerations (
    parent_Id varchar,
    idx numeric,
    Key varchar,
    Operator integer,
    Value varchar,
    TaintEffect integer,
    PRIMARY KEY(parent_Id, idx),
    CONSTRAINT fk_parent_table FOREIGN KEY (parent_Id) REFERENCES deployments(Id) ON DELETE CASCADE
)
`

	_, err := db.Exec(context.Background(), table)
	if err != nil {
		panic("error creating table: " + table)
	}

}

func createTableDeploymentsPorts(db *pgxpool.Pool) {
	// hack for testing, remove
	db.Exec(context.Background(), "DROP TABLE deployments_Ports CASCADE")

	table := `
create table if not exists deployments_Ports (
    parent_Id varchar,
    idx numeric,
    Name varchar,
    ContainerPort numeric,
    Protocol varchar,
    Exposure integer,
    ExposedPort numeric,
    PRIMARY KEY(parent_Id, idx),
    CONSTRAINT fk_parent_table FOREIGN KEY (parent_Id) REFERENCES deployments(Id) ON DELETE CASCADE
)
`

	_, err := db.Exec(context.Background(), table)
	if err != nil {
		panic("error creating table: " + table)
	}

	createTableDeploymentsPortsExposureInfos(db)
}

func createTableDeploymentsPortsExposureInfos(db *pgxpool.Pool) {
	// hack for testing, remove
	db.Exec(context.Background(), "DROP TABLE deployments_Ports_ExposureInfos CASCADE")

	table := `
create table if not exists deployments_Ports_ExposureInfos (
    parent_parent_Id varchar,
    parent_idx numeric,
    idx numeric,
    Level integer,
    ServiceName varchar,
    ServiceId varchar,
    ServiceClusterIp varchar,
    ServicePort numeric,
    NodePort numeric,
    ExternalIps text[],
    ExternalHostnames text[],
    PRIMARY KEY(parent_parent_Id, parent_idx, idx),
    CONSTRAINT fk_parent_table FOREIGN KEY (parent_parent_Id, parent_idx) REFERENCES deployments_Ports(parent_Id, idx) ON DELETE CASCADE
)
`

	_, err := db.Exec(context.Background(), table)
	if err != nil {
		panic("error creating table: " + table)
	}

}

func insertIntoDeployments(db *pgxpool.Pool, obj *storage.Deployment) error {

	serialized, marshalErr := obj.Marshal()
	if marshalErr != nil {
		return marshalErr
	}

	values := []interface{}{
		// parent primary keys start

		obj.GetId(),
		obj.GetName(),
		obj.GetHash(),
		obj.GetType(),
		obj.GetNamespace(),
		obj.GetNamespaceId(),
		obj.GetOrchestratorComponent(),
		obj.GetReplicas(),
		obj.GetLabels(),
		obj.GetPodLabels(),
		obj.GetLabelSelector().GetMatchLabels(),
		obj.GetCreated(),
		obj.GetClusterId(),
		obj.GetClusterName(),
		obj.GetAnnotations(),
		obj.GetPriority(),
		obj.GetInactive(),
		obj.GetImagePullSecrets(),
		obj.GetServiceAccount(),
		obj.GetServiceAccountPermissionLevel(),
		obj.GetAutomountServiceAccountToken(),
		obj.GetHostNetwork(),
		obj.GetHostPid(),
		obj.GetHostIpc(),
		obj.GetRuntimeClass(),
		obj.GetStateTimestamp(),
		obj.GetRiskScore(),
		obj.GetProcessTags(),
		serialized,
	}

	finalStr := "INSERT INTO deployments (Id, Name, Hash, Type, Namespace, NamespaceId, OrchestratorComponent, Replicas, Labels, PodLabels, LabelSelector_MatchLabels, Created, ClusterId, ClusterName, Annotations, Priority, Inactive, ImagePullSecrets, ServiceAccount, ServiceAccountPermissionLevel, AutomountServiceAccountToken, HostNetwork, HostPid, HostIpc, RuntimeClass, StateTimestamp, RiskScore, ProcessTags, serialized) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29) ON CONFLICT(Id) DO UPDATE SET Id = EXCLUDED.Id, Name = EXCLUDED.Name, Hash = EXCLUDED.Hash, Type = EXCLUDED.Type, Namespace = EXCLUDED.Namespace, NamespaceId = EXCLUDED.NamespaceId, OrchestratorComponent = EXCLUDED.OrchestratorComponent, Replicas = EXCLUDED.Replicas, Labels = EXCLUDED.Labels, PodLabels = EXCLUDED.PodLabels, LabelSelector_MatchLabels = EXCLUDED.LabelSelector_MatchLabels, Created = EXCLUDED.Created, ClusterId = EXCLUDED.ClusterId, ClusterName = EXCLUDED.ClusterName, Annotations = EXCLUDED.Annotations, Priority = EXCLUDED.Priority, Inactive = EXCLUDED.Inactive, ImagePullSecrets = EXCLUDED.ImagePullSecrets, ServiceAccount = EXCLUDED.ServiceAccount, ServiceAccountPermissionLevel = EXCLUDED.ServiceAccountPermissionLevel, AutomountServiceAccountToken = EXCLUDED.AutomountServiceAccountToken, HostNetwork = EXCLUDED.HostNetwork, HostPid = EXCLUDED.HostPid, HostIpc = EXCLUDED.HostIpc, RuntimeClass = EXCLUDED.RuntimeClass, StateTimestamp = EXCLUDED.StateTimestamp, RiskScore = EXCLUDED.RiskScore, ProcessTags = EXCLUDED.ProcessTags, serialized = EXCLUDED.serialized"
	_, err := db.Exec(context.Background(), finalStr, values...)
	if err != nil {
		return err
	}

	var query string

	for childIdx, child := range obj.GetLabelSelector().GetRequirements() {
		if err := insertIntoDeploymentsRequirements(db, child, obj.GetId(), childIdx); err != nil {
			return err
		}
	}

	query = "delete from deployments_Requirements where parent_Id = $1 AND idx >= $2"
	_, err = db.Exec(context.Background(), query, obj.GetId(), len(obj.GetLabelSelector().GetRequirements()))
	if err != nil {
		return err
	}
	for childIdx, child := range obj.GetContainers() {
		if err := insertIntoDeploymentsContainers(db, child, obj.GetId(), childIdx); err != nil {
			return err
		}
	}

	query = "delete from deployments_Containers where parent_Id = $1 AND idx >= $2"
	_, err = db.Exec(context.Background(), query, obj.GetId(), len(obj.GetContainers()))
	if err != nil {
		return err
	}
	for childIdx, child := range obj.GetTolerations() {
		if err := insertIntoDeploymentsTolerations(db, child, obj.GetId(), childIdx); err != nil {
			return err
		}
	}

	query = "delete from deployments_Tolerations where parent_Id = $1 AND idx >= $2"
	_, err = db.Exec(context.Background(), query, obj.GetId(), len(obj.GetTolerations()))
	if err != nil {
		return err
	}
	for childIdx, child := range obj.GetPorts() {
		if err := insertIntoDeploymentsPorts(db, child, obj.GetId(), childIdx); err != nil {
			return err
		}
	}

	query = "delete from deployments_Ports where parent_Id = $1 AND idx >= $2"
	_, err = db.Exec(context.Background(), query, obj.GetId(), len(obj.GetPorts()))
	if err != nil {
		return err
	}
	return nil
}

func insertIntoDeploymentsRequirements(db *pgxpool.Pool, obj *storage.LabelSelector_Requirement, parent_Id string, idx int) error {

	values := []interface{}{
		// parent primary keys start

		parent_Id,
		idx,
		obj.GetKey(),
		obj.GetOp(),
		obj.GetValues(),
	}

	finalStr := "INSERT INTO deployments_Requirements (parent_Id, idx, Key, Op, Values) VALUES($1, $2, $3, $4, $5) ON CONFLICT(parent_Id, idx) DO UPDATE SET parent_Id = EXCLUDED.parent_Id, idx = EXCLUDED.idx, Key = EXCLUDED.Key, Op = EXCLUDED.Op, Values = EXCLUDED.Values"
	_, err := db.Exec(context.Background(), finalStr, values...)
	if err != nil {
		return err
	}

	return nil
}

func insertIntoDeploymentsContainers(db *pgxpool.Pool, obj *storage.Container, parent_Id string, idx int) error {

	values := []interface{}{
		// parent primary keys start

		parent_Id,
		idx,
		obj.GetId(),
		obj.GetConfig().GetCommand(),
		obj.GetConfig().GetArgs(),
		obj.GetConfig().GetDirectory(),
		obj.GetConfig().GetUser(),
		obj.GetConfig().GetUid(),
		obj.GetConfig().GetAppArmorProfile(),
		obj.GetImage().GetId(),
		obj.GetImage().GetName().GetRegistry(),
		obj.GetImage().GetName().GetRemote(),
		obj.GetImage().GetName().GetTag(),
		obj.GetImage().GetName().GetFullName(),
		obj.GetImage().GetNotPullable(),
		obj.GetSecurityContext().GetPrivileged(),
		obj.GetSecurityContext().GetSelinux().GetUser(),
		obj.GetSecurityContext().GetSelinux().GetRole(),
		obj.GetSecurityContext().GetSelinux().GetType(),
		obj.GetSecurityContext().GetSelinux().GetLevel(),
		obj.GetSecurityContext().GetDropCapabilities(),
		obj.GetSecurityContext().GetAddCapabilities(),
		obj.GetSecurityContext().GetReadOnlyRootFilesystem(),
		obj.GetSecurityContext().GetSeccompProfile().GetType(),
		obj.GetSecurityContext().GetSeccompProfile().GetLocalhostProfile(),
		obj.GetResources().GetCpuCoresRequest(),
		obj.GetResources().GetCpuCoresLimit(),
		obj.GetResources().GetMemoryMbRequest(),
		obj.GetResources().GetMemoryMbLimit(),
		obj.GetName(),
		obj.GetLivenessProbe().GetDefined(),
		obj.GetReadinessProbe().GetDefined(),
	}

	finalStr := "INSERT INTO deployments_Containers (parent_Id, idx, Id, Config_Command, Config_Args, Config_Directory, Config_User, Config_Uid, Config_AppArmorProfile, Image_Id, Image_Name_Registry, Image_Name_Remote, Image_Name_Tag, Image_Name_FullName, Image_NotPullable, SecurityContext_Privileged, SecurityContext_Selinux_User, SecurityContext_Selinux_Role, SecurityContext_Selinux_Type, SecurityContext_Selinux_Level, SecurityContext_DropCapabilities, SecurityContext_AddCapabilities, SecurityContext_ReadOnlyRootFilesystem, SecurityContext_SeccompProfile_Type, SecurityContext_SeccompProfile_LocalhostProfile, Resources_CpuCoresRequest, Resources_CpuCoresLimit, Resources_MemoryMbRequest, Resources_MemoryMbLimit, Name, LivenessProbe_Defined, ReadinessProbe_Defined) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32) ON CONFLICT(parent_Id, idx) DO UPDATE SET parent_Id = EXCLUDED.parent_Id, idx = EXCLUDED.idx, Id = EXCLUDED.Id, Config_Command = EXCLUDED.Config_Command, Config_Args = EXCLUDED.Config_Args, Config_Directory = EXCLUDED.Config_Directory, Config_User = EXCLUDED.Config_User, Config_Uid = EXCLUDED.Config_Uid, Config_AppArmorProfile = EXCLUDED.Config_AppArmorProfile, Image_Id = EXCLUDED.Image_Id, Image_Name_Registry = EXCLUDED.Image_Name_Registry, Image_Name_Remote = EXCLUDED.Image_Name_Remote, Image_Name_Tag = EXCLUDED.Image_Name_Tag, Image_Name_FullName = EXCLUDED.Image_Name_FullName, Image_NotPullable = EXCLUDED.Image_NotPullable, SecurityContext_Privileged = EXCLUDED.SecurityContext_Privileged, SecurityContext_Selinux_User = EXCLUDED.SecurityContext_Selinux_User, SecurityContext_Selinux_Role = EXCLUDED.SecurityContext_Selinux_Role, SecurityContext_Selinux_Type = EXCLUDED.SecurityContext_Selinux_Type, SecurityContext_Selinux_Level = EXCLUDED.SecurityContext_Selinux_Level, SecurityContext_DropCapabilities = EXCLUDED.SecurityContext_DropCapabilities, SecurityContext_AddCapabilities = EXCLUDED.SecurityContext_AddCapabilities, SecurityContext_ReadOnlyRootFilesystem = EXCLUDED.SecurityContext_ReadOnlyRootFilesystem, SecurityContext_SeccompProfile_Type = EXCLUDED.SecurityContext_SeccompProfile_Type, SecurityContext_SeccompProfile_LocalhostProfile = EXCLUDED.SecurityContext_SeccompProfile_LocalhostProfile, Resources_CpuCoresRequest = EXCLUDED.Resources_CpuCoresRequest, Resources_CpuCoresLimit = EXCLUDED.Resources_CpuCoresLimit, Resources_MemoryMbRequest = EXCLUDED.Resources_MemoryMbRequest, Resources_MemoryMbLimit = EXCLUDED.Resources_MemoryMbLimit, Name = EXCLUDED.Name, LivenessProbe_Defined = EXCLUDED.LivenessProbe_Defined, ReadinessProbe_Defined = EXCLUDED.ReadinessProbe_Defined"
	_, err := db.Exec(context.Background(), finalStr, values...)
	if err != nil {
		return err
	}

	var query string

	for childIdx, child := range obj.GetConfig().GetEnv() {
		if err := insertIntoDeploymentsContainersEnv(db, child, parent_Id, idx, childIdx); err != nil {
			return err
		}
	}

	query = "delete from deployments_Containers_Env where parent_parent_Id = $1 AND parent_idx = $2 AND idx >= $3"
	_, err = db.Exec(context.Background(), query, parent_Id, idx, len(obj.GetConfig().GetEnv()))
	if err != nil {
		return err
	}
	for childIdx, child := range obj.GetVolumes() {
		if err := insertIntoDeploymentsContainersVolumes(db, child, parent_Id, idx, childIdx); err != nil {
			return err
		}
	}

	query = "delete from deployments_Containers_Volumes where parent_parent_Id = $1 AND parent_idx = $2 AND idx >= $3"
	_, err = db.Exec(context.Background(), query, parent_Id, idx, len(obj.GetVolumes()))
	if err != nil {
		return err
	}
	for childIdx, child := range obj.GetPorts() {
		if err := insertIntoDeploymentsContainersPorts(db, child, parent_Id, idx, childIdx); err != nil {
			return err
		}
	}

	query = "delete from deployments_Containers_Ports where parent_parent_Id = $1 AND parent_idx = $2 AND idx >= $3"
	_, err = db.Exec(context.Background(), query, parent_Id, idx, len(obj.GetPorts()))
	if err != nil {
		return err
	}
	for childIdx, child := range obj.GetSecrets() {
		if err := insertIntoDeploymentsContainersSecrets(db, child, parent_Id, idx, childIdx); err != nil {
			return err
		}
	}

	query = "delete from deployments_Containers_Secrets where parent_parent_Id = $1 AND parent_idx = $2 AND idx >= $3"
	_, err = db.Exec(context.Background(), query, parent_Id, idx, len(obj.GetSecrets()))
	if err != nil {
		return err
	}
	return nil
}

func insertIntoDeploymentsContainersEnv(db *pgxpool.Pool, obj *storage.ContainerConfig_EnvironmentConfig, parent_parent_Id string, parent_idx int, idx int) error {

	values := []interface{}{
		// parent primary keys start

		parent_parent_Id,
		parent_idx,
		idx,
		obj.GetKey(),
		obj.GetValue(),
		obj.GetEnvVarSource(),
	}

	finalStr := "INSERT INTO deployments_Containers_Env (parent_parent_Id, parent_idx, idx, Key, Value, EnvVarSource) VALUES($1, $2, $3, $4, $5, $6) ON CONFLICT(parent_parent_Id, parent_idx, idx) DO UPDATE SET parent_parent_Id = EXCLUDED.parent_parent_Id, parent_idx = EXCLUDED.parent_idx, idx = EXCLUDED.idx, Key = EXCLUDED.Key, Value = EXCLUDED.Value, EnvVarSource = EXCLUDED.EnvVarSource"
	_, err := db.Exec(context.Background(), finalStr, values...)
	if err != nil {
		return err
	}

	return nil
}

func insertIntoDeploymentsContainersVolumes(db *pgxpool.Pool, obj *storage.Volume, parent_parent_Id string, parent_idx int, idx int) error {

	values := []interface{}{
		// parent primary keys start

		parent_parent_Id,
		parent_idx,
		idx,
		obj.GetName(),
		obj.GetSource(),
		obj.GetDestination(),
		obj.GetReadOnly(),
		obj.GetType(),
		obj.GetMountPropagation(),
	}

	finalStr := "INSERT INTO deployments_Containers_Volumes (parent_parent_Id, parent_idx, idx, Name, Source, Destination, ReadOnly, Type, MountPropagation) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT(parent_parent_Id, parent_idx, idx) DO UPDATE SET parent_parent_Id = EXCLUDED.parent_parent_Id, parent_idx = EXCLUDED.parent_idx, idx = EXCLUDED.idx, Name = EXCLUDED.Name, Source = EXCLUDED.Source, Destination = EXCLUDED.Destination, ReadOnly = EXCLUDED.ReadOnly, Type = EXCLUDED.Type, MountPropagation = EXCLUDED.MountPropagation"
	_, err := db.Exec(context.Background(), finalStr, values...)
	if err != nil {
		return err
	}

	return nil
}

func insertIntoDeploymentsContainersPorts(db *pgxpool.Pool, obj *storage.PortConfig, parent_parent_Id string, parent_idx int, idx int) error {

	values := []interface{}{
		// parent primary keys start

		parent_parent_Id,
		parent_idx,
		idx,
		obj.GetName(),
		obj.GetContainerPort(),
		obj.GetProtocol(),
		obj.GetExposure(),
		obj.GetExposedPort(),
	}

	finalStr := "INSERT INTO deployments_Containers_Ports (parent_parent_Id, parent_idx, idx, Name, ContainerPort, Protocol, Exposure, ExposedPort) VALUES($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT(parent_parent_Id, parent_idx, idx) DO UPDATE SET parent_parent_Id = EXCLUDED.parent_parent_Id, parent_idx = EXCLUDED.parent_idx, idx = EXCLUDED.idx, Name = EXCLUDED.Name, ContainerPort = EXCLUDED.ContainerPort, Protocol = EXCLUDED.Protocol, Exposure = EXCLUDED.Exposure, ExposedPort = EXCLUDED.ExposedPort"
	_, err := db.Exec(context.Background(), finalStr, values...)
	if err != nil {
		return err
	}

	var query string

	for childIdx, child := range obj.GetExposureInfos() {
		if err := insertIntoDeploymentsContainersPortsExposureInfos(db, child, parent_parent_Id, parent_idx, idx, childIdx); err != nil {
			return err
		}
	}

	query = "delete from deployments_Containers_Ports_ExposureInfos where parent_parent_parent_Id = $1 AND parent_parent_idx = $2 AND parent_idx = $3 AND idx >= $4"
	_, err = db.Exec(context.Background(), query, parent_parent_Id, parent_idx, idx, len(obj.GetExposureInfos()))
	if err != nil {
		return err
	}
	return nil
}

func insertIntoDeploymentsContainersPortsExposureInfos(db *pgxpool.Pool, obj *storage.PortConfig_ExposureInfo, parent_parent_parent_Id string, parent_parent_idx int, parent_idx int, idx int) error {

	values := []interface{}{
		// parent primary keys start

		parent_parent_parent_Id,
		parent_parent_idx,
		parent_idx,
		idx,
		obj.GetLevel(),
		obj.GetServiceName(),
		obj.GetServiceId(),
		obj.GetServiceClusterIp(),
		obj.GetServicePort(),
		obj.GetNodePort(),
		obj.GetExternalIps(),
		obj.GetExternalHostnames(),
	}

	finalStr := "INSERT INTO deployments_Containers_Ports_ExposureInfos (parent_parent_parent_Id, parent_parent_idx, parent_idx, idx, Level, ServiceName, ServiceId, ServiceClusterIp, ServicePort, NodePort, ExternalIps, ExternalHostnames) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) ON CONFLICT(parent_parent_parent_Id, parent_parent_idx, parent_idx, idx) DO UPDATE SET parent_parent_parent_Id = EXCLUDED.parent_parent_parent_Id, parent_parent_idx = EXCLUDED.parent_parent_idx, parent_idx = EXCLUDED.parent_idx, idx = EXCLUDED.idx, Level = EXCLUDED.Level, ServiceName = EXCLUDED.ServiceName, ServiceId = EXCLUDED.ServiceId, ServiceClusterIp = EXCLUDED.ServiceClusterIp, ServicePort = EXCLUDED.ServicePort, NodePort = EXCLUDED.NodePort, ExternalIps = EXCLUDED.ExternalIps, ExternalHostnames = EXCLUDED.ExternalHostnames"
	_, err := db.Exec(context.Background(), finalStr, values...)
	if err != nil {
		return err
	}

	return nil
}

func insertIntoDeploymentsContainersSecrets(db *pgxpool.Pool, obj *storage.EmbeddedSecret, parent_parent_Id string, parent_idx int, idx int) error {

	values := []interface{}{
		// parent primary keys start

		parent_parent_Id,
		parent_idx,
		idx,
		obj.GetName(),
		obj.GetPath(),
	}

	finalStr := "INSERT INTO deployments_Containers_Secrets (parent_parent_Id, parent_idx, idx, Name, Path) VALUES($1, $2, $3, $4, $5) ON CONFLICT(parent_parent_Id, parent_idx, idx) DO UPDATE SET parent_parent_Id = EXCLUDED.parent_parent_Id, parent_idx = EXCLUDED.parent_idx, idx = EXCLUDED.idx, Name = EXCLUDED.Name, Path = EXCLUDED.Path"
	_, err := db.Exec(context.Background(), finalStr, values...)
	if err != nil {
		return err
	}

	return nil
}

func insertIntoDeploymentsTolerations(db *pgxpool.Pool, obj *storage.Toleration, parent_Id string, idx int) error {

	values := []interface{}{
		// parent primary keys start

		parent_Id,
		idx,
		obj.GetKey(),
		obj.GetOperator(),
		obj.GetValue(),
		obj.GetTaintEffect(),
	}

	finalStr := "INSERT INTO deployments_Tolerations (parent_Id, idx, Key, Operator, Value, TaintEffect) VALUES($1, $2, $3, $4, $5, $6) ON CONFLICT(parent_Id, idx) DO UPDATE SET parent_Id = EXCLUDED.parent_Id, idx = EXCLUDED.idx, Key = EXCLUDED.Key, Operator = EXCLUDED.Operator, Value = EXCLUDED.Value, TaintEffect = EXCLUDED.TaintEffect"
	_, err := db.Exec(context.Background(), finalStr, values...)
	if err != nil {
		return err
	}

	return nil
}

func insertIntoDeploymentsPorts(db *pgxpool.Pool, obj *storage.PortConfig, parent_Id string, idx int) error {

	values := []interface{}{
		// parent primary keys start

		parent_Id,
		idx,
		obj.GetName(),
		obj.GetContainerPort(),
		obj.GetProtocol(),
		obj.GetExposure(),
		obj.GetExposedPort(),
	}

	finalStr := "INSERT INTO deployments_Ports (parent_Id, idx, Name, ContainerPort, Protocol, Exposure, ExposedPort) VALUES($1, $2, $3, $4, $5, $6, $7) ON CONFLICT(parent_Id, idx) DO UPDATE SET parent_Id = EXCLUDED.parent_Id, idx = EXCLUDED.idx, Name = EXCLUDED.Name, ContainerPort = EXCLUDED.ContainerPort, Protocol = EXCLUDED.Protocol, Exposure = EXCLUDED.Exposure, ExposedPort = EXCLUDED.ExposedPort"
	_, err := db.Exec(context.Background(), finalStr, values...)
	if err != nil {
		return err
	}

	var query string

	for childIdx, child := range obj.GetExposureInfos() {
		if err := insertIntoDeploymentsPortsExposureInfos(db, child, parent_Id, idx, childIdx); err != nil {
			return err
		}
	}

	query = "delete from deployments_Ports_ExposureInfos where parent_parent_Id = $1 AND parent_idx = $2 AND idx >= $3"
	_, err = db.Exec(context.Background(), query, parent_Id, idx, len(obj.GetExposureInfos()))
	if err != nil {
		return err
	}
	return nil
}

func insertIntoDeploymentsPortsExposureInfos(db *pgxpool.Pool, obj *storage.PortConfig_ExposureInfo, parent_parent_Id string, parent_idx int, idx int) error {

	values := []interface{}{
		// parent primary keys start

		parent_parent_Id,
		parent_idx,
		idx,
		obj.GetLevel(),
		obj.GetServiceName(),
		obj.GetServiceId(),
		obj.GetServiceClusterIp(),
		obj.GetServicePort(),
		obj.GetNodePort(),
		obj.GetExternalIps(),
		obj.GetExternalHostnames(),
	}

	finalStr := "INSERT INTO deployments_Ports_ExposureInfos (parent_parent_Id, parent_idx, idx, Level, ServiceName, ServiceId, ServiceClusterIp, ServicePort, NodePort, ExternalIps, ExternalHostnames) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) ON CONFLICT(parent_parent_Id, parent_idx, idx) DO UPDATE SET parent_parent_Id = EXCLUDED.parent_parent_Id, parent_idx = EXCLUDED.parent_idx, idx = EXCLUDED.idx, Level = EXCLUDED.Level, ServiceName = EXCLUDED.ServiceName, ServiceId = EXCLUDED.ServiceId, ServiceClusterIp = EXCLUDED.ServiceClusterIp, ServicePort = EXCLUDED.ServicePort, NodePort = EXCLUDED.NodePort, ExternalIps = EXCLUDED.ExternalIps, ExternalHostnames = EXCLUDED.ExternalHostnames"
	_, err := db.Exec(context.Background(), finalStr, values...)
	if err != nil {
		return err
	}

	return nil
}

// New returns a new Store instance using the provided sql instance.
func New(db *pgxpool.Pool) Store {
	globaldb.RegisterTable(table, "storage.Deployment")

	createTableDeployments(db)

	// TBD(index creation)

	return &storeImpl{
		db: db,
	}
}

func (s *storeImpl) Upsert(obj *storage.Deployment) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Upsert, "storage.Deployment")

	return insertIntoDeployments(s.db, obj)
}

func (s *storeImpl) UpsertMany(objs []*storage.Deployment) error {
	for _, obj := range objs {
		if err := insertIntoDeployments(s.db, obj); err != nil {
			return err
		}
	}
	return nil
}

// Count returns the number of objects in the store
func (s *storeImpl) Count() (int, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Count, "storage.Deployment")

	row := s.db.QueryRow(context.Background(), countStmt)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// Exists returns if the id exists in the store
func (s *storeImpl) Exists(id string) (bool, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Exists, "storage.Deployment")

	row := s.db.QueryRow(context.Background(), existsStmt, id)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, pgutils.ErrNilIfNoRows(err)
	}
	return exists, nil
}

// Get returns the object, if it exists from the store
func (s *storeImpl) Get(id string) (*storage.Deployment, bool, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Get, "storage.Deployment")

	conn, release := s.acquireConn(ops.Get, "storage.Deployment")
	defer release()

	row := conn.QueryRow(context.Background(), getStmt, id)
	var data []byte
	if err := row.Scan(&data); err != nil {
		return nil, false, pgutils.ErrNilIfNoRows(err)
	}

	var msg storage.Deployment
	if err := proto.Unmarshal(data, &msg); err != nil {
		return nil, false, err
	}
	return &msg, true, nil
}

func (s *storeImpl) acquireConn(op ops.Op, typ string) (*pgxpool.Conn, func()) {
	defer metrics.SetAcquireDuration(time.Now(), op, typ)
	conn, err := s.db.Acquire(context.Background())
	if err != nil {
		panic(err)
	}
	return conn, conn.Release
}

// Delete removes the specified ID from the store
func (s *storeImpl) Delete(id string) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Remove, "storage.Deployment")

	conn, release := s.acquireConn(ops.Remove, "storage.Deployment")
	defer release()

	if _, err := conn.Exec(context.Background(), deleteStmt, id); err != nil {
		return err
	}
	return nil
}

// GetIDs returns all the IDs for the store
func (s *storeImpl) GetIDs() ([]string, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.GetAll, "storage.DeploymentIDs")

	rows, err := s.db.Query(context.Background(), getIDsStmt)
	if err != nil {
		return nil, pgutils.ErrNilIfNoRows(err)
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// GetMany returns the objects specified by the IDs or the index in the missing indices slice
func (s *storeImpl) GetMany(ids []string) ([]*storage.Deployment, []int, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.GetMany, "storage.Deployment")

	conn, release := s.acquireConn(ops.GetMany, "storage.Deployment")
	defer release()

	rows, err := conn.Query(context.Background(), getManyStmt, ids)
	if err != nil {
		if err == pgx.ErrNoRows {
			missingIndices := make([]int, 0, len(ids))
			for i := range ids {
				missingIndices = append(missingIndices, i)
			}
			return nil, missingIndices, nil
		}
		return nil, nil, err
	}
	defer rows.Close()
	elems := make([]*storage.Deployment, 0, len(ids))
	foundSet := make(map[string]struct{})
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, nil, err
		}
		var msg storage.Deployment
		if err := proto.Unmarshal(data, &msg); err != nil {
			return nil, nil, err
		}
		foundSet[msg.GetId()] = struct{}{}
		elems = append(elems, &msg)
	}
	missingIndices := make([]int, 0, len(ids)-len(foundSet))
	for i, id := range ids {
		if _, ok := foundSet[id]; !ok {
			missingIndices = append(missingIndices, i)
		}
	}
	return elems, missingIndices, nil
}

// Delete removes the specified IDs from the store
func (s *storeImpl) DeleteMany(ids []string) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.RemoveMany, "storage.Deployment")

	conn, release := s.acquireConn(ops.RemoveMany, "storage.Deployment")
	defer release()
	if _, err := conn.Exec(context.Background(), deleteManyStmt, ids); err != nil {
		return err
	}
	return nil
}

// Walk iterates over all of the objects in the store and applies the closure
func (s *storeImpl) Walk(fn func(obj *storage.Deployment) error) error {
	rows, err := s.db.Query(context.Background(), walkStmt)
	if err != nil {
		return pgutils.ErrNilIfNoRows(err)
	}
	defer rows.Close()
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return err
		}
		var msg storage.Deployment
		if err := proto.Unmarshal(data, &msg); err != nil {
			return err
		}
		return fn(&msg)
	}
	return nil
}

//// Stubs for satisfying legacy interfaces

// AckKeysIndexed acknowledges the passed keys were indexed
func (s *storeImpl) AckKeysIndexed(keys ...string) error {
	return nil
}

// GetKeysToIndex returns the keys that need to be indexed
func (s *storeImpl) GetKeysToIndex() ([]string, error) {
	return nil, nil
}
