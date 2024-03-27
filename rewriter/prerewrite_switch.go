package rewriter

import (
	"go/ast"
	"go/token"

	"github.com/goghcrow/go-matcher"
)

// golang spec
// https://go.dev/ref/spec#Expression_switches
// In an expression switch, the switch expression is evaluated and the case expressions,
// which need not be constants, are evaluated left-to-right and top-to-bottom;
// the first one that equals the switch expression triggers execution of the statements of the associated case;
// the other cases are skipped. If no case matches and there is a "default" case, its statements are executed.
// There can be at most one default case and it may appear anywhere in the "switch" statement.
// A missing switch expression is equivalent to the boolean value true.

// 先求值 SwitchStmt.Init,
// 再求值 SwitchStmt.Tag (switch expr)
// 再从左到右, 从上到下求值 case expr, 直到找到第一个 与 Switch.Tag 相等的 case, 执行 CaseClause.List

// SwitchStmt.Tag 如果不存在则默认为 true, (实际上就是匹配为 true 的 case)
// 最多一个 default, default 可以出现在任意 case 列表, 并不遵循从上到下的顺序的规则

// SwitchStmt.Tag 不能为常量
// SwitchStmt.Tag 如果是 untype, 会先隐式转换为defaylt type
// CaseClause.List 中的 expr 如果为 untyped, 会先隐式转换为 SwitchStmt.Tag 的类型
// SwitchStmt.Tag 必须可比较

// a := 1
//	switch a {
//	default:
//		println("default")
//	case 1:
//		println("1")
//	}
// => 1

// 比较吊诡的是, default 不遵循自然源码顺序执行, 但是 default 分支的 fallthrough 尊重书写顺序
// (感觉 default 分支 放最后可能是个正常的设计, 最后的分支不支持 fallthrough)
// a := 1
// switch a {
//	default:
//		println("default")
//		fallthrough
//	case 2:
//		println("2")
//	}
// => default 2

// In a case or default clause, the last non-empty statement may be a (possibly labeled) "fallthrough" statement
// to indicate that control should flow from the end of this clause to the first statement of the next clause.
// Otherwise control flows to the end of the "switch" statement.
// A "fallthrough" statement may appear as the last statement of all but the last clause of an expression switch.

// fallthrough 只能出现在最后一个case 分支最后一个 stmt, 或者 最后一个 stmt 是 LabedStmt, 且 LabedStmt.Stmt 是 fallthrough
// fallthrouth 进入下一个子句会跳过下一个子句的 case expr, 直接执行下一个子句的 case body

// e.g.
// switch {
//	default:
//		println("default")
//		fallthrough
//	case false:
//		println("false")
//	}
// => default false

type (
	branchSt = ast.BranchStmt
	labelSt  = ast.LabeledStmt
	switchSt = ast.SwitchStmt
	casec    = ast.CaseClause
)

// preRewrite 中新创建的节点都需要维护 tryNodes 信息,
// 后续 rewrite 需要依赖快速判断节点或者子节点是否包含 try 调用

