# Deployment
<!-- TOC -->

- [Deployment](#deployment)
    - [Prerequisites](#prerequisites)
    - [Deploy Contour](#deploy-contour)
    - [Deploy Discoverers](#deploy-discoverers)
        - [Kubernetes](#kubernetes)
        - [Openstack](#openstack)
    - [Deploy Prometheus](#deploy-prometheus)
        - [Stateful Deployment](#stateful-deployment)
        - [Quick start](#quick-start)
        - [Access the Prometheus web UI](#access-the-prometheus-web-ui)
        - [Access the Alertmanager web UI](#access-the-alertmanager-web-ui)
    - [Deploy Grafana](#deploy-grafana)
        - [Quick start](#quick-start-1)
        - [Access the Grafana UI](#access-the-grafana-ui)
        - [Configure Grafana](#configure-grafana)
            - [Configure datasource](#configure-datasource)
            - [Dashboards](#dashboards)
                - [Add Sample Kubernetes Dashboard](#add-sample-kubernetes-dashboard)
    - [Validation](#validation)
        - [Discovery cluster](#discovery-cluster)
        - [Gimbal cluster](#gimbal-cluster)

<!-- /TOC -->

## Prerequisites

- A copy of this repository. Download, or clone: 

  ```sh
  $ git clone git@github.com:heptio/gimbal.git
  ```

- A single Kubernetes cluster to deploy Gimbal
- Kubernetes or Openstack clusters with flat networking. That is, each Pod has a route-able IP address on the network.

## Deploy Contour

For additional information about Contour, see [the Gimbal architecture doc](../docs/gimbal-architecture.md).

```sh
# Navigate to deployment directory
$ cd deployment

# Deploy Contour
$ kubectl create -f contour/
```

The deployment also includes sample Network Policies which restrict access to Contour and Envoy as well as allow access from Prometheus to scrape for metrics. 

**NOTE**: The current configuration exposes the `/stats` path from the Envoy Admin UI so that Prometheus can scrape for metrics.

## Deploy Discoverers

Service discovery is enabled with the Discoverers, which have both Kubernetes and Openstack implementations.

```sh
# Create gimbal-discoverer namespace
kubectl create -f gimbal-discoverer/01-common.yaml
```

### Kubernetes

The Kubernetes Discoverer is responsible for looking at all services and endpoints in a Kubernetes cluster and synchronizing them to the host cluster. 

[Credentials](../docs/kubernetes-discoverer.md#credentials) to the remote cluster must be created as a Secret.

```sh
# Kubernetes secret
$ kubectl -n gimbal-discovery create secret generic remote-discover-kubecfg \
    --from-file=./config \
    --from-literal=cluster-name=node02

# Deploy Discoverer
$ kubectl apply -f gimbal-discoverer/02-kubernetes-discoverer.yaml
```

For more information, see [the Kubernetes Discoverer doc](../docs/kubernetes-discoverer.md).

### Openstack

The Openstack Discoverer is responsible for looking at all LBaaS and members in an Openstack cluster and synchronizing them to the host cluster. 
 
[Credentials](../docs/openstack-discoverer.md#credentials) to the remote cluster must be created as a secret.

```sh
# Openstack secret
$ kubectl -n gimbal-discovery create secret generic remote-discover-openstack \
    --from-file=certificate-authority-data=./ca.pem \
    --from-literal=cluster-name=openstack \
    --from-literal=username=admin \
    --from-literal=password=abc123 \
    --from-literal=auth-url=https://api.openstack:5000/ \
    --from-literal=tenant-name=heptio

# Deploy Discoverer
$ kubectl apply -f gimbal-discoverer/02-openstack-discoverer.yaml
```

For more information, see [the OpenStack Discoverer doc](../docs/openstack-discoverer.md).

## Deploy Prometheus

Included in the Gimbal repo is a sample development deployment of Prometheus and Alertmanager using temporary storage and may not be suitable for all environments.

### Stateful Deployment

 A stateful deployment of Prometheus should utilize persistent storage within your Kubernetes cluster. This is accomplished by utilizing [Persistent Volumes and Persistent Volume Claims](https://kubernetes.io/docs/concepts/storage/persistent-volumes/) to maintain a correlation between a data volume and the Prometheus pod. Persistent volumes can be `static` or `dynamic` and depends on the backend storage implementation utilized in environment in which the cluster is deployed. Please reference the [documentation](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#types-of-persistent-volumes) which best matches your environment & needs.

### Quick start

```sh
# Deploy 
$ kubectl apply -f prometheus

# Deploy kube-state-metrics to gather cluster information
$ git clone https://github.com/kubernetes/kube-state-metrics.git
$ cd kube-state-metrics
$ kubectl apply -f kubernetes/
```

### Access the Prometheus web UI

```sh
$ kubectl -n gimbal-monitoring port-forward $(kubectl -n gimbal-monitoring get pods -l app=prometheus -l component=server -o jsonpath='{.items[0].metadata.name}') 9090:9090
```

then go to [http://localhost:9090](http://localhost:9090) in your browser

### Access the Alertmanager web UI

```sh
$ kubectl -n gimbal-monitoring port-forward $(kubectl -n gimbal-monitoring get pods -l app=prometheus -l component=alertmanager -o jsonpath='{.items[0].metadata.name}') 9093:9093
```

then go to [http://localhost:9093](http://localhost:9093) in your browser

## Deploy Grafana

Sample development deployment of Grafana using temporary storage.

### Quick start

```sh
# Deploy
$ kubectl apply -f grafana/

# Create secret with grafana credentials
$ kubectl create secret generic grafana -n gimbal-monitoring \
    --from-literal=grafana-admin-password=admin \
    --from-literal=grafana-admin-user=admin 
```

### Access the Grafana UI

```sh
$ kubectl port-forward $(kubectl get pods -l app=grafana -n gimbal-monitoring -o jsonpath='{.items[0].metadata.name}') 3000 -n gimbal-monitoring
```

then go to [http://localhost:3000](http://localhost:3000) in your browser, with `admin` as the username.

### Configure Grafana

Grafana requires some configuration after it's deployed. These steps configure a datasource and import a dashboard to validate the connection. 

#### Configure datasource

1. On the main Grafana page, click **Add Datasource**
2. For **Name** enter _prometheus_
3. In `Type` selector, choose _Prometheus_ 
4. For the URL, enter `http://prometheus:9090`
5. Click **Save & Test**
6. Look for the message box in green stating _Data source is working_

#### Dashboards

Dashboards for Envoy and the Discovery components are included as part of the sample Grafana deployment.

##### Add Sample Kubernetes Dashboard

Add sample dashboard to validate that the data source is collecting data:

1. On the main page, click the plus icon and choose **Import dashboard**
2. Enter _1621_ in the first box
3. In the **prometheus** section, choose the datasource that you just created
4. Click **Import**

## Validation

Now you can verify the deployment:

### Discovery cluster

This example deploys a sample application into the default namespace of [the discovered Kubernetes cluster that you created](#kubernetes).

```sh
# Deploy sample apps
$ kubectl apply -f example-workload/deployment.yaml
```

### Gimbal cluster

Run the following commands on the Gimbal cluster to verify its components:

```sh
# Verify Discoverer Components
$ kubectl get po -n gimbal-discovery
NAME                                         READY     STATUS    RESTARTS   AGE
k8s-kubernetes-discoverer-55899dcb66-lgvnk   1/1       Running   0          5m

# Verify Contour
$ kubectl get po -n gimbal-contour
NAME            READY     STATUS    RESTARTS   AGE
contour-lq6mm   2/2       Running   0          5h

# Verify discovered services
$ kubectl get svc -l gimbal.heptio.com/cluster=node02 
NAME                TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)   AGE
nginx-node02        ClusterIP   None         <none>        80/TCP    17m
kuard-node02        ClusterIP   None         <none>        80/TCP    17m

# Deploy an Ingress route
$ kubectl apply -f example-workload/ingressroute.yaml

# Port forward to the Contour pod
$ kubectl port-forward $(kubectl get pods -n gimbal-contour -l app=contour -o jsonpath='{.items[0].metadata.name}') 9000:80 -n gimbal-contour

# Make a request to Gimbal cluster which will proxy traffic to the secondary cluster
# The curl should respond with the kuard and nginx alternating between requests
$ curl -i -H "Host: mixed.local" localhost:9000
```
