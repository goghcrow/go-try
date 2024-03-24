package rewriter

import (
	"go/ast"
	"go/token"

	"github.com/goghcrow/go-matcher"
	"github.com/goghcrow/go-matcher/combinator"
)

// 1. 处理 labeled for/range/switch/select
// 		prewrite_labeled.go
// 2. 处理 其他 block/body 中的 labeledStmt
// 		rewrite_labeled.go
// 		rewrite_others.go#rewriteStmtList
// 共同处理
// 		try 展开向前提取的 assign stmt 不能越过 label
//		因为会导致 goto 跳转到 label skip 掉 try 产生的新 stmt
//
// e.g.
// 	L:
//		var a = Try(ret1Err[int]())
//		goto L
//		println(a)
// ==> ❌
//		𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := ret1Err[int]()
//		if 𝗲𝗿𝗿𝟭 != nil {
//			return 𝗲𝗿𝗿𝟭
//		}
// 	L:
//		var a = 𝘃𝗮𝗹𝟭
//		goto L
//		println(a)
// ==> ✅
// 	L:
//		𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := ret1Err[int]()
//		if 𝗲𝗿𝗿𝟭 != nil {
//			return 𝗲𝗿𝗿𝟭
//		}
//		var a = 𝘃𝗮𝗹𝟭
//		goto L
//		println(a)
//
// e.g.
// 	L:
//		for range Try(ret1Err[[]int]()) {
//			break L
//			continue L
//			goto L
//		}
// 		println(42)
// ==> ❌
// 		𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := ret1Err[[]int]()
//		if 𝗲𝗿𝗿𝟭 != nil {
//			return 𝗲𝗿𝗿𝟭
//		}
// 	L:
//		for range 𝘃𝗮𝗹𝟭 {
//			break L
//			continue L
//			goto L
//		}
// 		println(42)
// 改写后,
// 	break L 之后跳转到 println(42) 继续执行
// 	continue L 要么重新执行 for-body 要么跳到 println(42) 继续执行
// 	goto L 之后会重新执行 L.Stmt, 上面是重新执行 Try(f)
//	但是改写后, Try(f) 被提取到 label 之前, 跳转后不会重新执行 Try(f), 语义错误
//  所以, 对于 同样是 labeledStmt, goto 和 continue/break 的处理方式不一样
// ==> ✅
// 	𝗟_𝗚𝗼𝘁𝗼_𝗟𝟭:
//		{
//			𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := ret1Err[[]int]()
//			if 𝗲𝗿𝗿𝟭 != nil {
//				return 𝗲𝗿𝗿𝟭
//			}
//		L:
//			for range 𝘃𝗮𝗹𝟭 {
//				break L
//				goto 𝗟_𝗚𝗼𝘁𝗼_𝗟𝟭
//				continue L
//			}
//		}

// golang spec
// https://go.dev/ref/spec#Label_scopes
// Labels are declared by labeled statements and are used in the "break", "continue", and "goto" statements.
// It is illegal to define a label that is never used. In contrast to other identifiers,
// labels are not block scoped and do not conflict with identifiers that are not labels.
// The scope of a label is the body of the function in which it is declared and excludes the body of any nested function.
//
// https://go.dev/ref/spec#Break_statements
// If there is a label, it must be that of an enclosing "for", "switch", or "select" statement, and that is the one whose execution terminates.
// https://go.dev/ref/spec#Continue_statements
// If there is a label, it must be that of an enclosing "for" statement, and that is the one whose execution advances.
//
// prerewriteLabeled 处理 goto/break/continue -> labeld for/range/switch/select
// rewriteLabeled 处理 其他 goto -> label
//
//
// try 改写的赋值语句必须提到前面
// 但是规范要求 break/continue label 必须 enclosing for/range/switch/select
// 换句话说, , break/continue 目标只能是 labeldStmt {stmt: for/range/switch/select}
// 		label 和 for/range/switch/select 之前不能有其他 stmt
// goto 只收到 block 限制, 不受到 enclosing 限制
//
// 所以
// 	1. 存在 try 语句改写的 stmt 需要往前放
// 	2. 同时存在 goto 和 break/continue 跳转到相同目标 label
// 满足这两条, 因为 try 改写产生了新的 stmt 需要需要插入到 for/range/switch/select 之前, labeled 之后,
// 产生冲突, 如果放到 labeled 之前, 语义错误, 放到之间 break/continue 语法错误
// 必须改写
//
// 	1. L 只有 goto  =>  LabledStmt, 套个 Block,
// 		将提取的assign 和 for/range/switch/select 打包作为新的 labelelStmt.Stmt
//		e.g.  range 为例
//		L:
//			for range try(f()){ goto L }
//		==> ✅
//		L:
//			{
//				v, err = f()
//				if err != nil { return err }
//				for range v { goto L }
//			}
//		这里不会产生符号冲突/缺失的问题, 因为 range-body 有自己独立的作用域, 提出的 assign sym 本身是 uniq 的
//		但是 rewriteLabeled() 如果也采用套个新的 BlockStmt{} 则会出现符号冲突/缺失问题
//		e.g.
// 		L:
//			var a = Try(ret1Err[int]())
//			goto L
//			println(a)
// 		==> ❌
// 		L:
//			{
//				𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := ret1Err[int]()
//				if 𝗲𝗿𝗿𝟭 != nil {
//					return 𝗲𝗿𝗿𝟭
//				}
//				var a = 𝘃𝗮𝗹𝟭
//			}
//			goto L
//			println(a)
//		L: var a = Try(ret1Err[int]()) 构成一条 LabeledStmt, 套个 Block, 会改变作用域 !!!
// 	2. L 有 goto, 同时有 break 或 continue
//		block 必须套, 否则 goto 的位置不对, 会 skip 提取的 assign 的再次执行 (e.g. range value)
//		但是套了 block (改 label.Stmt), break / continue 的跳转就错了, 因为 labeld for/range/switch/select 结构被破坏
//		所以, 要么改 break/continue, 要么改 goto
//		2.1 改 goto, 外层套个 block, 生成一个新的 goto 专用Label, goto 到这个新 label
//		2.2 改 break/continue, 外层套 block 的同时
//			break 可以改成 goto, 人工指向循环结构结尾,
//			continue, for 需要跳转到 post, range 则先要修改为 for 循环, 比较麻烦
// 		所以采用冲突则修改 goto 的方案