func (r *fileRewriter) preRewriteSwitch() {
	var (
		doRewrite = func(inFn fnNode, label *ast.Ident, sStmt *switchSt) (*ast.BlockStmt, map[*branchSt]bool) {
			block, brkToL, brk2go := r.preRewriteBreakToGotoInSwitch(label, sStmt, inFn)
			block = r.doPreRewriteSwitch(sStmt, inFn)

			// 如果存在从 break 改成 goto 的 branchStmt, 则需要在 block 后面加上对应 label
			if brkToL != nil {
				block.List = append(block.List, &labelSt{
					Label: brkToL,
					Stmt:  &ast.EmptyStmt{}, // 注意, labelStmt 的 stmt 不能为空
				})
			}
			return block, brk2go
		}
		// f 内部是否有跳转到 switchStmt 的 stmt
		labelUsedInFn = func(sStmt *switchSt, f fnNode, exclude map[*branchSt]bool) bool {
			jTbl := r.jmpTbl(f)
			for b := range jTbl {
				if !exclude[b] && jTbl.JumpTo(b, sStmt) {
					return true
				}
			}
			return false
		}
	)

	// case 中 (非 case body) 包含 try 调用的 switch
	switchPtn := &switchSt{
		Body: matcher.MkPattern[matcher.BlockStmtPattern](r.m, func(n ast.Node, ctx mctx) bool {
			for _, cs := range n.(*ast.BlockStmt).List {
				for _, it := range cs.(*casec).List {
					if r.tryNodes[it] {
						return true
					}
				}
			}
			return false
		}),
	}

	labeledSwitchPtn := &labelSt{Stmt: switchPtn}

	// 1. 先处理 labeled 的 swtich
	switchCache := map[*switchSt]bool{}
	r.match(labeledSwitchPtn, func(c cursor, ctx mctx) {
		lStmt := c.Node().(*labelSt)
		sStmt := lStmt.Stmt.(*switchSt)
		switchCache[sStmt] = true

		inFn := ctx.EnclosingFunc()
		assert(inFn != nil)

		block, brk2go := doRewrite(inFn, lStmt.Label, sStmt)

		// switch 附带的 label 除了 break 有使用, 可能还有 goto 在使用
		// 因为 break 都改成跳转到 改写后的 block 末尾 (label 名称不同)
		// 如果还存在 goto 到原 label, 不能一并改写, 语义不同, brk 跳转到结尾, goto 跳转到开头
		// 如果 label 没有 goto 在引用, 则删除, 有 goto 在引用, 则保留
		if labelUsedInFn(sStmt, inFn, brk2go) {
			// 不是 break 在用, 且存在跳转到 switch 的 label
			labelBlock := &labelSt{
				Label: lStmt.Label,
				Stmt:  block,
			}
			r.tryNodes[labelBlock] = true
			c.Replace(labelBlock)
			r.clearJmpTbl(inFn)
		} else {
			c.Replace(block)
		}
	})

	// 2. 再处理 labeledless switch
	r.match(switchPtn, func(c cursor, ctx mctx) {
		sStmt := c.Node().(*switchSt)
		if switchCache[sStmt] {
			return
		}
		inFn := ctx.EnclosingFunc()
		assert(inFn != nil)
		block, _ := doRewrite(inFn, nil, sStmt)
		c.Replace(block)
	})
}

func (r *fileRewriter) preRewriteBreakToGotoInSwitch(
	label *ast.Ident,
	sStmt *switchSt,
	inFn fnNode,
) (
	block *ast.BlockStmt,      // 改写结果
	brkToL *ast.Ident,         // break 跳转 label
	brk2go map[*branchSt]bool, // 记录改写的 break branch
) {
	// 如果是 labeled switch
	// 则所有 break 到该 label 的 stmt (包括嵌套的 子 block 的 break ),
	// 需要改成成使用 goto 跳转到 switch 改写的 if stmt 的结尾
	//
	// label:
	//	switch {
	//	case true:
	//		break label
	//		for { break label }
	//	}
	// =>
	//	{
	//		if true {
	//			goto brk
	//			for { goto brk }
	//		}
	// brk:
	// }
	//
	// 如果 switch 改写成了 if, 则 labeless 的 break 也要改写成 goto,
	//	switch {
	//		case true:
	//			break
	//	}
	// =>
	//	{
	//		if true {
	//			goto brk
	//		}
	//		brk:
	//	}

	brk2go = map[*branchSt]bool{}
	jTbl := r.jmpTbl(inFn)

	// 因为 switch 被改写成了 if, 所以需要把 break 改写成 goto
	// 匹配所有跳转到 switch 的 break, 包括
	// 	1. 有 label (自身 或者 嵌套 控制流内的 break)
	// 	2. 和 无 label (自身 case 内部的 break)
	// 生成新 label, 新建 goto branch 指过去,
	// 并将 break 所在的 branchStmt  记录在 brk2go 中,
	// 用来计算, 改写后是否还有其他 控制流跳转到 switch 的 label (如果 label != nil )
	// 进而判断是否将该 label 移除, 否则未使用 label 会报错

	brkPtn := &branchSt{
		Tok: token.BREAK,
		Label: matcher.MkPattern[matcher.IdentPattern](r.m, func(n ast.Node, ctx *matcher.MatchCtx) bool {
			isNilNode := matcher.IsNilNode(n)
			if label == nil {
				// break
				return isNilNode
			} else {
				// break Label
				return !isNilNode && n.(*ast.Ident).Name == label.Name
			}
		}),
	}

	r.m.Match(r.pkg.Package, brkPtn, r.f, func(c cursor, ctx mctx) {
		brk := c.Node().(*branchSt)
		if jTbl.JumpTo(brk, sStmt) {
			brk2go[brk] = true

			if brkToL == nil {
				if label == nil {
					brkToL = r.genSym(inFn, labelPrefix+"BrkTo")
				} else {
					brkToL = r.genSym(inFn, labelPrefix+"BrkTo_"+label.Name)
				}
			}

			// 这里原地修改, 否则 jumpTable 的 branchStmt: Stmt 映射类型会不对
			c.Replace(&branchSt{
				Tok:   token.GOTO,
				Label: brkToL,
			})
			r.clearJmpTbl(inFn)
		}
	})
	return
}

