package zk

import (
	"errors"
	log "github.com/hzhzh007/RoomChat/common/log"
	"github.com/samuel/go-zookeeper/zk"
	"os"
	"path"
	"strings"
	"syscall"
	"time"
)

var (
	// error
	ErrNoChild      = errors.New("zk: children is nil")
	ErrNodeNotExist = errors.New("zk: node not exist")
)

// Connect connect to zookeeper, and start a goroutine log the event.
func Connect(addr []string, timeout string) (*zk.Conn, error) {
	timeoutDuration, err := time.ParseDuration(timeout)
	if err != nil {
		log.Error("zk timeout error:%v", err)
		return nil, err
	}
	conn, session, err := zk.Connect(addr, timeoutDuration)
	if err != nil {
		log.Error("zk.Connect(\"%v\", %d) error(%v)", addr, timeoutDuration, err)
		return nil, err
	}
	go func() {
		for {
			event := <-session
			log.Debug("zookeeper get a event: %s", event.State.String())
		}
	}()
	return conn, nil
}

// Create create zookeeper path, if path exists ignore error
func Create(conn *zk.Conn, fpath string) error {
	// create zk root path
	tpath := ""
	for _, str := range strings.Split(fpath, "/")[1:] {
		tpath = path.Join(tpath, "/", str)
		log.Debug("create zookeeper path: \"%s\"", tpath)
		_, err := conn.Create(tpath, []byte(""), 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			if err == zk.ErrNodeExists {
				log.Warning("zk.create(\"%s\") exists", tpath)
			} else {
				log.Error("zk.create(\"%s\") error(%v)", tpath, err)
				return err
			}
		}
	}

	return nil
}

// RegisterTmp create a ephemeral node, and watch it, if node droped then send a SIGQUIT to self.
func RegisterTemp(conn *zk.Conn, fpath string, data []byte) error {
	tpath, err := conn.Create(path.Join(fpath), data, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		log.Error("conn.Create(\"%s\", \"%s\", zk.FlagEphemeral) error(%v)", fpath, string(data), err)
		return err
	}
	log.Debug("create a zookeeper node:%s", tpath)
	return nil
}

func DeleteNode(conn *zk.Conn, fpath string) error {
	log.Debug("delete znode %s", fpath)
	err := conn.Delete(fpath, -1)
	return err
}

// GetNodesW get all child from zk path with a watch.
func GetNodesW(conn *zk.Conn, path string) ([]string, <-chan zk.Event, error) {
	nodes, stat, watch, err := conn.ChildrenW(path)
	if err != nil {
		if err == zk.ErrNoNode {
			return nil, nil, ErrNodeNotExist
		}
		log.Error("zk.ChildrenW(\"%s\") error(%v)", path, err)
		return nil, nil, err
	}
	if stat == nil {
		return nil, nil, ErrNodeNotExist
	}
	if len(nodes) == 0 {
		return nil, nil, ErrNoChild
	}
	return nodes, watch, nil
}

// GetNodes get all child from zk path.
func GetNodes(conn *zk.Conn, path string) ([]string, error) {
	nodes, stat, err := conn.Children(path)
	if err != nil {
		if err == zk.ErrNoNode {
			return nil, ErrNodeNotExist
		}
		log.Error("zk.Children(\"%s\") error(%v)", path, err)
		return nil, err
	}
	if stat == nil {
		return nil, ErrNodeNotExist
	}
	if len(nodes) == 0 {
		return nil, ErrNoChild
	}
	return nodes, nil
}

// killSelf send a SIGQUIT to self.
func killSelf() {
	if err := syscall.Kill(os.Getpid(), syscall.SIGQUIT); err != nil {
		log.Error("syscall.Kill(%d, SIGQUIT) error(%v)", os.Getpid(), err)
	}
}
