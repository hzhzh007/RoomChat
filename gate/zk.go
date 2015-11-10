package main

import (
	. "github.com/hzhzh007/RoomChat/common"
	log "github.com/hzhzh007/RoomChat/common/log"
	myzk "github.com/hzhzh007/RoomChat/common/zk"
	"github.com/samuel/go-zookeeper/zk"
	"path"
	"time"
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

	fpath := path.Join(ZOOKEEPER_GATE_PATH, config.Zk.ZookeeperNode)

	err := myzk.DeleteNode(conn, fpath)
	if err != nil {
		log.Error("myzk delete node error:%v", err)
	}

	data := []byte{}
	log.Debug("myzk node:\"%s\" registe data: \"%s\"", fpath, string(data))
	if err = myzk.RegisterTemp(conn, fpath, data); err != nil {
		log.Error("myzk.RegisterTemp() error(%v)", err)
		return conn, err
	}
	go WatchConnChange(conn)

	return conn, nil
}

func WatchConnChange(conn *zk.Conn) {
	fpath := ZOOKEEPER_CONN_PATH
	for {
		log.Info("zk path:%s set a watch", fpath)
		_, _, watch, err := conn.ChildrenW(fpath)
		if err != nil {
			log.Info("path:%s get error ,try later", fpath)
			time.Sleep(10 * time.Second)
			continue
		}
		event := <-watch
		log.Info("zk path:%s receive a event %v", fpath, event)
		UpdateConnectorList()
	}
}
