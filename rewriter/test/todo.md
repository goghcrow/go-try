error propagating

// todo ast 中所有支持 Expr 和 Stmt (ExprStmt) 的地方都要检查一遍, 是否影响语义 !!!


optimize
Ø()
a, b = II(x, y)
a, b, c = III(x, y, z)

todo ast 中所有支持 Expr 和 Stmt (ExprStmt) 的地方都要检查一遍, 是否影响语义 !!!

call 不在 stmt 中的, edit1 的 indexOf 不能用, 比如 switch Try() {}

todo 泛型参数的 Try
import name 别名的 Try

todo 可以 hook 个日志打印 !!!
todo 检查 非 assign 的 call
todo Try[](xxx)
todo Try(T, error)
todo Try(f())
 todo 测试 自定义 error 类型, error 子类型
todo 测试 untyped nil
todo 测试类型报错, 返回值数量不对, 返回值类型不对

todo 全部 ast 节点试一下, 还有哪些 ???
func (s nodeStack) nearestStmt() *stmts {

todo try(...)...泛型版本
todo go try(...) / defer try(...)
