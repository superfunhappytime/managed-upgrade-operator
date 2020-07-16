// Package servicelog is a hacked together POC to integration between an operator and OCM Service Log API
//
// This POC is a quick and dirty use of the OCM sdk to emit status updates using the OCM Service Log API
// It requires the following environment variables defined:
// - `OCM_TOKEN` == ocm login token
// - `OCM_API_URL` == ocm api endpoint in use
//
package servicelog

import (
	"context"
	"fmt"
	"os"

	sdk "github.com/openshift-online/ocm-sdk-go"
	servicelogsv1 "github.com/openshift-online/ocm-sdk-go/servicelogs/v1"
	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/managed-upgrade-operator/pkg/metrics"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ServiceLogger provides log creation
type ServiceLogger interface {
	AttachToMetricsClient(c client.Client, m metrics.Metrics) (metrics.Metrics, error)
	CreateLog(summary, description string) error
	SetUpgradeEventID(id string)
}

type serviceLogger struct {
	scb         *sdk.ConnectionBuilder
	clusterUUID configv1.ClusterID
	eventID     string
	serviceName string
}

// NewServiceLogClient creates a new instance of ServiceLogger
func NewServiceLogClient(c client.Client, serviceName string) (ServiceLogger, error) {
	// Query for the cluster uuid
	cvList := &configv1.ClusterVersionList{}
	if err := c.List(context.TODO(), cvList); err != nil {
		return nil, err
	}
	if len(cvList.Items) == 0 {
		return nil, fmt.Errorf("Unable to find cluster version resource")
	}
	clusterUUID := cvList.Items[0].Spec.ClusterID

	// create debug logger for POC
	logger, err := sdk.NewGoLoggerBuilder().
		Debug(true).
		Build()
	if err != nil {
		return nil, fmt.Errorf("Can't build logger: %v", err)
	}
	// Create the connection buidler
	token := os.Getenv("OCM_TOKEN")
	connectionBuilder := sdk.NewConnectionBuilder().
		Logger(logger).
		Tokens(token).
		URL(os.Getenv("OCM_API_URL"))

	return &serviceLogger{
		scb:         connectionBuilder,
		clusterUUID: clusterUUID,
		serviceName: serviceName,
	}, nil
}

// AttachToMetricsClient wraps a Metric interface allowing a ServiceLogger to intercept calls
// expects m metrics.Metrics to be an *metrics.Counter
func (s *serviceLogger) AttachToMetricsClient(c client.Client, m metrics.Metrics) (metrics.Metrics, error) {
	if cntr, ok := m.(*metrics.Counter); ok {
		return &metricsWrapper{
			Counter:       cntr,
			serviceLogger: s,
		}, nil
	}

	return nil, fmt.Errorf("Unexpected Metrics type")
}

// SetUpgradeEventID to set the current eventID written into the summary field
func (s *serviceLogger) SetUpgradeEventID(id string) {
	s.eventID = id
}

// CreateLog will call the OSL api endpoint and create a log entry
// summary, description as defined and will set all other required fields
//
// really would like the output similar to this
// #id upgrade from x.y.z to x2.y2.z2
// for now this poc replaces the summary with:
// #id - summary
func (s *serviceLogger) CreateLog(summary, description string) error {
	ctx := context.Background()
	connection, err := s.scb.BuildContext(ctx)
	if err != nil {
		return err
	}
	// make sure to close connection
	defer connection.Close()

	logsCollection := connection.ServiceLogs().V1().ClusterLogs()

	logResource := logsCollection.Add()
	logEntry, err := servicelogsv1.NewLogEntry().
		ClusterUUID(fmt.Sprintf("%s", s.clusterUUID)).
		InternalOnly(true).
		Description(description).
		ServiceName(s.serviceName).
		Summary(fmt.Sprintf("Event #%s - %s", s.eventID, summary)).Build()

	cLogResp, err := logResource.Body(logEntry).SendContext(ctx)
	if err != nil {
		return err
	}
	if cLogResp.Status() != 201 {
		return fmt.Errorf("Unable to create service log")
	}

	return nil
}
