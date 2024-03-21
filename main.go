package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	webhookURL        = os.Getenv("WEBHOOK_URL")
	debug             = os.Getenv("DEBUG") == "true"
	intervalInMinutes = 10
)

type Report struct {
	Problems []Problem
}

type Problem struct {
	Summary     string
	Description string
}

func init() {
	var err error
	if os.Getenv("INTERVAL_IN_MINUTES") != "" {
		intervalInMinutes, err = strconv.Atoi(os.Getenv("INTERVAL_IN_MINUTES"))
		if err != nil {
			panic(fmt.Errorf("failed to parse environment variable INTERVAL_IN_MINUTES: %w", err))
		}
		if intervalInMinutes < 1 {
			log.Println("[init] INTERVAL_IN_MINUTES must be 1 or higher, defaulting to 1")
			intervalInMinutes = 1
		}
	}
	if len(webhookURL) == 0 {
		panic("environment variable WEBHOOK_URL is not defined")
	}
}

func main() {
	kubernetesClient, err := CreateClients()
	if err != nil {
		panic(fmt.Errorf("failed to initialize Kubernetes client: %w", err))
	}
	for {
		report := Report{}
		checkPods(kubernetesClient, &report)
		SendNotification(report)
		time.Sleep(time.Duration(intervalInMinutes) * time.Minute)
	}
}

func checkPods(kubernetesClient kubernetes.Interface, report *Report) {
	// TODO: Support pagination
	pods, err := kubernetesClient.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{ResourceVersion: "0"})
	if err != nil {
		log.Println("[checkPods] Failed to retrieve pods:", err.Error())
		return
	}
	for _, pod := range pods.Items {
		minutesSincePodWasCreated := time.Since(pod.GetCreationTimestamp().Time).Minutes()
		// If the pod has been pending for more than 2 minutes, we need to send an alert
		if pod.Status.Phase == v1.PodPending && minutesSincePodWasCreated > 2 {
			log.Printf("[checkPods] Pod with name=%s is in state=%s with reason=%s and message=%s", pod.Name, pod.Status.Phase, pod.Status.Reason, pod.Status.Message)
			problem := Problem{
				Summary: fmt.Sprintf("Pod `%s` in `%s` is stuck in a `Pending` state", pod.GetName(), pod.GetNamespace()),
			}
			// Try to get more information on why it's in a pending state if possible
			for _, containerStatus := range pod.Status.ContainerStatuses {
				// If the container is not ready, it's either waiting or terminated, so we need to send an alert
				if !containerStatus.Ready {
					state, reason, message := extractNameReasonMessageFromContainerState(containerStatus.State)
					log.Printf("[checkPods] Pod with name=%s has container=%s in state=%s for reason=%s with message=%s", pod.Name, containerStatus.Name, state, reason, message)
					problem.Description += fmt.Sprintf("Container `%s` is in state `%s` because of reason `%s`:\n```%s```", containerStatus.Name, state, reason, message)
				}
			}
			report.Problems = append(report.Problems, problem)
		} else if pod.Status.Phase == v1.PodRunning {
			for _, containerStatus := range pod.Status.ContainerStatuses {
				// If the container is not ready, it's either waiting or terminated, so we need to send an alert
				if !containerStatus.Ready {
					state, reason, message := extractNameReasonMessageFromContainerState(containerStatus.State)
					if containerStatus.RestartCount > 0 || (state == "Waiting" && (reason == "ErrImagePull" || reason == "ImagePullBackOff")) {
						log.Printf("[checkPods] Pod with name=%s has a container in state=%s for reason=%s with message=%s", pod.Name, state, reason, message)
						report.Problems = append(report.Problems, Problem{
							Summary:     fmt.Sprintf("Pod `%s` in `%s` has restarted `%d` times", pod.GetName(), pod.GetNamespace(), containerStatus.RestartCount),
							Description: fmt.Sprintf("Container `%s` is in state `%s` because of reason `%s`:\n```%s```", containerStatus.Name, state, reason, message),
						})
					}
				}
			}
		}
	}
}

func extractNameReasonMessageFromContainerState(state v1.ContainerState) (string, string, string) {
	if state.Waiting != nil {
		return "Waiting", state.Waiting.Reason, state.Waiting.Message
	} else if state.Running != nil {
		return "Running", "", ""
	} else if state.Terminated != nil {
		return "Terminated", state.Terminated.Reason, fmt.Sprintf("%s (exit code %d)", state.Terminated.Message, state.Terminated.ExitCode)
	}
	return "Unknown", "Unknown", "None"
}
