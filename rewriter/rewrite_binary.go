package rewriter

import (
	"go/ast"
	"go/token"
)

func (r *fileRewriter) rewriteBinaryExpr(ctx *rewriteCtx, n *ast.BinaryExpr) (ast.Expr, []ast.Stmt) {
	switch n.Op {
	case token.LOR, token.LAND:
		return r.rewriteShortCircuitExpr(ctx, n)
	default:
		x, xs := r.rewriteExpr(ctx, n.X)
		y, ys := r.rewriteExpr(ctx, n.Y)
		n1 := &ast.BinaryExpr{
			X:  x,
			Op: n.Op,
			Y:  y,
		}
		if n.Op == token.QUO || n.Op == token.SHL {
			// 可能会发生 runtime panic, 需要 fix 求值顺序
			// e.g.
			// 42/zero + Try(ret1Err[int]())
			// ----------------  wrong ----------------
			// 	𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := ret1Err[int]()
			//	if 𝗲𝗿𝗿𝟭 != nil {
			//		err = 𝗲𝗿𝗿𝟭
			//		return
			//	}
			//	42/zero + 𝘃𝗮𝗹𝟭 // 改变了语义, 会首先发生 panic, 不会先对 rhs 求值
			// ----------------  right ----------------
			// 	𝘃𝗮𝗹𝟭 := 42 / zero
			//	𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟭 := ret1Err[int]()
			//	if 𝗲𝗿𝗿𝟭 != nil {
			//		err = 𝗲𝗿𝗿𝟭
			//		return
			//	}
			//	𝘃𝗮𝗹𝟭 + 𝘃𝗮𝗹𝟮
			id, assign := r.assign(ctx, n1)
			return id, concat(xs, ys, sliceOf(assign))
		}
		return n1, append(xs, ys...)
	}
}

func (r *fileRewriter) rewriteShortCircuitExpr(ctx *rewriteCtx, n *ast.BinaryExpr) (ast.Expr, []ast.Stmt) {
	lhs, xs := r.rewriteExpr(ctx, n.X)
	rhs, ys := r.rewriteExpr(ctx, n.Y)

	// lhs 不影响求值顺序, 不需要处理短路逻辑
	if len(xs) == 0 {
		return &ast.BinaryExpr{
			X:  lhs,
			Op: n.Op,
			Y:  rhs,
		}, ys
	}

	// lhs 如果已经是 ident, 就不用另外生成 ident
	// 如果不是, 需要生成 ident, 提取赋值语句
	lhsId, ok := lhs.(*ast.Ident)
	if !ok {
		lhsId = r.genValId(ctx.fun)
		las := r.simpleAssign(lhsId, lhs, token.DEFINE)
		xs = append(xs, las)
	}

	// rhs 不影响求值顺序, 不需要处理短路逻辑
	if len(ys) == 0 {
		return &ast.BinaryExpr{
			X:  lhsId,
			Op: n.Op,
			Y:  rhs,
		}, xs
	}

	//		var 𝘃𝗮𝗹𝟯 bool
	//		𝘃𝗮𝗹𝟭 := id(true)
	//		if 𝘃𝗮𝗹𝟭 {
	//			𝘃𝗮𝗹𝟮 := id(false)
	//			𝘃𝗮𝗹𝟯 = 𝘃𝗮𝗹𝟮
	//		}
	// ===>
	// 		var 𝘃𝗮𝗹𝟮 bool
	//		𝘃𝗮𝗹𝟭 := id(true)
	//		if 𝘃𝗮𝗹𝟭 {
	//			𝘃𝗮𝗹𝟮 = id(false)
	//		}
	// 特殊 case 变量复用, 避免生成过多临时变量, 多余的赋值语句
	genRhsAssign := true
	if rhsId, ok := rhs.(*ast.Ident); ok && len(ys) == 1 {
		if assign, ok := ys[0].(*ast.AssignStmt); ok {
			if assign.Tok == token.DEFINE {
				if len(assign.Lhs) == 1 && assign.Lhs[0] == rhsId {
					assign.Tok = token.ASSIGN
					assign.Lhs[0] = lhsId
					genRhsAssign = false
				}
			}
		}
	}

	if genRhsAssign {
		// 𝘃𝗮𝗹𝟲 := 𝘃𝗮𝗹𝟰 && answer > 100
		//	if 𝘃𝗮𝗹𝟲 {
		//		𝘃𝗮𝗹𝟱, 𝗲𝗿𝗿𝟮 := ret1Err[bool]()
		//		if 𝗲𝗿𝗿𝟮 != nil {
		//			return 𝗲𝗿𝗿𝟮
		//		}
		//		𝘃𝗮𝗹𝟲 = 𝘃𝗮𝗹𝟱
		//	}
		// 这里的赋值不能省略, 因为
		// ys stmts 会被放在 if 内, 产生新的 block
		// 子作用域 block 内的变量会 shadow
		// 所以没法用 𝘃𝗮𝗹𝟲 替代 𝘃𝗮𝗹𝟱, 需要生成赋值语句
		ras := r.simpleAssign(lhsId, rhs, token.ASSIGN)
		ys = append(ys, ras)
	}

	// 处理短路逻辑, ys 需要把 wrap if lhsId
	var cond ast.Expr = lhsId
	if n.Op == token.LOR {
		cond = &ast.UnaryExpr{
			Op: token.NOT,
			X:  lhsId,
		}
	}

	xs = append(xs, &ast.IfStmt{
		Cond: cond,
		Body: &ast.BlockStmt{
			List: ys,
		},
	})

	return lhsId, xs
}