func (r *fileRewriter) doPreRewriteSwitch(sStmt *switchSt, inFn ast.Node) *ast.BlockStmt {
	// 因为 case 有 try, 才需要改写, 所以肯定有至少一个 case 分支
	assert(len(sStmt.Body.List) > 0)

	var xs []ast.Stmt
	var initHasTry, tagHasTry, bodyHasTry bool

	// 改写 init
	if sStmt.Init != nil {
		xs = append(xs, sStmt.Init)
		if r.tryNodes[sStmt.Init] {
			initHasTry = true
		}
	}

	// 改写 tag
	var tag *ast.Ident
	if sStmt.Tag != nil {
		tag = r.genValId(inFn)
		declTag := r.simpleAssign(tag, sStmt.Tag, token.DEFINE)
		xs = append(xs, declTag)
		if r.tryNodes[sStmt.Tag] {
			r.tryNodes[declTag] = true
			tagHasTry = true
		}
	}

	// 改写 body
	csr := r.mkCaseRewriter(sStmt.Body.List, tag)
	iff, bodyHasTry := csr.rewriteCasesToIfStmt()

	block := &ast.BlockStmt{
		List: append(xs, iff),
	}
	if initHasTry || tagHasTry || bodyHasTry {
		r.tryNodes[block] = true
	}
	return block
}

type caseRewriter struct {
	fr         *fileRewriter
	cs         []*casec
	bodyHasTry []bool
	conds      []ast.Expr
}

func (r *fileRewriter) mkCaseRewriter(xs []ast.Stmt, tag *ast.Ident) *caseRewriter {
	csr := &caseRewriter{
		fr: r,
		cs: sliceMap(xs, func(_ int, n ast.Stmt) *casec { return n.(*casec) }),
	}
	csr.bodyHasTry = sliceMap(csr.cs, func(i int, c *casec) bool {
		fst := findFirst(c.Body, func(n ast.Stmt) bool { return csr.fr.tryNodes[n] })
		return fst != nil
	})
	csr.conds = sliceMap(csr.cs, func(i int, c *casec) ast.Expr {
		return csr.rewriteCond(c.List, tag)
	})
	csr.eliminationFallthrough()
	return csr
}

func (r *caseRewriter) rewriteCasesToIfStmt() (_ ast.Stmt, outHasTry bool) {
	iff := &ast.IfStmt{}
	lastIf := iff
	var defaultClause *ast.BlockStmt
	for i, c := range r.cs {
		body := &ast.BlockStmt{
			List: c.Body,
		}
		if r.bodyHasTry[i] {
			r.fr.tryNodes[body] = true
			outHasTry = true
		}

		isDefault := c.List == nil
		if isDefault {
			defaultClause = body
		} else {
			cond := r.conds[i]
			els := &ast.IfStmt{
				Cond: cond,
				Body: body,
			}
			if r.fr.tryNodes[cond] || r.bodyHasTry[i] {
				r.fr.tryNodes[els] = true
				outHasTry = true
			}
			lastIf.Else = els
			lastIf = els
		}
	}
	if defaultClause != nil {
		lastIf.Else = defaultClause
	}

	assert(iff.Else != nil)
	return iff.Else, outHasTry
}

