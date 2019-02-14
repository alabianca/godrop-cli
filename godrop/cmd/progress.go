package cmd

import "fmt"

type Bar struct {
	Total   int
	Current int
	update  chan int
	done    chan int
}

const (
	zero  = `[                    ]`
	one   = `[==                  ]`
	two   = `[====                ]`
	three = `[======              ]`
	four  = `[========            ]`
	five  = `[==========          ]`
	six   = `[============        ]`
	seven = `[==============      ]`
	eight = `[================    ]`
	nine  = `[==================  ]`
	ten   = `[====================]`
)

func NewProgressBar(total int) *Bar {
	return &Bar{
		Total:   total,
		Current: 0,
		update:  make(chan int, 10),
		done:    make(chan int),
	}
}

func (b *Bar) Init() {
	b.draw() // initial draw

	go func() {
		for {
			select {
			case <-b.update:
				//update the bar
				b.draw()
			case <-b.done:
				fmt.Println()
				return
			}
		}
	}()
}

func (b *Bar) Update(val int) {
	b.Current = val
	b.update <- val
}

func (b *Bar) Done() {
	b.done <- 1
}

func (b *Bar) draw() {
	percent := (float32(b.Current) / float32(b.Total)) * float32(100)

	if percent < 10 {
		print(zero, percent)
	}

	if percent >= 10 && percent < 20 {
		print(one, percent)
	}

	if percent >= 20 && percent < 30 {
		print(two, percent)
	}

	if percent >= 30 && percent < 40 {
		print(three, percent)
	}

	if percent >= 40 && percent < 50 {
		print(four, percent)
	}

	if percent >= 50 && percent < 60 {
		print(five, percent)
	}

	if percent >= 60 && percent < 70 {
		print(six, percent)
	}

	if percent >= 80 && percent < 90 {
		print(seven, percent)
	}

	if percent >= 90 && percent < 100 {
		print(eight, percent)
	}

	if percent >= 100 {
		print(nine, percent)
	}

}

func print(progress string, percent float32) {
	fmt.Printf("\r %s %.2f%%", progress, percent)
}
