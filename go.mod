module coding/golang/gomodplay

go 1.13

require (
	github.com/imdario/mergo v0.3.9 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	k8s.io/api v0.16.4
	k8s.io/apimachinery v0.16.4
	k8s.io/client-go v0.16.4
	k8s.io/utils v0.0.0-20200327001022-6496210b90e8 // indirect
)

replace (
	k8s.io/api => k8s.io/api v0.16.4
	k8s.io/apimachinery => k8s.io/apimachinery v0.16.4
	k8s.io/client-go => k8s.io/client-go v0.16.4
	k8s.io/utils => k8s.io/utils v0.0.0-20200327001022-6496210b90e8
)
