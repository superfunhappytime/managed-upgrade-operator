package servicelog

import (
	"fmt"
	"time"

	"github.com/openshift/managed-upgrade-operator/pkg/metrics"
)

type metricsWrapper struct {
	*metrics.Counter
	*serviceLogger
}

func (s *metricsWrapper) UpdateMetricClusterCheckFailed(upgradeConfigName string) {
	s.Counter.UpdateMetricClusterCheckFailed(upgradeConfigName)
	s.CreateLog("Upgrade failed", fmt.Sprintf("Failed cluster health check - %s", upgradeConfigName))
}

// UpdateMetricValidationFailed
// called on failure to process a CR
// You probably wouldn't want to create a service log for this
//func (s *metricsWrapper) UpdateMetricValidationFailed(upgradeConfigName string) {
//	s.Counter.UpdateMetricValidationFailed(upgradeConfigName)
//	s.CreateLog("Upgrade failed", fmt.Sprintf("Failed validation of an upgrade config - %s", upgradeConfigName))
//}

func (s *metricsWrapper) UpdateMetricValidationSucceeded(upgradeConfigName string) {
	s.Counter.UpdateMetricValidationSucceeded(upgradeConfigName)
	s.CreateLog("Upgrade starting", fmt.Sprintf("Passed validation of an upgrade config - %s", upgradeConfigName))
}

func (s *metricsWrapper) UpdateMetricClusterCheckSucceeded(upgradeConfigName string) {
	s.Counter.UpdateMetricClusterCheckSucceeded(upgradeConfigName)
	s.CreateLog("Upgrade starting", fmt.Sprintf("Passed cluster health check - %s", upgradeConfigName))
}

func (s *metricsWrapper) UpdateMetricUpgradeStartTime(time time.Time, upgradeConfigName string, version string) {
	s.Counter.UpdateMetricUpgradeStartTime(time, upgradeConfigName, version)
	s.CreateLog("Upgrade progressing", fmt.Sprintf("Upgrade started at %s - %s", time, upgradeConfigName))
}

func (s *metricsWrapper) UpdateMetricControlPlaneEndTime(time time.Time, upgradeConfigName string, version string) {
	s.Counter.UpdateMetricControlPlaneEndTime(time, upgradeConfigName, version)
	s.CreateLog("Upgrade progressing", fmt.Sprintf("Upgrade of control plane completed at %s - %s", time, upgradeConfigName))
}

func (s *metricsWrapper) UpdateMetricNodeUpgradeEndTime(time time.Time, upgradeConfigName string, version string) {
	s.Counter.UpdateMetricNodeUpgradeEndTime(time, upgradeConfigName, version)
	s.CreateLog("Upgrade progressing", fmt.Sprintf("Upgrade of worker nodes completed at %s - %s", time, upgradeConfigName))
}

func (s *metricsWrapper) UpdateMetricClusterVerificationFailed(upgradeConfigName string) {
	s.Counter.UpdateMetricClusterVerificationFailed(upgradeConfigName)
	s.CreateLog("Upgrade failed", fmt.Sprintf("Post-upgrade verification has failed - %s", upgradeConfigName))
}

func (s *metricsWrapper) UpdateMetricClusterVerificationSucceeded(upgradeConfigName string) {
	s.Counter.UpdateMetricClusterVerificationSucceeded(upgradeConfigName)
	s.CreateLog("Upgrade progressing", fmt.Sprintf("Passed post-upgrade verification - %s", upgradeConfigName))
}
