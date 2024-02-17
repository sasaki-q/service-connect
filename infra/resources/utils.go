package resource

import "github.com/aws/jsii-runtime-go"

func vToP(e []string) *[]*string {
	tmp := []*string{}
	for _, v := range e {
		tmp = append(tmp, jsii.String(v))
	}

	return &tmp
}
