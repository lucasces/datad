package main

import "encoding/json"
import "log"
import "errors"

import "github.com/boltdb/bolt"

type NodeService struct {
	ctx   *Context
	nodes map[string]Node
}

var nodesBucketName = []byte("nodes")

func (self *NodeService) AddNode(node Node) error {
	log.Printf("Adding node %s", node.Name)
	data, err := json.Marshal(node)
	if err != nil {
		return nil
	}

	err = self.ctx.Database.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(nodesBucketName)
		if err != nil {
			return err
		}
		return b.Put([]byte(node.Id), data)
	})

	return err
}

func (self *NodeService) NodeExists(id string) (bool, error) {
	exists := false
	err := self.ctx.Database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(nodesBucketName)

		data := b.Get([]byte(id))
		exists = (data != nil)
		return nil
	})

	return exists, err
}

func (self *NodeService) RemoveNode(id string) {
	delete(self.nodes, id)
}

func (self *NodeService) GetNode(id string) (Node, error) {
	node := Node{}
	err := self.ctx.Database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(nodesBucketName)
		data := b.Get([]byte(id))
		if data == nil {
			return errors.New("Node does not exist")
		}

		json.Unmarshal(data, &node)
		return nil
	})

	return node, err
}

func (self *NodeService) ListNodes(offset int, limit int) ([]Node, error) {
	nodes := []Node{}
	err := self.ctx.Database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(nodesBucketName)
		n := 0
		err := b.ForEach(func(k, v []byte) error {
			switch {
			case n > (offset + limit - 1):
				return NewLimitReachedError()
			case n >= offset:
				node := Node{}
				err := json.Unmarshal(v, &node)
				if err != nil {
					return err
				}
				nodes = append(nodes, node)
			}
			n += 1
			return nil
		})
		if _, ok := err.(LimitReachedError); ok {
			return nil
		}
		return err
	})
	return nodes, err
}

func CreateNodeService(ctx *Context) {
	nodeService := NodeService{ctx, make(map[string]Node)}

	err := ctx.Database.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(nodesBucketName)
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	ctx.NodeService = &nodeService
}
