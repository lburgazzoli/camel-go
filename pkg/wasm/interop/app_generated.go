package interop

import (
	"unsafe"

	karmem "karmem.org/golang"
)

var _ unsafe.Pointer

var _Null = make([]byte, 112)
var _NullReader = karmem.NewReader(_Null)

type (
	PacketIdentifier uint64
)

const (
	PacketIdentifierPair         = 7704677971900589564
	PacketIdentifierMessage      = 14302180353067076632
	PacketIdentifierHttpRequest  = 8832126455728058533
	PacketIdentifierHttpResponse = 13025749761452334274
)

type Pair struct {
	Key string
	Val string
}

func NewPair() Pair {
	return Pair{}
}

func (x *Pair) PacketIdentifier() PacketIdentifier {
	return PacketIdentifierPair
}

func (x *Pair) Reset() {
	x.Read((*PairViewer)(unsafe.Pointer(&_Null)), _NullReader)
}

func (x *Pair) WriteAsRoot(writer *karmem.Writer) (offset uint, err error) {
	return x.Write(writer, 0)
}

func (x *Pair) Write(writer *karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(32)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	__KeySize := uint(1 * len(x.Key))
	__KeyOffset, err := writer.Alloc(__KeySize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+0, uint32(__KeyOffset))
	writer.Write4At(offset+0+4, uint32(__KeySize))
	writer.Write4At(offset+0+4+4, 1)
	__KeySlice := [3]uint{*(*uint)(unsafe.Pointer(&x.Key)), __KeySize, __KeySize}
	writer.WriteAt(__KeyOffset, *(*[]byte)(unsafe.Pointer(&__KeySlice)))
	__ValSize := uint(1 * len(x.Val))
	__ValOffset, err := writer.Alloc(__ValSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+12, uint32(__ValOffset))
	writer.Write4At(offset+12+4, uint32(__ValSize))
	writer.Write4At(offset+12+4+4, 1)
	__ValSlice := [3]uint{*(*uint)(unsafe.Pointer(&x.Val)), __ValSize, __ValSize}
	writer.WriteAt(__ValOffset, *(*[]byte)(unsafe.Pointer(&__ValSlice)))

	return offset, nil
}

func (x *Pair) ReadAsRoot(reader *karmem.Reader) {
	x.Read(NewPairViewer(reader, 0), reader)
}

func (x *Pair) Read(viewer *PairViewer, reader *karmem.Reader) {
	__KeyString := viewer.Key(reader)
	if x.Key != __KeyString {
		__KeyStringCopy := make([]byte, len(__KeyString))
		copy(__KeyStringCopy, __KeyString)
		x.Key = *(*string)(unsafe.Pointer(&__KeyStringCopy))
	}
	__ValString := viewer.Val(reader)
	if x.Val != __ValString {
		__ValStringCopy := make([]byte, len(__ValString))
		copy(__ValStringCopy, __ValString)
		x.Val = *(*string)(unsafe.Pointer(&__ValStringCopy))
	}
}

type Message struct {
	ID            string
	Source        string
	Type          string
	Subject       string
	ContentType   string
	ContentSchema string
	Time          int64
	Content       []byte
	Annotations   []Pair
}

func NewMessage() Message {
	return Message{}
}

func (x *Message) PacketIdentifier() PacketIdentifier {
	return PacketIdentifierMessage
}

func (x *Message) Reset() {
	x.Read((*MessageViewer)(unsafe.Pointer(&_Null)), _NullReader)
}

func (x *Message) WriteAsRoot(writer *karmem.Writer) (offset uint, err error) {
	return x.Write(writer, 0)
}

