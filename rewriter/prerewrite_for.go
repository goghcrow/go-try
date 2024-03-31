package rewriter

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/goghcrow/go-imports"
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
	// 无法表达 try 的展开形式 (需要支持 block expr {...})
	// 所以, 如果 post 包含 try 调用, 也需要参照 cond 改写位置并保持语义
	// 参见语言规范 for 部分:
	// "the post statement is executed after each execution of the block (and only if the block was executed)."
	//
	// 处理 post 有两种思路
	// body 后面加个 post-label, post 放到 label stmt 中,
	// 	- 顺序控制流不发生异常, 一定会执行 post-label stmt
	// 	- break 不会执行 post, 不需要处理
	//	- continue 需要人工调整成 goto 到 post-label
	// 		- 特别的, continue 是支持label (指向 for stmt), 准确的说,
	// 			无论在当前 for stmt, 还是在 nested 的 for stmt 中的所有跳转到
	//			当前 for stmt label 的 continue 都要修改成 goto, e.g.
	//
	// outer:
	// for ; ; ~$post~ {
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
	// 之后, 都变成了 forward goto, forward goto 的限制是:  goto 和 label 之前不能存在变量声明,
	// 即 decl(var) 或者 assign(:=), 否则会报错 "jumps over declaration"
	//
	// 所以, 需要把 forward goto 转变成 normal goto, 说人话, 就是不能把 post 放到
	// for body 的末尾, 而是 for body 的开头部分, 同时需要加个 flag, 每次 goto 之前
	// 需要先设置 flag, 允许执行 post 代码, e.g.
	//
	// 	for {
	//		postFlag := false
	//	L_post:
	//		if postFlag {
	//			Try0(ret0())
	//			postFlag = false
	//			continue
	//		}
	//
	//		postFlag = true
	//		goto L_post
	//	}
	//
	// 测试了一下, 可以实现, 处理繁琐, 代码分散, 控制流绕, 生成代码可读性差
	// 好处是作用域不会产生副作用, post 迁移到 if body, 父作用域是 for stmt 本身的隐式作用域
	// 即 post 原来所在的作用域, 又因为 post 不会产生新的符号定义, 所以作用域可以看成一样
	//
	// 第二种思路:
	// 借鉴了 java finally 字节码上的处理, finally 代码采用复制, 而不是 jmp
	// 按照语义, 需要把 post 复制到所以 for block 所有 terminated 节点末尾
	// 	1. body 末尾
	// 		特别的, 如果 body 不会终止, 那 post 其实是死代码
	// 		如果把死代码删除 (不复制 post), 可能会导致包导入没被使用, 优化掉 import 可能会导致,
	// 		对应 init 不执行, 代码改写影响了语义, 所以要保留死代码!!!
	// 	2. 当前 forStmt 的 continue 之前
	// 	3. 跳转到当前 forStmt labeled continue 之前
	// 注意:
	//	- break 不会执行 post 不加
	//	- labeled continue 到 外层 for,  当前 for 不执行 post, 只需要复制外层 for 的 post
	// 	- return / panic / exit 之类自然也不需要
	// 复制的方式, 控制流简单, 生成代码可读性高 (后续采用了这种简单粗暴的方式)
	// 但是, 作用域会产生问题, 先看 spec 对 for scope 的说明
	//
	// https://go.dev/ref/spec#Blocks
	// block (也就是 scope)
	// A block is a possibly empty sequence of declarations and statements within matching brace brackets.
	//		Block = "{" StatementList "}" .
	//		StatementList = { Statement ";" } .
	// In addition to explicit blocks in the source code, there are implicit blocks:
	//	The universe block encompasses all Go source text.
	//	Each package has a package block containing all Go source text for that package.
	//	Each file has a file block containing all Go source text in that file.
	//	Each "if", "for", and "switch" statement is considered to be in its own implicit block.
	//	Each clause in a "switch" or "select" statement acts as an implicit block.
	// Blocks nest and influence scoping.
	//
	// scope 链 universe -> package -> file
	// 隐式block: if/for/switch/switch clause/select clause
	//
	// 看一个简单的例子
	// 	for a := 1; a > 0; a = Try(ret1Err[int]()) {
	//		a := 42
	//		_ = a
	//	}
	// ==> ❌ 复制 post 产生的作用域问题
	// 	for a := 1; a > 0; {
	//		a := 42
	//		_ = a
	//		𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := ret1Err[int]()
	//		if 𝗲𝗿𝗿𝟭 != nil {
	//			return 𝗲𝗿𝗿𝟭
	//		}
	//		a = 𝘃𝗮𝗹𝟭 // a resolve 错误!!!
	//	}
	// 错误原因
	// for scope 的 a 是 a1 (1), body 中的 a 是 a2 (42)
	// for-post 的 a resolve 到 a1
	// for-body 的 define 在 for-body 子作用域 shadow 了 parent for 作用域的 a1
	// 但是 for-post 被复制到 for-body 作用域的 shadow 之后, 所以出问题
	// 目的需要保持 scope 父子关系
	// 且for-post 和 for-body 隔离独立, 共同指向隐式的 for scope
	//
	// post 复制到 body 末尾, 分别 wrap block 可以满足语义
	//	for {
	//		{ body }
	//		{ post }
	//	}
	// 但是, 子作用域如果采取单纯的复制 post, 没法保持 body 和 post 作用域独立
	// e.g.
	// 	for a := 1; a > 0; a = Try(ret1Err[int]()) {
	//		a := 42
	//		_ = a
	//		if true {
	//			continue
	//		}
	//		println(a)
	//	}
	// ==> ❌
	// 	for a := 1; a > 0; {
	//		a := 42
	//		_ = a
	//		if true {
	//			𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := ret1Err[int]()
	//			if 𝗲𝗿𝗿𝟭 != nil {
	//				return 𝗲𝗿𝗿𝟭
	//			}
	//			a = 𝘃𝗮𝗹𝟭 // resolve 错 a
	//			continue
	//		}
	//		println(a)
	//		𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟮 := ret1Err[int]()
	//		if 𝗲𝗿𝗿𝟮 != nil {
	//			return 𝗲𝗿𝗿𝟮
	//		}
	//		a = 𝘃𝗮𝗹𝟮 // resolve 错 a
	//	}
	//
	// 所以, 不能直接复制 post 代码, 而是需要保持 post scope 的同时, 复制 post 的计算
	// 所以, 需要把包含 try 的 for-post 打包成作用域正确的闭包放到 for-body 开头, 并在需要复制 post 的地方调用闭包执行 post
	//
	// ==> ✅
	// for a := 1; a > 0; {
	//		𝗽𝗼𝘀𝘁𝟭 := func() error {
	//			𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := ret1Err[int]()
	//			if 𝗲𝗿𝗿𝟭 != nil {
	//				return 𝗲𝗿𝗿𝟭
	//			}
	//			a = 𝘃𝗮𝗹𝟭
	//			return nil
	//		}
	//		a := 42
	//		_ = a
	//		if true {
	//			𝗲𝗿𝗿𝟭 := 𝗽𝗼𝘀𝘁𝟭()
	//			if 𝗲𝗿𝗿𝟭 != nil {
	//				err = 𝗲𝗿𝗿𝟭
	//				return
	//			}
	//			continue
	//		}
	//		println(a)
	//		𝗲𝗿𝗿𝟮 := 𝗽𝗼𝘀𝘁𝟭()
	//		if 𝗲𝗿𝗿𝟮 != nil {
	//			err = 𝗲𝗿𝗿𝟮
	//			return
	//		}
	//	}

	varPostFn, postFnAssign := r.mkPostFnLit(forStmt.Post, enclosingFn)
	tryCallPost := r.mkTryCallPost(varPostFn)

	forStmt.Post = nil
	forStmt.Body.List = prepend[ast.Stmt](forStmt.Body.List, postFnAssign)

	// 处理上述第二种思路 2 和 3 两种 case
	// 遍历 enclosing func 的 body 中所有跳转到当前 for 的 continue 节点
	// 并在之前 copy 一份 post 调用
	cpCnt := 0
	jTbl := r.jmpTbl(enclosingFn)
	_, stmt := unpackFunc(enclosingFn)
	r.m.Match(r.pkg.Package, &ast.BranchStmt{Tok: token.CONTINUE}, stmt, func(c cursor, ctx mctx) {
		n := c.Node().(*ast.BranchStmt)
		if jTbl.JumpTo(n, forStmt) {
			// 这里有两种思路
			//	1. 复制 post 节点
			//	 	c.InsertBefore(cloneNode(tryCallPost))
			//		但是, 只复制 ast 结构不行, tryNodes 信息也需要复制
			//		否则 rewrite 其他节点的 tryNodes 判断不准确
			//	2. 复制 post 引用
			//		带来的问题是, 后续 rewrite 只能默认节点是 immutable
			//		而不能 in_place 替换, 每次修改都需要新建节点, 否则, 会改写所有被复制的节点
			//	这里采用第二种方式
			c.InsertBefore(tryCallPost)
			// 更新parent tryNode
			for _, x := range ctx.Stack {
				r.tryNodes[x] = true
			}
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
	// 往 body 末尾复制一份 post 调用, 并保留了死代码的语义
	if isTerminatingForBody(forStmt.Body) {
		forStmt.Body.List = append(forStmt.Body.List, tryCallPost)

		// 这里不用继续向上更新了, 因为 for.Post 包含 try 调用, 所以 parent 已经是正确的
		r.tryNodes[forStmt.Body] = true
	} else if cpCnt == 0 {
		// 处理 post 死代码的情况
		deadCode := &ast.IfStmt{
			Cond: constFalse(),
			Body: &ast.BlockStmt{
				List: []ast.Stmt{tryCallPost},
			},
		}
		forStmt.Body.List = append(forStmt.Body.List, deadCode)

		// 这里不用继续向上更新了, 因为 for.Post 包含 try 调用, 所以 parent 已经是正确的
		r.tryNodes[forStmt.Body] = true
		r.tryNodes[deadCode] = true
		r.tryNodes[deadCode.Body] = true
	}
}

func (r *fileRewriter) mkPostFnLit(post ast.Stmt, enclosingFn fnNode) (varPostFn ast.Expr, postFnAssign ast.Stmt) {
	r.importRT = true // for rt.𝗘𝗿𝗿𝗼𝗿

	varPostFn = r.genPostId(enclosingFn)
	postFnLit := &ast.FuncLit{ // try!
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{{
					// 这里 func() (_ error) {} 处理成 _ 是为了 body 最后不用 return nil
					// 从而避免检查 nil 是否 shadow
					Names: []*ast.Ident{ast.NewIdent("_")},
					// error 可能被 shadow 重新定义, e.g.
					// type error = int
					// 𝗽𝗼𝘀𝘁𝟭 := func() (_ error) { ... }
					// 所以这里 ref rt 中的 error 别名
					Type: ast.NewIdent(rtErrorTyName), // rt.E𝗿𝗿𝗼𝗿
				}},
			},
		},
		Body: &ast.BlockStmt{ // try!
			List: []ast.Stmt{
				post, // try!
				&ast.ReturnStmt{ /*Results: []ast.Expr{ ast.NewIdent("nil") }*/ },
			},
		},
	}
	postFnAssign = &ast.AssignStmt{ // try!
		Lhs: []ast.Expr{varPostFn},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{postFnLit},
	}

	r.tryNodes[postFnLit] = true
	r.tryNodes[postFnLit.Body] = true
	r.tryNodes[postFnAssign] = true

	postFnSig := types.NewSignatureType(
		nil,
		nil,
		nil,
		nil,
		types.NewTuple(types.NewVar(
			token.NoPos,
			r.pkg.Types,
			"",
			r.errTy,
		)),
		false,
	)
	r.pkg.UpdateType(varPostFn, postFnSig)
	r.pkg.UpdateType(postFnLit, postFnSig)
	return
}

