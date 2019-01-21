# Go in Kubernetes
Single node cluster orchestration for go application with minimum effort

## Kubernetes Installation
Installion instruction for 
[kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) and 
[minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/)

Start minikube:
```
minikube start
```

To install and run Kubernetes on Mac follow [this](https://rominirani.com/tutorial-getting-started-with-kubernetes-with-docker-on-mac-7f58467203fd)

## Docker Image
[Dockerfile](https://github.com/shoeb240/go-in-kubernetes/blob/master/Dockerfile)
```dockerfile
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: go-app
  labels:
        app: go-app

spec:
  replicas: 3
  template:
    metadata:
      labels:
        app: go-app

    spec:
      containers:
      - image: go-built-app:1.0
        name: go-built-app
        volumeMounts:
          - name: my-host-path
            mountPath: /go/src

      volumes:
      - name: my-host-path
        hostPath:
          path: /Users/shoeb240/go/src/github.com/shoeb240/go-in-kubernetes
          type: Directory
```

Create docker image from Dockerfile
```text
$ docker build -t go-built-app:1.0 .
Sending build context to Docker daemon  60.93kB
Step 1/11 : FROM golang:alpine AS builder
alpine: Pulling from library/golang
Digest: sha256:198cb8c94b9ee6941ce6d58f29aadb855f64600918ce602cdeacb018ad77d647
Status: Downloaded newer image for golang:alpine
 ---> f56365ec0638
Step 2/11 : ADD . /go/src/
 ---> c2fc6a922ba0
Step 3/11 : RUN apk add git
 ---> Running in 527f60150853
fetch http://dl-cdn.alpinelinux.org/alpine/v3.8/main/x86_64/APKINDEX.tar.gz
fetch http://dl-cdn.alpinelinux.org/alpine/v3.8/community/x86_64/APKINDEX.tar.gz
(1/6) Installing nghttp2-libs (1.32.0-r0)
(2/6) Installing libssh2 (1.8.0-r3)
(3/6) Installing libcurl (7.61.1-r1)
(4/6) Installing expat (2.2.5-r0)
(5/6) Installing pcre2 (10.31-r0)
(6/6) Installing git (2.18.1-r0)
Executing busybox-1.28.4-r2.trigger
OK: 19 MiB in 20 packages
Removing intermediate container 527f60150853
 ---> 8855f35e2204
Step 4/11 : RUN git config --global core.autocrlf false
 ---> Running in a6630a8f9e41
Removing intermediate container a6630a8f9e41
 ---> 19e15d367f2b
Step 5/11 : RUN go get github.com/gorilla/mux
 ---> Running in 4174834b8f6e
Removing intermediate container 4174834b8f6e
 ---> 5d56aad41d92
Step 6/11 : WORKDIR /go/src
 ---> Running in 0f978cca307c
Removing intermediate container 0f978cca307c
 ---> baec6fad4026
Step 7/11 : RUN go build -o main .
 ---> Running in acb05153cb3c
Removing intermediate container acb05153cb3c
 ---> f668e9a60eb7
Step 8/11 : FROM alpine
latest: Pulling from library/alpine
Digest: sha256:46e71df1e5191ab8b8034c5189e325258ec44ea739bba1e5645cff83c9048ff1
Status: Downloaded newer image for alpine:latest
 ---> 3f53bb00af94
Step 9/11 : COPY --from=builder /go/src/main /app/
 ---> 7fc3723cd54a
Step 10/11 : WORKDIR /app
 ---> Running in be5f4e013dca
Removing intermediate container be5f4e013dca
 ---> 0a8db822adfe
Step 11/11 : CMD ["./main"]
 ---> Running in 81309e6bec41
Removing intermediate container 81309e6bec41
 ---> 2a34bf9f52fe
Successfully built 2a34bf9f52fe
Successfully tagged go-built-app:1.0
```

Check created docker images:
```text
Shoeb-Mac:go-in-kubernetes shoeb240$ docker images
REPOSITORY                                 TAG                 IMAGE ID            CREATED             SIZE
go-built-app                               1.0                 2a34bf9f52fe        2 minutes ago       11.3MB
```

## Deployment
One of the most common Kubernetes object is the deployment object. The deployment object defines the container spec required, along with the name and labels used by other parts of Kubernetes to discover and connect to the application.

We will create container using image go-built-app in 3 pods.

[deployment.yaml](https://github.com/shoeb240/go-in-kubernetes/blob/master/deployment.yaml)
```
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: go-app
  labels:
        app: go-app

spec:
  replicas: 3
  template:
    metadata:
      labels:
        app: go-app

    spec:
      containers:
      - image: go-built-app:1.0
        name: go-built-app
        volumeMounts:
          - name: my-host-path
            mountPath: /go/src

      volumes:
      - name: my-host-path
        hostPath:
          path: /Users/shoeb240/go/src/github.com/shoeb240/go-in-kubernetes
          type: Directory
```
*Note: You must change spec.template.spec.volumes.hostPath.path according to your codebase absolute path.

Run kubectl command:
```
$ kubectl create -f deployment.yaml
deployment.extensions/go-app created
```

Checking deployment and service
```
$ kubectl get deployments
NAME     DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
go-app   3         3         3            3           25s
$ kubectl get pods
NAME                      READY   STATUS    RESTARTS   AGE
go-app-5f55f4c977-6dkhr   1/1     Running   0          32s
go-app-5f55f4c977-6tg46   1/1     Running   0          32s
go-app-5f55f4c977-sprfb   1/1     Running   0          32s
Shoeb-Mac:go-in-kubernetes shoeb240$ 
```

We cannot yet access the app because port is not exposed outside world. We need to create service object for that.

## Service
Kubernetes has powerful networking capabilities that control how applications communicate. Service is assigned a unique IP address (also called clusterIP). 

[service.yaml](https://github.com/shoeb240/go-in-kubernetes/blob/master/service.yaml)
```
apiVersion: v1
kind: Service
metadata:
  name: go-app-service
  labels:
    app: go-app
spec:
  type: NodePort
  ports:
  - port: 8081
    nodePort: 30081
  selector:
    app: go-app
```

Run kubectl command:
```
$ kubectl create -f service.yaml
service/go-app-service created
```

The Service selects all applications with the label go-app. As multiple replicas are deployed, they will be automatically load balanced based on this common label. The Service makes the application available via a NodePort.

Checking services
```
$ kubectl get service
NAME             TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
go-app-service   NodePort    10.109.45.113   <none>        8081:30081/TCP   7s
kubernetes       ClusterIP   10.96.0.1       <none>        443/TCP          10d
```


Now we should be able to curl the go Service on <CLUSTER-IP>:<PORT>
```
$ curl localhost:30081/users
  I am responding to your API call
```


## Scalling
We can scale up to 4 replicas from 3 modifying deployment.yaml replicas field as "replicas: 4" and running the following kubectl command
```
$ kubectl apply -f deployment.yaml
deployment.extensions/go-app configured
```

Lets check deployment and pods
```
$ kubectl get deployments
NAME     DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
go-app   4         4         4            4           25m
$ kubectl get pods
NAME                      READY   STATUS    RESTARTS   AGE
go-app-5f55f4c977-6dkhr   1/1     Running   0          25m
go-app-5f55f4c977-6tg46   1/1     Running   0          25m
go-app-5f55f4c977-pzzfr   1/1     Running   0          41s
go-app-5f55f4c977-sprfb   1/1     Running   0          25m
```


## Troubleshooting
Lets see how to find out problem when our deployment is not successful. For instance, we can put wrong image version in deployment.yaml. Our image version is 1.0, we can change it to 1.1 and create deployment again. Donot forget to delete previous deployment.

```
$ kubectl delete -f deployment.yaml
deployment.extensions "go-app" deleted
$ kubectl create -f deployment.yaml
deployment.extensions/go-app created
```

Checking deployment. We see that available deployment is 0.
```
$ kubectl get deployments
NAME     DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
go-app   4         4         4            0           13s
```

Checking pods, they are not ready and status shows ErrImagePull and ImagePullBackOff. So we realize that we have issue with image.
```
$ kubectl get pods
NAME                      READY   STATUS             RESTARTS   AGE
go-app-56577bc9d4-8gmhc   0/1     ErrImagePull       0          21s
go-app-56577bc9d4-gbhkt   0/1     ImagePullBackOff   0          21s
go-app-56577bc9d4-lb9c7   0/1     ErrImagePull       0          21s
go-app-56577bc9d4-r5274   0/1     ErrImagePull       0          21s
```

To get more insight, we write command to describe pod
```
$ kubectl describe pod go-app-56577bc9d4-8gmhc
Name:           go-app-56577bc9d4-8gmhc
Namespace:      default
Node:           docker-for-desktop/192.168.65.3
Start Time:     Sun, 20 Jan 2019 22:04:45 +0600
Labels:         app=go-app
                pod-template-hash=1213367580
Annotations:    <none>
Status:         Pending
IP:             10.1.0.137
Controlled By:  ReplicaSet/go-app-56577bc9d4
Containers:
  go-built-app:
    Container ID:   
    Image:          go-built-app:1.1
    Image ID:       
    Port:           <none>
    Host Port:      <none>
    State:          Waiting
      Reason:       ImagePullBackOff
    Ready:          False
    Restart Count:  0
    Environment:    <none>
    Mounts:
      /go/src from my-host-path (rw)
      /var/run/secrets/kubernetes.io/serviceaccount from default-token-2q6zg (ro)
Conditions:
  Type           Status
  Initialized    True 
  Ready          False 
  PodScheduled   True 
Volumes:
  my-host-path:
    Type:          HostPath (bare host directory volume)
    Path:          /Users/shoeb240/go/src/github.com/shoeb240/go-in-kubernetes
    HostPathType:  Directory
  default-token-2q6zg:
    Type:        Secret (a volume populated by a Secret)
    SecretName:  default-token-2q6zg
    Optional:    false
QoS Class:       BestEffort
Node-Selectors:  <none>
Tolerations:     node.kubernetes.io/not-ready:NoExecute for 300s
                 node.kubernetes.io/unreachable:NoExecute for 300s
Events:
  Type     Reason                 Age                 From                         Message
  ----     ------                 ----                ----                         -------
  Normal   Scheduled              103s                default-scheduler            Successfully assigned go-app-56577bc9d4-8gmhc to docker-for-desktop
  Normal   SuccessfulMountVolume  103s                kubelet, docker-for-desktop  MountVolume.SetUp succeeded for volume "my-host-path"
  Normal   SuccessfulMountVolume  103s                kubelet, docker-for-desktop  MountVolume.SetUp succeeded for volume "default-token-2q6zg"
  Normal   Pulling                43s (x3 over 101s)  kubelet, docker-for-desktop  pulling image "go-built-app:1.1"
  Warning  Failed                 36s (x3 over 88s)   kubelet, docker-for-desktop  Failed to pull image "go-built-app:1.1": rpc error: code = Unknown desc = Error response from daemon: pull access denied for go-built-app, repository does not exist or may require 'docker login'
  Warning  Failed                 36s (x3 over 88s)   kubelet, docker-for-desktop  Error: ErrImagePull
  Normal   BackOff                10s (x4 over 88s)   kubelet, docker-for-desktop  Back-off pulling image "go-built-app:1.1"
  Warning  Failed                 10s (x4 over 88s)   kubelet, docker-for-desktop  Error: ImagePullBackOff
```

At the bottom of the output we see - Failed to pull image "go-built-app:1.1"
