package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	K8S_LABEL_POD_TO_DEL = "delete-pod"
)

var clientset *kubernetes.Clientset

func main() {
	log.Print("Shared Informer app started")

	// TODO: creates the in-cluster config
	// kubeconfig := os.Getenv("KUBECONFIG")
	// config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	// if err != nil {
	// 	log.Panic(err.Error())
	// }

	config, err := rest.InClusterConfig()
	if err != nil {
		var kubeconfig *string
		if home := homeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		// use the current context in kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}

	// create clientset upon the specific config local vs in-cluster
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic(err.Error())
	}

	factory := informers.NewSharedInformerFactory(clientset, 0)
	informer := factory.Core().V1().Pods().Informer()
	stopper := make(chan struct{})
	defer close(stopper)
	defer runtime.HandleCrash()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    onAdd,
		UpdateFunc: onUpdate,
	})
	go informer.Run(stopper)
	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}
	<-stopper
}

// onAdd is the function executed when the kubernetes informer notified the
// presence of a new kubernetes node in the cluster
func onAdd(obj interface{}) {
	// Cast the obj as node
	pod := obj.(*corev1.Pod)
	_, ok := pod.GetLabels()[K8S_LABEL_POD_TO_DEL]
	if ok {
		fmt.Printf(
			"The label %s will not be honored on creation, you must update the pod to take effect!\n",
			K8S_LABEL_POD_TO_DEL)
	} else {
		fmt.Printf(
			"The pod %s/%s does not have the label %s\n",
			pod.GetNamespace(), pod.GetName(), K8S_LABEL_POD_TO_DEL)
	}
}

func onUpdate(oldObj interface{}, newObj interface{}) {
	// Cast the obj as node
	pod := newObj.(*corev1.Pod)
	_, ok := pod.GetLabels()[K8S_LABEL_POD_TO_DEL]
	if ok {
		fmt.Printf("Pod has been labeleld w/ %s so it will be deleted\n", K8S_LABEL_POD_TO_DEL)
		go deletePod(pod)
	} else {
		fmt.Printf("label %s is missing\n", K8S_LABEL_POD_TO_DEL)
	}
}

func deletePod(pod *corev1.Pod) {
	clientset.CoreV1().Pods(pod.GetNamespace()).Delete(pod.GetName(), &metav1.DeleteOptions{})
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
