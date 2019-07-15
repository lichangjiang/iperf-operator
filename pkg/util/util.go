package util

import (
	"reflect"

	iperfalpha1 "github.com/lichangjiang/iperf-operator/pkg/apis/iperf.test.svc/alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	Version = iperfalpha1.Version
	Group   = iperfalpha1.CustomResourceGroup
)

var Kind = reflect.TypeOf(iperfalpha1.IperfTask{}).Name()

func Int32Ptr(i int32) *int32 {
	return &i
}

func IperfTaskOwnRef(namespace, uid string) metav1.OwnerReference {
	blockOwner := true
	return metav1.OwnerReference{
		APIVersion:         Version,
		Kind:               Kind,
		Name:               namespace,
		UID:                types.UID(uid),
		BlockOwnerDeletion: &blockOwner,
	}
}
