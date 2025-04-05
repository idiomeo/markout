package markout

import (
	"bytes"

	"github.com/common-nighthawk/go-figure"
	"github.com/yuin/goldmark"
	ast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type Renderer struct {
	buf      *bytes.Buffer
	source   []byte // 新增：保存原始文本
	inHeader bool   // 新增：标题状态跟踪
}

func RenderToTerminal(md string) (string, error) {
	var buf bytes.Buffer
	source := []byte(md) // 保存原始数据

	parser := goldmark.DefaultParser()
	doc := parser.Parse(text.NewReader(source))

	renderer := &Renderer{
		buf:    &buf,
		source: source, // 传入原始数据
	}

	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		return renderer.renderNode(n, entering)
	})

	return buf.String(), err
}

// 修复后的节点处理方法
func (r *Renderer) renderNode(n ast.Node, entering bool) (ast.WalkStatus, error) {
	switch v := n.(type) {
	case *ast.Heading:
		if entering {
			r.inHeader = true
		} else {
			r.buf.WriteString("\n\n")
			r.inHeader = false
		}
		return ast.WalkContinue, nil

	case *ast.Text:
		if entering {
			// 正确获取文本段
			segment := v.Segment
			text := segment.Value(r.source)

			if r.inHeader {
				f := figure.NewFigure(string(text), "mini", true)
				r.buf.WriteString(f.String())
			} else {
				r.buf.Write(text)
			}
		}
		return ast.WalkContinue, nil

	case *ast.Emphasis:
		if entering {
			if v.Level == 2 {
				r.buf.WriteString("\033[1m")
			} else {
				r.buf.WriteString("\033[3m")
			}
		} else {
			r.buf.WriteString("\033[0m")
		}
		return ast.WalkContinue, nil

	case *ast.Link:
		if entering {
			r.buf.WriteString("\033[4;34m")
		} else {
			r.buf.WriteString("\033[0m")
		}
		return ast.WalkContinue, nil

	default:
		return ast.WalkContinue, nil
	}
}
