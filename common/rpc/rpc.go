package rpc

import (
	"log"
	"net/rpc"
)

var (
	clientMap map[string]*rpc.Client = make(map[string]*rpc.Client)
)

func GetRpcClient(clientAddr string) (client *rpc.Client, err error) {
	client, ok := clientMap[clientAddr]
	if ok == false {
		log.Println("start create new client to:" + clientAddr)
		client, err = rpc.DialHTTP("tcp", clientAddr)
		if err != nil {
			log.Println("dialing err:", err)
			return client, err
		}
		log.Printf("new created client is:%v", client)
		clientMap[clientAddr] = client
	}
	return client, nil
}

func RemoveRpcClient(clientAddr string) error {
	client, ok := clientMap[clientAddr]
	if ok == false {
		return nil
	}

	client.Close()
	delete(clientMap, clientAddr)
	return nil
}
