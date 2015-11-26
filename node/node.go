package node

import "encoding/json"
import "log"
import "errors"

import "github.com/boltdb/bolt"

import "datad/defs"
import "datad/context"

type DefaultNodeService struct {
	ctx   context.Context
	nodes map[string]defs.Node
}

var bucketName = []byte("nodes")

func (self DefaultNodeService) AddNode(node defs.Node) error {
	log.Printf("Adding node %s", node.Name)
	data, err := json.Marshal(node)
	if err != nil {
		return nil
	}

	err = self.ctx.Database.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
		}
		return b.Put([]byte(node.Id), data)
	})

	return err
}

func (self DefaultNodeService) NodeExists(id string) (bool, error) {
	exists := false
	err := self.ctx.Database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)

		data := b.Get([]byte(id))
		exists = (data != nil)
		return nil
	})

	return exists, err
}

func (self DefaultNodeService) RemoveNode(id string) {
	delete(self.nodes, id)
}

func (self DefaultNodeService) GetNode(id string) (defs.Node, error) {
	node := defs.Node{}
	err := self.ctx.Database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		data := b.Get([]byte(id))
		if data == nil {
			return errors.New("Node does not exist")
		}

		json.Unmarshal(data, &node)
		return nil
	})

	return node, err
}

func (self DefaultNodeService) NodeInfo() defs.Node {
	return self.ctx.Config.NodeInfo
}

func NewNodeService(ctx context.Context) DefaultNodeService {
	nodeService := DefaultNodeService{ctx, make(map[string]defs.Node)}

	err := ctx.Database.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return nodeService
}

func NewNode(id string, name string, addr string) defs.Node {
	return defs.Node{id, name, addr}
}
