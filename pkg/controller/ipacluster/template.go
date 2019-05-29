// Copyright 2019 The FreeIPA Operator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ipacluster

const Template = `apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: foo-statefulset
  namespace: default
spec:
  selector:
    matchLabels:
      statefulset: foo-statefulset
  replicas: 1
  template:
    metadata:
      labels:
        statefulset: foo-statefulset
    spec:
      containers:
      - name: nginx
        image: nginx
`

//const Template = `apiVersion: apps/v1
//kind: StatefulSet
//metadata:
//  name: foo-statefulset
//  namespace: default
//spec:
//  selector:
//    matchLabels:
//      app: freeipa-server
//  serviceName: ns
//  replicas: 1
//  template:
//    metadata:
//      labels:
//        app: freeipa-server
//      annotations:
//        seccomp.security.alpha.kubernetes.io/pod: docker/default
//    spec:
//      priorityClassName: high-priority
//      terminationGracePeriodSeconds: 300
//      containers:
//      - name: freeipa-server
//        imagePullPolicy: IfNotPresent
//        image: freeipa/freeipa-server:centos-7
//        args: ["ipa-replica-install"]
//        volumeMounts:
//        - name: rbd-data
//          mountPath: /data
//        - mountPath: /sys/fs/cgroup
//          name: cgroups
//          readOnly: true
//        - mountPath: /run
//          name: run
//        - mountPath: /run/systemd
//          name: run-sysd
//        - mountPath: /tmp
//          name: tmp
//        ports:
//        - containerPort: 53
//          protocol: TCP
//        - containerPort: 53
//          protocol: UDP
//        - containerPort: 80
//          protocol: TCP
//        - containerPort: 443
//          protocol: TCP
//        - containerPort: 88
//          protocol: TCP
//        - containerPort: 88
//          protocol: UDP
//        - containerPort: 389
//          protocol: TCP
//        - containerPort: 636
//          protocol: TCP
//        - containerPort: 464
//          protocol: TCP
//        - containerPort: 7389
//          protocol: TCP
//        - containerPort: 9443
//          protocol: TCP
//        - containerPort: 9444
//          protocol: TCP
//        - containerPort: 9445
//          protocol: TCP
//        env:
//#        - name: DEBUG_TRACE
//#          value: "1"
//        - name: DEBUG_NO_EXIT
//          value: "1"
//      volumes:
//      - name: cgroups
//        hostPath:
//          path: /sys/fs/cgroup
//      - name: run
//        emptyDir:
//          medium: Memory
//      - name: run-sysd
//        emptyDir:
//          medium: Memory
//      - name: tmp
//        emptyDir:
//          medium: Memory
//  volumeClaimTemplates:
//  - metadata:
//      name: rbd-data
//    spec:
//      accessModes:
//      - ReadWriteOnce
//      storageClassName: cephrbd-dmz
//      resources:
//        requests:
//          storage: 10Gi
//---
//apiVersion: v1
//kind: Service
//metadata:
//  name: ns
//  namespace: dmz
//spec:
//  selector:
//    app: freeipa-server
//  type: ClusterIP
//  clusterIP:  None
//---
//apiVersion: v1
//kind: Service
//metadata:
//  name: ns-0a
//  namespace: dmz
//  annotations:
//    metallb.universe.tf/allow-shared-ip: ns-0
//spec:
//  selector:
//    app: freeipa-server
//    statefulset.kubernetes.io/pod-name: ns-0
//  type: LoadBalancer
//  externalTrafficPolicy: Local
//  loadBalancerIP: 204.152.96.10
//  ports:
//  - name: dns-udp
//    port: 53
//    protocol: UDP
//    targetPort: 53
//---
//apiVersion: v1
//kind: Service
//metadata:
//  name: ns-0b
//  namespace: dmz
//  annotations:
//    metallb.universe.tf/allow-shared-ip: ns-0
//spec:
//  selector:
//    app: freeipa-server
//    statefulset.kubernetes.io/pod-name: ns-0
//  type: LoadBalancer
//  externalTrafficPolicy: Local
//  loadBalancerIP: 204.152.96.10
//  ports:
//  - name: dns-tcp
//    port: 53
//    targetPort: 53
//  - name: http
//    port: 80
//    targetPort: 80
//  - name: https
//    port: 443
//    targetPort: 443
//  - name: kerberos-tcp
//    port: 88
//    protocol: TCP
//    targetPort: 88
//  - name: kerberos-kpasswd
//    port: 464
//    protocol: TCP
//    targetPort: 464
//  - name: ldap
//    port: 389
//    protocol: TCP
//    targetPort: 389
//  - name: ldaps
//    port: 636
//    protocol: TCP
//    targetPort: 636
//---
//apiVersion: v1
//kind: Service
//metadata:
//  name: ns-1a
//  namespace: dmz
//  annotations:
//    metallb.universe.tf/allow-shared-ip: ns-1
//spec:
//  selector:
//    app: freeipa-server
//    statefulset.kubernetes.io/pod-name: ns-1
//  type: LoadBalancer
//  externalTrafficPolicy: Local
//  loadBalancerIP: 204.152.96.11
//  ports:
//  - name: dns-udp
//    port: 53
//    protocol: UDP
//    targetPort: 53
//---
//apiVersion: v1
//kind: Service
//metadata:
//  name: ns-1b
//  namespace: dmz
//  annotations:
//    metallb.universe.tf/allow-shared-ip: ns-1
//spec:
//  selector:
//    app: freeipa-server
//    statefulset.kubernetes.io/pod-name: ns-1
//  type: LoadBalancer
//  externalTrafficPolicy: Local
//  loadBalancerIP: 204.152.96.11
//  ports:
//  - name: dns-tcp
//    port: 53
//    targetPort: 53
//  - name: http
//    port: 80
//    targetPort: 80
//  - name: https
//    port: 443
//    targetPort: 443
//  - name: kerberos-tcp
//    port: 88
//    protocol: TCP
//    targetPort: 88
//  - name: kerberos-kpasswd
//    port: 464
//    protocol: TCP
//    targetPort: 464
//  - name: ldap
//    port: 389
//    protocol: TCP
//    targetPort: 389
//  - name: ldaps
//    port: 636
//    protocol: TCP
//    targetPort: 636
//---
//apiVersion: v1
//kind: Service
//metadata:
//  name: ns-2a
//  namespace: dmz
//  annotations:
//    metallb.universe.tf/allow-shared-ip: ns-2
//spec:
//  selector:
//    app: freeipa-server
//    statefulset.kubernetes.io/pod-name: ns-2
//  type: LoadBalancer
//  externalTrafficPolicy: Local
//  loadBalancerIP: 204.152.96.12
//  ports:
//  - name: dns-udp
//    port: 53
//    protocol: UDP
//    targetPort: 53
//---
//apiVersion: v1
//kind: Service
//metadata:
//  name: ns-2b
//  namespace: dmz
//  annotations:
//    metallb.universe.tf/allow-shared-ip: ns-2
//spec:
//  selector:
//    app: freeipa-server
//    statefulset.kubernetes.io/pod-name: ns-2
//  type: LoadBalancer
//  externalTrafficPolicy: Local
//  loadBalancerIP: 204.152.96.12
//  ports:
//  - name: dns-tcp
//    port: 53
//    targetPort: 53
//  - name: http
//    port: 80
//    targetPort: 80
//  - name: https
//    port: 443
//    targetPort: 443
//  - name: kerberos-tcp
//    port: 88
//    protocol: TCP
//    targetPort: 88
//  - name: kerberos-kpasswd
//    port: 464
//    protocol: TCP
//    targetPort: 464
//  - name: ldap
//    port: 389
//    protocol: TCP
//    targetPort: 389
//  - name: ldaps
//    port: 636
//    protocol: TCP
//    targetPort: 636
//`