func (x *Message) Write(writer *karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(112)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	writer.Write4At(offset, uint32(108))
	__IDSize := uint(1 * len(x.ID))
	__IDOffset, err := writer.Alloc(__IDSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+4, uint32(__IDOffset))
	writer.Write4At(offset+4+4, uint32(__IDSize))
	writer.Write4At(offset+4+4+4, 1)
	__IDSlice := [3]uint{*(*uint)(unsafe.Pointer(&x.ID)), __IDSize, __IDSize}
	writer.WriteAt(__IDOffset, *(*[]byte)(unsafe.Pointer(&__IDSlice)))
	__SourceSize := uint(1 * len(x.Source))
	__SourceOffset, err := writer.Alloc(__SourceSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+16, uint32(__SourceOffset))
	writer.Write4At(offset+16+4, uint32(__SourceSize))
	writer.Write4At(offset+16+4+4, 1)
	__SourceSlice := [3]uint{*(*uint)(unsafe.Pointer(&x.Source)), __SourceSize, __SourceSize}
	writer.WriteAt(__SourceOffset, *(*[]byte)(unsafe.Pointer(&__SourceSlice)))
	__TypeSize := uint(1 * len(x.Type))
	__TypeOffset, err := writer.Alloc(__TypeSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+28, uint32(__TypeOffset))
	writer.Write4At(offset+28+4, uint32(__TypeSize))
	writer.Write4At(offset+28+4+4, 1)
	__TypeSlice := [3]uint{*(*uint)(unsafe.Pointer(&x.Type)), __TypeSize, __TypeSize}
	writer.WriteAt(__TypeOffset, *(*[]byte)(unsafe.Pointer(&__TypeSlice)))
	__SubjectSize := uint(1 * len(x.Subject))
	__SubjectOffset, err := writer.Alloc(__SubjectSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+40, uint32(__SubjectOffset))
	writer.Write4At(offset+40+4, uint32(__SubjectSize))
	writer.Write4At(offset+40+4+4, 1)
	__SubjectSlice := [3]uint{*(*uint)(unsafe.Pointer(&x.Subject)), __SubjectSize, __SubjectSize}
	writer.WriteAt(__SubjectOffset, *(*[]byte)(unsafe.Pointer(&__SubjectSlice)))
	__ContentTypeSize := uint(1 * len(x.ContentType))
	__ContentTypeOffset, err := writer.Alloc(__ContentTypeSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+52, uint32(__ContentTypeOffset))
	writer.Write4At(offset+52+4, uint32(__ContentTypeSize))
	writer.Write4At(offset+52+4+4, 1)
	__ContentTypeSlice := [3]uint{*(*uint)(unsafe.Pointer(&x.ContentType)), __ContentTypeSize, __ContentTypeSize}
	writer.WriteAt(__ContentTypeOffset, *(*[]byte)(unsafe.Pointer(&__ContentTypeSlice)))
	__ContentSchemaSize := uint(1 * len(x.ContentSchema))
	__ContentSchemaOffset, err := writer.Alloc(__ContentSchemaSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+64, uint32(__ContentSchemaOffset))
	writer.Write4At(offset+64+4, uint32(__ContentSchemaSize))
	writer.Write4At(offset+64+4+4, 1)
	__ContentSchemaSlice := [3]uint{*(*uint)(unsafe.Pointer(&x.ContentSchema)), __ContentSchemaSize, __ContentSchemaSize}
	writer.WriteAt(__ContentSchemaOffset, *(*[]byte)(unsafe.Pointer(&__ContentSchemaSlice)))
	__TimeOffset := offset + 76
	writer.Write8At(__TimeOffset, *(*uint64)(unsafe.Pointer(&x.Time)))
	__ContentSize := uint(1 * len(x.Content))
	__ContentOffset, err := writer.Alloc(__ContentSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+84, uint32(__ContentOffset))
	writer.Write4At(offset+84+4, uint32(__ContentSize))
	writer.Write4At(offset+84+4+4, 1)
	__ContentSlice := *(*[3]uint)(unsafe.Pointer(&x.Content))
	__ContentSlice[1] = __ContentSize
	__ContentSlice[2] = __ContentSize
	writer.WriteAt(__ContentOffset, *(*[]byte)(unsafe.Pointer(&__ContentSlice)))
	__AnnotationsSize := uint(32 * len(x.Annotations))
	__AnnotationsOffset, err := writer.Alloc(__AnnotationsSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+96, uint32(__AnnotationsOffset))
	writer.Write4At(offset+96+4, uint32(__AnnotationsSize))
	writer.Write4At(offset+96+4+4, 32)
	for i := range x.Annotations {
		if _, err := x.Annotations[i].Write(writer, __AnnotationsOffset); err != nil {
			return offset, err
		}
		__AnnotationsOffset += 32
	}

	return offset, nil
}

func (x *Message) ReadAsRoot(reader *karmem.Reader) {
	x.Read(NewMessageViewer(reader, 0), reader)
}

