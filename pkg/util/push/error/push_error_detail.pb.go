// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: push_error_detail.proto

package pushpb

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
	reflect "reflect"
	strconv "strconv"
	strings "strings"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type ErrorType int32

const (
	INGESTION_RATE_LIMIT_ERROR ErrorType = 0
	REQUEST_RATE_LIMIT_ERROR   ErrorType = 1
	INSTANCE_LIMIT_ERROR       ErrorType = 2
	TENANT_LIMIT_ERROR         ErrorType = 3
	DATA_ERROR                 ErrorType = 4
)

var ErrorType_name = map[int32]string{
	0: "INGESTION_RATE_LIMIT_ERROR",
	1: "REQUEST_RATE_LIMIT_ERROR",
	2: "INSTANCE_LIMIT_ERROR",
	3: "TENANT_LIMIT_ERROR",
	4: "DATA_ERROR",
}

var ErrorType_value = map[string]int32{
	"INGESTION_RATE_LIMIT_ERROR": 0,
	"REQUEST_RATE_LIMIT_ERROR":   1,
	"INSTANCE_LIMIT_ERROR":       2,
	"TENANT_LIMIT_ERROR":         3,
	"DATA_ERROR":                 4,
}

func (ErrorType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_096c13f3e2c4c43c, []int{0}
}

type OptionalFlag int32

const (
	SERVICE_OVERLOAD_STATUS_ON_RATE_LIMIT_ENABLED OptionalFlag = 0
)

var OptionalFlag_name = map[int32]string{
	0: "SERVICE_OVERLOAD_STATUS_ON_RATE_LIMIT_ENABLED",
}

var OptionalFlag_value = map[string]int32{
	"SERVICE_OVERLOAD_STATUS_ON_RATE_LIMIT_ENABLED": 0,
}

func (OptionalFlag) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_096c13f3e2c4c43c, []int{1}
}

type PushErrorDetails struct {
	ErrorType     ErrorType      `protobuf:"varint,1,opt,name=ErrorType,proto3,enum=push.ErrorType" json:"ErrorType,omitempty"`
	OptionalFlags []OptionalFlag `protobuf:"varint,3,rep,packed,name=OptionalFlags,proto3,enum=push.OptionalFlag" json:"OptionalFlags,omitempty"`
}

func (m *PushErrorDetails) Reset()      { *m = PushErrorDetails{} }
func (*PushErrorDetails) ProtoMessage() {}
func (*PushErrorDetails) Descriptor() ([]byte, []int) {
	return fileDescriptor_096c13f3e2c4c43c, []int{0}
}
func (m *PushErrorDetails) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PushErrorDetails) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PushErrorDetails.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PushErrorDetails) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PushErrorDetails.Merge(m, src)
}
func (m *PushErrorDetails) XXX_Size() int {
	return m.Size()
}
func (m *PushErrorDetails) XXX_DiscardUnknown() {
	xxx_messageInfo_PushErrorDetails.DiscardUnknown(m)
}

var xxx_messageInfo_PushErrorDetails proto.InternalMessageInfo

func (m *PushErrorDetails) GetErrorType() ErrorType {
	if m != nil {
		return m.ErrorType
	}
	return INGESTION_RATE_LIMIT_ERROR
}

func (m *PushErrorDetails) GetOptionalFlags() []OptionalFlag {
	if m != nil {
		return m.OptionalFlags
	}
	return nil
}

func init() {
	proto.RegisterEnum("push.ErrorType", ErrorType_name, ErrorType_value)
	proto.RegisterEnum("push.OptionalFlag", OptionalFlag_name, OptionalFlag_value)
	proto.RegisterType((*PushErrorDetails)(nil), "push.PushErrorDetails")
}

func init() { proto.RegisterFile("push_error_detail.proto", fileDescriptor_096c13f3e2c4c43c) }

