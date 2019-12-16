# TaskPool

TaskPool is a simple worker pool to run functions in the background. 

### Usage

Define a Taskpool:

```go
pool := taskpool.NewInMemoryPool(2, 2)
pool.Start()
defer pool.Stop()
```

Create a Task ...

```go
task := taskpool.NewRunnableTask(
    func(outChan chan interface{}, errChan chan error) {
        errChan <- testError
    }
)
```
... and put it into the pool

```go
err := pool.Put(task)
if err != nil {
    panic(err)
}
```

To get a status, a result or an error call Get(). 
The call returns a result and an error.
In case this task wasn't processed yet, err returns taskpool.ErrInProgress

```go
result, err := task.Get()
```
or to wait for a result 
```go
result, err := task.GetWait(1 * time.Second)
```

### Install

```go
go get github.com/hazeglide/taskpool
```

### License

See [LICENSE](LICENSE)