func (x *Message) Read(viewer *MessageViewer, reader *karmem.Reader) {
	__IDString := viewer.ID(reader)
	if x.ID != __IDString {
		__IDStringCopy := make([]byte, len(__IDString))
		copy(__IDStringCopy, __IDString)
		x.ID = *(*string)(unsafe.Pointer(&__IDStringCopy))
	}
	__SourceString := viewer.Source(reader)
	if x.Source != __SourceString {
		__SourceStringCopy := make([]byte, len(__SourceString))
		copy(__SourceStringCopy, __SourceString)
		x.Source = *(*string)(unsafe.Pointer(&__SourceStringCopy))
	}
	__TypeString := viewer.Type(reader)
	if x.Type != __TypeString {
		__TypeStringCopy := make([]byte, len(__TypeString))
		copy(__TypeStringCopy, __TypeString)
		x.Type = *(*string)(unsafe.Pointer(&__TypeStringCopy))
	}
	__SubjectString := viewer.Subject(reader)
	if x.Subject != __SubjectString {
		__SubjectStringCopy := make([]byte, len(__SubjectString))
		copy(__SubjectStringCopy, __SubjectString)
		x.Subject = *(*string)(unsafe.Pointer(&__SubjectStringCopy))
	}
	__ContentTypeString := viewer.ContentType(reader)
	if x.ContentType != __ContentTypeString {
		__ContentTypeStringCopy := make([]byte, len(__ContentTypeString))
		copy(__ContentTypeStringCopy, __ContentTypeString)
		x.ContentType = *(*string)(unsafe.Pointer(&__ContentTypeStringCopy))
	}
	__ContentSchemaString := viewer.ContentSchema(reader)
	if x.ContentSchema != __ContentSchemaString {
		__ContentSchemaStringCopy := make([]byte, len(__ContentSchemaString))
		copy(__ContentSchemaStringCopy, __ContentSchemaString)
		x.ContentSchema = *(*string)(unsafe.Pointer(&__ContentSchemaStringCopy))
	}
	x.Time = viewer.Time()
	__ContentSlice := viewer.Content(reader)
	__ContentLen := len(__ContentSlice)
	if __ContentLen > cap(x.Content) {
		x.Content = append(x.Content, make([]byte, __ContentLen-len(x.Content))...)
	}
	if __ContentLen > len(x.Content) {
		x.Content = x.Content[:__ContentLen]
	}
	copy(x.Content, __ContentSlice)
	x.Content = x.Content[:__ContentLen]
	__AnnotationsSlice := viewer.Annotations(reader)
	__AnnotationsLen := len(__AnnotationsSlice)
	if __AnnotationsLen > cap(x.Annotations) {
		x.Annotations = append(x.Annotations, make([]Pair, __AnnotationsLen-len(x.Annotations))...)
	}
	if __AnnotationsLen > len(x.Annotations) {
		x.Annotations = x.Annotations[:__AnnotationsLen]
	}
	for i := 0; i < __AnnotationsLen; i++ {
		x.Annotations[i].Read(&__AnnotationsSlice[i], reader)
	}
	x.Annotations = x.Annotations[:__AnnotationsLen]
}

type HttpRequest struct {
	URL     string
	Method  string
	Headers []Pair
	Params  []Pair
	Content []byte
}

func NewHttpRequest() HttpRequest {
	return HttpRequest{}
}

func (x *HttpRequest) PacketIdentifier() PacketIdentifier {
	return PacketIdentifierHttpRequest
}

func (x *HttpRequest) Reset() {
	x.Read((*HttpRequestViewer)(unsafe.Pointer(&_Null)), _NullReader)
}

func (x *HttpRequest) WriteAsRoot(writer *karmem.Writer) (offset uint, err error) {
	return x.Write(writer, 0)
}