var fileDescriptor_096c13f3e2c4c43c = []byte{
	// 354 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x91, 0x3f, 0x4f, 0x02, 0x31,
	0x18, 0x87, 0x5b, 0x21, 0x24, 0x36, 0x8a, 0xa4, 0x31, 0x7a, 0x21, 0xa6, 0x21, 0x4e, 0x84, 0x84,
	0x23, 0xe2, 0xe2, 0xe0, 0x52, 0xb8, 0x6a, 0x2e, 0xc1, 0x3b, 0xed, 0x15, 0x06, 0x97, 0x86, 0x53,
	0x3c, 0x48, 0xd0, 0xbb, 0xdc, 0x9f, 0xc1, 0xb8, 0xb8, 0xbb, 0xf8, 0x31, 0xfc, 0x28, 0x8e, 0x8c,
	0x8c, 0x72, 0x2c, 0x8e, 0x7c, 0x04, 0x73, 0xc5, 0x28, 0xe8, 0xd6, 0xe7, 0x7d, 0xde, 0xf7, 0xd7,
	0x37, 0x2d, 0xda, 0x0f, 0x92, 0x68, 0x28, 0x07, 0x61, 0xe8, 0x87, 0xf2, 0x76, 0x10, 0xf7, 0x47,
	0x63, 0x3d, 0x08, 0xfd, 0xd8, 0xc7, 0xf9, 0x4c, 0x94, 0xeb, 0xde, 0x28, 0x1e, 0x26, 0xae, 0x7e,
	0xe3, 0xdf, 0x37, 0x3c, 0xdf, 0xf3, 0x1b, 0x4a, 0xba, 0xc9, 0x9d, 0x22, 0x05, 0xea, 0xb4, 0x1c,
	0x3a, 0x7c, 0x42, 0xa5, 0xcb, 0x24, 0x1a, 0xb2, 0x2c, 0xce, 0x50, 0x69, 0x11, 0xae, 0xa3, 0x4d,
	0xc5, 0xe2, 0x31, 0x18, 0x68, 0xb0, 0x02, 0xab, 0xc5, 0xe6, 0x8e, 0x9e, 0x85, 0xeb, 0x3f, 0x65,
	0xfe, 0xdb, 0x81, 0x4f, 0xd0, 0xb6, 0x1d, 0xc4, 0x23, 0xff, 0xa1, 0x3f, 0x3e, 0x1b, 0xf7, 0xbd,
	0x48, 0xcb, 0x55, 0x72, 0xd5, 0x62, 0x13, 0x2f, 0x47, 0x56, 0x15, 0x5f, 0x6f, 0xac, 0xbd, 0xc0,
	0x95, 0x9b, 0x30, 0x41, 0x65, 0xd3, 0x3a, 0x67, 0x8e, 0x30, 0x6d, 0x4b, 0x72, 0x2a, 0x98, 0xec,
	0x98, 0x17, 0xa6, 0x90, 0x8c, 0x73, 0x9b, 0x97, 0x00, 0x3e, 0x40, 0x1a, 0x67, 0x57, 0x5d, 0xe6,
	0x88, 0xff, 0x16, 0x62, 0x0d, 0xed, 0x9a, 0x96, 0x23, 0xa8, 0xd5, 0x5e, 0x37, 0x1b, 0x78, 0x0f,
	0x61, 0xc1, 0x2c, 0x6a, 0x89, 0xb5, 0x7a, 0x0e, 0x17, 0x11, 0x32, 0xa8, 0xa0, 0xdf, 0x9c, 0xaf,
	0x51, 0xb4, 0xb5, 0xba, 0x1e, 0x3e, 0x42, 0x75, 0x87, 0xf1, 0x9e, 0xd9, 0x66, 0xd2, 0xee, 0x31,
	0xde, 0xb1, 0xa9, 0x21, 0x1d, 0x41, 0x45, 0xd7, 0x91, 0x7f, 0x16, 0xb4, 0x68, 0xab, 0xc3, 0x8c,
	0x12, 0x68, 0x9d, 0x4e, 0x66, 0x04, 0x4c, 0x67, 0x04, 0x2c, 0x66, 0x04, 0x3e, 0xa7, 0x04, 0xbe,
	0xa5, 0x04, 0xbe, 0xa7, 0x04, 0x4e, 0x52, 0x02, 0x3f, 0x52, 0x02, 0x3f, 0x53, 0x02, 0x16, 0x29,
	0x81, 0xaf, 0x73, 0x02, 0x26, 0x73, 0x02, 0xa6, 0x73, 0x02, 0xae, 0x0b, 0xd9, 0x43, 0x05, 0xae,
	0x5b, 0x50, 0x5f, 0x72, 0xfc, 0x15, 0x00, 0x00, 0xff, 0xff, 0x3b, 0xe6, 0xe8, 0x0b, 0xe2, 0x01,
	0x00, 0x00,
}

