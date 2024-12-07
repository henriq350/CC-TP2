package tcp_handler

import (
	"bytes"
	"encoding/gob"
	"fmt"
)



func SerializeAlert(alert AlertMessage) ([]byte, error) {
	var buffer bytes.Buffer
	
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(alert)
	
	if err != nil {
		fmt.Println("Error encoding alert message:", err)
		return nil, err
	}
	
	return buffer.Bytes(), nil
}


func DeserializeAlert(data []byte) (AlertMessage, error) {
	var alert AlertMessage

	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	
	if err := decoder.Decode(&alert); err != nil {
		fmt.Println("Error decoding alert message:", err)
		return AlertMessage{}, err
	
	}
	
	return alert, nil
}
