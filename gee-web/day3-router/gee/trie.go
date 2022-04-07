package gee

import (
	"fmt"
	"strings"
)

// Trie 节点
type node struct {
	pattern  string  // 待匹配路由
	part     string  // 当前节点对应内容
	children []*node // 子节点
	isWild   bool    // 通配标志
}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

// 插入节点
func (n *node) insert(pattern string, parts []string, height int) {
	// 递归出口
	// 当高度和 parts 长度相同
	if len(parts) == height {
		// 将待匹配路由存在节点中
		n.pattern = pattern
		return
	}

	// 获取当前高度对应的部分
	// 由上至下地存储 part
	part := parts[height]
	// 查找第一个匹配的子节点
	child := n.matchChild(part)
	// 没有则新增
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	// 插入下一层节点
	child.insert(pattern, parts, height+1)
}

// 查找节点
func (n *node) search(parts []string, height int) *node {
	// 递归出口
	// 高度和切片长度相同 或 遇到通配节点
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		// 节点中的 pattern 为空则表示未找到
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	// 获取匹配的子节点
	children := n.matchChildren(part)

	for _, child := range children {
		// 在子节点中递归查找
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}

// 遍历所有节点，获取 pattern 不为空的节点
func (n *node) travel(list *([]*node)) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}

// 寻找第一个匹配的节点
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 获取所有匹配节点
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}