func (r *fileRewriter) mkTryCallPost(varPostFn ast.Expr) (tryCallPost *ast.ExprStmt) {
	// 不能直接展开成 err:=post();  if err != nil { return err } 形式
	// 需要处理 defer err 的逻辑, 所以改成行, try.Try0(post()), 而 Try0 展开统一到 rewrite 阶段处理
	// 但是带来另一个问题, 新生成的代码需要标记
	// 1. tryNode 信息, 后续 rewrite 会使用
	// 2. caller 节点需要更新类型信息, 后续 rewrite 会使用

	// 这里不怕 try 或者 Try0 已经被重新定义
	// 这里会手动 Try0 的 type.Object 信息, rewrite 会被改写消除
	// e.g.,
	// 	for a := 1; a > 0; a = Try(ret1Err[int]()) {
	//		Try0 := 1
	//		_ = Try0
	//	}

	var tryFn ast.Expr
	try0Id := ast.NewIdent("Try0")
	x := imports.ImportSpec(r.f, pkgTryPath).Name
	switch {
	case x == nil:
		tryFn = &ast.SelectorExpr{
			X:   ast.NewIdent("try"),
			Sel: try0Id,
		}
	case x.Name == ".":
		tryFn = try0Id
	default:
		tryFn = &ast.SelectorExpr{
			X:   ast.NewIdent(x.Name),
			Sel: try0Id,
		}
	}

	callPost := &ast.CallExpr{Fun: varPostFn}
	tryCallPost = &ast.ExprStmt{ //try!
		X: &ast.CallExpr{ // try!
			Fun: tryFn, // try!
			Args: []ast.Expr{
				callPost,
			},
		},
	}

	r.tryNodes[tryCallPost] = true
	r.tryNodes[tryCallPost.X] = true

	// 后续改写 try call 需要检查返回值, ref retCntOfFnExpr
	// 所以这里需要更新下生成的 post 调用的类型信息
	r.pkg.UpdateType(callPost, r.errTy)

	// 后续改写 try call 需要检查 callee 是否是 try func call
	// 所以这里需要更新 try func 的类型信息
	try0Obj := r.tryCallObject(tryFnNames[0])
	if false { // 暂时没用到
		r.pkg.UpdateType(tryFn, try0Obj.Type())
		r.pkg.UpdateType(try0Id, try0Obj.Type())
	}
	r.pkg.UpdateUses(tryFn, try0Obj)
	return
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
