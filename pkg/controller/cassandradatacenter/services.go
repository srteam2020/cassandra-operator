package cassandradatacenter

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func createOrUpdateNodesService(rctx *reconciliationRequestContext) (*corev1.Service, error) {
	nodesService := &corev1.Service{ObjectMeta: DataCenterResourceMetadata(rctx.cdc, "nodes")}

	logger := rctx.logger.WithValues("Service.Name", nodesService.Name)

	opresult, err := controllerutil.CreateOrUpdate(context.TODO(), rctx.client, nodesService, func(_ runtime.Object) error {
		nodesService.Spec = corev1.ServiceSpec{
			ClusterIP: "None",
			Ports:     ports{cqlPort, jmxPort}.asServicePorts(),
			Selector:  DataCenterLabels(rctx.cdc),
		}

		if rctx.cdc.Spec.PrometheusSupport {
			nodesService.Spec.Ports = append(nodesService.Spec.Ports, prometheusPort.asServicePort())
		}

		if err := controllerutil.SetControllerReference(rctx.cdc, nodesService, rctx.scheme); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Only log if something has changed
	if opresult != controllerutil.OperationResultNone {
		logger.Info(fmt.Sprintf("Service %s %s.", nodesService.Name, opresult))
	}

	return nodesService, err
}

func createOrUpdateSeedNodesService(rctx *reconciliationRequestContext) (*corev1.Service, error) {
	seedNodesService := &corev1.Service{ObjectMeta: DataCenterResourceMetadata(rctx.cdc, "seeds")}

	logger := rctx.logger.WithValues("Service.Name", seedNodesService.Name)

	opresult, err := controllerutil.CreateOrUpdate(context.TODO(), rctx.client, seedNodesService, func(_ runtime.Object) error {
		seedNodesService.Spec = corev1.ServiceSpec{
			ClusterIP:                "None",
			Ports:                    internodePort.asServicePorts(),
			Selector:                 DataCenterLabels(rctx.cdc),
			PublishNotReadyAddresses: true,
		}

		seedNodesService.Annotations = map[string]string{
			"service.alpha.kubernetes.io/tolerate-unready-endpoints": "true",
		}

		if err := controllerutil.SetControllerReference(rctx.cdc, seedNodesService, rctx.scheme); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Only log if something has changed
	if opresult != controllerutil.OperationResultNone {
		logger.Info(fmt.Sprintf("Service %s %s.", seedNodesService.Name, opresult))
	}

	return seedNodesService, err
}