func (r *fileRewriter) preRewriteLabeled() {
	// golang spec
	// A "break" statement terminates execution of the innermost
	// "for", "switch", or "select" statement within the same function.
	// A "continue" statement begins the next iteration of the innermost enclosing "for" loop
	// by advancing control to the end of the loop block. The "for" loop must be within the same function.
	// break 用在 for/range/switch/select 中, continue 用在 for/range 中
	ptn := &labelSt{
		Stmt: matcher.MkPattern[matcher.StmtPattern](r.m, func(n ast.Node, _ *matcher.MatchCtx) bool {
			// 是否需要在 for/ range / switch / select 前插入 stmt
			switch n := n.(type) {
			case *ast.ForStmt:
				// 参见 rewrite_others.go 对 *ast.ForStmt 改写
				// 只有 n.Init 包含 try 调用, 才会产生新 stmt, n.Cond, n.Body 只会影响 n.body
				return r.tryNodes[n.Init]
			case *ast.RangeStmt:
				return r.tryNodes[n.X] // 参见 rewriteRange, 同上
			case *ast.SwitchStmt:
				return r.tryNodes[n.Init] || r.tryNodes[n.Tag] // 参见 rewriteSwitch, 同上
			case *ast.SelectStmt:
				return r.tryNodes[n.Body] // 参见 rewriteSelect, 同上
			default:
				return false
			}
		}),
	}

	r.match(ptn, func(c cursor, ctx mctx) {
		lbStmt := c.Node().(*labelSt)
		inFn := ctx.EnclosingFunc()

		// collect 跳转到 label 的 branch 语句, 并记录 goto branch 个数
		branchTo := map[*branchSt]bool{}
		gotoCnt := 0
		r.m.Match(r.pkg.Package, &branchSt{
			Tok: combinator.Wildcard[combinator.TokenPattern](r.m),
		}, inFn, func(c *matcher.Cursor, ctx *matcher.MatchCtx) {
			br := c.Node().(*branchSt)
			jTbl := r.jmpTbl(inFn)
			if jTbl.JumpTo(br, lbStmt.Stmt) || jTbl.JumpTo(br, lbStmt) {
				branchTo[br] = true
				if br.Tok == token.GOTO {
					gotoCnt++
				}
			}
		})

		// 三种情况
		switch {
		case gotoCnt == 0:
			// 1. 没有 goto 跳过去, 只有 break/continue 不要该
		case len(branchTo)-gotoCnt == 0:
			// 2. 只有 goto 跳过去, 只需要套个 block
			lbStmt.Stmt = &ast.BlockStmt{
				List: []ast.Stmt{
					lbStmt.Stmt,
				},
			}
			// 更新 tryNodes 信息, labelStmt 包含 try, 只可能是 .Stmt 包含 try
			r.tryNodes[lbStmt.Stmt] = true
		case len(branchTo)-gotoCnt > 0:
			// 3. goto / break / continue 都会跳过去
			// 则需要生成新的 gotoL, 并替换跳转到原 label 的 branch (可能存在嵌套子结构)
			gotoL := r.genSym(inFn, labelPrefix+"Goto_"+lbStmt.Label.Name)
			n1 := &labelSt{
				Label: gotoL,
				Stmt: &ast.BlockStmt{
					List: []ast.Stmt{
						lbStmt,
					},
				},
			}
			// 更新 tryNodes 信息
			r.tryNodes[n1] = true
			c.Replace(n1)

			for br := range branchTo {
				if br.Tok == token.GOTO {
					br.Label = gotoL
				}
			}

			r.clearJmpTbl(inFn) // 修改了跳转控制流, 清空缓存
		}
	})
}