func (x ErrorType) String() string {
	s, ok := ErrorType_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
func (x OptionalFlag) String() string {
	s, ok := OptionalFlag_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
func (this *PushErrorDetails) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*PushErrorDetails)
	if !ok {
		that2, ok := that.(PushErrorDetails)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.ErrorType != that1.ErrorType {
		return false
	}
	if len(this.OptionalFlags) != len(that1.OptionalFlags) {
		return false
	}
	for i := range this.OptionalFlags {
		if this.OptionalFlags[i] != that1.OptionalFlags[i] {
			return false
		}
	}
	return true
}
func (this *PushErrorDetails) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 6)
	s = append(s, "&pushpb.PushErrorDetails{")
	s = append(s, "ErrorType: "+fmt.Sprintf("%#v", this.ErrorType)+",\n")
	s = append(s, "OptionalFlags: "+fmt.Sprintf("%#v", this.OptionalFlags)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringPushErrorDetail(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func (m *PushErrorDetails) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PushErrorDetails) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PushErrorDetails) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.OptionalFlags) > 0 {
		dAtA2 := make([]byte, len(m.OptionalFlags)*10)
		var j1 int
		for _, num := range m.OptionalFlags {
			for num >= 1<<7 {
				dAtA2[j1] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j1++
			}
			dAtA2[j1] = uint8(num)
			j1++
		}
		i -= j1
		copy(dAtA[i:], dAtA2[:j1])
		i = encodeVarintPushErrorDetail(dAtA, i, uint64(j1))
		i--
		dAtA[i] = 0x1a
	}
	if m.ErrorType != 0 {
		i = encodeVarintPushErrorDetail(dAtA, i, uint64(m.ErrorType))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintPushErrorDetail(dAtA []byte, offset int, v uint64) int {
	offset -= sovPushErrorDetail(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *PushErrorDetails) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ErrorType != 0 {
		n += 1 + sovPushErrorDetail(uint64(m.ErrorType))
	}
	if len(m.OptionalFlags) > 0 {
		l = 0
		for _, e := range m.OptionalFlags {
			l += sovPushErrorDetail(uint64(e))
		}
		n += 1 + sovPushErrorDetail(uint64(l)) + l
	}
	return n
}

func sovPushErrorDetail(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozPushErrorDetail(x uint64) (n int) {
	return sovPushErrorDetail(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *PushErrorDetails) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&PushErrorDetails{`,
		`ErrorType:` + fmt.Sprintf("%v", this.ErrorType) + `,`,
		`OptionalFlags:` + fmt.Sprintf("%v", this.OptionalFlags) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringPushErrorDetail(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *PushErrorDetails) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPushErrorDetail
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: PushErrorDetails: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PushErrorDetails: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ErrorType", wireType)
			}
			m.ErrorType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPushErrorDetail
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ErrorType |= ErrorType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType == 0 {
				var v OptionalFlag
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowPushErrorDetail
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= OptionalFlag(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.OptionalFlags = append(m.OptionalFlags, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowPushErrorDetail
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthPushErrorDetail
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthPushErrorDetail
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				if elementCount != 0 && len(m.OptionalFlags) == 0 {
					m.OptionalFlags = make([]OptionalFlag, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v OptionalFlag
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowPushErrorDetail
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= OptionalFlag(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.OptionalFlags = append(m.OptionalFlags, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field OptionalFlags", wireType)
			}
		default:
			iNdEx = preIndex
			skippy, err := skipPushErrorDetail(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPushErrorDetail
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPushErrorDetail
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipPushErrorDetail(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowPushErrorDetail
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowPushErrorDetail
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowPushErrorDetail
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthPushErrorDetail
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthPushErrorDetail
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowPushErrorDetail
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipPushErrorDetail(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthPushErrorDetail
				}
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthPushErrorDetail = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowPushErrorDetail   = fmt.Errorf("proto: integer overflow")
)
