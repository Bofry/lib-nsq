package nsq

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

var (
	_MessageContentSignature         = []byte{0x1b, 0x4e, 0x53, 0x51}
	_MessageContentStateEndDelimiter = []byte{'\r', '\n'}
)

type MessageContent struct {
	State MessageState
	Body  []byte
}

func NewMessageContent() *MessageContent {
	return &MessageContent{
		State: MessageState{},
		Body:  nil,
	}
}

func (c *MessageContent) WriteTo(w io.Writer) (int64, error) {
	var (
		total int64
	)

	n, err := w.Write([]byte{
		0x1b, 0x4e, 0x53, 0x51,
		0x01,
		0x00, 0x00,
	})
	total += int64(n)
	if err != nil {
		return total, err
	}

	// write state and tags
	{
		var (
			statesize int = c.State.byteSize() + (c.State.Len() * 3)
			buf           = make([]byte, statesize+2)
			offset    int = 0
			bytes     int = 0
		)

		bytes = 2
		binary.BigEndian.PutUint16(buf[offset:offset+bytes], uint16(statesize))
		offset += bytes

		if c.State.Len() > 0 {
			for k, v := range c.State.values {
				var size int = len(k) + len(v) + 1

				bytes = 2
				binary.BigEndian.PutUint16(buf[offset:offset+bytes], uint16(size))
				offset += bytes

				offset += copy(buf[offset:], k)
				buf[offset] = ':'
				offset++
				offset += copy(buf[offset:], v)
			}
		}

		n, err = w.Write(buf)
		total += int64(n)
		if err != nil {
			return total, err
		}
	}

	// write state end delimiter
	n, err = w.Write([]byte{'\r', '\n'})
	total += int64(n)
	if err != nil {
		return total, err
	}

	// write body
	n, err = w.Write(c.Body)
	total += int64(n)
	if err != nil {
		return total, err
	}

	return total, nil
}

/* content format:
 *                    2-byte
 *                    unused
 *                      v
 *  [esc][N][S][Q][x][x][x][x][x][x][x][x][x][x][x][x][x][x]...[x][x][x][x][x][x][x][x][x]...[\r][\n][x][x][x][x]...
 *  |  4-byte    || ||    ||    ||    ||  N-byte              ||    ||  N-byte              ||      ||  N-byte
 *  ------------------------------------------------------------------------------------------------------
 *    signature   ^          ^     ^^    tag key and value       ^^    tag key and value        ^^     message body
 *              1-byte       |   2-byte                        2-byte                         2-byte
 *             version       |  tag size                      tag size                  state end delimiter
 *                           |
 *                         2-byte
 *                     tag total size
 */
func DecodeMessageContent(source []byte) (*MessageContent, error) {
	var (
		state *MessageState = &MessageState{}
		body  []byte
	)

	r := bytes.NewReader(source)
	// check signature
	{
		var b [4]byte
		_, err := r.Read(b[:])
		if err != nil {
			return nil, err
		}
		if !bytes.Equal(b[:], _MessageContentSignature) {
			return nil, fmt.Errorf("invalid MessageContent signature")
		}
	}
	// check version
	{
		version, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		if !isSupportedMessageContentVersion(version) {
			return nil, fmt.Errorf("unsupported MessageContent version")
		}
	}
	// skip reserved field
	{
		_, err := r.Seek(2, io.SeekCurrent)
		if err != nil {
			return nil, err
		}
	}
	// read all tags
	{
		var tagReader *bytes.Reader
		// get tag bytes
		{
			var b [2]byte
			_, err := r.Read(b[:])
			if err != nil {
				return nil, err
			}
			size := int(binary.BigEndian.Uint16(b[:]))
			if size > r.Len() {
				return nil, fmt.Errorf("tag total size exceeds")
			}
			offset := int(r.Size()) - r.Len()
			tagReader = bytes.NewReader(source[offset : offset+size])
		}

		for lastLen := -1; tagReader.Len() > 0; {
			if lastLen == tagReader.Len() {
				break
			}
			lastLen = tagReader.Len()

			var b [2]byte
			_, err := tagReader.Read(b[:])
			if err != nil {
				return nil, err
			}
			size := int(binary.BigEndian.Uint16(b[:]))
			if size > tagReader.Len() {
				return nil, fmt.Errorf("tag size exceeds")
			}

			var tagBytes []byte = make([]byte, size)
			_, err = tagReader.Read(tagBytes)
			if err != nil {
				return nil, err
			}

			tagKeyValue := bytes.SplitN(tagBytes, []byte{':'}, 2)
			if tagKeyValue == nil {
				return nil, fmt.Errorf("invalid tag")
			} else {
				if len(tagKeyValue) == 2 {
					key := tagKeyValue[0]
					val := tagKeyValue[1]
					state.Set(string(key), val)
				}
			}
		}
		r.Seek(tagReader.Size(), io.SeekCurrent)
	}

	// read state end delimiter
	{
		var b [2]byte
		_, err := r.Read(b[:])
		if err != nil {
			return nil, err
		}

		if !bytes.Equal(b[:], _MessageContentStateEndDelimiter) {
			return nil, fmt.Errorf("invalid MessageContent state end delimiter")
		}
	}

	// read body
	{
		body = make([]byte, r.Len())
		_, err := r.Read(body)
		if err != nil {
			return nil, err
		}
	}
	return &MessageContent{
		State: *state,
		Body:  body,
	}, nil
}

func isSupportedMessageContentVersion(ver byte) bool {
	return ver == 0x01
}
