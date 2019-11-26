# spinlock（自旋锁）


```go
type SpinLock struct {
	f uint32
}
func (sl *SpinLock) Lock() {
	for !sl.TryLock() {
		runtime.Gosched()
	}
}
func (sl *SpinLock) Unlock() {
	atomic.StoreUint32(&sl.f, 0)
}
func (sl *SpinLock) TryLock() bool {
	return atomic.CompareAndSwapUint32(&sl.f, 0, 1)
}
```