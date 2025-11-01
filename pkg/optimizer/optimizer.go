package optimizer

import (
	"context"
	"runtime"
	"sync"
	"time"

	"github.com/projectdiscovery/gologger"
)

// WorkerPool manages concurrent processing of tasks
type WorkerPool struct {
	workers    int
	taskQueue  chan Task
	resultChan chan Result
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

// Task represents a unit of work
type Task interface {
	Execute() Result
	GetID() string
}

// Result represents the result of a task
type Result interface {
	GetError() error
	GetData() interface{}
	GetTaskID() string
}

// BasicTask implements the Task interface
type BasicTask struct {
	ID       string
	Function func() (interface{}, error)
}

func (t *BasicTask) Execute() Result {
	data, err := t.Function()
	return &BasicResult{
		TaskID: t.ID,
		Data:   data,
		Error:  err,
	}
}

func (t *BasicTask) GetID() string {
	return t.ID
}

// BasicResult implements the Result interface
type BasicResult struct {
	TaskID string
	Data   interface{}
	Error  error
}

func (r *BasicResult) GetError() error {
	return r.Error
}

func (r *BasicResult) GetData() interface{} {
	return r.Data
}

func (r *BasicResult) GetTaskID() string {
	return r.TaskID
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(workers int, queueSize int) *WorkerPool {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}
	
	ctx, cancel := context.WithCancel(context.Background())
	
	return &WorkerPool{
		workers:    workers,
		taskQueue:  make(chan Task, queueSize),
		resultChan: make(chan Result, queueSize),
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Start begins processing tasks
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
	
	gologger.Debug().Msgf("Started %d workers", wp.workers)
}

// Stop gracefully shuts down the worker pool
func (wp *WorkerPool) Stop() {
	close(wp.taskQueue)
	wp.wg.Wait()
	close(wp.resultChan)
	wp.cancel()
	
	gologger.Debug().Msg("Worker pool stopped")
}

// Submit adds a task to the queue
func (wp *WorkerPool) Submit(task Task) {
	select {
	case wp.taskQueue <- task:
	case <-wp.ctx.Done():
		gologger.Warning().Msgf("Cannot submit task %s: worker pool is shutting down", task.GetID())
	}
}

// Results returns the result channel
func (wp *WorkerPool) Results() <-chan Result {
	return wp.resultChan
}

// worker processes tasks from the queue
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	
	gologger.Debug().Msgf("Worker %d started", id)
	
	for {
		select {
		case task, ok := <-wp.taskQueue:
			if !ok {
				gologger.Debug().Msgf("Worker %d: task queue closed", id)
				return
			}
			
			gologger.Debug().Msgf("Worker %d processing task %s", id, task.GetID())
			result := task.Execute()
			
			select {
			case wp.resultChan <- result:
			case <-wp.ctx.Done():
				gologger.Debug().Msgf("Worker %d: context cancelled", id)
				return
			}
			
		case <-wp.ctx.Done():
			gologger.Debug().Msgf("Worker %d: context cancelled", id)
			return
		}
	}
}

// RateLimiter controls the rate of operations
type RateLimiter struct {
	limiter chan struct{}
	rate    time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requestsPerSecond int) *RateLimiter {
	if requestsPerSecond <= 0 {
		requestsPerSecond = 10 // Default rate
	}
	
	return &RateLimiter{
		limiter: make(chan struct{}, requestsPerSecond),
		rate:    time.Second / time.Duration(requestsPerSecond),
	}
}

// Wait blocks until it's safe to proceed
func (rl *RateLimiter) Wait() {
	rl.limiter <- struct{}{}
	go func() {
		time.Sleep(rl.rate)
		<-rl.limiter
	}()
}

// BatchProcessor processes items in batches
type BatchProcessor struct {
	batchSize int
	timeout   time.Duration
}

// NewBatchProcessor creates a new batch processor
func NewBatchProcessor(batchSize int, timeout time.Duration) *BatchProcessor {
	if batchSize <= 0 {
		batchSize = 100
	}
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	
	return &BatchProcessor{
		batchSize: batchSize,
		timeout:   timeout,
	}
}

// Process processes items in batches
func (bp *BatchProcessor) Process(items []interface{}, processor func([]interface{}) error) error {
	for i := 0; i < len(items); i += bp.batchSize {
		end := i + bp.batchSize
		if end > len(items) {
			end = len(items)
		}
		
		batch := items[i:end]
		
		ctx, cancel := context.WithTimeout(context.Background(), bp.timeout)
		
		done := make(chan error, 1)
		go func() {
			done <- processor(batch)
		}()
		
		select {
		case err := <-done:
			cancel()
			if err != nil {
				return err
			}
		case <-ctx.Done():
			cancel()
			gologger.Warning().Msgf("Batch processing timeout for batch %d-%d", i, end-1)
		}
	}
	
	return nil
}

// MemoryOptimizer helps manage memory usage
type MemoryOptimizer struct {
	maxMemoryMB int64
}

// NewMemoryOptimizer creates a new memory optimizer
func NewMemoryOptimizer(maxMemoryMB int64) *MemoryOptimizer {
	if maxMemoryMB <= 0 {
		maxMemoryMB = 512 // Default 512MB
	}
	
	return &MemoryOptimizer{
		maxMemoryMB: maxMemoryMB,
	}
}

// CheckMemoryUsage returns current memory usage
func (mo *MemoryOptimizer) CheckMemoryUsage() (int64, bool) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	currentMB := int64(m.Alloc) / 1024 / 1024
	exceeded := currentMB > mo.maxMemoryMB
	
	if exceeded {
		gologger.Warning().Msgf("Memory usage (%d MB) exceeded limit (%d MB)", currentMB, mo.maxMemoryMB)
	}
	
	return currentMB, exceeded
}

// ForceGC forces garbage collection
func (mo *MemoryOptimizer) ForceGC() {
	runtime.GC()
	gologger.Debug().Msg("Forced garbage collection")
}

// OptimizeForSpeed configures runtime for speed
func OptimizeForSpeed() {
	// Set GOMAXPROCS to use all available CPUs
	runtime.GOMAXPROCS(runtime.NumCPU())
	
	gologger.Debug().Msgf("Optimized for speed: GOMAXPROCS=%d", runtime.NumCPU())
}

// OptimizeForMemory configures runtime for memory efficiency
func OptimizeForMemory() {
	// Force initial GC for memory efficiency
	runtime.GC()
	
	gologger.Debug().Msg("Optimized for memory: forced GC")
}
