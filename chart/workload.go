/*
 * @Author: lwnmengjing
 * @Date: 2021/10/29 11:21 下午
 * @Last Modified by: lwnmengjing
 * @Last Modified time: 2021/10/29 11:21 下午
 */

package chart

import (
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s"
	"github.com/lwnmengjing/cd-template/imports/k8s"
	"github.com/lwnmengjing/cd-template/pkg/config"
	"os"
	"strconv"
)

func NewWorkloadChart(scope constructs.Construct, id string, props *cdk8s.ChartProps) cdk8s.Chart {
	chart := cdk8s.NewChart(scope, jsii.String(id), props)
	ports := make([]*k8s.ContainerPort, 0)
	//port
	for i := range config.Cfg.Ports {
		ports = append(ports, &k8s.ContainerPort{
			ContainerPort: jsii.Number(float64(config.Cfg.Ports[i].Port)),
		})
	}
	//env
	env := make([]*k8s.EnvVar, 0)
	for i := range config.Cfg.ImportEnvNames {
		if config.Cfg.ImportEnvNames[i] == "" {
			continue
		}
		v := os.Getenv(config.Cfg.ImportEnvNames[i])
		env = append(env, &k8s.EnvVar{
			Name:  &config.Cfg.ImportEnvNames[i],
			Value: &v,
		})
	}
	env = append(env, &k8s.EnvVar{
		Name: jsii.String("NODE_NAME"),
		ValueFrom: &k8s.EnvVarSource{
			FieldRef: &k8s.ObjectFieldSelector{
				FieldPath: jsii.String("metadata.name"),
			},
		},
	}, &k8s.EnvVar{
		Name: jsii.String("STAGE"),
		ValueFrom: &k8s.EnvVarSource{
			FieldRef: &k8s.ObjectFieldSelector{
				FieldPath: jsii.String("metadata.namespace"),
			},
		},
	})
	//config
	volumeMounts := make([]*k8s.VolumeMount, 0)
	volumes := make([]*k8s.Volume, 0)
	if len(config.Cfg.Config) > 0 {
		readOnly := true
		for i := range config.Cfg.Config {
			volumes = append(volumes, &k8s.Volume{
				Name: &config.Cfg.Config[i].Name,
				ConfigMap: &k8s.ConfigMapVolumeSource{
					Name: &config.Cfg.Config[i].Name,
				},
			})

			volumeMounts = append(volumeMounts, &k8s.VolumeMount{
				MountPath: &config.Cfg.Config[i].Path,
				Name:      &config.Cfg.Config[i].Name,
				ReadOnly:  &readOnly,
			})
		}
	}

	var serviceAccountName *string
	if config.Cfg.ServiceAccount {
		serviceAccountName = jsii.String(config.Cfg.App + "-" + config.Cfg.Service)
	}
	if config.Cfg.ServiceAccountName != "" {
		serviceAccountName = jsii.String(config.Cfg.ServiceAccountName)
	}
	var command *[]*string
	if len(config.Cfg.Command) > 0 {

		command = &config.Cfg.Command
	}
	var args *[]*string
	if len(config.Cfg.Args) > 0 {
		args = &config.Cfg.Args
	}

	var resources k8s.ResourceRequirements
	if len(config.Cfg.Resources) > 0 {
		for k, r := range config.Cfg.Resources {
			switch k {
			case "limits":
				resources.Limits = &map[string]k8s.Quantity{
					"cpu":    k8s.Quantity_FromString(&r.CPU),
					"memory": k8s.Quantity_FromString(&r.Memory),
				}
			case "requests":
				resources.Requests = &map[string]k8s.Quantity{
					"cpu":    k8s.Quantity_FromString(&r.CPU),
					"memory": k8s.Quantity_FromString(&r.Memory),
				}
			}
		}
	}
	annotations := make(map[string]*string)
	if config.Cfg.Metrics.Scrape {
		annotations["prometheus.io/scrape"] = jsii.String("true")
		annotations["prometheus.io/port"] = jsii.String(strconv.Itoa(int(config.Cfg.Metrics.Port)))
		annotations["prometheus.io/path"] = jsii.String(config.Cfg.Metrics.Path)
	}
	var replicas *float64
	if !config.Cfg.Hpa {
		replicas = jsii.Number(float64(config.Cfg.Replicas))
	}
	switch config.Cfg.WorkloadType {
	case "statefulset":
		k8s.NewKubeStatefulSet(chart, jsii.String("statefulset"), &k8s.KubeStatefulSetProps{
			Metadata: &k8s.ObjectMeta{
				Name:   &config.Cfg.Service,
				Labels: props.Labels,
			},
			Spec: &k8s.StatefulSetSpec{
				ServiceName: &config.Cfg.Service,
				Replicas:    replicas,
				Selector: &k8s.LabelSelector{
					MatchLabels: props.Labels,
				},
				Template: &k8s.PodTemplateSpec{
					Metadata: &k8s.ObjectMeta{
						Labels:      props.Labels,
						Annotations: &annotations,
					},
					Spec: &k8s.PodSpec{
						ServiceAccountName: serviceAccountName,
						Containers: &[]*k8s.Container{{
							Name:         jsii.String(config.Cfg.Service),
							Image:        jsii.String(config.Cfg.Image.String()),
							Ports:        &ports,
							Env:          &env,
							VolumeMounts: &volumeMounts,
							Command:      command,
							Args:         args,
							Resources:    &resources,
						}},
						Volumes: &volumes,
					},
				},
			},
		})
	default:
		k8s.NewKubeDeployment(chart, jsii.String("deployment"), &k8s.KubeDeploymentProps{
			Metadata: &k8s.ObjectMeta{
				Name:   &config.Cfg.Service,
				Labels: props.Labels,
			},
			Spec: &k8s.DeploymentSpec{
				Replicas: replicas,
				Selector: &k8s.LabelSelector{
					MatchLabels: props.Labels,
				},
				Template: &k8s.PodTemplateSpec{
					Metadata: &k8s.ObjectMeta{
						Labels:      props.Labels,
						Annotations: &annotations,
					},
					Spec: &k8s.PodSpec{
						ServiceAccountName: serviceAccountName,
						Containers: &[]*k8s.Container{{
							Name:         jsii.String(config.Cfg.Service),
							Image:        jsii.String(config.Cfg.Image.String()),
							Ports:        &ports,
							Env:          &env,
							VolumeMounts: &volumeMounts,
							Command:      command,
							Args:         args,
							Resources:    &resources,
						}},
						Volumes: &volumes,
					},
				},
			},
		})
	}
	return chart
}
