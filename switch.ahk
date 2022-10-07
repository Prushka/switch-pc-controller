#SingleInstance, force
#NoEnv
#MaxHotkeysPerInterval 99000000
#HotkeyInterval 99000000
#KeyHistory 0
#include MouseDelta.ahk

ListLines Off
Process, Priority, , A
SetBatchLines, -1
SetKeyDelay, -1, -1
SetMouseDelay, -1
SetDefaultMouseSpeed, 0
SendMode Input
CoordMode, Mouse, Screen


global output
global output_1 := ""
global output_2 := ""
global output_3 := ""
global debugRows := 3
global debug_launched := false

global mouseX := 0
global mouseY := 0
global prevXDiff := 0
global prevYDiff := 0
global paused := false
global lastMouseCheck := A_TickCount


Sysget, totalWidth, 78
Sysget, totalHeight, 79
global centerX := Round(totalWidth//2)
global centerY := Round(totalHeight//2)

DllCall("Winmm.dll\timeBeginPeriod", UInt, 1)
; launchDebug()

md := new MouseDelta("MouseEvent")
md.SetState(true)

SetTimer,CheckMouseEvent,8
SystemCursor("Off")
ClipCursor()

UnclipCursor() {
  print("unclipping")
  DllCall( "ClipCursor" )
}

ClipCursor() {
  print("clipping")
  VarSetCapacity(R,16,0),  NumPut(300,&R+0),NumPut(500,&R+4),NumPut(1000,&R+8),NumPut(1000,&R+12)
  DllCall( "ClipCursor", UInt,&R )
}

CheckMouseEvent() {
    if (paused) {
        return
    }
    if (A_TickCount - lastMouseCheck > 50 and (prevXDiff != 0 or prevYDiff != 0)) {
        sendMouseDiff(0, 0)
        prevXDiff := 0
        prevYDiff := 0
        return
    }
}

MouseEvent(MouseID, x := 0, y := 0){
	t := A_TickCount
  if t - lastMouseCheck < 10 {
    return
  }
	text := "x: " x ", y: " y (lastMouseCheck ? (", Delta Time: " t - lastMouseCheck " ms, MouseID: " MouseID) : "")
	print(text)
	lastMouseCheck := t
  prevXDiff := x
  prevYDiff := y
  sendMouseDiff(x, y)
}

CheckMouse() {
  if paused{
    return
  }
  MouseGetPos,newX,newY
  if (newX != mouseX or newY != mouseY or prevXDiff != 0 or prevYDiff != 0) {
    xDiff := newX - mouseX
    yDiff := newY - mouseY
    if (xDiff == 0 and yDiff == 0) {
      MouseMove, centerX, centerY
      print(centerX "," centerY)
      mouseX := centerX
      mouseY := centerY
    }else{
      mouseX := newX
      mouseY := newY
    }
    prevXDiff := xDiff
    prevYDiff := yDiff
    sendMouseDiff(xDiff, yDiff)
  }
}

AccurateSleep(ms) {
    DllCall("Sleep", UInt, ms)
}

launchDebug(){
	Gui, 99: -Caption +AlwaysOnTop
	Gui, 99:Margin, 0, 0
	Gui, 99:Add, Edit, w500 r%debugRows% voutput
	x := A_ScreenWidth - 500
	Gui, 99:Show, x%x% y0 NoActivate
  debug_launched := true
}

SystemCursor(OnOff=1)   ; INIT = "I","Init"; OFF = 0,"Off"; TOGGLE = -1,"T","Toggle"; ON = others
{
    static AndMask, XorMask, $, h_cursor
        ,c0,c1,c2,c3,c4,c5,c6,c7,c8,c9,c10,c11,c12,c13 ; system cursors
        , b1,b2,b3,b4,b5,b6,b7,b8,b9,b10,b11,b12,b13   ; blank cursors
        , h1,h2,h3,h4,h5,h6,h7,h8,h9,h10,h11,h12,h13   ; handles of default cursors
    if (OnOff = "Init" or OnOff = "I" or $ = "")       ; init when requested or at first call
    {
        $ = h                                          ; active default cursors
        VarSetCapacity( h_cursor,4444, 1 )
        VarSetCapacity( AndMask, 32*4, 0xFF )
        VarSetCapacity( XorMask, 32*4, 0 )
        system_cursors = 32512,32513,32514,32515,32516,32642,32643,32644,32645,32646,32648,32649,32650
        StringSplit c, system_cursors, `,
        Loop %c0%
        {
            h_cursor   := DllCall( "LoadCursor", "uint",0, "uint",c%A_Index% )
            h%A_Index% := DllCall( "CopyImage",  "uint",h_cursor, "uint",2, "int",0, "int",0, "uint",0 )
            b%A_Index% := DllCall("CreateCursor","uint",0, "int",0, "int",0
                , "int",32, "int",32, "uint",&AndMask, "uint",&XorMask )
        }
    }
    if (OnOff = 0 or OnOff = "Off" or $ = "h" and (OnOff < 0 or OnOff = "Toggle" or OnOff = "T"))
        $ = b  ; use blank cursors
    else
        $ = h  ; use the saved cursors

    Loop %c0%
    {
        h_cursor := DllCall( "CopyImage", "uint",%$%%A_Index%, "uint",2, "int",0, "int",0, "uint",0 )
        DllCall( "SetSystemCursor", "uint",h_cursor, "uint",c%A_Index% )
    }
}

print(input){
  printNoLineBreak(input "`n")
}

printNoLineBreak(input){
  if !debug_launched {
    return
  }
    output_3 := output_2
    output_2 := output_1
    output_1 := input
	GuiControlGet, output, 99:
	GuiControl, 99: , output, %output_1%%output_2%%output_3%
}

$F9::
Suspend
paused := !paused
if paused{
UnclipCursor()
SystemCursor("On")
md.SetState(false)
}else{
ClipCursor()
SystemCursor("Off")
md.SetState(true)
}
return

$F8::
Gui, 99:Destroy
md.Delete()
md := ""
releaseAll()
SystemCursor("On")
Reload
return

$Home::
holdButton("Home")
return

$Home Up::
releaseButton("Home")
return

~F7::
releaseAll()
; Using 'true' above and the call below allows the script to remain responsive.
; you can get the response data either in raw or text format
; raw: hObject.responseBody
; text: hObject.responseText
return

$1::
while(GetKeyState("1", "P")){
    holdButton("ZL")
    AccurateSleep(17)
    releaseButton("ZL")
    AccurateSleep(17)
}
return

cancel() {
releaseButton("ZL")
AccurateSleep(17)
holdButton("LR")
AccurateSleep(17)
holdButton("ZL")
AccurateSleep(170)
releaseButton("LR")
AccurateSleep(17)
}

$Q::
cancel()
AccurateSleep(17)
holdButton("LDown")
AccurateSleep(300)
releaseButton("LDown")
AccurateSleep(17)
holdButton("LUp")
AccurateSleep(1)
holdButton("B")
AccurateSleep(34)
releaseButton("B")
return

$E::
cancel()
return

$Z::
holdButton("LL")
return

$Z Up::
releaseButton("LL")
return

$C::
holdButton("LR")
return

$C Up::
releaseButton("LR")
return

$ESC::
holdButton("B")
return

$ESC Up::
releaseButton("B")
return

$Enter::
holdButton("A")
return

$Enter Up::
releaseButton("A")
return

$x::
holdButton("X")
return

$x Up::
releaseButton("X")
return

$w::
holdButton("LUp")
return

$w Up::
releaseButton("LUp")
return

$a::
holdButton("LLeft")
return

$a Up::
releaseButton("LLeft")
return

$s::
holdButton("LDown")
return

$s Up::
releaseButton("LDown")
return

$d::
holdButton("LRight")
return

$d Up::
releaseButton("LRight")
return

$r::
holdButton("RClick")
return

$r Up::
releaseButton("RClick")
return

$-::
holdButton("Minus")
return

$- Up::
releaseButton("Minus")
return

$=::
holdButton("Plus")
return

$= Up::
releaseButton("Plus")
return

$LButton::
holdButton("ZR")
return

$LButton Up::
releaseButton("ZR")
return

$MButton::
holdButton("Y")
return

$MButton Up::
releaseButton("Y")
return

$RButton::
holdButton("ZL")
return

$RButton Up::
releaseButton("ZL")
return

$Shift::
holdButton("ZL")
return

$Shift Up::
releaseButton("ZL")
return

$Right::
holdButton("R")
return

$Right Up::
releaseButton("R")
return

$Left::
holdButton("L")
return

$Left Up::
releaseButton("L")
return

$Up::
holdButton("U")
return

$Up Up::
releaseButton("U")
return

$Down::
holdButton("D")
return

$Down Up::
releaseButton("D")
return

$Space::
holdButton("B")
return

$Space Up::
releaseButton("B")
return

$xbutton1::
holdButton("LR")
return

$xbutton1 Up::
releaseButton("LR")
return

$xbutton2::
holdButton("A")
return

$xbutton2 Up::
releaseButton("A")
return

sendMouseDiff(x, y) {
  sendCommand("D", x "," y)
}

holdButton(btn){
  sendCommand("H", btn)
}

releaseButton(btn) {
    sendCommand("R", btn)
}

releaseAll() {
    sendCommand("A", "s")
}

sendCommand(action, key) {
    whr := ComObjCreate("WinHttp.WinHttpRequest.5.1")
    whr.Open("GET", "http://localhost/" action "/" key, false)
    whr.Send()
}