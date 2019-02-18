package cmd

import (
	"fmt"
	"io"
)

// Bar is the Progress Bar for download
type Bar struct {
	Total   int64
	Current int64
	update  chan int
	done    chan int
	writer  io.Writer
}

const (
	zero  = `[                    ]`
	one   = `[===>                ]`
	two   = `[=====>              ]`
	three = `[=======>            ]`
	four  = `[=========>          ]`
	five  = `[===========>        ]`
	six   = `[=============>      ]`
	seven = `[===============>    ]`
	eight = `[=================>  ]`
	nine  = `[====================]`
)

func NewProgressBar(target io.Writer) *Bar {
	return &Bar{
		Current: 0,
		update:  make(chan int, 10),
		done:    make(chan int),
		writer:  target,
	}
}

func (b *Bar) Init(total int64) {
	b.Total = total
}

func (b *Bar) Write(p []byte) (n int, err error) {
	b.Current += int64(len(p))
	progress := b.draw()
	return b.writer.Write([]byte(progress))
}

func (b *Bar) Done() {
	b.writer.Write([]byte("\n"))
}

func (b *Bar) draw() string {
	percent := (float32(b.Current) / float32(b.Total)) * float32(100)

	if percent > float32(100) {
		percent = float32(100)
	}
	//fmt.Println("Precent: ", percent)
	var progress string

	if percent < 10 {
		progress = getProgress(zero, percent)
	}

	if percent >= float32(10) && percent < float32(20) {
		progress = getProgress(one, percent)
	}

	if percent >= float32(20) && percent < float32(30) {
		progress = getProgress(two, percent)
	}

	if percent >= float32(30) && percent < float32(40) {
		progress = getProgress(three, percent)
	}

	if percent >= float32(40) && percent < float32(50) {
		progress = getProgress(four, percent)
	}

	if percent >= float32(50) && percent < float32(60) {
		progress = getProgress(five, percent)
	}

	if percent >= float32(60) && percent < float32(70) {
		progress = getProgress(six, percent)
	}

	if percent >= float32(80) && percent < float32(90) {
		progress = getProgress(seven, percent)
	}

	if percent >= float32(90) && percent < float32(100) {
		progress = getProgress(eight, percent)
	}

	if percent >= float32(100) {
		progress = getProgress(nine, percent)
	}

	return progress

}

func getProgress(progress string, percent float32) string {
	return fmt.Sprintf("\r %s %.2f%%", progress, percent)
}