func (r *caseRewriter) eliminationFallthrough() {
	// 有两种思路来消除 fallthrough
	// 1. 最早采用的思路
	//		switch translate 成 if {} else if {} ... else {} 结构
	//		包含 fallthroug 的 case, 不生成 else 表示控制继续  if {}; if {}
	//		default 分支生成 else if true { } 结构
	//	但是犯了两个错误:
	//		第一个错误是, fallthroug 之后会直接跳过下一个 case 分支的判断,
	//			直接执行下一个 case 分支的 body, 语义错了
	//			translate 成 label 跳转不行, 因为 label 遵循 scope 规则, body 在子作用域
	//			if { goto L} else if { L: } 结构是不行的, 因为 前一个 if-body 是访问不到后一个 if-body 的符号表的
	//		第二个错误是 default 分支的处理
	//			golang default 可以出现在任意位置, 但是总是最后一个 test
	//			所以, 生成的 if 结构, 如果有 default 分支, 不管在什么位置, 都要调整到最后
	//			但是, fallthroug 到 default 分支, 控制流却遵循书写位置 (有点吊诡的特性)
	//			test 和 fallthroug 在行为上并不一致!!!
	//			所以按顺序生成 if 控制流的方式, 如果兼容语义, 务必会引入大量的 flag 以及跳转结构
	// 2. 现在的思路, 跟处理 for-post 一样, 直接无脑复制 fallthroug 代码块
	//			但要注意 case: fallthroug; case: fallthrough 这种连续 fallthroug 的情况
	//			需要从后往前复制, 保证 pre-body 复制的 post-post 已经是复制好的
	//			另外, 复制产生的作用域问题, 不同 case 分支的 body 原来的 scope 是兄弟关系, 隔离且拥有共同 parent
	// 			注意: 不能直接 copy fallthrough 下一个 case clause 的 body
	// 			一个错误的示例:
	//			  var x int
	//				switch {
	//				case Try(ret1Err[int]()) == 1:
	//					x := 1
	//					println(x)
	//					fallthrough
	//				case true:
	//					x = 2
	//				}
	//				println(x)
	//			 ========================>
	//			 	var x int
	//				{
	//					𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := ret1Err[int]()
	//					if 𝗲𝗿𝗿𝟭 != nil {
	//						return 𝗲𝗿𝗿𝟭
	//					}
	//					if 𝘃𝗮𝗹𝟭 == 1 {
	//						x := 1
	//						println(x)
	//						x = 2 // x = 2 本来是要修改 外层 x 的值, 这里修改的是内层 x 的值
	//					} else if true {
	//						x = 2
	//					}
	//				}
	//				println(x)
	//			 出现问题的原因
	//			 复制 case body 破坏了 scope 的父子关系, 兄弟的作用域 mixed 在一起
	//			 会产生符号冲突, 以及符号 shadow 问题
	//			 所以, copy case body 需要保持原来 scope 关系

	n := len(r.cs)

	bodies := make([][]ast.Stmt /*BlockStmt*/, n)
	bodiesHasTry := make([]bool, n)

	// 复制阶段:
	// 因为 fallthrough 是正向不中断控制流, 所以
	// 通过逆向遍历, 来把 fallthrough 的代码块都向跳转来源复制
	// 为了保证 scope 正确, 不直接复制 []Stmt, 而是都打包成 BlockStmt 向上复制
	for i := n - 1; i >= 0; i-- {
		endWithFall, body := r.rtrimFallthrough(r.cs[i].Body)

		if len(body) > 0 {
			block := &ast.BlockStmt{
				List: body,
			}
			if r.bodyHasTry[i] {
				r.fr.tryNodes[block] = true
				bodiesHasTry[i] = true
			}
			bodies[i] = []ast.Stmt{block}
		}

		if endWithFall {
			// 假设语法正常, 最后一个 case clause 不能是 fallthrough
			assert(i+1 < n)
			// 合并 body 并更新 try 状态
			// try = body[n+1].try? || body[n].try?
			bodies[i] = concat(bodies[i], bodies[i+1])
			if bodiesHasTry[i+1] {
				r.bodyHasTry[i] = true
			}
		}
	}

	// 合并阶段:
	for i, xs := range bodies {
		if len(xs) == 1 {
			r.cs[i].Body = xs[0].(*ast.BlockStmt).List
			delete(r.fr.tryNodes, xs[0])
		} else {
			// 不同的 case clause body 之间
			// 都 wrap block 保证 scope 隔离, 且父子关系正确
			r.cs[i].Body = xs
		}
		r.bodyHasTry[i] = bodiesHasTry[i]
	}

	// 上面错误示例的正确转换如下
	// 	var x int // Scope1
	//	{
	//		𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := ret1Err[int]()
	//		if 𝗲𝗿𝗿𝟭 != nil {
	//			return 𝗲𝗿𝗿𝟭
	//		}
	//		if 𝘃𝗮𝗹𝟭 == 1 {
	//			{
	//				x := 1
	//				println(x)
	//			}
	//			{
	//				x = 2 // x 引用 scope1
	//			}
	//		} else if true {
	//			x = 2 // x 引用 scope1
	//		}
	//	}
	//	println(x)
}

