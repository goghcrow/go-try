package rewriter

import (
	"go/ast"

	"github.com/goghcrow/go-matcher"
	"github.com/goghcrow/go-matcher/combinator"
)

// 改写 else if try() { } 避免 if 分支惰性求值的问题
func (r *fileRewriter) preRewriteIf() {
	type (
		stmtPtn = matcher.StmtPattern
		exprPtn = matcher.ExprPattern
		ifStmt  = ast.IfStmt
	)

	// else if 的 init 或者 cond 不分包含 try 调用
	// if ... { } else if try?(); try?() { }
	ptn := &ifStmt{
		Else: combinator.OrEx[stmtPtn](r.m, &ifStmt{
			Init: matcher.MkPattern[stmtPtn](r.m, r.containsTryCall),
		}, &ifStmt{
			Cond: matcher.MkPattern[exprPtn](r.m, r.containsTryCall),
		}),
	}
	r.match(ptn, func(c cursor, ctx mctx) {
		iff := c.Node().(*ifStmt)
		elif := iff.Else.(*ifStmt)
		// 将 else if 改写成 else { if ... { } }
		iff.Else = &ast.BlockStmt{
			List: []ast.Stmt{elif},
		}
		// 更新 tryNodes 信息
		r.tryNodes[iff.Else] = true
	})
}
