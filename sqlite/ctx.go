package sqlite

import (
	"context"
	"sync"
)

// Context 执行数据库的上下文
type Context struct {
	context.Context                    // 上下文
	sync.WaitGroup                     // 工作控制
	data            chan []interface{} // 数据通道
	cancelFunc      context.CancelFunc // 退出函数
	err             error              // 执行错误
}

// NewContext 上下文构造函数
func NewContext(parent context.Context) *Context {
	data := make(chan []interface{})
	if parent == nil {
		parent = context.Background()
	}
	context, cancelFunc := context.WithCancel(parent)
	return &Context{context, sync.WaitGroup{}, data, cancelFunc, nil}
}

// SetError 设置错误，如设置错误则所该任务取消
func (ctx *Context) SetError(err error) {
	ctx.err = err
	ctx.cancelFunc()
}

// CheckError 检查错误，如有错误则直接 panic
func (ctx *Context) CheckError() {
	if ctx.err != nil {
		panic(ctx.err)
	}
}

// Error 获取错误
func (ctx *Context) Error() error {
	return ctx.err
}

// CloseData 关闭数据通道，发送数据的协程结束后应调用此方法
func (ctx *Context) CloseData() {
	close(ctx.data)
}

// Data 获取数据通道
func (ctx *Context) Data() chan []interface{} {
	return ctx.data
}

// Done 任务完成
func (ctx *Context) Done() {
	ctx.WaitGroup.Done()
}

// Cancel 获取任务取消通道
func (ctx *Context) Cancel() <-chan struct{} {
	return ctx.Context.Done()
}
