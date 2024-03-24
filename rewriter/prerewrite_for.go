package rewriter

import (
	"go/ast"
	"go/token"

	"github.com/goghcrow/go-matcher"
	"github.com/goghcrow/go-matcher/combinator"
)

func (r *fileRewriter) preRewriteFor() {
	type (
		stmtPtn = matcher.StmtPattern
		exprPtn = matcher.ExprPattern
	)
	var ptn = combinator.OrEx[stmtPtn](r.m,
		&ast.ForStmt{
			Cond: matcher.MkPattern[exprPtn](r.m, r.containsTryCall),
		},
		&ast.ForStmt{
			Post: matcher.MkPattern[stmtPtn](r.m, r.containsTryCall),
		},
	)
	r.match(ptn, func(c cursor, ctx mctx) {
		forStmt := c.Node().(*ast.ForStmt)
		if r.tryNodes[forStmt.Cond] {
			r.preRewriteForCond(forStmt)
		}
		if r.tryNodes[forStmt.Post] {
			enclosing := ctx.EnclosingFunc()
			assert(enclosing != nil)
			r.preRewriteForPost(forStmt, enclosing)
		}
	})
}

func (r *fileRewriter) preRewriteForCond(forStmt *ast.ForStmt) {
	cond := forStmt.Cond

	// cond 是个 expr
	// 但是 try 展开形式 if err != nil { return } 是 stmt
	// 所以, 如果 cond 包含 try 调用
	// 则将 cond 修改成 body 开头的 if !cond { break }
	ifBrk := &ast.IfStmt{
		Cond: negativeCondExpr(cond),
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.BranchStmt{
					Tok: token.BREAK,
				},
			},
		},
	}
	forStmt.Body.List = prepend(forStmt.Body.List, ast.Stmt(ifBrk))
	forStmt.Cond = nil

	// 更新 tryNodes
	r.tryNodes[ifBrk.Cond] = true
	r.tryNodes[ifBrk] = true
	r.tryNodes[forStmt.Body] = true
}

