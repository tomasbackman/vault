// Code generated by protoc-gen-go. DO NOT EDIT.
// source: sdk/physical/encrypted_blob.proto

package physical

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// EncryptedBlobInfo contains information about the encrypted value along with
// information about the key used to encrypt it
type EncryptedBlobInfo struct {
	// Ciphertext is the encrypted bytes
	Ciphertext []byte `sentinel:"" protobuf:"bytes,1,opt,name=ciphertext,proto3" json:"ciphertext,omitempty"`
	// IV is the initialization value used during encryption
	IV []byte `sentinel:"" protobuf:"bytes,2,opt,name=iv,proto3" json:"iv,omitempty"`
	// HMAC is the bytes of the HMAC, if any
	HMAC []byte `sentinel:"" protobuf:"bytes,3,opt,name=hmac,proto3" json:"hmac,omitempty"`
	// Wrapped can be used by the client to indicate whether Ciphertext
	// actually contains wrapped data or not. This can be useful if you want to
	// reuse the same struct to pass data along before and after wrapping.
	Wrapped bool `sentinel:"" protobuf:"varint,4,opt,name=wrapped,proto3" json:"wrapped,omitempty"`
	// KeyInfo contains information about the key that was used to create this value
	KeyInfo *KeyInfo `sentinel:"" protobuf:"bytes,5,opt,name=key_info,json=keyInfo,proto3" json:"key_info,omitempty"`
	// ValuePath can be used by the client to store information about where the
	// value came from
	ValuePath            string   `sentinel:"" protobuf:"bytes,6,opt,name=ValuePath,proto3" json:"ValuePath,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EncryptedBlobInfo) Reset()         { *m = EncryptedBlobInfo{} }
func (m *EncryptedBlobInfo) String() string { return proto.CompactTextString(m) }
func (*EncryptedBlobInfo) ProtoMessage()    {}
func (*EncryptedBlobInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_74ed4b518ec7cf88, []int{0}
}

func (m *EncryptedBlobInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EncryptedBlobInfo.Unmarshal(m, b)
}
func (m *EncryptedBlobInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EncryptedBlobInfo.Marshal(b, m, deterministic)
}
func (m *EncryptedBlobInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EncryptedBlobInfo.Merge(m, src)
}
func (m *EncryptedBlobInfo) XXX_Size() int {
	return xxx_messageInfo_EncryptedBlobInfo.Size(m)
}
func (m *EncryptedBlobInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_EncryptedBlobInfo.DiscardUnknown(m)
}

var xxx_messageInfo_EncryptedBlobInfo proto.InternalMessageInfo

func (m *EncryptedBlobInfo) GetCiphertext() []byte {
	if m != nil {
		return m.Ciphertext
	}
	return nil
}

func (m *EncryptedBlobInfo) GetIV() []byte {
	if m != nil {
		return m.IV
	}
	return nil
}

func (m *EncryptedBlobInfo) GetHMAC() []byte {
	if m != nil {
		return m.HMAC
	}
	return nil
}

func (m *EncryptedBlobInfo) GetWrapped() bool {
	if m != nil {
		return m.Wrapped
	}
	return false
}

func (m *EncryptedBlobInfo) GetKeyInfo() *KeyInfo {
	if m != nil {
		return m.KeyInfo
	}
	return nil
}

func (m *EncryptedBlobInfo) GetValuePath() string {
	if m != nil {
		return m.ValuePath
	}
	return ""
}

// KeyInfo contains information regarding which Wrapper key was used to
// encrypt the entry
type KeyInfo struct {
	// Mechanism is the method used by the wrapper to encrypt and sign the
	// data as defined by the wrapper.
	Mechanism     uint64 `sentinel:"" protobuf:"varint,1,opt,name=Mechanism,proto3" json:"Mechanism,omitempty"`
	HMACMechanism uint64 `sentinel:"" protobuf:"varint,2,opt,name=HMACMechanism,proto3" json:"HMACMechanism,omitempty"`
	// This is an opaque ID used by the wrapper to identify the specific
	// key to use as defined by the wrapper.  This could be a version, key
	// label, or something else.
	KeyID     string `sentinel:"" protobuf:"bytes,3,opt,name=KeyID,proto3" json:"KeyID,omitempty"`
	HMACKeyID string `sentinel:"" protobuf:"bytes,4,opt,name=HMACKeyID,proto3" json:"HMACKeyID,omitempty"`
	// These value are used when generating our own data encryption keys
	// and encrypting them using the wrapper
	WrappedKey []byte `sentinel:"" protobuf:"bytes,5,opt,name=WrappedKey,proto3" json:"WrappedKey,omitempty"`
	// Mechanism specific flags
	Flags                uint64   `sentinel:"" protobuf:"varint,6,opt,name=Flags,proto3" json:"Flags,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *KeyInfo) Reset()         { *m = KeyInfo{} }
func (m *KeyInfo) String() string { return proto.CompactTextString(m) }
func (*KeyInfo) ProtoMessage()    {}
func (*KeyInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_74ed4b518ec7cf88, []int{1}
}

func (m *KeyInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KeyInfo.Unmarshal(m, b)
}
func (m *KeyInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KeyInfo.Marshal(b, m, deterministic)
}
func (m *KeyInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KeyInfo.Merge(m, src)
}
func (m *KeyInfo) XXX_Size() int {
	return xxx_messageInfo_KeyInfo.Size(m)
}
func (m *KeyInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_KeyInfo.DiscardUnknown(m)
}

var xxx_messageInfo_KeyInfo proto.InternalMessageInfo

func (m *KeyInfo) GetMechanism() uint64 {
	if m != nil {
		return m.Mechanism
	}
	return 0
}

func (m *KeyInfo) GetHMACMechanism() uint64 {
	if m != nil {
		return m.HMACMechanism
	}
	return 0
}

func (m *KeyInfo) GetKeyID() string {
	if m != nil {
		return m.KeyID
	}
	return ""
}

func (m *KeyInfo) GetHMACKeyID() string {
	if m != nil {
		return m.HMACKeyID
	}
	return ""
}

func (m *KeyInfo) GetWrappedKey() []byte {
	if m != nil {
		return m.WrappedKey
	}
	return nil
}

func (m *KeyInfo) GetFlags() uint64 {
	if m != nil {
		return m.Flags
	}
	return 0
}

func init() {
	proto.RegisterType((*EncryptedBlobInfo)(nil), "github.com.hashicorp.vault.physical.encryption.EncryptedBlobInfo")
	proto.RegisterType((*KeyInfo)(nil), "github.com.hashicorp.vault.physical.encryption.KeyInfo")
}

func init() { proto.RegisterFile("sdk/physical/encrypted_blob.proto", fileDescriptor_74ed4b518ec7cf88) }

var fileDescriptor_74ed4b518ec7cf88 = []byte{
	// 325 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x91, 0xbd, 0x4e, 0xfb, 0x30,
	0x14, 0xc5, 0xe5, 0xfc, 0xd3, 0xaf, 0xfb, 0x2f, 0x48, 0x58, 0x0c, 0x1e, 0x10, 0x0a, 0x15, 0x43,
	0x26, 0x57, 0x82, 0x81, 0x99, 0xf2, 0x21, 0x50, 0x55, 0x09, 0x79, 0x00, 0x89, 0xa5, 0x72, 0x5c,
	0xb7, 0xb6, 0x9a, 0xc6, 0x56, 0xe2, 0x16, 0xf2, 0x64, 0x3c, 0x10, 0x2f, 0x82, 0xe2, 0x10, 0x52,
	0x46, 0x36, 0xdf, 0x5f, 0xce, 0x3d, 0x3a, 0xe7, 0x06, 0xce, 0x8a, 0xc5, 0x7a, 0x6c, 0x55, 0x59,
	0x68, 0xc1, 0xd3, 0xb1, 0xcc, 0x44, 0x5e, 0x5a, 0x27, 0x17, 0xf3, 0x24, 0x35, 0x09, 0xb5, 0xb9,
	0x71, 0x06, 0xd3, 0x95, 0x76, 0x6a, 0x9b, 0x50, 0x61, 0x36, 0x54, 0xf1, 0x42, 0x69, 0x61, 0x72,
	0x4b, 0x77, 0x7c, 0x9b, 0x3a, 0xda, 0x6c, 0xd2, 0xef, 0x4d, 0x6d, 0xb2, 0xd1, 0x27, 0x82, 0xa3,
	0xbb, 0xc6, 0x68, 0x92, 0x9a, 0xe4, 0x31, 0x5b, 0x1a, 0x7c, 0x0a, 0x20, 0xb4, 0x55, 0x32, 0x77,
	0xf2, 0xdd, 0x11, 0x14, 0xa1, 0x78, 0xc8, 0xf6, 0x08, 0x3e, 0x84, 0x40, 0xef, 0x48, 0xe0, 0x79,
	0xa0, 0x77, 0x18, 0x43, 0xa8, 0x36, 0x5c, 0x90, 0x7f, 0x9e, 0xf8, 0x37, 0x26, 0xd0, 0x7b, 0xcb,
	0xb9, 0xb5, 0x72, 0x41, 0xc2, 0x08, 0xc5, 0x7d, 0xd6, 0x8c, 0x98, 0x41, 0x7f, 0x2d, 0xcb, 0xb9,
	0xce, 0x96, 0x86, 0x74, 0x22, 0x14, 0xff, 0xbf, 0xb8, 0xfa, 0x63, 0x6c, 0x3a, 0x95, 0x65, 0x15,
	0x94, 0xf5, 0xd6, 0xf5, 0x03, 0x9f, 0xc0, 0xe0, 0x99, 0xa7, 0x5b, 0xf9, 0xc4, 0x9d, 0x22, 0xdd,
	0x08, 0xc5, 0x03, 0xd6, 0x82, 0xd1, 0x07, 0x82, 0xde, 0xb4, 0x55, 0xce, 0xa4, 0x50, 0x3c, 0xd3,
	0xc5, 0xc6, 0x57, 0x0b, 0x59, 0x0b, 0xf0, 0x39, 0x1c, 0x3c, 0xcc, 0xae, 0x6f, 0x5a, 0x45, 0xe0,
	0x15, 0xbf, 0x21, 0x3e, 0x86, 0x4e, 0x65, 0x77, 0xeb, 0x0b, 0x0f, 0x58, 0x3d, 0x54, 0xce, 0x95,
	0xac, 0xfe, 0x12, 0xd6, 0x19, 0x7e, 0x40, 0x75, 0xd3, 0x97, 0xfa, 0x00, 0x53, 0x59, 0xfa, 0xde,
	0x43, 0xb6, 0x47, 0x2a, 0xcf, 0xfb, 0x94, 0xaf, 0x0a, 0x9f, 0x3e, 0x64, 0xf5, 0x30, 0x81, 0xd7,
	0x7e, 0xd3, 0x3f, 0xe9, 0xfa, 0x5f, 0x7c, 0xf9, 0x15, 0x00, 0x00, 0xff, 0xff, 0xac, 0x92, 0x5a,
	0x15, 0x07, 0x02, 0x00, 0x00,
}
