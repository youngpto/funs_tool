package algorithm

import (
	"container/list"
	"github.com/youngpto/funs_tool/async"
	"math"
	"sync"
	"time"
)

var timeWheel *TimeWheel

// @author qiang.ou<qingqianludao@gmail.com>

// Job 延时任务回调函数
type Job func(arge interface{})

// TaskData 回调函数参数类型
type TaskData interface{}

// TimeWheel 时间轮
type TimeWheel struct {
	sync.Mutex
	interval          time.Duration       // 指针每隔多久往前移动一格
	accumulator       int64               // 累加器
	ticker            *time.Ticker        // 滴答计时器
	slots             []*list.List        // 时间轮槽
	timer             map[interface{}]int // key: 定时器唯一标识 value: 定时器所在的槽, 主要用于删除定时器
	currentPos        int                 // 当前指针指向哪一个槽
	slotNum           int                 // 槽数量
	addTaskChannel    chan Task           // 新增任务channel
	removeTaskChannel chan interface{}    // 删除任务channel
	stopChannel       chan bool           // 停止定时器channel
}

// Task 延时任务
type Task struct {
	delay  time.Duration // 延迟时间
	circle int           // 时间轮需要转动几圈
	key    interface{}   // 定时器唯一标识, 用于删除定时器
	job    Job           // 定时器回调函数
	data   TaskData      // 回调函数参数
}

func SetTimeWheel(t *TimeWheel) {
	timeWheel = t
}

func GetTimeWheel() *TimeWheel {
	return timeWheel
}

// NewTimeWheel 创建时间轮
func NewTimeWheel(interval time.Duration, slotNum int) *TimeWheel {
	if interval <= 0 || slotNum <= 0 {
		return nil
	}
	tw := &TimeWheel{
		interval:          interval,
		accumulator:       0,
		slots:             make([]*list.List, slotNum),
		timer:             make(map[interface{}]int),
		currentPos:        0,
		slotNum:           slotNum,
		addTaskChannel:    make(chan Task),
		removeTaskChannel: make(chan interface{}),
		stopChannel:       make(chan bool),
	}

	tw.initSlots()

	return tw
}

// 初始化槽，每个槽指向一个双向链表
func (tw *TimeWheel) initSlots() {
	for i := 0; i < tw.slotNum; i++ {
		tw.slots[i] = list.New()
	}
}

// Start 启动时间轮
func (tw *TimeWheel) Start(closeSig chan bool) {
	tw.ticker = time.NewTicker(tw.interval)
	for {
		select {
		case <-tw.stopChannel:
			tw.ticker.Stop()
			return
		case <-closeSig:
			tw.ticker.Stop()
			return
		case <-tw.ticker.C:
			tw.tickHandler()
		case task := <-tw.addTaskChannel:
			tw.addTask(&task)
		case key := <-tw.removeTaskChannel:
			tw.removeTask(key)

		}
	}
}

// Stop 停止时间轮
func (tw *TimeWheel) Stop() {
	tw.stopChannel <- true
}

// AddTimer 添加定时器 不可撤销
func (tw *TimeWheel) AddTimer(delay time.Duration, data TaskData, job Job) {
	if delay <= 0 {
		return
	}
	tw.addTaskChannel <- Task{delay: delay, data: data, job: job}
}

// AddTimerCustom 可以通过key来撤销一个未执行的定时器
func (tw *TimeWheel) AddTimerCustom(delay time.Duration, key interface{}, data TaskData, job Job) {
	if delay <= 0 {
		return
	}
	tw.addTaskChannel <- Task{delay: delay, key: key, data: data, job: job}
}

// RemoveTimer 删除定时器 key为添加定时器时传递的定时器唯一标识
func (tw *TimeWheel) RemoveTimer(key interface{}) {
	if key == nil {
		return
	}
	tw.removeTaskChannel <- key
}

func (tw *TimeWheel) tickHandler() {
	l := tw.slots[tw.currentPos]
	if tw.currentPos == tw.slotNum-1 {
		tw.currentPos = 0
	} else {
		tw.currentPos++
	}
	tw.scanAndRunTask(l)
}

// 扫描链表中过期定时器, 并执行回调函数
func (tw *TimeWheel) scanAndRunTask(l *list.List) {
	tw.Lock()
	defer tw.Unlock()
	for e := l.Front(); e != nil; {
		task := e.Value.(*Task)
		if task.circle > 0 {
			task.circle--
			e = e.Next()
			continue
		}

		next := e.Next()
		l.Remove(e)
		delete(tw.timer, task.key)
		e = next

		async.Go(func() {
			task.job(task.data)
		})
	}
}

// 新增任务到链表中
func (tw *TimeWheel) addTask(task *Task) {
	tw.Lock()
	defer tw.Unlock()
	if task.key == nil {
		if tw.accumulator >= math.MaxInt64 {
			tw.accumulator = 0
		}
		tw.accumulator++
		task.key = tw.accumulator
	}
	pos, circle := tw.getPositionAndCircle(task.delay)
	task.circle = circle

	tw.slots[pos].PushBack(task)

	tw.timer[task.key] = pos
}

// 获取定时器在槽中的位置, 时间轮需要转动的圈数
func (tw *TimeWheel) getPositionAndCircle(d time.Duration) (int, int) {
	delayMillisecond := int(d.Nanoseconds() / 1e6)
	intervalMillisecond := int(tw.interval.Nanoseconds() / 1e6)
	circle := int(delayMillisecond / intervalMillisecond / tw.slotNum)
	pos := int(tw.currentPos+delayMillisecond/intervalMillisecond) % tw.slotNum
	return pos, circle
}

// 从链表中删除任务
func (tw *TimeWheel) removeTask(key interface{}) {
	tw.Lock()
	defer tw.Unlock()
	// 获取定时器所在的槽
	position, ok := tw.timer[key]
	if !ok {
		return
	}
	// 获取槽指向的链表
	l := tw.slots[position]
	for e := l.Front(); e != nil; {
		task := e.Value.(*Task)
		if task.key == key {
			delete(tw.timer, task.key)
			l.Remove(e)
		}

		e = e.Next()
	}
}

var stopTimeWheel = make(chan bool, 1)

func init() {
	SetTimeWheel(NewTimeWheel(10*time.Millisecond, 100))
	async.Go(timeWheel.Start, stopTimeWheel)
}

func StopTimeWheel() {
	stopTimeWheel <- true
}
