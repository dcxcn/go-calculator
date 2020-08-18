// SPDX-License-Identifier: Unlicense OR MIT

package main

// Multiple windows in Gio.

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/event"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"go-calculator/controls"
	"go-calculator/engine"
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
)

type window struct {
	*app.Window
	calEditor widget.Editor
	calLogger widget.Editor
	more      widget.Clickable
	btn_1     widget.Clickable
	btn_2     widget.Clickable
	btn_3     widget.Clickable
	btn_4     widget.Clickable
	btn_5     widget.Clickable
	btn_6     widget.Clickable
	btn_7     widget.Clickable
	btn_8     widget.Clickable
	btn_9     widget.Clickable
	btn_0     widget.Clickable

	btn_00       widget.Clickable
	btn_dot      widget.Clickable
	btn_plus     widget.Clickable
	btn_minus    widget.Clickable
	btn_multiply widget.Clickable
	btn_divide   widget.Clickable

	btn_xkh_l widget.Clickable
	btn_xkh_r widget.Clickable
	btn_mod   widget.Clickable
	btn_sqrt  widget.Clickable
	btn_cbrt  widget.Clickable
	btn_cfang widget.Clickable

	btn_sin widget.Clickable
	btn_cos widget.Clickable
	btn_tan widget.Clickable
	btn_cot widget.Clickable

	btn_abs   widget.Clickable
	btn_ceil  widget.Clickable
	btn_floor widget.Clickable
	btn_round widget.Clickable

	btn_backspace widget.Clickable
	btn_clean     widget.Clickable
	btn_equal     widget.Clickable
	btn_savelog   widget.Clickable
	btn_readlog   widget.Clickable
	btn_about     widget.Clickable
	close         widget.Clickable
}
type (
	D = layout.Dimensions
	C = layout.Context
)

func main() {
	newWindow()
	app.Main()
}

var (
	windowCount int32
	listV       = &layout.List{
		Axis: layout.Vertical,
	}
	listH = &layout.List{
		Axis: layout.Horizontal,
	}
)

func newWindow() {
	atomic.AddInt32(&windowCount, +1)
	go func() {
		w := new(window)
		w.Window = app.NewWindow(
			app.Size(unit.Dp(500), unit.Dp(500)),
			app.MinSize(unit.Dp(500), unit.Dp(500)),
			app.MaxSize(unit.Dp(500), unit.Dp(500)),
			app.Title("Go-Calculator"),
		)
		if err := w.loop(w.Events()); err != nil {
			log.Fatal(err)
		}
		if c := atomic.AddInt32(&windowCount, -1); c == 0 {
			os.Exit(0)
		}
	}()
}