func (r *caseRewriter) rtrimFallthrough(xs []ast.Stmt) (_ bool, stmtsWithoutFallthrough []ast.Stmt) {
	xs = trimTrailingEmptyStmts(xs)
	if len(xs) == 0 {
		return false, xs
	}

	init, lst := xs[:len(xs)-1], xs[len(xs)-1]
	switch n := lst.(type) {
	case *branchSt:
		if n.Tok == token.FALLTHROUGH {
			return true, init
		}
	case *labelSt:
		// case ...:
		//	L:
		//		println()
		// 		fallthrough // last
		// ----------------
		// case ...:
		//	L:				// last
		// 		fallthrough // fallthrough 是 LabelStmt 的 child
		b, ok := n.Stmt.(*branchSt) // 假设语法正确
		// 这里 label 不能去掉, 因为 body 中有跳转的来源
		if ok && b.Tok == token.FALLTHROUGH {
			n.Stmt = &ast.EmptyStmt{}
			return true, xs
		}
	}
	return false, xs
}

// xs: case 子句的表达式列表 或 类型列表
// tag: switch 的 tag expr
func (r *caseRewriter) rewriteCond(xs []ast.Expr, tag *ast.Ident) ast.Expr {
	if xs == nil { // default clause
		return nil
	}

	// x: case 子句的表达式列表 或 类型列表 中的单个元素
	tagWith := func(caseV ast.Expr) ast.Expr {
		if tag == nil {
			// golang spec
			// A missing switch expression is equivalent to the boolean value true.
			// 省略 tag, tag 为 true,
			// case 表达式 实际在对 case 的值与 tag 进行判等操作不变
			// 不能理解成 switch { caes bool_expr: }
			// 而是 swithc { case expr == true: }
			// 所以 case expr 是可以为 any 类型的, 而不一定是 true

			// golang spec
			// https://go.dev/ref/spec#Comparison_operators
			// Interface types that are not type parameters are comparable.
			// Two interface values are equal if they have identical dynamic types and equal dynamic values or if both have value nil.
			// A comparison of two interface values with identical dynamic types causes a run-time panic if that type is not comparable.
			// any 类型判等编译期也没法检查类型, 所以变成运行时检查, 类型不对 runtime panic

			// switch on implicit bool converted to interface
			// was broken: see issue 3980
			//
			//	switch i := interface{}(true); {
			//	case i:
			//		assert(true, "true")
			//	case false:
			//		assert(false, "i should be true")
			//	default:
			//		assert(false, "i should be true")
			//	}
			// 错误的 translate
			//  i := interface{}(true)
			//	𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := lit(i)
			//	if 𝗲𝗿𝗿𝟭 != nil {  return 𝗲𝗿𝗿𝟭 }
			//	if 𝘃𝗮𝗹𝟭 { // 𝘃𝗮𝗹𝟭 是 any 类型 !!!
			//		assert(true, "true")
			//	} else { ...
			ty := r.fr.pkg.TypeOf(caseV)
			// 只有 bool 类型, 可以省略判等
			if !isAny(ty, true) {
				// 这里不支持别名
				assert(isBool(ty, false))
				return caseV
			}
		}

		//	swtich tag {
		//		case a,b: 	==> 	tag == a || tag == b
		//	}
		//	switch {
		//		case non_nil_expr1, non_nil_expr2:	==> non_nil_expr1 || non_nil_expr2
		//	}
		//	switch {
		//		case nil_expr1, nil_expr2:	==> expr1 == true || expr2 == true
		//	}

		var x ast.Expr
		if tag == nil {
			x = &ast.ParenExpr{
				X: constTrue(),
			}
		} else {
			x = tag
		}

		// 否则, 转换成 tag == x 形式
		eq := &ast.BinaryExpr{
			X:  x,
			Op: token.EQL,
			Y:  caseV,
		}

		// tag是 id, 不需要判断 tryNodes
		// 顺便更新 tryNodes 信息
		if r.fr.tryNodes[caseV] {
			r.fr.tryNodes[eq] = true
		}
		return eq
	}

	e := tagWith(xs[0])
	for i := 1; i < len(xs); i++ {
		x, y := e, tagWith(xs[i])
		e = &ast.BinaryExpr{X: x, Op: token.LOR, Y: y}
		if r.fr.tryNodes[x] || r.fr.tryNodes[y] {
			r.fr.tryNodes[e] = true
		}
	}
	return e
}
