// Package js includes shortcuts to manipulate everything javascript and DOM related.
package js

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/gowebapi/webapi"
	"github.com/gowebapi/webapi/core/js"
	"github.com/gowebapi/webapi/css/cssom"
	"github.com/gowebapi/webapi/dom"
	"github.com/gowebapi/webapi/dom/domcore"
	"github.com/gowebapi/webapi/html"
	"github.com/gowebapi/webapi/html/canvas"
	"github.com/gowebapi/webapi/html/htmlevent"
	"log"
	"net/url"
	"reflect"
	"strings"
)

var (
	// Doc is a shortcut to the pages Document.
	Doc *webapi.Document
	Win *webapi.Window

	// URI is the parsed URL/URI (a net/url.URL object). See URI.Query() to get
	// access to the GET parameters.
	URI *url.URL
)

// Initialize document information: URI, GetParameters, etc.
func init() {
	if flag.Lookup("test.v") != nil {
		// If used in test, it is not being run in a browser, so no document is available.
		return
	}

	// Get document and parse URI and query parameters.
	Doc = webapi.GetDocument()
	Win = webapi.GetWindow()

	var err error
	if URI, err = url.Parse(Doc.DocumentURI()); err != nil {
		log.Fatalf("Failed to parse URI (%s): %v", Doc.DocumentURI(), err)
	}
	glog.V(1).Infof("URI: %v", URI)
	glog.V(1).Infof("\tHostname=%s", URI.Hostname())
	glog.V(1).Infof("\tQuery=%v", URI.Query())
}

// Compatible interface takes anything that has JSValue().
type Compatible interface {
	JSValue() js.Value
}

// EventTargetCompatible is the interface that all Element/Node types of the DOM implement.
type EventTargetCompatible interface {
	Compatible
	AddEventListener(_type string, callback *domcore.EventListenerValue, options *domcore.Union)
}

type NodeCompatible interface {
	EventTargetCompatible
	AppendChild(node *dom.Node) (_result *dom.Node)
	ChildNodes() *dom.NodeList
	RemoveChild(child *dom.Node) (_result *dom.Node)
}

// ElementCompatible is anything that behaves like an Element.
type ElementCompatible interface {
	NodeCompatible
	SetAttribute(qualifiedName string, value string)
	RemoveAttribute(qualifiedName string)
}

// HtmlElementCompatible is anything that behaves like an HtmlElement.
type HtmlElementCompatible interface {
	ElementCompatible
	Style() *cssom.CSSStyleDeclaration
}

type EventCompatible interface {
	Compatible
	Bubbles() bool
	PreventDefault()
	StopPropagation()
	StopImmediatePropagation()
}

type MouseEventCompatible interface {
	EventCompatible
	Button() int

	AltKey() bool
	CtrlKey() bool
	ShiftKey() bool

	OffsetX() float64
	OffsetY() float64
}

type KeyboardEventCompatible interface {
	EventCompatible
	Key() string
	Code() string
	KeyCode() uint
	CharCode() uint

	AltKey() bool
	CtrlKey() bool
	ShiftKey() bool
	MetaKey() bool

	Repeat() bool
	IsComposing() bool
}

// Key codes that can be used when comparing with a KeyboardEventCompatible.KeyCode()
const (
	KeyCodeNone  = uint(0)
	KeyCodeLeft  = 37
	KeyCodeUp    = 38
	KeyCodeRight = 39
	KeyCodeDown  = 40
)

// IsNil checks for nil pointers as interface.
func IsNil(e Compatible) bool {
	if e == nil || (reflect.ValueOf(e).Kind() == reflect.Ptr && reflect.ValueOf(e).IsNil()) {
		return true
	}
	return false
}