func (x *HttpRequest) Write(writer *karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(72)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	writer.Write4At(offset, uint32(64))
	__URLSize := uint(1 * len(x.URL))
	__URLOffset, err := writer.Alloc(__URLSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+4, uint32(__URLOffset))
	writer.Write4At(offset+4+4, uint32(__URLSize))
	writer.Write4At(offset+4+4+4, 1)
	__URLSlice := [3]uint{*(*uint)(unsafe.Pointer(&x.URL)), __URLSize, __URLSize}
	writer.WriteAt(__URLOffset, *(*[]byte)(unsafe.Pointer(&__URLSlice)))
	__MethodSize := uint(1 * len(x.Method))
	__MethodOffset, err := writer.Alloc(__MethodSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+16, uint32(__MethodOffset))
	writer.Write4At(offset+16+4, uint32(__MethodSize))
	writer.Write4At(offset+16+4+4, 1)
	__MethodSlice := [3]uint{*(*uint)(unsafe.Pointer(&x.Method)), __MethodSize, __MethodSize}
	writer.WriteAt(__MethodOffset, *(*[]byte)(unsafe.Pointer(&__MethodSlice)))
	__HeadersSize := uint(32 * len(x.Headers))
	__HeadersOffset, err := writer.Alloc(__HeadersSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+28, uint32(__HeadersOffset))
	writer.Write4At(offset+28+4, uint32(__HeadersSize))
	writer.Write4At(offset+28+4+4, 32)
	for i := range x.Headers {
		if _, err := x.Headers[i].Write(writer, __HeadersOffset); err != nil {
			return offset, err
		}
		__HeadersOffset += 32
	}
	__ParamsSize := uint(32 * len(x.Params))
	__ParamsOffset, err := writer.Alloc(__ParamsSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+40, uint32(__ParamsOffset))
	writer.Write4At(offset+40+4, uint32(__ParamsSize))
	writer.Write4At(offset+40+4+4, 32)
	for i := range x.Params {
		if _, err := x.Params[i].Write(writer, __ParamsOffset); err != nil {
			return offset, err
		}
		__ParamsOffset += 32
	}
	__ContentSize := uint(1 * len(x.Content))
	__ContentOffset, err := writer.Alloc(__ContentSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+52, uint32(__ContentOffset))
	writer.Write4At(offset+52+4, uint32(__ContentSize))
	writer.Write4At(offset+52+4+4, 1)
	__ContentSlice := *(*[3]uint)(unsafe.Pointer(&x.Content))
	__ContentSlice[1] = __ContentSize
	__ContentSlice[2] = __ContentSize
	writer.WriteAt(__ContentOffset, *(*[]byte)(unsafe.Pointer(&__ContentSlice)))

	return offset, nil
}

func (x *HttpRequest) ReadAsRoot(reader *karmem.Reader) {
	x.Read(NewHttpRequestViewer(reader, 0), reader)
}

func (x *HttpRequest) Read(viewer *HttpRequestViewer, reader *karmem.Reader) {
	__URLString := viewer.URL(reader)
	if x.URL != __URLString {
		__URLStringCopy := make([]byte, len(__URLString))
		copy(__URLStringCopy, __URLString)
		x.URL = *(*string)(unsafe.Pointer(&__URLStringCopy))
	}
	__MethodString := viewer.Method(reader)
	if x.Method != __MethodString {
		__MethodStringCopy := make([]byte, len(__MethodString))
		copy(__MethodStringCopy, __MethodString)
		x.Method = *(*string)(unsafe.Pointer(&__MethodStringCopy))
	}
	__HeadersSlice := viewer.Headers(reader)
	__HeadersLen := len(__HeadersSlice)
	if __HeadersLen > cap(x.Headers) {
		x.Headers = append(x.Headers, make([]Pair, __HeadersLen-len(x.Headers))...)
	}
	if __HeadersLen > len(x.Headers) {
		x.Headers = x.Headers[:__HeadersLen]
	}
	for i := 0; i < __HeadersLen; i++ {
		x.Headers[i].Read(&__HeadersSlice[i], reader)
	}
	x.Headers = x.Headers[:__HeadersLen]
	__ParamsSlice := viewer.Params(reader)
	__ParamsLen := len(__ParamsSlice)
	if __ParamsLen > cap(x.Params) {
		x.Params = append(x.Params, make([]Pair, __ParamsLen-len(x.Params))...)
	}
	if __ParamsLen > len(x.Params) {
		x.Params = x.Params[:__ParamsLen]
	}
	for i := 0; i < __ParamsLen; i++ {
		x.Params[i].Read(&__ParamsSlice[i], reader)
	}
	x.Params = x.Params[:__ParamsLen]
	__ContentSlice := viewer.Content(reader)
	__ContentLen := len(__ContentSlice)
	if __ContentLen > cap(x.Content) {
		x.Content = append(x.Content, make([]byte, __ContentLen-len(x.Content))...)
	}
	if __ContentLen > len(x.Content) {
		x.Content = x.Content[:__ContentLen]
	}
	copy(x.Content, __ContentSlice)
	x.Content = x.Content[:__ContentLen]
}

