package rpc

import (
	"log"
	"net/rpc"
)

var (
	clientMap map[string]*rpc.Client
)

func GetRpcClient(clientAddr string) (client *rpc.Client, err error) {
	client, ok := clientMap[clientAddr]
	if ok == false {
		client, err := rpc.DialHTTP("tcp", clientAddr)
		if err != nil {
			log.Println("dialing err:", err)
			return client, err
		}
		clientMap[clientAddr] = client
	}
	return
}
