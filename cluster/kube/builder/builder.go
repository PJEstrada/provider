package builder

import (
	"fmt"
	"strconv"

	mani "github.com/akash-network/akash-api/go/manifest/v2beta2"
	"github.com/tendermint/tendermint/libs/log"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/util/intstr"

	mtypes "github.com/akash-network/akash-api/go/node/market/v1beta3"

	clusterUtil "github.com/akash-network/provider/cluster/util"
	crd "github.com/akash-network/provider/pkg/apis/akash.network/v2beta2"
)

const (
	AkashManagedLabelName         = "akash.network"
	AkashManifestServiceLabelName = "akash.network/manifest-service"
	AkashNetworkStorageClasses    = "akash.network/storageclasses"
	AkashServiceTarget            = "akash.network/service-target"
	AkashMetalLB                  = "metal-lb"
	akashDeploymentPolicyName     = "akash-deployment-restrictions"

	akashNetworkNamespace       = "akash.network/namespace"
	AkashLeaseOwnerLabelName    = "akash.network/lease.id.owner"
	AkashLeaseDSeqLabelName     = "akash.network/lease.id.dseq"
	AkashLeaseGSeqLabelName     = "akash.network/lease.id.gseq"
	AkashLeaseOSeqLabelName     = "akash.network/lease.id.oseq"
	AkashLeaseProviderLabelName = "akash.network/lease.id.provider"
	AkashLeaseManifestVersion   = "akash.network/manifest.version"
)

const (
	runtimeClassNoneValue = "none"
	runtimeClassNvidia    = "nvidia"
)

const (
	envVarAkashGroupSequence         = "AKASH_GROUP_SEQUENCE"
	envVarAkashDeploymentSequence    = "AKASH_DEPLOYMENT_SEQUENCE"
	envVarAkashOrderSequence         = "AKASH_ORDER_SEQUENCE"
	envVarAkashOwner                 = "AKASH_OWNER"
	envVarAkashProvider              = "AKASH_PROVIDER"
	envVarAkashClusterPublicHostname = "AKASH_CLUSTER_PUBLIC_HOSTNAME"
)

var (
	dnsPort     = intstr.FromInt(53)
	udpProtocol = corev1.Protocol("UDP")
	tcpProtocol = corev1.Protocol("TCP")
)

type builderBase interface {
	NS() string
	Name() string
	Validate() error
}

type builder struct {
	log      log.Logger
	settings Settings
	lid      mtypes.LeaseID
	group    *mani.Group
	sparams  crd.ParamsServices
}

var _ builderBase = (*builder)(nil)

func (b *builder) NS() string {
	return LidNS(b.lid)
}

func (b *builder) Name() string {
	return b.NS()
}

func (b *builder) labels() map[string]string {
	return map[string]string{
		AkashManagedLabelName: "true",
		akashNetworkNamespace: LidNS(b.lid),
	}
}

func (b *builder) Validate() error {
	return nil
}

func addIfNotPresent(envVarsAlreadyAdded map[string]int, env []corev1.EnvVar, key string, value interface{}) []corev1.EnvVar {
	_, exists := envVarsAlreadyAdded[key]
	if exists {
		return env
	}

	env = append(env, corev1.EnvVar{Name: key, Value: fmt.Sprintf("%v", value)})
	return env
}

const SuffixForNodePortServiceName = "-np"

func makeGlobalServiceNameFromBasename(basename string) string {
	return fmt.Sprintf("%s%s", basename, SuffixForNodePortServiceName)
}

// LidNS generates a unique sha256 sum for identifying a provider's object name.
func LidNS(lid mtypes.LeaseID) string {
	return clusterUtil.LeaseIDToNamespace(lid)
}

func AppendLeaseLabels(lid mtypes.LeaseID, labels map[string]string) map[string]string {
	labels[AkashLeaseOwnerLabelName] = lid.Owner
	labels[AkashLeaseDSeqLabelName] = strconv.FormatUint(lid.DSeq, 10)
	labels[AkashLeaseGSeqLabelName] = strconv.FormatUint(uint64(lid.GSeq), 10)
	labels[AkashLeaseOSeqLabelName] = strconv.FormatUint(uint64(lid.OSeq), 10)
	labels[AkashLeaseProviderLabelName] = lid.Provider
	return labels
}
