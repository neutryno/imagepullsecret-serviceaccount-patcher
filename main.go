package main

import (
	"context"
	"encoding/json"
	"flag"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
	"strings"
	"time"
)

type PrivateRegistrySecret struct {
	name string
}

func main() {

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the client
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	var privateRegistrySecretNames string
	flag.StringVar(&privateRegistrySecretNames, "registrysecretnames", LookupEnv("REGISTRY_SECRET_NAMES"), "Comma separated names of secrets that shall be used as ImagePullSecrets.")

	var privateRegistries []PrivateRegistrySecret // an empty list
	for _, i := range strings.Split(privateRegistrySecretNames, ",") {
		privateRegistries = append(privateRegistries, PrivateRegistrySecret{name: i})
	}

	log.Info("Service-Account Patcher started")
	var patched = 0
	for {
		// get all service accounts in all namespaces
		serviceAccounts, err := client.CoreV1().ServiceAccounts("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		for _, sa := range serviceAccounts.Items {
			log.WithField("serviceaccount", sa.Name).WithField("namesapce", sa.Namespace).
				Debug("Processing ServiceAccount")

			for _, privateReg := range privateRegistries {
				log.WithField("registryname", privateReg.name).Debug("Processing Secret")

				if includeImagePullSecret(&sa, privateReg.name) {
					log.WithField("serviceaccount", sa.Name).
						WithField("namesapce", sa.Namespace).
						WithField("imagePullSecrets", sa.ImagePullSecrets).
						Debug("ServiceAccount has ImagePullSecrets")
				} else {
					log.WithField("serviceaccount", sa.Name).
						WithField("namesapce", sa.Namespace).
						WithField("imagePullSecret", privateReg.name).
						Info("ServiceAccount does not have ImagePullSecret")

					patch, err := getPatchString(&sa, privateReg.name)
					if err != nil {
						panic(err.Error())
					}

					result, err := client.CoreV1().ServiceAccounts(sa.Namespace).
						Patch(context.TODO(), sa.Name, types.StrategicMergePatchType, patch, metav1.PatchOptions{})
					log.WithField("serviceaccount", result.Name).
						WithField("namesapce", sa.Namespace).
						WithField("imagePullSecrets", result.ImagePullSecrets).
						Info("ServiceAccount patched")
					if err != nil {
						panic(err.Error())
					}
					patched++
				}
			}
		}

		log.WithField("numberOfPatchedServiceAccounts", patched).
			Info("Patched ServiceAccounts")
		patched = 0
		time.Sleep(10 * time.Second)
	}
}

// the below is taken from https://github.com/titansoft-pte-ltd/imagepullsecret-patcher
type patch struct {
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`
}

func includeImagePullSecret(sa *corev1.ServiceAccount, secretName string) bool {
	for _, imagePullSecret := range sa.ImagePullSecrets {
		if imagePullSecret.Name == secretName {
			return true
		}
	}
	return false
}
func getPatchString(sa *corev1.ServiceAccount, secretName string) ([]byte, error) {
	saPatch := patch{
		// copy the slice
		ImagePullSecrets: append([]corev1.LocalObjectReference(nil), sa.ImagePullSecrets...),
	}
	if !includeImagePullSecret(sa, secretName) {
		saPatch.ImagePullSecrets = append(saPatch.ImagePullSecrets, corev1.LocalObjectReference{Name: secretName})
	}
	return json.Marshal(saPatch)
}

// LookupEnvOrString lookup ENV string with given key,
func LookupEnv(key string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	panic("Please provide environment variable REGISTRY_SECRET_NAMES, with comma separated secret names.")
}