type HttpResponse struct {
	Code    int32
	Headers []Pair
	Content []byte
}

func NewHttpResponse() HttpResponse {
	return HttpResponse{}
}

func (x *HttpResponse) PacketIdentifier() PacketIdentifier {
	return PacketIdentifierHttpResponse
}

func (x *HttpResponse) Reset() {
	x.Read((*HttpResponseViewer)(unsafe.Pointer(&_Null)), _NullReader)
}

func (x *HttpResponse) WriteAsRoot(writer *karmem.Writer) (offset uint, err error) {
	return x.Write(writer, 0)
}

func (x *HttpResponse) Write(writer *karmem.Writer, start uint) (offset uint, err error) {
	offset = start
	size := uint(40)
	if offset == 0 {
		offset, err = writer.Alloc(size)
		if err != nil {
			return 0, err
		}
	}
	writer.Write4At(offset, uint32(32))
	__CodeOffset := offset + 4
	writer.Write4At(__CodeOffset, *(*uint32)(unsafe.Pointer(&x.Code)))
	__HeadersSize := uint(32 * len(x.Headers))
	__HeadersOffset, err := writer.Alloc(__HeadersSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+8, uint32(__HeadersOffset))
	writer.Write4At(offset+8+4, uint32(__HeadersSize))
	writer.Write4At(offset+8+4+4, 32)
	for i := range x.Headers {
		if _, err := x.Headers[i].Write(writer, __HeadersOffset); err != nil {
			return offset, err
		}
		__HeadersOffset += 32
	}
	__ContentSize := uint(1 * len(x.Content))
	__ContentOffset, err := writer.Alloc(__ContentSize)
	if err != nil {
		return 0, err
	}
	writer.Write4At(offset+20, uint32(__ContentOffset))
	writer.Write4At(offset+20+4, uint32(__ContentSize))
	writer.Write4At(offset+20+4+4, 1)
	__ContentSlice := *(*[3]uint)(unsafe.Pointer(&x.Content))
	__ContentSlice[1] = __ContentSize
	__ContentSlice[2] = __ContentSize
	writer.WriteAt(__ContentOffset, *(*[]byte)(unsafe.Pointer(&__ContentSlice)))

	return offset, nil
}

func (x *HttpResponse) ReadAsRoot(reader *karmem.Reader) {
	x.Read(NewHttpResponseViewer(reader, 0), reader)
}

func (x *HttpResponse) Read(viewer *HttpResponseViewer, reader *karmem.Reader) {
	x.Code = viewer.Code()
	__HeadersSlice := viewer.Headers(reader)
	__HeadersLen := len(__HeadersSlice)
	if __HeadersLen > cap(x.Headers) {
		x.Headers = append(x.Headers, make([]Pair, __HeadersLen-len(x.Headers))...)
	}
	if __HeadersLen > len(x.Headers) {
		x.Headers = x.Headers[:__HeadersLen]
	}
	for i := 0; i < __HeadersLen; i++ {
		x.Headers[i].Read(&__HeadersSlice[i], reader)
	}
	x.Headers = x.Headers[:__HeadersLen]
	__ContentSlice := viewer.Content(reader)
	__ContentLen := len(__ContentSlice)
	if __ContentLen > cap(x.Content) {
		x.Content = append(x.Content, make([]byte, __ContentLen-len(x.Content))...)
	}
	if __ContentLen > len(x.Content) {
		x.Content = x.Content[:__ContentLen]
	}
	copy(x.Content, __ContentSlice)
	x.Content = x.Content[:__ContentLen]
}

type PairViewer struct {
	_data [32]byte
}

func NewPairViewer(reader *karmem.Reader, offset uint32) (v *PairViewer) {
	if !reader.IsValidOffset(offset, 32) {
		return (*PairViewer)(unsafe.Pointer(&_Null))
	}
	v = (*PairViewer)(unsafe.Add(reader.Pointer, offset))
	return v
}

