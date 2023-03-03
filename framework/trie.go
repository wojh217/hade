package framework

import (
	"fmt"
	"strings"
)

type Tree struct {
	root * node
}

type node struct {
	isLast bool    // 是否是代表一个完整路由的节点 比如路径里的中间节点就不是完整节点
	segment string
	handlers []ControllerHandler // 一个handler数组，中间件也是handler
	childs []*node
	parent *node
}

func NewNode() *node {
	return &node{
		isLast: false,
		segment: "",
		handlers: []ControllerHandler{},
		childs: []*node{},
	}
}

func NewTree() *Tree {
	return &Tree{
		root: NewNode(),
	}
}

func PrintNode(prefix string, n *node) {
	new_prefix := prefix + n.segment
	fmt.Printf("%s-%s", prefix, n.segment)
	if len(n.childs) == 0 {
		fmt.Printf("\n")
	}
	for _, node := range n.childs {
		PrintNode(new_prefix, node)
	}
}



// 是否是动态段
func isWildSegment(segment string) bool {
	return strings.HasPrefix(segment, ":")
}

// 在n的子节点中找出匹配segment段的节点， 只找一层
func (n *node) filterChildNodes(segment string) []*node {
	if len(n.childs) == 0 {
		return nil
	}

	// 如果是动态段，则如果有子节点，就肯定匹配
	if isWildSegment(segment) {
		return n.childs
	}

	nodes := make([]*node, 0, len(n.childs))
	for _, cnode := range n.childs {
		// 如果节点的段是动态段，那肯定匹配
		if isWildSegment(cnode.segment) {
			nodes = append(nodes, cnode)
		} else if segment == cnode.segment {
			nodes = append(nodes, cnode)
		}
	}
	return nodes
}

// 从某个节点开始，查找树中是否有匹配uri的节点了
func (n *node) matchNode(uri string) *node {
	segments := strings.SplitN(uri, "/", 2)

	first := segments[0]
	// 从当前节点的子节点中找匹配first的节点
	cnodes := n.filterChildNodes(first)
	if len(cnodes) == 0 {
		return nil
	}

	// 如果段是最终断，则不用再递归往下找
	if len(segments) == 1 {
		for _, tn := range cnodes {
			// 用isLast判断是否是代表路径的节点， 因为可能存在/user/:id 和/user/name/zhang这种情况，那么:id和name是存在一个节点上的
			//
			if tn.isLast {
				return tn
			}
		}
		return nil
	}

	// 如果不是最终断，则从cnodes中接续寻找匹配全断的剩余部分
	for _, tn := range cnodes {
		tnMatch := tn.matchNode(segments[1])
		if tnMatch != nil {
			fmt.Printf("tnMatch != nil: %p\n", tnMatch)
			return tnMatch
		}
	}

	return nil
}

/*
AddRouter
/book/list
/book/:id(冲突)
/book/:id/name
/book/:student/age
/:user/name
/:user/name/:age(冲突)
*/
func (tree *Tree) AddRouter(uri string, handlers []ControllerHandler) error {

	n := tree.root
	// 如果发现uri已经存在于树形结构中，则返回报错
	//if n.matchNode(uri) != nil {
	//	return errors.New("route exists: " + uri)
	//}

	segments := strings.Split(uri, "/")
	for idx, segment := range segments {
		isLast := idx == len(segments) - 1
		var objNode *node

		//childNodes := n.filterChildNodes(segment)
		//// 匹配上了，但名字完全一样才算真的拥有此segment的节点
		//if len(childNodes) > 0 {
		//	for _, cnode := range childNodes {
		//		if cnode.segment == segment {
		//			objNode = cnode
		//			break
		//		}
		//	}
		//}

		// 创建新的子节点
		if objNode == nil {
			cnode := NewNode()
			cnode.segment = segment
			if isLast {
				cnode.isLast = true
				cnode.handlers = handlers

			}
			cnode.parent = n
			n.childs = append(n.childs, cnode)
			objNode = cnode
		}
		n = objNode
	}


	return nil
}

func (tree *Tree) FindHandlers(uri string) []ControllerHandler {
	matchNode := tree.root.matchNode(uri)
	if matchNode == nil {
		return nil
	}
	return matchNode.handlers
}

// n为uri路径匹配的节点，
// uri包含路由参数 如/books/:id   /user/:name/list
func (n *node) parseParamsFromEndNode(uri string) map[string]string {
	result := make(map[string]string)
	segments := strings.Split(uri, "/")
	cnode := n
	// 从后往前找父节点
	for i := len(segments)-1; i >= 0; i-- {
		if cnode.segment == "" {
			break
		}
		if isWildSegment(cnode.segment) {
			// key为去掉冒号后的字符串
			result[cnode.segment[1:]] = segments[i]
		}
		cnode = cnode.parent
	}
	return result
}