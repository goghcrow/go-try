ast 中所有支持 Expr 和 Stmt (ExprStmt) 的地方都要检查一遍, 是否影响语义 !!!
泛型参数的 Try
可以 hook 个日志打印 !!!
检查 非 assign 的 call
Try[](xxx)
Try(T, error)
Try(f())
测试 untyped nil
测试类型报错, 返回值数量不对, 返回值类型不对
try(...)...泛型版本