// AsNode converts any pointer to a struct that extends dom.Node, back to a *dom.Node.
func AsNode(e NodeCompatible) *dom.Node {
	if IsNil(e) {
		return nil
	}
	if n, ok := e.(*dom.Node); ok {
		fmt.Println("Converted easy!")
		return n
	}
	if reflect.ValueOf(e).Kind() != reflect.Ptr {
		return nil
	}
	val := reflect.Indirect(reflect.ValueOf(e))
	val = val.FieldByName("Node") // Get dom.Node Value.
	if !val.IsValid() {
		return nil
	}
	val = val.Addr() // Get the *dom.Node Value.
	return val.Interface().(*dom.Node)
}

// ById returns the element with the given Id.
func ById(id string) *html.HTMLElement {
	e := Doc.GetElementById(id)
	if e == nil {
		return nil
	}
	return html.HTMLElementFromWrapper(e)
}

// MustById returns the element with the given Id, and glog.Fatal if not found.
func MustById(id string) *html.HTMLElement {
	e := ById(id)
	if e == nil {
		glog.Fatalf("MustById didn't find element %q in HTML page.", id)
	}
	return e
}

// ByTag returns the elements with of the given Tag
func ByTag(tag string) *dom.HTMLCollection {
	return Doc.GetElementsByTagName(tag)
}

// ByTag1 returns the first element with the given Tag.
func ByTag1(tag string) *html.HTMLElement {
	l := ByTag(tag)
	if l.Length() == 0 {
		return nil
	}
	return html.HTMLElementFromWrapper(l.Item(0))
}

var _ = ByTag1

// Elem creates a new DOM element with the tag / attributes given. Notice that attributes with space won't work,
// those will need to be set by using e.SetAttribute().
func Elem(tagAttributes string) *dom.Element {
	splits := strings.Split(tagAttributes, " ")
	e := Doc.CreateElement(splits[0], nil)
	for ii := 1; ii < len(splits); ii++ {
		split := splits[ii]
		idx := strings.Index(split, "=")
		if idx == -1 {
			e.SetAttribute(split, "")
		} else {
			e.SetAttribute(split[0:idx], split[idx+1:])
		}
	}
	return e
}

// Append will append the child to the parent node.
func Append(parent, child NodeCompatible) {
	cn := AsNode(child)
	if cn == nil {
		fmt.Println("ChildNode failed to convert to node!!!")
		return
	}
	parent.AppendChild(cn)
}

// RemoveChildren will remove all children from the given parent node.
func RemoveChildren(parent NodeCompatible) {
	htmlNodeList := parent.ChildNodes()
	nodes := make([]*dom.Node, 0, htmlNodeList.Length())
	for ii := int(htmlNodeList.Length() - 1); ii >= 0; ii-- {
		nodes = append(nodes, htmlNodeList.Item(uint(ii)))
	}
	for _, node := range nodes {
		parent.RemoveChild(node)
	}
}

// InnerText returns the first TextNode within ElementNode. If none is found one is created.
func InnerText(element *html.HTMLElement) *dom.Node {
	children := element.ChildNodes()
	for ii := uint(0); ii < children.Length(); ii++ {
		child := children.Index(ii)
		if child.NodeName() == "#text" {
			return child
		}
	}
	return nil
}

// EventOptions are the options passed to AddEventListener used in the `OnWith` function.
type EventOptions struct {
	Capture, Once, Passive bool
}

// On adds a callback to the element when event type happens, using `element.AddEventListener`.
// All event options default to false. See OnWith() to set options.
func On(element EventTargetCompatible, eventType string, callback func(ev *domcore.Event)) {
	OnWith(element, eventType, callback, EventOptions{})
}

// OnWith adds a callback to the element when event type happens, using `element.AddEventListener`.
// There is no easy way to remove the event listener from Go :(
func OnWith(element EventTargetCompatible, eventType string, callback func(ev *domcore.Event), options EventOptions) {
	if len(eventType) < 2 || eventType[:2] == "on" {
		log.Fatalf("Suspicious eventType name: %s -- REMOVE THE \"ON\" PREFIX from the eventType", eventType)
	}
	jsOptions := map[string]interface{}{
		"Capture": options.Capture,
		"Once":    options.Once,
		"Passive": options.Passive,
	}
	element.AddEventListener(eventType, domcore.NewEventListenerFunc(callback),
		domcore.UnionFromJS(js.ValueOf(jsOptions)))
}

