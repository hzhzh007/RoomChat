package main

import (
	"encoding/json"
	. "github.com/hzhzh007/RoomChat/common"
	log "github.com/hzhzh007/RoomChat/common/log"
	myzk "github.com/hzhzh007/RoomChat/common/zk"
	"github.com/samuel/go-zookeeper/zk"
	"path"
)

func InitZK() (*zk.Conn, error) {
	config := GetConfig()
	if config == nil {
		log.Fatal("get config error")
	}
	conn, err := myzk.Connect(config.Zk.ZookeeperAddr, config.Zk.ZookeeperTimeout)
	if err != nil {
		log.Error("myzk.Connect() error(%v)", err)
		return nil, err
	}

	fpath := path.Join(ZOOKEEPER_CONN_PATH, config.Zk.ZookeeperNode)

	nodeInfo := ConnectorConf{
		Hostname:  config.Zk.ZookeeperNode,
		ServeAddr: config.ServeAddr,
		RpcAddr:   config.RpcAddr,
	}
	data, err := json.Marshal(nodeInfo)

	if err != nil {
		log.Error("json.Marshal() error(%v)", err)
		return conn, err
	}
	log.Debug("myzk node:\"%s\" registe data: \"%s\"", fpath, string(data))
	err = myzk.DeleteNode(conn, fpath)
	if err != nil {
		log.Error("myzk delete node error:%v", err)
	}
	if err = myzk.RegisterTemp(conn, fpath, data); err != nil {
		log.Error("myzk.RegisterTemp() error(%v)", err)
		return conn, err
	}
	return conn, nil
}

func ZkAddRoomConn(conn *zk.Conn, roomid string) error {
	config := GetConfig()
	fpath := path.Join(ZOOKEEPER_ROOM_PATH, roomid, config.RpcAddr)
	data := []byte(config.RpcAddr)
	log.Debug("myzk node:\"%s\" registe data: \"%s\"", fpath, string(data))
	err := myzk.RegisterTemp(conn, fpath, data)
	if err != nil {
		log.Error("myzk.RegisterTemp() error(%v)", err)
	}
	return err
}

func ZkDeleteRoomConn(conn *zk.Conn, roomid string) error {
	config := GetConfig()
	fpath := path.Join(ZOOKEEPER_ROOM_PATH, roomid, config.RpcAddr)
	log.Debug("myzk delete znode:%s", fpath)
	err := myzk.DeleteNode(conn, fpath)
	if err != nil {
		log.Error("myzk delete node error:%v", err)
	}
	return err
}
