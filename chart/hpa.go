package chart

import (
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s"
	"github.com/lwnmengjing/cd-template-go/imports/k8s"
	"github.com/lwnmengjing/cd-template-go/pkg/config"
)

func NewHpaChart(scope constructs.Construct, id string, props *cdk8s.ChartProps) cdk8s.Chart {
	chart := cdk8s.NewChart(scope, jsii.String(id), props)

	k8s.NewKubeHorizontalPodAutoscalerV2Beta2(chart, jsii.String("hpa"), &k8s.KubeHorizontalPodAutoscalerV2Beta2Props{
		Metadata: &k8s.ObjectMeta{
			Labels: props.Labels,
			Name:   &config.Cfg.Service,
		},
		Spec: &k8s.HorizontalPodAutoscalerSpecV2Beta2{
			MinReplicas: jsii.Number(float64(config.Cfg.Replicas)),
			MaxReplicas: jsii.Number(float64(config.Cfg.MaxReplicas)),
			ScaleTargetRef: &k8s.CrossVersionObjectReferenceV2Beta2{
				Kind:       jsii.String(config.Cfg.WorkloadType),
				Name:       jsii.String(config.Cfg.Service),
				ApiVersion: jsii.String("apps/v1"),
			},
			Metrics: &[]*k8s.MetricSpecV2Beta2{
				{
					Type: jsii.String("Resource"),
					Resource: &k8s.ResourceMetricSourceV2Beta2{
						Name: jsii.String("memory"),
						Target: &k8s.MetricTargetV2Beta2{
							Type:               jsii.String("Utilization"),
							AverageUtilization: jsii.Number(80),
						},
					},
				},
				{
					Type: jsii.String("Resource"),
					Resource: &k8s.ResourceMetricSourceV2Beta2{
						Name: jsii.String("cpu"),
						Target: &k8s.MetricTargetV2Beta2{
							Type:               jsii.String("Utilization"),
							AverageUtilization: jsii.Number(80),
						},
					},
				},
			},
		},
	})

	return chart
}
