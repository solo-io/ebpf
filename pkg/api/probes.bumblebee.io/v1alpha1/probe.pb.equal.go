// Code generated by protoc-gen-ext. DO NOT EDIT.
// source: github.com/solo-io/bumblebee/api/bumblebee.io/probes/v1alpha1/probe.proto

package v1alpha1

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
	equality "github.com/solo-io/protoc-gen-ext/pkg/equality"
)

// ensure the imports are used
var (
	_ = errors.New("")
	_ = fmt.Print
	_ = binary.LittleEndian
	_ = bytes.Compare
	_ = strings.Compare
	_ = equality.Equalizer(nil)
	_ = proto.Message(nil)
)

// Equal function
func (m *ProbeSpec) Equal(that interface{}) bool {
	if that == nil {
		return m == nil
	}

	target, ok := that.(*ProbeSpec)
	if !ok {
		that2, ok := that.(ProbeSpec)
		if ok {
			target = &that2
		} else {
			return false
		}
	}
	if target == nil {
		return m == nil
	} else if m == nil {
		return false
	}

	if strings.Compare(m.GetImage(), target.GetImage()) != 0 {
		return false
	}

	if len(m.GetNodeSelector()) != len(target.GetNodeSelector()) {
		return false
	}
	for k, v := range m.GetNodeSelector() {

		if strings.Compare(v, target.GetNodeSelector()[k]) != 0 {
			return false
		}

	}

	if m.GetImagePullPolicy() != target.GetImagePullPolicy() {
		return false
	}

	return true
}
