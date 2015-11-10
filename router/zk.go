package main

import (
	. "github.com/hzhzh007/RoomChat/common"
	log "github.com/hzhzh007/RoomChat/common/log"
	myzk "github.com/hzhzh007/RoomChat/common/zk"
	"github.com/samuel/go-zookeeper/zk"
	"path"
	"strings"
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

	fpath := path.Join(ZOOKEEPER_ROUTER_PATH, config.ServeAddr)

	err = myzk.DeleteNode(conn, fpath)
	if err != nil {
		log.Error("myzk delete node error:%v", err)
	}

	data := []byte{}
	log.Debug("myzk node:\"%s\" registe data: \"%s\"", fpath, string(data))
	if err = myzk.RegisterTemp(conn, fpath, data); err != nil {
		log.Error("myzk.RegisterTemp() error(%v)", err)
		return conn, err
	}
	return conn, nil
}

func WatchRoomConnChange(conn *zk.Conn, roomid string, cancel <-chan int) {
	fpath := path.Join(ZOOKEEPER_ROOM_PATH, roomid)
	for {
		log.Info("zk path:%s set a watch", fpath)
		_, _, watch, err := conn.ChildrenW(fpath)
		if err != nil {
			log.Info("path:%s get error ,try later", fpath)
			time.Sleep(10 * time.Second)
			continue
		}
		select {
		case event := <-watch:
			log.Info("zk path:%s receive a event %v", fpath, event)
			err = GetRoomConnectorFromZk(roomid)
			if err != nil {
				log.Error("get room:%s from zk error", roomid)
			}
		case <-cancel:
			log.Info("cancel room:%s watch", roomid)
			break
		}
	}
}

func GetRoomConnector(conn *zk.Conn, roomid string) []string {
	fpath := path.Join(ZOOKEEPER_ROOM_PATH, roomid)
	nodes, _, _, err := conn.ChildrenW(fpath)
	if err != nil {
		log.Error("get roomid:%s from zk error:%v", roomid, err)
		return []string{}
	}
	return nodes
}
