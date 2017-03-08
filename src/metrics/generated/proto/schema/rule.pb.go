// Copyright (c) 2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Code generated by protoc-gen-go.
// source: rule.proto
// DO NOT EDIT!

package schema

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type MappingRule struct {
	Name       string            `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	TagFilters map[string]string `protobuf:"bytes,2,rep,name=tag_filters,json=tagFilters" json:"tag_filters,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Policies   []*Policy         `protobuf:"bytes,3,rep,name=policies" json:"policies,omitempty"`
}

func (m *MappingRule) Reset()                    { *m = MappingRule{} }
func (m *MappingRule) String() string            { return proto.CompactTextString(m) }
func (*MappingRule) ProtoMessage()               {}
func (*MappingRule) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *MappingRule) GetTagFilters() map[string]string {
	if m != nil {
		return m.TagFilters
	}
	return nil
}

func (m *MappingRule) GetPolicies() []*Policy {
	if m != nil {
		return m.Policies
	}
	return nil
}

type RollupTarget struct {
	Name     string    `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Tags     []string  `protobuf:"bytes,2,rep,name=tags" json:"tags,omitempty"`
	Policies []*Policy `protobuf:"bytes,3,rep,name=policies" json:"policies,omitempty"`
}

func (m *RollupTarget) Reset()                    { *m = RollupTarget{} }
func (m *RollupTarget) String() string            { return proto.CompactTextString(m) }
func (*RollupTarget) ProtoMessage()               {}
func (*RollupTarget) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

func (m *RollupTarget) GetPolicies() []*Policy {
	if m != nil {
		return m.Policies
	}
	return nil
}

type RollupRule struct {
	Name       string            `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	TagFilters map[string]string `protobuf:"bytes,2,rep,name=tag_filters,json=tagFilters" json:"tag_filters,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Targets    []*RollupTarget   `protobuf:"bytes,3,rep,name=targets" json:"targets,omitempty"`
}

func (m *RollupRule) Reset()                    { *m = RollupRule{} }
func (m *RollupRule) String() string            { return proto.CompactTextString(m) }
func (*RollupRule) ProtoMessage()               {}
func (*RollupRule) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2} }

func (m *RollupRule) GetTagFilters() map[string]string {
	if m != nil {
		return m.TagFilters
	}
	return nil
}

func (m *RollupRule) GetTargets() []*RollupTarget {
	if m != nil {
		return m.Targets
	}
	return nil
}

type RuleSet struct {
	Namespace     string         `protobuf:"bytes,1,opt,name=namespace" json:"namespace,omitempty"`
	CreatedAt     int64          `protobuf:"varint,2,opt,name=created_at,json=createdAt" json:"created_at,omitempty"`
	LastUpdatedAt int64          `protobuf:"varint,3,opt,name=last_updated_at,json=lastUpdatedAt" json:"last_updated_at,omitempty"`
	Tombstoned    bool           `protobuf:"varint,4,opt,name=tombstoned" json:"tombstoned,omitempty"`
	Version       int32          `protobuf:"varint,5,opt,name=version" json:"version,omitempty"`
	Cutover       int64          `protobuf:"varint,6,opt,name=cutover" json:"cutover,omitempty"`
	MappingRules  []*MappingRule `protobuf:"bytes,7,rep,name=mapping_rules,json=mappingRules" json:"mapping_rules,omitempty"`
	RollupRules   []*RollupRule  `protobuf:"bytes,8,rep,name=rollup_rules,json=rollupRules" json:"rollup_rules,omitempty"`
}

func (m *RuleSet) Reset()                    { *m = RuleSet{} }
func (m *RuleSet) String() string            { return proto.CompactTextString(m) }
func (*RuleSet) ProtoMessage()               {}
func (*RuleSet) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{3} }

func (m *RuleSet) GetMappingRules() []*MappingRule {
	if m != nil {
		return m.MappingRules
	}
	return nil
}

func (m *RuleSet) GetRollupRules() []*RollupRule {
	if m != nil {
		return m.RollupRules
	}
	return nil
}

type Namespaces struct {
	Namespaces     []string `protobuf:"bytes,1,rep,name=namespaces" json:"namespaces,omitempty"`
	RulesetCutover int64    `protobuf:"varint,2,opt,name=ruleset_cutover,json=rulesetCutover" json:"ruleset_cutover,omitempty"`
}

func (m *Namespaces) Reset()                    { *m = Namespaces{} }
func (m *Namespaces) String() string            { return proto.CompactTextString(m) }
func (*Namespaces) ProtoMessage()               {}
func (*Namespaces) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{4} }

func init() {
	proto.RegisterType((*MappingRule)(nil), "schema.MappingRule")
	proto.RegisterType((*RollupTarget)(nil), "schema.RollupTarget")
	proto.RegisterType((*RollupRule)(nil), "schema.RollupRule")
	proto.RegisterType((*RuleSet)(nil), "schema.RuleSet")
	proto.RegisterType((*Namespaces)(nil), "schema.Namespaces")
}

func init() { proto.RegisterFile("rule.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 449 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xa4, 0x53, 0x4d, 0x8f, 0xd3, 0x30,
	0x10, 0x55, 0x9a, 0x7e, 0x4e, 0xbb, 0x5b, 0x34, 0xec, 0xc1, 0x5a, 0x01, 0xaa, 0x82, 0x04, 0x15,
	0x87, 0x1c, 0x40, 0x48, 0x2b, 0x24, 0x0e, 0xab, 0x05, 0x6e, 0x20, 0x64, 0xb6, 0xe7, 0xc8, 0x4d,
	0x4d, 0x88, 0x48, 0xe2, 0xc8, 0x9e, 0x54, 0xea, 0xef, 0xe2, 0x9f, 0x20, 0x7e, 0x10, 0x8a, 0x5d,
	0xa7, 0x05, 0x0a, 0x12, 0xda, 0xdb, 0xcc, 0xf3, 0xcc, 0xf3, 0xbc, 0xe7, 0x31, 0x80, 0x6e, 0x0a,
	0x19, 0xd7, 0x5a, 0x91, 0xc2, 0xa1, 0x49, 0xbf, 0xc8, 0x52, 0x5c, 0xce, 0x6a, 0x55, 0xe4, 0xe9,
	0xce, 0xa1, 0xd1, 0xf7, 0x00, 0xa6, 0xef, 0x45, 0x5d, 0xe7, 0x55, 0xc6, 0x9b, 0x42, 0x22, 0x42,
	0xbf, 0x12, 0xa5, 0x64, 0xc1, 0x22, 0x58, 0x4e, 0xb8, 0x8d, 0xf1, 0x0d, 0x4c, 0x49, 0x64, 0xc9,
	0xe7, 0xbc, 0x20, 0xa9, 0x0d, 0xeb, 0x2d, 0xc2, 0xe5, 0xf4, 0xf9, 0xe3, 0xd8, 0xf1, 0xc5, 0x47,
	0xdd, 0xf1, 0xad, 0xc8, 0xde, 0xb9, 0xaa, 0xb7, 0x15, 0xe9, 0x1d, 0x07, 0xea, 0x00, 0x7c, 0x06,
	0x63, 0x7b, 0x73, 0x2e, 0x0d, 0x0b, 0x2d, 0xc5, 0xb9, 0xa7, 0xf8, 0x68, 0x27, 0xe2, 0xdd, 0xf9,
	0xe5, 0x6b, 0x98, 0xff, 0x46, 0x85, 0xf7, 0x20, 0xfc, 0x2a, 0x77, 0xfb, 0xb9, 0xda, 0x10, 0x2f,
	0x60, 0xb0, 0x15, 0x45, 0x23, 0x59, 0xcf, 0x62, 0x2e, 0x79, 0xd5, 0xbb, 0x0a, 0xa2, 0x35, 0xcc,
	0xb8, 0x2a, 0x8a, 0xa6, 0xbe, 0x15, 0x3a, 0x93, 0x74, 0x52, 0x14, 0x42, 0x9f, 0x44, 0xe6, 0xd4,
	0x4c, 0xb8, 0x8d, 0xff, 0x67, 0xc4, 0xe8, 0x47, 0x00, 0xe0, 0x2e, 0xf9, 0xab, 0x6f, 0x37, 0xa7,
	0x7c, 0x8b, 0x3c, 0xe3, 0xa1, 0xf9, 0x9f, 0xb6, 0xc5, 0x30, 0x22, 0xab, 0xc2, 0x8f, 0x74, 0xf1,
	0x2b, 0x81, 0x93, 0xc8, 0x7d, 0xd1, 0x5d, 0xad, 0xfb, 0xd6, 0x83, 0x51, 0x3b, 0xd3, 0x27, 0x49,
	0xf8, 0x00, 0x26, 0xad, 0x0e, 0x53, 0x8b, 0xd4, 0x0b, 0x3b, 0x00, 0xf8, 0x10, 0x20, 0xd5, 0x52,
	0x90, 0xdc, 0x24, 0x82, 0x2c, 0x51, 0xc8, 0x27, 0x7b, 0xe4, 0x9a, 0xf0, 0x09, 0xcc, 0x0b, 0x61,
	0x28, 0x69, 0xea, 0x8d, 0xaf, 0x09, 0x6d, 0xcd, 0x59, 0x0b, 0xaf, 0x1c, 0x7a, 0x4d, 0xf8, 0x08,
	0x80, 0x54, 0xb9, 0x36, 0xa4, 0x2a, 0xb9, 0x61, 0xfd, 0x45, 0xb0, 0x1c, 0xf3, 0x23, 0x04, 0x19,
	0x8c, 0xb6, 0x52, 0x9b, 0x5c, 0x55, 0x6c, 0xb0, 0x08, 0x96, 0x03, 0xee, 0xd3, 0xf6, 0x24, 0x6d,
	0x48, 0x6d, 0xa5, 0x66, 0x43, 0xcb, 0xec, 0x53, 0xbc, 0x82, 0xb3, 0xd2, 0x6d, 0x65, 0xd2, 0x7e,
	0x00, 0xc3, 0x46, 0xd6, 0xb9, 0xfb, 0x27, 0x56, 0x96, 0xcf, 0xca, 0x43, 0x62, 0xf0, 0x25, 0xcc,
	0xb4, 0xb5, 0x75, 0xdf, 0x38, 0xb6, 0x8d, 0xf8, 0xe7, 0x9b, 0xf1, 0xa9, 0xee, 0x62, 0x13, 0xad,
	0x00, 0x3e, 0x78, 0x63, 0x4c, 0x2b, 0xa9, 0xb3, 0xc9, 0xb0, 0xc0, 0x2e, 0xd8, 0x11, 0x82, 0x4f,
	0x61, 0x6e, 0xd9, 0x25, 0x25, 0x5e, 0x80, 0xb3, 0xef, 0x7c, 0x0f, 0xdf, 0x38, 0x74, 0x3d, 0xb4,
	0x7f, 0xf4, 0xc5, 0xcf, 0x00, 0x00, 0x00, 0xff, 0xff, 0xeb, 0xfa, 0xe5, 0x5b, 0xc7, 0x03, 0x00,
	0x00,
}