func (x *PairViewer) size() uint32 {
	return 32
}
func (x *PairViewer) Key(reader *karmem.Reader) (v string) {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(&x._data), 0))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(&x._data), 0+4))
	if !reader.IsValidOffset(offset, size) {
		return ""
	}
	length := uintptr(size / 1)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*string)(unsafe.Pointer(&slice))
}
func (x *PairViewer) Val(reader *karmem.Reader) (v string) {
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(&x._data), 12))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(&x._data), 12+4))
	if !reader.IsValidOffset(offset, size) {
		return ""
	}
	length := uintptr(size / 1)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*string)(unsafe.Pointer(&slice))
}

type MessageViewer struct {
	_data [112]byte
}

func NewMessageViewer(reader *karmem.Reader, offset uint32) (v *MessageViewer) {
	if !reader.IsValidOffset(offset, 8) {
		return (*MessageViewer)(unsafe.Pointer(&_Null))
	}
	v = (*MessageViewer)(unsafe.Add(reader.Pointer, offset))
	if !reader.IsValidOffset(offset, v.size()) {
		return (*MessageViewer)(unsafe.Pointer(&_Null))
	}
	return v
}

func (x *MessageViewer) size() uint32 {
	return *(*uint32)(unsafe.Pointer(&x._data))
}
func (x *MessageViewer) ID(reader *karmem.Reader) (v string) {
	if 4+12 > x.size() {
		return v
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(&x._data), 4))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(&x._data), 4+4))
	if !reader.IsValidOffset(offset, size) {
		return ""
	}
	length := uintptr(size / 1)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*string)(unsafe.Pointer(&slice))
}
func (x *MessageViewer) Source(reader *karmem.Reader) (v string) {
	if 16+12 > x.size() {
		return v
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(&x._data), 16))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(&x._data), 16+4))
	if !reader.IsValidOffset(offset, size) {
		return ""
	}
	length := uintptr(size / 1)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*string)(unsafe.Pointer(&slice))
}
func (x *MessageViewer) Type(reader *karmem.Reader) (v string) {
	if 28+12 > x.size() {
		return v
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(&x._data), 28))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(&x._data), 28+4))
	if !reader.IsValidOffset(offset, size) {
		return ""
	}
	length := uintptr(size / 1)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*string)(unsafe.Pointer(&slice))
}
func (x *MessageViewer) Subject(reader *karmem.Reader) (v string) {
	if 40+12 > x.size() {
		return v
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(&x._data), 40))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(&x._data), 40+4))
	if !reader.IsValidOffset(offset, size) {
		return ""
	}
	length := uintptr(size / 1)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*string)(unsafe.Pointer(&slice))
}
func (x *MessageViewer) ContentType(reader *karmem.Reader) (v string) {
	if 52+12 > x.size() {
		return v
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(&x._data), 52))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(&x._data), 52+4))
	if !reader.IsValidOffset(offset, size) {
		return ""
	}
	length := uintptr(size / 1)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*string)(unsafe.Pointer(&slice))
}
func (x *MessageViewer) ContentSchema(reader *karmem.Reader) (v string) {
	if 64+12 > x.size() {
		return v
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(&x._data), 64))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(&x._data), 64+4))
	if !reader.IsValidOffset(offset, size) {
		return ""
	}
	length := uintptr(size / 1)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*string)(unsafe.Pointer(&slice))
}
func (x *MessageViewer) Time() (v int64) {
	if 76+8 > x.size() {
		return v
	}
	return *(*int64)(unsafe.Add(unsafe.Pointer(&x._data), 76))
}
func (x *MessageViewer) Content(reader *karmem.Reader) (v []byte) {
	if 84+12 > x.size() {
		return []byte{}
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(&x._data), 84))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(&x._data), 84+4))
	if !reader.IsValidOffset(offset, size) {
		return []byte{}
	}
	length := uintptr(size / 1)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]byte)(unsafe.Pointer(&slice))
}
func (x *MessageViewer) Annotations(reader *karmem.Reader) (v []PairViewer) {
	if 96+12 > x.size() {
		return []PairViewer{}
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(&x._data), 96))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(&x._data), 96+4))
	if !reader.IsValidOffset(offset, size) {
		return []PairViewer{}
	}
	length := uintptr(size / 32)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]PairViewer)(unsafe.Pointer(&slice))
}

type HttpRequestViewer struct {
	_data [72]byte
}

