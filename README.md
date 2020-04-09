# Pinning k8s subcomponents with go mod

- bootstrap go mod
```
go mod init
```

- get the main dependencies you need for you project
```
go get k8s.io/client-go@v0.16.4
```

- manually make fake null version for each component to the latest version within the *require* section
skip the one already pinned by go mod such as ```k8s.io/{SUBCOMPONENT} v0.0.0-20200327001022-6496210b90e8 // indirect```
```
k8s.io/{SUBCOMPONENT} => k8s.io/{SUBCOMPONENT} v0.0.0
```

- manually set the *replace* section as following
```
k8s.io/{SUBCOMPONENT} => k8s.io/{SUBCOMPONENT} v0.16.4
```

- run the pod killer
```
go run main.go
```

- test the killer
```
k run nginx --image=nginx --restart=Never
k label pods nginx app=crixo
k label pods delete-pod=yes
```