func (w *window) loop(events <-chan event.Event) error {
	th := material.NewTheme(gofont.Collection())
	var ops op.Ops
	for {
		e := <-events
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			for w.more.Clicked() {
				newWindow()
			}
			for w.close.Clicked() {
				w.Close()
			}
			for w.btn_readlog.Clicked() {
				logStr, _ := GetFileContentAsStringLines("./calculator.log")
				w.calLogger.SetText("")
				for _, str := range logStr {
					w.calLogger.Insert(str)
				}
			}
			for w.btn_savelog.Clicked() {
				out, _ := os.Create("./calculator.log")
				out.WriteString(w.calLogger.Text())
				defer out.Close()
			}
			for w.btn_about.Clicked() {
				controls.ShowMessageBox("About", "|developer: duanchunxiao\n|__country: china\n|_____city: beijing changping\n|____email: 837578856@qq.com")
			}
			for w.btn_backspace.Clicked() {
				txt := w.calEditor.Text()
				w.calEditor.SetText("")
				w.calEditor.Insert(txt[:len(txt)-1])
			}
			for w.btn_clean.Clicked() {
				w.calEditor.SetText("")
			}
			for w.btn_equal.Clicked() {
				txt := w.calEditor.Text()
				res, val := doCompute(txt)
				w.calEditor.SetText("")
				w.calEditor.Insert(val)
				w.calLogger.Insert(res)
			}
			for w.btn_1.Clicked() {
				w.calEditor.Insert("1")
			}
			for w.btn_2.Clicked() {
				w.calEditor.Insert("2")
			}
			for w.btn_3.Clicked() {
				w.calEditor.Insert("3")
			}
			for w.btn_4.Clicked() {
				w.calEditor.Insert("4")
			}
			for w.btn_5.Clicked() {
				w.calEditor.Insert("5")
			}
			for w.btn_6.Clicked() {
				w.calEditor.Insert("6")
			}
			for w.btn_7.Clicked() {
				w.calEditor.Insert("7")
			}
			for w.btn_8.Clicked() {
				w.calEditor.Insert("8")
			}
			for w.btn_9.Clicked() {
				w.calEditor.Insert("9")
			}
			for w.btn_0.Clicked() {
				w.calEditor.Insert("0")
			}
			for w.btn_00.Clicked() {
				w.calEditor.Insert("00")
			}
			for w.btn_dot.Clicked() {
				w.calEditor.Insert(".")
			}
			for w.btn_plus.Clicked() {
				w.calEditor.Insert("+")
			}
			for w.btn_minus.Clicked() {
				w.calEditor.Insert("-")
			}
			for w.btn_multiply.Clicked() {
				w.calEditor.Insert("*")
			}
			for w.btn_divide.Clicked() {
				w.calEditor.Insert("/")
			}
			for w.btn_xkh_l.Clicked() {
				w.calEditor.Insert("(")
			}
			for w.btn_xkh_r.Clicked() {
				w.calEditor.Insert(")")
			}
			for w.btn_mod.Clicked() {
				w.calEditor.Insert("%")
			}
			for w.btn_sqrt.Clicked() {
				txt := w.calEditor.Text()
				if IsNum(txt) {
					w.calEditor.SetText("")
					w.calEditor.Insert("sqrt(" + txt + ")")
				} else {
					w.calEditor.Insert("sqrt()")
				}
			}
			for w.btn_cbrt.Clicked() {
				txt := w.calEditor.Text()
				if IsNum(txt) {
					w.calEditor.SetText("")
					w.calEditor.Insert("cbrt(" + txt + ")")
				} else {
					w.calEditor.Insert("cbrt()")
				}
			}
			for w.btn_cfang.Clicked() {
				w.calEditor.Insert("^")
			}
			for w.btn_sin.Clicked() {
				txt := w.calEditor.Text()
				if IsNum(txt) {
					w.calEditor.SetText("")
					w.calEditor.Insert("sin(" + txt + ")")
				} else {
					w.calEditor.Insert("sin()")
				}
			}
			for w.btn_cos.Clicked() {
				txt := w.calEditor.Text()
				if IsNum(txt) {
					w.calEditor.SetText("")
					w.calEditor.Insert("cos(" + txt + ")")
				} else {
					w.calEditor.Insert("cos()")
				}
			}
			for w.btn_tan.Clicked() {
				txt := w.calEditor.Text()
				if IsNum(txt) {
					w.calEditor.SetText("")
					w.calEditor.Insert("tan(" + txt + ")")
				} else {
					w.calEditor.Insert("tan()")
				}
			}
			for w.btn_cot.Clicked() {
				txt := w.calEditor.Text()
				if IsNum(txt) {
					w.calEditor.SetText("")
					w.calEditor.Insert("cot(" + txt + ")")
				} else {
					w.calEditor.Insert("cot()")
				}
			}
			for w.btn_abs.Clicked() {
				txt := w.calEditor.Text()
				if IsNum(txt) {
					w.calEditor.SetText("")
					w.calEditor.Insert("abs(" + txt + ")")
				} else {
					w.calEditor.Insert("abs()")
				}
			}
			for w.btn_ceil.Clicked() {
				txt := w.calEditor.Text()
				if IsNum(txt) {
					w.calEditor.SetText("")
					w.calEditor.Insert("ceil(" + txt + ")")
				} else {
					w.calEditor.Insert("ceil()")
				}
			}
			for w.btn_floor.Clicked() {
				txt := w.calEditor.Text()
				if IsNum(txt) {
					w.calEditor.SetText("")
					w.calEditor.Insert("floor(" + txt + ")")
				} else {
					w.calEditor.Insert("floor()")
				}
			}
			for w.btn_round.Clicked() {
				txt := w.calEditor.Text()
				if IsNum(txt) {
					w.calEditor.SetText("")
					w.calEditor.Insert("round(" + txt + ")")
				} else {
					w.calEditor.Insert("round()")
				}
			}
			gtx := layout.NewContext(&ops, e)
			//搭建界面
			BuildUI(gtx, th, w)
			e.Frame(gtx.Ops)
		}
	}
}
func IsNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func GetFileContentAsStringLines(filePath string) ([]string, error) {
	fmt.Printf("get file content as lines: %v", filePath)
	result := []string{}
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("read file: %v error: %v", filePath, err)
		return result, err
	}
	s := string(b)
	for _, lineStr := range strings.Split(s, "\n") {
		lineStr = strings.TrimSpace(lineStr)
		if lineStr == "" {
			continue
		}
		result = append(result, lineStr)
	}
	fmt.Printf("get file content as lines: %v, size: %v", filePath, len(result))
	return result, nil
}
func doCompute(exp string) (res string, val string) {
	// input text -> []token
	toks, err := engine.Parse(exp)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		return
	}
	// []token -> AST Tree
	ast := engine.NewAST(toks, exp)
	if ast.Err != nil {
		fmt.Println("ERROR: " + ast.Err.Error())
		return
	}
	// AST builder
	ar := ast.ParseExpression()
	if ast.Err != nil {
		fmt.Println("ERROR: " + ast.Err.Error())
		return
	}
	fmt.Printf("ExprAST: %+v\n", ar)
	// catch runtime errors
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("ERROR: ", e)
		}
	}()
	// AST traversal -> result
	r := engine.ExprASTResult(ar)
	fmt.Println("progressing ...\t", r)
	fmt.Printf("%s = %v\n", exp, r)
	return fmt.Sprintf("%s = %v\n", exp, r), fmt.Sprintf("%v", r)
}
func BuildUI(gtx layout.Context, th *material.Theme, w *window) layout.Dimensions {

	widgets := []layout.Widget{
		func(gtx C) D {
			gtx.Constraints.Min.Y = gtx.Px(unit.Dp(100))
			gtx.Constraints.Max.Y = gtx.Px(unit.Dp(100))
			ce := material.Editor(th, &w.calEditor, "mathematical expression")
			border := widget.Border{Color: color.RGBA{A: 0xff}, CornerRadius: unit.Dp(8), Width: unit.Px(2)}
			return border.Layout(gtx, func(gtx C) D {
				return layout.UniformInset(unit.Dp(8)).Layout(gtx, ce.Layout)
			})

		},
		func(gtx C) D {
			gtx.Constraints.Min.Y = gtx.Px(unit.Dp(100))
			gtx.Constraints.Max.Y = gtx.Px(unit.Dp(100))
			ce := material.Editor(th, &w.calLogger, "log")
			border := widget.Border{Color: color.RGBA{R: 0xee, G: 0xee, B: 0xee, A: 0xff}, CornerRadius: unit.Dp(8), Width: unit.Px(2)}
			return border.Layout(gtx, func(gtx C) D {
				return layout.UniformInset(unit.Dp(8)).Layout(gtx, ce.Layout)
			})
		},
		func(gtx C) D {
			in := layout.UniformInset(unit.Dp(2))
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx C) D {
					return in.Layout(gtx, material.Button(th, &w.btn_savelog, "SaveLog").Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					return in.Layout(gtx, material.Button(th, &w.btn_readlog, "ReadLog").Layout)
				}),

				layout.Flexed(1, func(gtx C) D {
					return in.Layout(gtx, material.Button(th, &w.btn_about, "About").Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					btn := material.Button(th, &w.btn_xkh_l, "(")
					btn.Background = color.RGBA{A: 0xff, R: 0x6e, G: 0x6e, B: 0x6e}
					return in.Layout(gtx, btn.Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					btn := material.Button(th, &w.btn_xkh_r, ")")
					btn.Background = color.RGBA{A: 0xff, R: 0x6e, G: 0x6e, B: 0x6e}
					return in.Layout(gtx, btn.Layout)
				}),
			)
		},
		func(gtx C) D {
			in := layout.UniformInset(unit.Dp(2))
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx C) D {
					btn := material.Button(th, &w.btn_7, "7")
					btn.Background = color.RGBA{A: 0xff, R: 0x00, G: 0x00, B: 0x00}
					return in.Layout(gtx, btn.Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					btn := material.Button(th, &w.btn_8, "8")
					btn.Background = color.RGBA{A: 0xff, R: 0x00, G: 0x00, B: 0x00}
					return in.Layout(gtx, btn.Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					btn := material.Button(th, &w.btn_9, "9")
					btn.Background = color.RGBA{A: 0xff, R: 0x00, G: 0x00, B: 0x00}
					return in.Layout(gtx, btn.Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					btn := material.Button(th, &w.btn_divide, "÷")
					btn.Background = color.RGBA{A: 0xff, R: 0x6e, G: 0x6e, B: 0x6e}
					return in.Layout(gtx, btn.Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					btn := material.Button(th, &w.btn_sqrt, "sqrt")
					btn.Background = color.RGBA{A: 0xff, R: 0x6e, G: 0x6e, B: 0x6e}
					return in.Layout(gtx, btn.Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					btn := material.Button(th, &w.btn_sin, "sin")
					btn.Background = color.RGBA{A: 0xff, R: 0x6e, G: 0x6e, B: 0x6e}
					return in.Layout(gtx, btn.Layout)

				}),
				layout.Flexed(1, func(gtx C) D {
					btn := material.Button(th, &w.btn_abs, "abs")
					btn.Background = color.RGBA{A: 0xff, R: 0x6e, G: 0x6e, B: 0x6e}
					return in.Layout(gtx, btn.Layout)
				}),
			)
		},
		func(gtx C) D {
			in := layout.UniformInset(unit.Dp(2))
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx C) D {
					btn := material.Button(th, &w.btn_4, "4")
					btn.Background = color.RGBA{A: 0xff, R: 0x00, G: 0x00, B: 0x00}
					return in.Layout(gtx, btn.Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					btn := material.Button(th, &w.btn_5, "5")
					btn.Background = color.RGBA{A: 0xff, R: 0x00, G: 0x00, B: 0x00}
					return in.Layout(gtx, btn.Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					btn := material.Button(th, &w.btn_6, "6")
					btn.Background = color.RGBA{A: 0xff, R: 0x00, G: 0x00, B: 0x00}
					return in.Layout(gtx, btn.Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					btn := material.Button(th, &w.btn_multiply, "×")
					btn.Background = color.RGBA{A: 0xff, R: 0x6e, G: 0x6e, B: 0x6e}
					return in.Layout(gtx, btn.Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					btn := material.Button(th, &w.btn_cbrt, "cbrt")
					btn.Background = color.RGBA{A: 0xff, R: 0x6e, G: 0x6e, B: 0x6e}
					return in.Layout(gtx, btn.Layout)
				}),

				layout.Flexed(1, func(gtx C) D {
					btn := material.Button(th, &w.btn_cos, "cos")
					btn.Background = color.RGBA{A: 0xff, R: 0x6e, G: 0x6e, B: 0x6e}
					return in.Layout(gtx, btn.Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					btn := material.Button(th, &w.btn_ceil, "ceil")
					btn.Background = color.RGBA{A: 0xff, R: 0x6e, G: 0x6e, B: 0x6e}
					return in.Layout(gtx, btn.Layout)
				}),
			)
		},
		func(gtx C) D {
			in := layout.UniformInset(unit.Dp(2))
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx C) D {
					btn := material.Button(th, &w.btn_1, "1")
					btn.Background = color.RGBA{A: 0xff, R: 0x00, G: 0x00, B: 0x00}
					return in.Layout(gtx, btn.Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					btn := material.Button(th, &w.btn_2, "2")
					btn.Background = color.RGBA{A: 0xff, R: 0x00, G: 0x00, B: 0x00}
					return in.Layout(gtx, btn.Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					btn := material.Button(th, &w.btn_3, "3")
					btn.Background = color.RGBA{A: 0xff, R: 0x00, G: 0x00, B: 0x00}
					return in.Layout(gtx, btn.Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					btn := material.Button(th, &w.btn_minus, "-")
					btn.Background = color.RGBA{A: 0xff, R: 0x6e, G: 0x6e, B: 0x6e}
					return in.Layout(gtx, btn.Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					btn := material.Button(th, &w.btn_cfang, "^")
					btn.Background = color.RGBA{A: 0xff, R: 0x6e, G: 0x6e, B: 0x6e}
					return in.Layout(gtx, btn.Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					btn := material.Button(th, &w.btn_tan, "tan")
					btn.Background = color.RGBA{A: 0xff, R: 0x6e, G: 0x6e, B: 0x6e}
					return in.Layout(gtx, btn.Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					btn := material.Button(th, &w.btn_floor, "floor")
					btn.Background = color.RGBA{A: 0xff, R: 0x6e, G: 0x6e, B: 0x6e}
					return in.Layout(gtx, btn.Layout)
				}),
			)
		},
		func(gtx C) D {
			in := layout.UniformInset(unit.Dp(2))
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx C) D {
					btn := material.Button(th, &w.btn_0, "0")
					btn.Background = color.RGBA{A: 0xff, R: 0x00, G: 0x00, B: 0x00}
					return in.Layout(gtx, btn.Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					btn := material.Button(th, &w.btn_00, "00")
					btn.Background = color.RGBA{A: 0xff, R: 0x00, G: 0x00, B: 0x00}
					return in.Layout(gtx, btn.Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					btn := material.Button(th, &w.btn_dot, ".")
					btn.Background = color.RGBA{A: 0xff, R: 0x00, G: 0x00, B: 0x00}
					return in.Layout(gtx, btn.Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					btn := material.Button(th, &w.btn_plus, "+")
					btn.Background = color.RGBA{A: 0xff, R: 0x6e, G: 0x6e, B: 0x6e}
					return in.Layout(gtx, btn.Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					btn := material.Button(th, &w.btn_mod, "%")
					btn.Background = color.RGBA{A: 0xff, R: 0x6e, G: 0x6e, B: 0x6e}
					return in.Layout(gtx, btn.Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					btn := material.Button(th, &w.btn_cot, "cot")
					btn.Background = color.RGBA{A: 0xff, R: 0x6e, G: 0x6e, B: 0x6e}
					return in.Layout(gtx, btn.Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					btn := material.Button(th, &w.btn_round, "round")
					btn.Background = color.RGBA{A: 0xff, R: 0x6e, G: 0x6e, B: 0x6e}
					return in.Layout(gtx, btn.Layout)
				}),
			)
		},
		func(gtx C) D {
			in := layout.UniformInset(unit.Dp(2))
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx C) D {
					btn := material.Button(th, &w.btn_backspace, "<-")
					btn.Background = color.RGBA{A: 0xff, R: 0xda, G: 0x26, B: 0x26}
					return in.Layout(gtx, btn.Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					btn := material.Button(th, &w.btn_clean, "clean")
					btn.Background = color.RGBA{A: 0xff, R: 0xda, G: 0x26, B: 0x26}
					return in.Layout(gtx, btn.Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					btn := material.Button(th, &w.btn_equal, "=")
					btn.Background = color.RGBA{A: 0xff, R: 0x07, G: 0x74, B: 0x70}
					return in.Layout(gtx, btn.Layout)
				}),
			)
		},
	}
	return listV.Layout(gtx, len(widgets), func(gtx C, i int) D {
		return layout.UniformInset(unit.Dp(3)).Layout(gtx, widgets[i])
	})
	//layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
	//	return layout.Flex{
	//		Alignment: layout.Middle,
	//	}.Layout(gtx,
	//
	//		//RigidInset(material.Editor(th, &w.calEditor, "0").Layout),
	//		//RigidInset(material.Button(th, &w.more, "More!").Layout),
	//		//RigidInset(material.Button(th, &w.close, "Exit").Layout),
	//		//RigidInset(material.Button(th, &w.btn_1, "1").Layout),
	//		//RigidInset(material.Button(th, &w.btn_2, "2").Layout),
	//		//RigidInset(material.Button(th, &w.btn_3, "3").Layout),
	//		//RigidInset(material.Button(th, &w.btn_4, "4").Layout),
	//		//RigidInset(material.Button(th, &w.btn_5, "5").Layout),
	//		//RigidInset(material.Button(th, &w.btn_6, "6").Layout),
	//		//RigidInset(material.Button(th, &w.btn_7, "7").Layout),
	//		//RigidInset(material.Button(th, &w.btn_8, "8").Layout),
	//		//RigidInset(material.Button(th, &w.btn_9, "9").Layout),
	//		//RigidInset(material.Button(th, &w.btn_0, "0").Layout),
	//	)
	//})
}

func RigidInset(w layout.Widget) layout.FlexChild {
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Sp(5)).Layout(gtx, w)
	})
}