func NewHttpRequestViewer(reader *karmem.Reader, offset uint32) (v *HttpRequestViewer) {
	if !reader.IsValidOffset(offset, 8) {
		return (*HttpRequestViewer)(unsafe.Pointer(&_Null))
	}
	v = (*HttpRequestViewer)(unsafe.Add(reader.Pointer, offset))
	if !reader.IsValidOffset(offset, v.size()) {
		return (*HttpRequestViewer)(unsafe.Pointer(&_Null))
	}
	return v
}

func (x *HttpRequestViewer) size() uint32 {
	return *(*uint32)(unsafe.Pointer(&x._data))
}
func (x *HttpRequestViewer) URL(reader *karmem.Reader) (v string) {
	if 4+12 > x.size() {
		return v
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(&x._data), 4))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(&x._data), 4+4))
	if !reader.IsValidOffset(offset, size) {
		return ""
	}
	length := uintptr(size / 1)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*string)(unsafe.Pointer(&slice))
}
func (x *HttpRequestViewer) Method(reader *karmem.Reader) (v string) {
	if 16+12 > x.size() {
		return v
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(&x._data), 16))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(&x._data), 16+4))
	if !reader.IsValidOffset(offset, size) {
		return ""
	}
	length := uintptr(size / 1)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*string)(unsafe.Pointer(&slice))
}
func (x *HttpRequestViewer) Headers(reader *karmem.Reader) (v []PairViewer) {
	if 28+12 > x.size() {
		return []PairViewer{}
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(&x._data), 28))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(&x._data), 28+4))
	if !reader.IsValidOffset(offset, size) {
		return []PairViewer{}
	}
	length := uintptr(size / 32)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]PairViewer)(unsafe.Pointer(&slice))
}
func (x *HttpRequestViewer) Params(reader *karmem.Reader) (v []PairViewer) {
	if 40+12 > x.size() {
		return []PairViewer{}
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(&x._data), 40))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(&x._data), 40+4))
	if !reader.IsValidOffset(offset, size) {
		return []PairViewer{}
	}
	length := uintptr(size / 32)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]PairViewer)(unsafe.Pointer(&slice))
}
func (x *HttpRequestViewer) Content(reader *karmem.Reader) (v []byte) {
	if 52+12 > x.size() {
		return []byte{}
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(&x._data), 52))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(&x._data), 52+4))
	if !reader.IsValidOffset(offset, size) {
		return []byte{}
	}
	length := uintptr(size / 1)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]byte)(unsafe.Pointer(&slice))
}

type HttpResponseViewer struct {
	_data [40]byte
}

func NewHttpResponseViewer(reader *karmem.Reader, offset uint32) (v *HttpResponseViewer) {
	if !reader.IsValidOffset(offset, 8) {
		return (*HttpResponseViewer)(unsafe.Pointer(&_Null))
	}
	v = (*HttpResponseViewer)(unsafe.Add(reader.Pointer, offset))
	if !reader.IsValidOffset(offset, v.size()) {
		return (*HttpResponseViewer)(unsafe.Pointer(&_Null))
	}
	return v
}

func (x *HttpResponseViewer) size() uint32 {
	return *(*uint32)(unsafe.Pointer(&x._data))
}
func (x *HttpResponseViewer) Code() (v int32) {
	if 4+4 > x.size() {
		return v
	}
	return *(*int32)(unsafe.Add(unsafe.Pointer(&x._data), 4))
}
func (x *HttpResponseViewer) Headers(reader *karmem.Reader) (v []PairViewer) {
	if 8+12 > x.size() {
		return []PairViewer{}
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(&x._data), 8))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(&x._data), 8+4))
	if !reader.IsValidOffset(offset, size) {
		return []PairViewer{}
	}
	length := uintptr(size / 32)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]PairViewer)(unsafe.Pointer(&slice))
}
func (x *HttpResponseViewer) Content(reader *karmem.Reader) (v []byte) {
	if 20+12 > x.size() {
		return []byte{}
	}
	offset := *(*uint32)(unsafe.Add(unsafe.Pointer(&x._data), 20))
	size := *(*uint32)(unsafe.Add(unsafe.Pointer(&x._data), 20+4))
	if !reader.IsValidOffset(offset, size) {
		return []byte{}
	}
	length := uintptr(size / 1)
	slice := [3]uintptr{
		uintptr(unsafe.Add(reader.Pointer, offset)), length, length,
	}
	return *(*[]byte)(unsafe.Pointer(&slice))
}
