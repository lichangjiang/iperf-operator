package alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type IperfTask struct {
	metav1.TypeMeta `json:",inline"`

	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +optional
	Status IperfStatus `json:"status,omitempty"`
	Spec   IperfSpec   `json:"spec"`
}

type IperfSpec struct {
	IperfImage string          `json:"iperfImage"`
	ToEmail    string          `json:"toEmail"`
	ServerSpec IperfServerSpec `json:"serverSpec"`
	ClientSpec IperfClientSpec `json:"clientSpec"`
}

type IperfServerSpec struct {
	Port int32 `json:"port,omitempty"`
}

type IperfClientSpec struct {
	Mode     string `json:"mode,omitempty"`
	Udp      bool   `json:"udp,omitempty"`
	BwLimit  string `json:"bwLimit,omitempty"`
	Parallel int32  `json:"parallel,omitempty"`
	Interval int32  `json:"interval,omitempty"`
	Duration int32  `json:"duration,omitempty"`
}

type IperfStatus struct {
	State   string `json:"state,omitempty"`
	Deploy  string `json:"deploy,omitempty"`
	Message string `json:"message,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type IperfTaskList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IperfTask `json:"items"`
}
