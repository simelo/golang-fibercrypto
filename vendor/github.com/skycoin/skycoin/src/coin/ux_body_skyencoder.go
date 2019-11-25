// Code generated by github.com/skycoin/skyencoder. DO NOT EDIT.

package coin

import "github.com/skycoin/skycoin/src/cipher/encoder"

// encodeSizeUxBody computes the size of an encoded object of type UxBody
func encodeSizeUxBody(obj *UxBody) uint64 {
	i0 := uint64(0)

	// obj.SrcTransaction
	i0 += 32

	// obj.Address.Version
	i0++

	// obj.Address.Key
	i0 += 20

	// obj.Coins
	i0 += 8

	// obj.Hours
	i0 += 8

	return i0
}

// encodeUxBody encodes an object of type UxBody to a buffer allocated to the exact size
// required to encode the object.
func encodeUxBody(obj *UxBody) ([]byte, error) {
	n := encodeSizeUxBody(obj)
	buf := make([]byte, n)

	if err := encodeUxBodyToBuffer(buf, obj); err != nil {
		return nil, err
	}

	return buf, nil
}

// encodeUxBodyToBuffer encodes an object of type UxBody to a []byte buffer.
// The buffer must be large enough to encode the object, otherwise an error is returned.
func encodeUxBodyToBuffer(buf []byte, obj *UxBody) error {
	if uint64(len(buf)) < encodeSizeUxBody(obj) {
		return encoder.ErrBufferUnderflow
	}

	e := &encoder.Encoder{
		Buffer: buf[:],
	}

	// obj.SrcTransaction
	e.CopyBytes(obj.SrcTransaction[:])

	// obj.Address.Version
	e.Uint8(obj.Address.Version)

	// obj.Address.Key
	e.CopyBytes(obj.Address.Key[:])

	// obj.Coins
	e.Uint64(obj.Coins)

	// obj.Hours
	e.Uint64(obj.Hours)

	return nil
}

// decodeUxBody decodes an object of type UxBody from a buffer.
// Returns the number of bytes used from the buffer to decode the object.
// If the buffer not long enough to decode the object, returns encoder.ErrBufferUnderflow.
func decodeUxBody(buf []byte, obj *UxBody) (uint64, error) {
	d := &encoder.Decoder{
		Buffer: buf[:],
	}

	{
		// obj.SrcTransaction
		if len(d.Buffer) < len(obj.SrcTransaction) {
			return 0, encoder.ErrBufferUnderflow
		}
		copy(obj.SrcTransaction[:], d.Buffer[:len(obj.SrcTransaction)])
		d.Buffer = d.Buffer[len(obj.SrcTransaction):]
	}

	{
		// obj.Address.Version
		i, err := d.Uint8()
		if err != nil {
			return 0, err
		}
		obj.Address.Version = i
	}

	{
		// obj.Address.Key
		if len(d.Buffer) < len(obj.Address.Key) {
			return 0, encoder.ErrBufferUnderflow
		}
		copy(obj.Address.Key[:], d.Buffer[:len(obj.Address.Key)])
		d.Buffer = d.Buffer[len(obj.Address.Key):]
	}

	{
		// obj.Coins
		i, err := d.Uint64()
		if err != nil {
			return 0, err
		}
		obj.Coins = i
	}

	{
		// obj.Hours
		i, err := d.Uint64()
		if err != nil {
			return 0, err
		}
		obj.Hours = i
	}

	return uint64(len(buf) - len(d.Buffer)), nil
}

// decodeUxBodyExact decodes an object of type UxBody from a buffer.
// If the buffer not long enough to decode the object, returns encoder.ErrBufferUnderflow.
// If the buffer is longer than required to decode the object, returns encoder.ErrRemainingBytes.
func decodeUxBodyExact(buf []byte, obj *UxBody) error {
	if n, err := decodeUxBody(buf, obj); err != nil {
		return err
	} else if n != uint64(len(buf)) {
		return encoder.ErrRemainingBytes
	}

	return nil
}
