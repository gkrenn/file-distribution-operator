package controllers

import (
	"context"
	"fmt"
	"path"
	"strings"
	"time"

	krenndevv1alpha1 "github.com/gkrenn/file-distribution-operator/api/v1alpha1"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	sourceVolumePath      = "/source"
	destinationVolumePath = "/destination"

	defaultFilePermissions = "644"
)

var jobLog *zap.SugaredLogger

func init() {
	// create new zap json logger
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	jobLog = logger.Sugar()
}

type DistributionJob struct {
	client.Client
	ctx context.Context

	hostPathDir      string
	hostPathFileName string

	Namespace              string
	FileDistributionConfig krenndevv1alpha1.FileDistributionConfig
}

func NewDistributionJob(client client.Client, ctx context.Context, namespace string, fileDistributionConfig krenndevv1alpha1.FileDistributionConfig) DistributionJob {
	var hostPathDir string
	var hostPathFileName string

	if fileDistributionConfig.Spec.FilePermissions == "" {
		fileDistributionConfig.Spec.FilePermissions = defaultFilePermissions
	}

	// destination can be a directory or a file
	if isStringDirectory(fileDistributionConfig.Spec.Destination) {
		jobLog.Info("destination is a directory")
		hostPathDir = fileDistributionConfig.Spec.Destination
		hostPathFileName = path.Base(fileDistributionConfig.Spec.FileName)
	} else {
		jobLog.Info("destination is a file")
		hostPathDir = path.Dir(fileDistributionConfig.Spec.Destination)
		hostPathFileName = path.Base(fileDistributionConfig.Spec.Destination)
	}

	return DistributionJob{
		Client:                 client,
		ctx:                    ctx,
		Namespace:              namespace,
		FileDistributionConfig: fileDistributionConfig,
		hostPathDir:            hostPathDir,
		hostPathFileName:       hostPathFileName,
	}
}

func (distributionJob DistributionJob) BuildKubernetesJob(nodeName string) *batchv1.Job {
	// todo set mode and ownership
	// todo mount actual node filesystem

	sourceVolumeName := "file-mnt"
	destinationVolumeName := "node-mnt"

	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "file-distribution-job-" + nodeName,
			Namespace: distributionJob.Namespace,
			Labels: map[string]string{
				"app": "file-distribution",
			},
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					NodeName: nodeName,
					Containers: []corev1.Container{
						{
							Name:    "file-distribution",
							Image:   "ubuntu:latest",
							Command: distributionJob.buildCommand(),
							// set to true to allow for interactive shell i.e. debugging
							TTY:   false,
							Stdin: false,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      sourceVolumeName,
									MountPath: sourceVolumePath,
								},
								{
									Name:      destinationVolumeName,
									MountPath: destinationVolumePath,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: sourceVolumeName,
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: distributionJob.FileDistributionConfig.Spec.SecretName,
									Optional:   &[]bool{false}[0],
									Items: []corev1.KeyToPath{
										{
											Key:  distributionJob.FileDistributionConfig.Spec.FileName,
											Path: distributionJob.FileDistributionConfig.Spec.FileName,
										},
									},
								},
							},
						},
						{
							Name: destinationVolumeName,
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: distributionJob.hostPathDir,
								},
							},
						},
					},
					// only run once
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
			// prevent job from restarting
			BackoffLimit: &[]int32{0}[0],
		},
	}
}

func isStringDirectory(path string) bool {
	return strings.HasSuffix(path, "/")
}

func (distributionJob DistributionJob) setupJobOnAllNodes() ([]corev1.Node, error) {
	// get list of nodes from kubernetes
	var nodelist corev1.NodeList
	err := distributionJob.List(context.TODO(), &nodelist)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, node := range nodelist.Items {
		// create job
		job := distributionJob.BuildKubernetesJob(node.Name)
		err = distributionJob.Create(context.TODO(), job)
		if err != nil {
			return nodelist.Items, errors.WithStack(err)
		}
		jobLog.Info("created job for node: ", node.Name)
	}

	return nodelist.Items, nil
}

func (distributionJob DistributionJob) deleteAllJobsInNamespace() error {
	return errors.WithStack(distributionJob.DeleteAllOf(distributionJob.ctx, &batchv1.Job{}, client.InNamespace(distributionJob.Namespace), client.PropagationPolicy(metav1.DeletePropagationBackground)))
}

func (distributionJob DistributionJob) waitForAllJobsToBeCompleted() error {
	// could be implemented with watch => did not work during development => switched to polling for now

	for i := 0; i < 10; i++ {
		var jobList batchv1.JobList
		err := distributionJob.List(context.TODO(), &jobList, client.InNamespace(distributionJob.Namespace))
		if err != nil {
			return errors.WithStack(err)
		}

		// check if all jobs are completed
		allJobsCompleted := true
		for _, job := range jobList.Items {
			if job.Status.Succeeded == 0 {
				allJobsCompleted = false
			}
		}

		if allJobsCompleted {
			return nil
		}

		time.Sleep(10 * time.Second)
	}

	err := fmt.Errorf("timeout while waiting for jobs to be completed")
	jobLog.Error(err, err.Error())

	return err
}

func (distributionJob DistributionJob) buildCommand() []string {
	// kubernetes ensure that source and destination exist
	source := path.Join(sourceVolumePath, distributionJob.FileDistributionConfig.Spec.FileName)
	destination := path.Join(destinationVolumePath, distributionJob.hostPathFileName)
	return []string{
		"/bin/bash",
		"-c",
		fmt.Sprintf("cp -f '%s' '%s' && chmod %s '%s'", source, destination, distributionJob.FileDistributionConfig.Spec.FilePermissions, destination),
	}
}
