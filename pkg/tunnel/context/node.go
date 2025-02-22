/*
Copyright 2020 The SuperEdge Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package context

import (
	"github.com/superedge/superedge/pkg/tunnel/proto"
	"sync"
)

type node struct {
	name      string
	ch        chan *proto.StreamMsg
	conns     *[]string
	connsLock sync.RWMutex
	pairnodes map[string]string
	nodesLock sync.RWMutex
}

func (edge *node) BindNode(uuid string) {
	edge.connsLock.Lock()
	if edge.conns == nil {
		edge.conns = &[]string{uuid}
	}
	edge.connsLock.Unlock()
}

func (edge *node) UnbindNode(uuid string) {
	//删除连接绑定
	edge.connsLock.Lock()
	for k, v := range *edge.conns {
		if v == uuid {
			*edge.conns = append((*edge.conns)[:k], (*edge.conns)[k+1:len(*edge.conns)]...)
		}
	}
	edge.connsLock.Unlock()
	//删除节点绑定
	edge.nodesLock.Lock()
	delete(edge.pairnodes, uuid)
	edge.nodesLock.Unlock()
}

func (edge *node) Send2Node(msg *proto.StreamMsg) {
	edge.ch <- msg
}

func (edge *node) NodeRecv() <-chan *proto.StreamMsg {
	return edge.ch
}

func (edge *node) GetName() string {
	return edge.name
}

func (edge *node) GetBindConns() []string {
	edge.connsLock.RLock()
	defer edge.connsLock.RUnlock()
	if edge.conns == nil {
		return nil
	}
	return *edge.conns
}
func (edge *node) GetChan() chan *proto.StreamMsg {
	return edge.ch
}

func (edge *node) AddPairNode(uid, nodeName string) {
	edge.nodesLock.Lock()
	defer edge.nodesLock.Unlock()
	edge.pairnodes[uid] = nodeName
}
func (edge *node) RemovePairNode(uid string) {
	edge.nodesLock.Lock()
	defer edge.nodesLock.Unlock()
	delete(edge.pairnodes, uid)
}

func (edge *node) GetPairNode(uid string) string {
	edge.nodesLock.RLock()
	defer edge.nodesLock.RUnlock()
	return edge.pairnodes[uid]
}
