package pact

import (
	"bytes"
	"encoding/binary"
	"errors"

	"io"
	"log"
)

//简单协议

//读出一个完整的数据包
func Read(r io.Reader) ([]byte, error) {
	header := make([]byte, 4)
	n, err := io.ReadFull(r, header)

	if err != nil {
		return nil, err
	}
	if n != 4 {
		return nil, errors.New("can not read header ")
	}

	r2 := bytes.NewReader(header)
	var length int32
	binary.Read(r2, binary.BigEndian, &length)

	body := make([]byte, int(length))

	n2, err := io.ReadFull(r, body)
	if err != nil {
		return nil, err
	}
	if n2 == 0 {
		//长度是0的数据也是合法的
		return body, nil
	}
	return body, nil
}

func Write(w io.Writer, data []byte) error {
	length := len(data)
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, int32(length))
	if err != nil {
		return err
	}

	body := []byte{}
	body = append(body, data...)

	n, err := w.Write(buf.Bytes())

	if err != nil {
		return err
	}
	if n != len(buf.Bytes()) {
		log.Println("this should not happen")
	}
	n2, err := w.Write(body)
	if err != nil {
		return err
	}
	if n2 != len(body) {
		log.Println("this should not happening 2")
		return errors.New("write part of the body ")
	}
	return nil

}