// DiscardEvent is an event handler that stops propagation and the
// default browser handling. Can be used as a parameter to the `On` function.
func DiscardEvent(e *domcore.Event) {
	e.StopPropagation()
	e.PreventDefault()
}

// Cookie returns the named cookie. Returns ifMissing, if the cookie was not found.
func Cookie(name, ifMissing string) (value string) {
	allCookies := Doc.Cookie()
	cookies := strings.Split(allCookies, "; ")
	for _, kv := range cookies {
		if kv == "" {
			continue
		}
		idx := strings.Index(kv, "=")
		var key, value string
		if idx == -1 {
			key = kv
		} else {
			key = kv[0:idx]
			value = kv[idx+1:]
		}
		if key == name {
			return value
		}
	}
	return ifMissing
}

// SetCookie sets tne named cookie to the given value.
func SetCookie(name, value string) {
	Doc.SetCookie(fmt.Sprintf("%s=%s", name, value))
}

// DelCookie will delete the given cookie, by giving an expired time.
func DelCookie(name string) {
	SetCookie(name, "; expires=Thu, 01 Jan 1970 00:00:00 GMT")
}

var _ = DelCookie

// Alert opens a small dialog window with the alert message.
func Alert(msg string) {
	Win.Alert2(msg)
}

// StyleSet sets the style of an element.
func StyleSet(e HtmlElementCompatible, key, value string) {
	e.Style().SetProperty(key, value, nil)
}

// Display sets the display style (CSS) of an element.
func Display(e HtmlElementCompatible, mode string) {
	StyleSet(e, "display", mode)
}

// DisplayOff sets the display style (CSS) of an element to "none".
func DisplayOff(e HtmlElementCompatible) {
	Display(e, "none")
}

// DisplayInline sets the display style (CSS) of an element to "inline".
func DisplayInline(e HtmlElementCompatible) {
	Display(e, "inline")
}

// DisplayBlock sets the display style (CSS) of an element to "block".
func DisplayBlock(e HtmlElementCompatible) {
	Display(e, "block")
}

// DisableIf disables button conditioned on evaluation.
func DisableIf(e HtmlElementCompatible, disable bool) {
	if disable {
		Disable(e)
	} else {
		Enable(e)
	}
}

// Disable disables buttons.
func Disable(e HtmlElementCompatible) {
	e.SetAttribute("disable", "")
}

// Enable enables buttons.
func Enable(e HtmlElementCompatible) {
	e.RemoveAttribute("disable")
}

// AsHTML casts some type of element as an HTMLElement.
func AsHTML(e EventTargetCompatible) *html.HTMLElement {
	if IsNil(e) {
		return nil
	}
	return html.HTMLElementFromJS(e.JSValue())
}

// AsInput casts some type of element as an HTMLInputElement.
func AsInput(e EventTargetCompatible) *html.HTMLInputElement {
	if IsNil(e) {
		return nil
	}
	return html.HTMLInputElementFromJS(e.JSValue())
}

// AsButton casts some type of element as an HTMLButtonElement.
func AsButton(e EventTargetCompatible) *html.HTMLButtonElement {
	if IsNil(e) {
		return nil
	}
	return html.HTMLButtonElementFromJS(e.JSValue())
}

// AsTable casts some type of element as an HTMLTableElement.
func AsTable(e EventTargetCompatible) *html.HTMLTableElement {
	if IsNil(e) {
		return nil
	}
	return html.HTMLTableElementFromJS(e.JSValue())
}

