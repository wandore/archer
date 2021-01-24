package dag

import (
	"log"
)

type DAG struct {
	dag map[string][]string
}

func (d *DAG) AddEdge(root, target string) {
	_, exist := d.dag[root]
	if exist {
		flag := false
		for _, v := range d.dag[root] {
			if v == target {
				flag = true
				break
			}
		}
		if !flag {
			d.dag[root] = append(d.dag[root], target)
		}
	} else {
		d.dag[root] = append(d.dag[root], target)
	}
}

func (d *DAG) TopoSort() [][]string {
	res := make([][]string, 0)

	inDegree := make(map[string]int, 0)
	for _, arr := range d.dag {
		for _, v := range arr {
			inDegree[v]++
		}
	}

	q := make([]string, 0)
	for k, _ := range d.dag {
		_, exist := inDegree[k]
		if !exist {
			inDegree[k] = 0
		}
		if inDegree[k] == 0 {
			q = append(q, k)
		}
	}

	if len(q) == 0 {
		log.Fatal("toposort failed because of cyclic graph")
	}

	for len(q) > 0 {
		n := len(q)
		nodeList := make([]string, 0)
		for n > 0 {
			node := q[0]
			nodeList = append(nodeList, node)
			q = q[1:]
			for _, v := range d.dag[node] {
				inDegree[v]--
				if inDegree[v] == 0 {
					q = append(q, v)
				}
			}
			n--
		}
		res = append(res, nodeList)
	}

	n := 0
	for _, v := range res {
		n += len(v)
	}
	if n != len(inDegree) {
		log.Fatal("toposort failed because of cyclic graph")
	}

	return res
}

func New() *DAG {
	return &DAG{dag: map[string][]string{}}
}
