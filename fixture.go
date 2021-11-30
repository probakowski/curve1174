package curve1174

import (
	"github.com/dterei/gotsc"
	"math"
	"math/rand"
	"sort"
	"time"
)

const testsNumber = 1 + percentilesNumber + 1
const enoughMeasurements = 10000
const percentilesNumber = 100

type Timer struct {
	startCount uint64
	endCount   uint64
}

func (t *Timer) Start() {
	t.startCount = gotsc.BenchStart()
}

func (t *Timer) End() {
	t.endCount = gotsc.BenchEnd()
}

func (t *Timer) count() uint64 {
	return t.endCount - t.startCount
}

func preparePercentiles(ticks []uint64) (percentiles []uint64) {
	percentiles = make([]uint64, percentilesNumber)
	ticks = append([]uint64{}, ticks...)
	sort.Slice(ticks, func(i, j int) bool {
		return ticks[i] < ticks[j]
	})
	for i := 0; i < percentilesNumber; i++ {
		index := (1 - math.Pow(0.5, 10*(float64(i)+1)/percentilesNumber)) * float64(len(ticks))
		percentiles[i] = ticks[int(index)]
	}
	return percentiles
}

type Ctx struct {
	stats       [101]Stat
	percentiles []uint64
}

func (ctx *Ctx) Measure(testMethod func(timer *Timer, random bool, vector int), randomVectors, iterationCount int) float64 {
	rand.Seed(time.Now().UnixNano())
	//var stats [testsNumber]Stat
	ticks := make([]uint64, iterationCount)
	tests := make([]int, iterationCount)
	classes := make([]int, iterationCount)
	for i := 0; i < iterationCount; i++ {
		tests[i] = rand.Intn(randomVectors)
		classes[i] = rand.Intn(2)
	}
	for i := 0; i < iterationCount; i++ {
		timer := new(Timer)
		testMethod(timer, classes[i] == 1, tests[i])
		ticks[i] = timer.count()
	}
	if ctx.percentiles == nil {
		ctx.percentiles = preparePercentiles(ticks)
	}
	for i := 0; i < iterationCount; i++ {
		ctx.stats[0].push(ticks[i], classes[i])

		for crop := 0; crop < percentilesNumber; crop++ {
			if ticks[i] < ctx.percentiles[crop] {
				ctx.stats[crop+1].push(ticks[i], classes[i])
			}
		}
	}

	max := 0.0
	for i := 0; i < len(ctx.stats); i++ {
		if ctx.stats[i].count[0]+ctx.stats[i].count[1] < 10000 {
			continue
		}
		tt := math.Abs(ctx.stats[i].compute())
		if tt > max {
			max = tt
		}
	}

	return max
}

type Stat struct {
	mean  [2]float64
	m2    [2]float64
	count [2]float64
}

func (s *Stat) push(ticks uint64, class int) {
	s.count[class]++
	delta := float64(ticks) - s.mean[class]
	s.mean[class] += delta / s.count[class]
	s.m2[class] += delta * (float64(ticks) - s.mean[class])
}

func (s *Stat) compute() float64 {
	v0 := s.m2[0] / (s.count[0] - 1)
	v1 := s.m2[1] / (s.count[1] - 1)
	num := s.mean[0] - s.mean[1]
	den := math.Sqrt(v0/s.count[0] + v1/s.count[1])
	t := num / den
	return t
}
