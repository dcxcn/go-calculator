package controls

import (
	"log"
	"sync/atomic"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"gioui.org/font/gofont"
)

type MessageBox struct {
	*app.Window

	msgText  string
	btn_sure widget.Clickable
}
type (
	D = layout.Dimensions
	C = layout.Context
)

var (
	windowCount int32
	listV       = &layout.List{
		Axis: layout.Vertical,
	}
	listH = &layout.List{
		Axis: layout.Horizontal,
	}
)

func ShowMessageBox(title string, msg string) {
	newWindow(title, msg)
}

func newWindow(title string, msg string) {
	atomic.AddInt32(&windowCount, +1)
	go func() {
		w := new(MessageBox)
		w.msgText = msg
		w.Window = app.NewWindow(
			app.Size(unit.Dp(400), unit.Dp(200)),
			app.MinSize(unit.Dp(400), unit.Dp(200)),
			app.MaxSize(unit.Dp(400), unit.Dp(200)),
			app.Title(title),
		)
		if err := w.loop(w.Events()); err != nil {
			log.Fatal(err)
		}
		if c := atomic.AddInt32(&windowCount, -1); c == 0 {
			w.Close()
		}
	}()
}

func (w *MessageBox) loop(events <-chan event.Event) error {
	th := material.NewTheme(gofont.Collection())
	var ops op.Ops
	for {
		e := <-events
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:

			for w.btn_sure.Clicked() {
				w.Close()
			}
			gtx := layout.NewContext(&ops, e)

			//搭建界面
			BuildUI(gtx, th, w)
			e.Frame(gtx.Ops)
		}
	}
}
func BuildUI(gtx layout.Context, th *material.Theme, w *MessageBox) layout.Dimensions {
	widgets := []layout.Widget{
		material.H6(th, w.msgText).Layout,
		func(gtx C) D {
			in := layout.UniformInset(unit.Dp(2))
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx C) D {
					return in.Layout(gtx, material.Button(th, &w.btn_sure, "sure").Layout)
				}),
			)
		},
	}
	return listV.Layout(gtx, len(widgets), func(gtx C, i int) D {
		return layout.UniformInset(unit.Dp(3)).Layout(gtx, widgets[i])
	})
}
