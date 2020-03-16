# Kubernetes | Minikube Deployment + Service + Horizontal Autoscaler

In this post we are going to learn about how deploy containerized app into kubernetes cluster, enable the horizontal auto scaling on it, and create a service that make the application accessible from outside the cluster.

The application that we are going to use on the tutorial is a simple hello world app written in Go. The app is dockerized and the image is available on [Docker Hub](https://hub.docker.com/repository/docker/novalagung/hello-world).

You can also deploy your own app, just do push it into Docker Hub. This guide might help you [Docker - Push Image to hub.docker.com](/docker-push-image-to-hub.html).

## 1. Prerequisites

#### 1.1. Docker engine

Ensure Docker engine is running. If you haven't install it, then follow this guide [Docker Installation](docker-installation.md).

#### Ensure minikube is running

Run the minikube using command below. Do execute the command on PowerShell with admin privilege.

```bash
minikube start
```

If you haven't install minikube, do install it first. Follow this guide https://kubernetes.io/docs/setup/learning-environment/minikube/.

#### Ensure the `kubectl` command is available

If you haven't install `kubectl` then follow this guide https://kubernetes.io/docs/tasks/tools/install-kubectl/.

## 3. Guide

#### 3.1. For Windows user only, run PowerShell with admin privilege

CMD won't be helpful here. Run the PowerShell under administrator privilege.

#### 3.2. Create the kubernetes configuration file (in `.yml` format)

Configuration file defines the configuration for a Kubernetes object. Create one, name it `k8s.yaml` or it can be other name, it is fine. Open the file using any editor. Next we shall begin defining the configurations.

We are going to create a deployment config, service config, and auto scaler config. All of them will be written in one file (`k8s.yaml`).

#### 3.3. Deployment config


deployments represent a set of multiple, identical Pods with no unique identities.
a deployment runs multiple replicas of your application and automatically replaces any instances that fail or become unresponsive.
in this way, Deployments help ensure that one or more instances of your application are available to serve user requests.
deployments are managed by the Kubernetes Deployment controller.
deployments use a Pod template, which contains a specification for its Pods.
the Pod specification determines how each Pod should look like:
- what applications should run inside its containers
- which volumes the Pods should mount
- its labels, and more.

```yaml
---
# ======================== deployment
#
# deployments represent a set of multiple, identical Pods with no unique identities.
# a deployment runs multiple replicas of your application and automatically replaces any instances that fail or become unresponsive.
# in this way, Deployments help ensure that one or more instances of your application are available to serve user requests.
# deployments are managed by the Kubernetes Deployment controller.
#
# deployments use a Pod template, which contains a specification for its Pods.
# the Pod specification determines how each Pod should look like:
# - what applications should run inside its containers
# - which volumes the Pods should mount
# - its labels, and more.

# there are a lot of apis available in kubernetes (try in cli: `kubectl api-versions`).
# for this block of deployment code, we will use `apps/v1`.
apiVersion: apps/v1

# book this block of yaml for Deployment.
kind: Deployment

# name it `my-app-deployment`.
metadata:
  name: my-app-deployment

# specification of the desired behavior of the Deployment.
spec:

  # selector.matchLabels basically used to determine which pods are managed by the deployment.
  # this deployment will manage all pods that have labels matching the selector.
  #
  # refer to line 41 or 46
  selector:
    matchLabels:
      app: my-app

  # template describes the pods that will be created.
  #
  # pods are the smallest, most basic deployable objects in Kubernetes. A Pod represents a single instance of a running process in your cluster.
  # pods contain one or more containers, such as Docker containers.
  template:

    # put label on the pods as `my-app`.
    metadata:
      labels:
        app: my-app

    # specification of the desired behavior of the `my-app` pod.
    spec:

      # list of containers belonging to the `my-app` pod.
      # a single pod might contains multiple containers.
      containers:

          # allocate a container, name it as `hello-world`.
        - name: hello-world

          # the container image is on docker hub repo `novalagung/hello-world`.
          # if the particular image is not available locally, then it'll be pulled first.
          image: novalagung/hello-world

          # set the env vars during container build process.
          # for more details take a look at https://hub.docker.com/repository/docker/novalagung/hello-world.
          env:
            - name: PORT
              value: "8080"
            - name: INSTANCE_ID
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name

          # this pod only have one container (`hello-world`), and what this container do is  start a webserver that listen to port `8080`.
          # the port need to be exported, to make it accessible between the pods within the cluster.
          ports:
            - containerPort: 8080

          # compute resources required by this container `hello-world`.
          resources:
            limits:
              cpu: 250m
              memory: 32Mi

---
# ======================== service
#
# the idea of a Service is to group a set of Pod endpoints into a single resource.
# by default, you get a stable cluster IP address that clients inside the cluster can use to contact Pods in the Service.
# a client sends a request to the stable IP address, and the request is routed to one of the Pods in the Service.

# pick api version `v1` for service.
apiVersion: v1

# book this block of yaml for Service.
kind: Service

# name it `my-service`.
metadata:
  name: my-service

# spec of the desired behavior of the service.
spec:

  # pick LoadBalancer as the type of the service.
  # a LoadBalancer service is the standard way to expose a service to the internet.
  # this will spin up a Network Load Balancer that will give you a single IP address that will forward all traffic to your service.
  #
  # on cloud provider this will generate an external IP for public access.
  # in local usage (e.g. minikube), the service will be accessible trhough minikube exposed IP.
  type: LoadBalancer

  # route service traffic to pods with label keys and values matching this selector.
  #
  # refer to line 41 or 46
  selector:
    app: my-app

  # the list of ports that are exposed by this service.
  ports:

      # expose the service to outside of cluster, make it publicily accessible
      # via external IP or via cluster public IP (e.g minikube IP) using nodePort below.
      #
      # to get the exposed URL (with IP): `minikube service my-service --url`.
      #   => http://<cluster-public-ip>:<nodePort>
    - nodePort: 32199

      # the incoming external request into nodePort will be directed towards port 80 of this particular service, within the cluster.
      #
      # to get the exposed URL (with IP): `kubectl describe service my-service | findstr "IP"`.
      #   => http://<service-ip>:<port>
      port: 80

      # then from the service, it'll be directed to the available pods (in round-robin style), to pod IP with port 8080.
      #   => http://<pod-ip>:<targetPort>
      #
      # refer to line 75
      targetPort: 8080

---
# ======================== horizontal pod auto scaler
#
# the Horizontal Pod Autoscaler automatically scales the number of pods in a replication controller, deployment, replica set
# or stateful set based on observed CPU utilization (or, with custom metrics support, on some other application-provided metrics).

# pick api version `autoscaling/v2beta2` for auto scaler.
apiVersion: autoscaling/v2beta2

# book this block of yaml for HPA (HorizontalPodAutoscaler).
kind: HorizontalPodAutoscaler

# name it `my-auto-scaler`.
metadata:
  name: my-auto-scaler

# spec of the desired behavior of the auto scaler.
spec:

  # min replica allowed.
  minReplicas: 2

  # max replica allowed.
  maxReplicas: 10

  # the deployment that will be scalled is `my-app-deployment`.
  #
  # refer to line 24 or 27
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: my-app-deployment

  # metrics contains the specifications for which to use to calculate the desired replica count (the maximum replica count across all metrics).
  # the desired replica count is calculated multiplying the ratio between the target value and the current value by the current number of pods.
  metrics:

      # resource refers to a resource metric known to Kubernetes describing each pod in the current scale target (e.g. CPU or memory).
      # in below we define the scaling criteria as, if cpu utlization is changed between the amount of 50% utilization, then scaling process shall happen.
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 50
  
```


```bash
kubectl create -f k8s.yaml
kubectl get deployments
```

```bash
minikube service my-service --url # http://172.18.86.42:30642
minikube ip
```

```bash
kubectl get pods
kubectl describe pods
kubectl logs <pod name> // here to see error log
kubectl exec -it <pod name> -- /bin/sh


kubectl delete deployment my-app-deployment
kubectl delete service my-service
kubectl delete hpa my-auto-scaler
kubectl apply -f k8s.yaml



hey -c 50 -z 5m  http://172.18.86.40:32199/

kubectl get hpa # auto scale

kubectl delete clusterrolebinding kubernetes-dashboard
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.0.0-beta8/aio/deploy/recommended.yaml
minikube dashboard
```

<!-- 
kubectl apply -f https://raw.githubusercontent.com/kubernetes-sigs/metrics-server/master/deploy/kubernetes/aggregated-metrics-reader.yaml
kubectl apply -f https://raw.githubusercontent.com/kubernetes-sigs/metrics-server/master/deploy/kubernetes/auth-delegator.yaml
kubectl apply -f https://raw.githubusercontent.com/kubernetes-sigs/metrics-server/master/deploy/kubernetes/auth-reader.yaml
kubectl apply -f https://raw.githubusercontent.com/kubernetes-sigs/metrics-server/master/deploy/kubernetes/metrics-apiservice.yaml
kubectl apply -f https://raw.githubusercontent.com/kubernetes-sigs/metrics-server/master/deploy/kubernetes/metrics-server-deployment.yaml
kubectl apply -f https://raw.githubusercontent.com/kubernetes-sigs/metrics-server/master/deploy/kubernetes/metrics-server-service.yaml
kubectl apply -f https://raw.githubusercontent.com/kubernetes-sigs/metrics-server/master/deploy/kubernetes/resource-reader.yaml

kubectl edit deployment -n kube-system metrics-server

- --kubelet-insecure-tls
- --kubelet-preferred-address-types=InternalIP,Hostname,InternalDNS,ExternalDNS,ExternalIP -->
