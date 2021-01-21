package dag

import "archer/common"

type DAG struct {
	dag map[string][]string
}

func (d *DAG) AddEdge(root, target string) error {
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
	return nil
}

func (d *DAG) TopoSort() ([]string, error) {
	res := make([]string, 0)
	inDegree := make(map[string]int, 0)
	q := make([]string, 0)
	for k, v := range d.dag {
		inDegree[k] = 0
		for _, node := range v {
			inDegree[node] = 0
		}
	}
	for _, arr := range d.dag {
		for _, v := range arr {
			inDegree[v]++
		}
	}
	for k, _ := range d.dag {
		if inDegree[k] == 0 {
			q = append(res, k)
		}
	}
	if len(q) == 0 {
		return nil, common.ErrCyclicGraph
	}
	for len(q) > 0 {
		node := q[0]
		res = append(res, node)
		q = q[1:]
		for _, v := range d.dag[node] {
			inDegree[v]--
			if inDegree[v] == 0 {
				q = append(q, v)
			}
		}
	}
	if len(res) != len(inDegree) {
		return nil, common.ErrCyclicGraph
	}
	return res, nil
}

func New() *DAG {
	return &DAG{dag: map[string][]string{}}
}
