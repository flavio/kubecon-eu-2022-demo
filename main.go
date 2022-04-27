package main

import (
	demo "github.com/saschagrunert/demo"
	"os/exec"
)

const DEMO_NS = "test"

func main() {
	cleanupNamespace()
	if err := setupNamespace(); err != nil {
		panic(err)
	}

	d := demo.New()
	d.Add(k3sWasmRun(), "demo", "k3s running wasm demo")
	d.Run()

}

func k3sWasmRun() *demo.Run {
	r := demo.NewRun(
		"Kubernetes API server 💖 WebAssembly",
	)

	r.Step(demo.S(
		"A simple k3s cluster",
	), demo.S("kubectl get nodes"))
	r.Step(nil,
		demo.S("kubectl get pods -A"))

	r.Step(demo.S(
		"The configuration of the Kubewarden feature gate",
	), demo.S("bat /etc/rancher/k3s/admission-control-config.yaml"))

	r.Step(demo.S(
		"A privileged Pod",
	), demo.S("bat privileged.yaml"))
	r.StepCanFail(demo.S(
		"Attempt to create a privileged Pod inside of the default namespace",
	), demo.S("kubectl apply -n default -f privileged.yaml"))

	r.Step(demo.S(
		"A Pod using a container image pulled from the Docker Hub",
	), demo.S("bat alpine.yaml"))
	r.StepCanFail(demo.S(
		"Attempt to create a Pod that uses a container image from the Docker Hub",
	), demo.S("kubectl apply -n default -f alpine.yaml"))

	r.Step(demo.S(
		"Deploy the same privileged Pod in a namespace where the policies are not enforced",
	), demo.S("kubectl apply -n", DEMO_NS, "-f privileged.yaml"))

	return r
}

func setupNamespace() error {
	exec.Command("kubectl", "create", "namespace", DEMO_NS).Run()
	return nil
}

func cleanupNamespace() error {
	exec.Command("kubectl", "delete", "namespace", DEMO_NS).Run()
	return nil
}
