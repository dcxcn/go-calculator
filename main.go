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
	"go-calculator/engine"
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync/atomic"
)

type window struct {
	*app.Window
	calEditor     widget.Editor
	calLogger     widget.Editor
	more          widget.Clickable
	btn_1         widget.Clickable
	btn_2         widget.Clickable
	btn_3         widget.Clickable
	btn_4         widget.Clickable
	btn_5         widget.Clickable
	btn_6         widget.Clickable
	btn_7         widget.Clickable
	btn_8         widget.Clickable
	btn_9         widget.Clickable
	btn_0         widget.Clickable
	btn_00        widget.Clickable
	btn_dot       widget.Clickable
	btn_plus      widget.Clickable
	btn_minus     widget.Clickable
	btn_multiply  widget.Clickable
	btn_divide    widget.Clickable
	btn_xkh_l     widget.Clickable
	btn_xkh_r     widget.Clickable
	btn_mod       widget.Clickable
	btn_sqrt      widget.Clickable
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
			app.Size(unit.Dp(400), unit.Dp(500)),
			app.MinSize(unit.Dp(400), unit.Dp(500)),
			app.MaxSize(unit.Dp(400), unit.Dp(500)),
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
				w.calEditor.SetText("")
				w.calEditor.Insert("sqrt(" + txt + ")")
			}
			gtx := layout.NewContext(&ops, e)
			//搭建界面
			BuildUI(gtx, th, w)
			e.Frame(gtx.Ops)
		}
	}
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
			ce := material.Editor(th, &w.calEditor, "Hint")
			border := widget.Border{Color: color.RGBA{A: 0xff}, CornerRadius: unit.Dp(8), Width: unit.Px(2)}
			return border.Layout(gtx, func(gtx C) D {
				return layout.UniformInset(unit.Dp(8)).Layout(gtx, ce.Layout)
			})

		},
		func(gtx C) D {
			gtx.Constraints.Min.Y = gtx.Px(unit.Dp(100))
			gtx.Constraints.Max.Y = gtx.Px(unit.Dp(100))
			ce := material.Editor(th, &w.calLogger, "Hint")
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
			)
		},
		func(gtx C) D {
			in := layout.UniformInset(unit.Dp(2))
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx C) D {
					return in.Layout(gtx, material.Button(th, &w.btn_7, "7").Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					return in.Layout(gtx, material.Button(th, &w.btn_8, "8").Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					return in.Layout(gtx, material.Button(th, &w.btn_9, "9").Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					return in.Layout(gtx, material.Button(th, &w.btn_divide, "÷").Layout)
				}),

				layout.Flexed(1, func(gtx C) D {
					return in.Layout(gtx, material.Button(th, &w.btn_sqrt, "sqrt").Layout)
				}),
			)
		},
		func(gtx C) D {
			in := layout.UniformInset(unit.Dp(2))
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx C) D {
					return in.Layout(gtx, material.Button(th, &w.btn_4, "4").Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					return in.Layout(gtx, material.Button(th, &w.btn_5, "5").Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					return in.Layout(gtx, material.Button(th, &w.btn_6, "6").Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					return in.Layout(gtx, material.Button(th, &w.btn_multiply, "×").Layout)
				}),

				layout.Flexed(1, func(gtx C) D {
					return in.Layout(gtx, material.Button(th, &w.btn_mod, "%").Layout)
				}),
			)
		},
		func(gtx C) D {
			in := layout.UniformInset(unit.Dp(2))
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx C) D {
					return in.Layout(gtx, material.Button(th, &w.btn_1, "1").Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					return in.Layout(gtx, material.Button(th, &w.btn_2, "2").Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					return in.Layout(gtx, material.Button(th, &w.btn_3, "3").Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					return in.Layout(gtx, material.Button(th, &w.btn_minus, "-").Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					return in.Layout(gtx, material.Button(th, &w.btn_xkh_l, "(").Layout)
				}),
			)
		},
		func(gtx C) D {
			in := layout.UniformInset(unit.Dp(2))
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx C) D {
					return in.Layout(gtx, material.Button(th, &w.btn_0, "0").Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					return in.Layout(gtx, material.Button(th, &w.btn_00, "00").Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					return in.Layout(gtx, material.Button(th, &w.btn_dot, ".").Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					return in.Layout(gtx, material.Button(th, &w.btn_plus, "+").Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					return in.Layout(gtx, material.Button(th, &w.btn_xkh_r, ")").Layout)
				}),
			)
		},
		func(gtx C) D {
			in := layout.UniformInset(unit.Dp(2))
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx C) D {
					return in.Layout(gtx, material.Button(th, &w.btn_backspace, "<-").Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					return in.Layout(gtx, material.Button(th, &w.btn_clean, "C").Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					return in.Layout(gtx, material.Button(th, &w.btn_equal, "=").Layout)
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