func (r *fileRewriter) preRewriteForPost(forStmt *ast.ForStmt, enclosingFn fnNode) {
	// golang post 只支持 simple stmt
	// 无法表达 try 的展开形式 (如果支持 {...} ) 则可以
	// 所以, 如果 post 包含 try 调用, 也需要像 cond 一样换个位置
	// 但要保留 post 执行的语义, 参见语言规范 for 部分:
	// "the post statement is executed after each execution of the block (and only if the block was executed)."

	// 处理 post 有两种思路
	// body 后面加个 post-label, post 放到 label stmt 中,
	// 	- 顺序控制流不发生异常, 一定会执行 post-label stmt
	// 	- break 不会执行 post, 不需要处理
	//	- continue 需要人工调整成 goto 到 post-label
	// 		- 特别的, continue 是支持label (指向 for stmt), 准确的说,
	// 			所有跳转到 当前 for stmt label 的 continue 都要修改成 goto
	//			无论在当前 for stmt, 还是在 nested 的 for stmt 中
	//	e.g.
	//
	// outer:
	// for ; ; $post {
	//		continue => goto L_post
	//		for {
	//			continue outer: => goto L_Post // 内嵌的带 label 的 cont 也一样
	//		}
	// 		L_post:
	// 			$post
	// }
	//
	// 但是会有两个问题,
	// 1. 原来 continue 可能是带标签的, 改成 goto 之后,
	// 原来的标签可能就没用了需要移除, 可以加个 refCnt 来计数
	// 2. 因为存在 post 放到了 body 最后, continue (不管是否内嵌) 修改成 goto
	// 之后, 都变成了 forward goto, 这种 goto 和 label 之前不能存在变量声明,
	// 即 decl(var) 或者 assign(:=), 否则会报错 "jumps over declaration"
	//
	// 所以, 需要把 forward goto 转变成 normal goto, 说人话, 就是不能把 post 放到
	// for body 的末尾, 而是 for body 的开头部分, 同时需要加个 flag, 每次 goto 之前
	// 需要先设置 flag, 允许执行 post 代码, e.g.
	//
	// 	for {
	//
	//		postFlag := false
	//	L_post:
	//		if postFlag {
	//			Try0(ret0())
	//
	//			postFlag = false
	//			continue
	//		}
	//
	//		postFlag = true
	//		goto L_post
	//	}
	//
	// 测试了一下, 可以实现, 处理起来比较繁琐, 代码分散, 可读性差, 而且控制流比较绕

	// 第二种思路:
	// 借鉴了 java finally 字节码上的处理, finally 代码采用复制, 而不是 jmp
	// 按照语义, 需要把 post 复制到所以 for block 可能执行结束部分
	// 	1. body 最后,
	// 		特别的, 如果 body 不会终止, 那 post 其实是死代码
	// 		如果把死代码删除 (post 没有被复制过), 可能会导致包导入没被使用, 优化掉 import 可能会导致,
	// 		对应 init 不执行, 代码改写影响了语义, 所以要保留死代码!!!
	// 	2. continue 之前
	// 	3. labeled continue 跳转到当前 forStmt
	// 注意:
	//	- break 不会执行 post 不加
	//	- continue label 到 outer for 不加, 也不会执行内层 post
	// 	- return / panic / exit 之类自然也不需要
	// 这里代码采用了 第二种简单暴力的思路!!!

	post := forStmt.Post
	forStmt.Post = nil

	// 处理上述 2 和 3 两种 case
	// 遍历 enclosing func 的 body 中所有跳转到当前 for 的 continue 节点
	// 并在之前 copy 一份 post
	cpCnt := 0
	jTbl := r.jmpTbl(enclosingFn)
	_, stmt := unpackFunc(enclosingFn)
	r.m.Match(r.pkg.Package, &ast.BranchStmt{Tok: token.CONTINUE}, stmt, func(c cursor, ctx mctx) {
		n := c.Node().(*ast.BranchStmt)
		if jTbl.JumpTo(n, forStmt) {
			// 这里有两种思路
			//	1. 复制 post 节点
			//	 	c.InsertBefore(cloneNode(post))
			//		但是, 只复制 ast 结构不行, tryNodes 信息也需要复制
			//		否则 rewrite 其他节点的 tryNodes 判断不准确
			//	2. 复制 post 引用
			//		带来的问题是, 后续 rewrite 只能当节点是 immutable 是偶
			//		而不能 inplace 替换, 每次修改都需要新建节点, 否则, 会改写所有被复制的节点
			//	这里采用第二种方式, 复制 tryNodes 的工作量也不小, 基本是无差别遍历 ast
			c.InsertBefore(post)
			cpCnt++
		}
	})

	// for body 是否会终止
	// 1. 通用 terminating 检查
	// 2. continue 会退出当前 for, 所以还要加上 continue 的 case
	isTerminatingForBody := func(body *ast.BlockStmt) bool {
		// 去除空语句的 stmts 是否最后是 continue
		list := trimTrailingEmptyStmts(body.List)
		endWithCont := false
		if len(list) > 0 {
			if bs, ok := list[len(list)-1].(*ast.BranchStmt); ok {
				if bs.Tok == token.CONTINUE {
					endWithCont = true
				}
			}
		}
		isForBodyTerminating := endWithCont || r.tChecker.IsTerminating(&ast.BlockStmt{List: list}, "")
		return !isForBodyTerminating
	}

	// 处理上述 case 1
	// 往 body 末尾复制一份 post, 并保留了 post 如果是死代码的语义
	if isTerminatingForBody(forStmt.Body) {
		forStmt.Body.List = append(forStmt.Body.List, post)
	} else if cpCnt == 0 { // 处理 post 死代码的情况
		deadcode := &ast.IfStmt{
			Cond: constFalse(),
			Body: &ast.BlockStmt{
				List: []ast.Stmt{post},
			},
		}
		forStmt.Body.List = append(forStmt.Body.List, deadcode)

		// 更新 tryNodes 信息
		r.tryNodes[deadcode] = true
		r.tryNodes[deadcode.Body] = true
	}
}

func negativeCondExpr(cond ast.Expr) ast.Expr {
	switch bin := cond.(type) {
	case *ast.BinaryExpr:
		switch bin.Op {
		case token.EQL:
			bin.Op = token.NEQ
			return bin
		case token.NEQ:
			bin.Op = token.EQL
			return bin

		case token.LSS:
			bin.Op = token.GEQ
			return bin
		case token.GEQ:
			bin.Op = token.LSS
			return bin

		case token.LEQ:
			bin.Op = token.GTR
			return bin
		case token.GTR:
			bin.Op = token.LEQ
			return bin

		case token.LAND:
			bin.Op = token.LOR
			bin.X = negativeCondExpr(bin.X)
			bin.Y = negativeCondExpr(bin.Y)
			return bin

		case token.LOR:
			bin.Op = token.LAND
			bin.X = negativeCondExpr(bin.X)
			bin.Y = negativeCondExpr(bin.Y)
			return bin
		default: // make ide happier
		}
	case *ast.UnaryExpr:
		if bin.Op == token.NOT {
			return bin.X
		}
	}
	return &ast.UnaryExpr{
		Op: token.NOT,
		X:  cond,
	}
}