// AsTR casts some type of element as an HTMLTableRowElement.
func AsTR(e EventTargetCompatible) *html.HTMLTableRowElement {
	if IsNil(e) {
		return nil
	}
	return html.HTMLTableRowElementFromJS(e.JSValue())
}

func TRAddValue(tr *html.HTMLTableRowElement, value interface{}) *html.HTMLTableColElement {
	td := html.HTMLTableColElementFromJS(Elem("td").JSValue())
	td.SetInnerText(fmt.Sprintf("%s", value))
	Append(tr, td)
	return td
}

// AsSpan casts some type of element as an HTMLSpanElement.
func AsSpan(e EventTargetCompatible) *html.HTMLSpanElement {
	if IsNil(e) {
		return nil
	}
	return html.HTMLSpanElementFromJS(e.JSValue())
}

// AsSelect casts some type of element as an HTMLSelectElement.
func AsSelect(e EventTargetCompatible) *html.HTMLSelectElement {
	if IsNil(e) {
		return nil
	}
	return html.HTMLSelectElementFromJS(e.JSValue())
}

// AsOption casts some type of element as an HTMLOptionElement.
func AsOption(e EventTargetCompatible) *html.HTMLOptionElement {
	if IsNil(e) {
		return nil
	}
	return html.HTMLOptionElementFromJS(e.JSValue())
}

// AsCanvas casts some type of element as an HTMLCanvasElement.
func AsCanvas(e EventTargetCompatible) *canvas.HTMLCanvasElement {
	if IsNil(e) {
		return nil
	}
	return canvas.HTMLCanvasElementFromJS(e.JSValue())
}

// AsImage casts some type of element as an HTMLImageElement.
func AsImage(e EventTargetCompatible) *html.HTMLImageElement {
	if IsNil(e) {
		return nil
	}
	return html.HTMLImageElementFromJS(e.JSValue())
}

// AsScript casts some type of element as an HTMLScriptElement.
func AsScript(e EventTargetCompatible) *html.HTMLScriptElement {
	if IsNil(e) {
		return nil
	}
	return html.HTMLScriptElementFromJS(e.JSValue())
}

// AsMouseEvent casts some type of element as an HTMLMouseEventElement.
func AsMouseEvent(e EventCompatible) *htmlevent.MouseEvent {
	if IsNil(e) {
		return nil
	}
	return htmlevent.MouseEventFromJS(e.JSValue())
}

// AsKeyboardEvent casts some type of element as an HTMLKeyboardEventElement.
func AsKeyboardEvent(e EventCompatible) *htmlevent.KeyboardEvent {
	if IsNil(e) {
		return nil
	}
	return htmlevent.KeyboardEventFromJS(e.JSValue())
}

// AsAudio casts some type of element as an audio.AudioNode.
func AsAudio(e *html.HTMLElement) *Audio {
	if IsNil(e) {
		return nil
	}
	return &Audio{e.JSValue()}
}

// SprintDuration pretty-prints duration given in milliseconds
func SprintDuration(milliseconds int64) string {
	// Show seconds.
	secs := milliseconds / 1000
	if secs < 60 {
		return fmt.Sprintf("%ds", secs)
	}

	// Show minutes.
	minutes := secs / 60
	secs -= minutes * 60
	if minutes < 60 {
		if secs == 0 {
			return fmt.Sprintf("%dmin", minutes)
		} else {
			return fmt.Sprintf("%dmin %ds", minutes, secs)
		}
	}

	// Show hours.
	hours := minutes / 60
	minutes -= hours * 60
	if minutes == 0 {
		return fmt.Sprintf("%dh", hours)
	} else {
		return fmt.Sprintf("%dh %dmin", hours, minutes)
	}
}

// Audio is a simple wrapper over JS value that should represent a "<audio ...>" node.
type Audio struct {
	value js.Value
}

// Play will play the audio object using Javascript functionality.
func (a *Audio) Play() { a.value.Call("play") }
