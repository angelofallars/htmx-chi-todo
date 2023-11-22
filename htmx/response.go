package htmx

import "net/http"

// Request headers
const (
	HeaderBoosted               = "HX-Boosted"
	HeaderCurrentURL            = "HX-Current-URL"
	HeaderHistoryRestoreRequest = "HX-History-Restore-Request"
	HeaderPrompt                = "HX-Prompt"
	HeaderRequest               = "Hx-Request"
	HeaderTarget                = "HX-Target"
	HeaderTriggerName           = "Hx-Trigger-Name"
)

// Common headers
const (
	HeaderTrigger = "HX-Trigger"
)

// Response headers
const (
	HeaderLocation           = "HX-Location"
	HeaderPushURL            = "HX-Push-Url"
	HeaderRedirect           = "HX-Redirect"
	HeaderRefresh            = "HX-Refresh"
	HeaderReplaceUrl         = "HX-Replace-Url"
	HeaderReswap             = "HX-Reswap"
	HeaderRetarget           = "HX-Retarget"
	HeaderReselect           = "HX-Reselect"
	HeaderTriggerAfterSettle = "HX-Trigger-After-Settle"
	HeaderTriggerAfterSwap   = "HX-Trigger-After-Swap"
)

type (
	// Interface to define 'hx-swap' values.
	swapValue interface {
		Swap() string
	}
	// Concrete 'hx-swap' values.
	swap string
)

func (s swap) Swap() string {
	return string(s)
}

var (
	// Replace the inner html of the target element
	SwapInnerHTML swap = "innerHTML"

	// Replace the entire target element with the response
	SwapOuterHTML swap = "outerHTML"

	// Insert the response before the target element
	SwapBeforeBegin swap = "beforebegin"

	// Insert the response before the first child of the target element
	SwapAfterBegin swap = "afterbegin"

	// Insert the response after the last child of the target element
	SwapBeforeEnd swap = "beforeend"

	// Insert the response after the target element
	SwapAfterEnd swap = "afterend"

	// Deletes the target element regardless of the response
	SwapDelete swap = "delete"

	// Does not append content from response (out of band items will still be processed).
	SwapNone swap = "none"
)

var falseString = "false"

type response struct {
	headers map[string]string
}

// NewResponse returns a new HTMX response header writer.
func NewResponse() response {
	return response{
		headers: make(map[string]string),
	}
}

// Write applies the defined HTMX headers to a given response writer.
func (r response) Write(w http.ResponseWriter) error {
	header := w.Header()
	for key, value := range r.headers {
		header.Add(key, value)
	}

	return nil
}

// Location	allows you to do a client-side redirect that does not do a full page reload.
//
// Sets the 'HX-Location' header.
func (r response) Location(url string) response {
	r.headers[HeaderLocation] = url
	return r
}

// PushURL pushes a new URL into the browser location history.
//
// Sets the 'HX-Push-Url' header.
func (r response) PushURL(url string) response {
	r.headers[HeaderPushURL] = url
	return r
}

// PreventPushURL prevents the browser’s history from being updated.
//
// Sets the 'HX-Push-Url' header.
func (r response) PreventPushURL(url string) response {
	r.headers[HeaderPushURL] = falseString
	return r
}

// Redirect does a client-side redirect to a new location.
//
// Sets the 'HX-Redirect' header.
func (r response) Redirect(url string) response {
	r.headers[HeaderRedirect] = url
	return r
}

// If set to true, Refresh makes the client-side do a full refresh of the page.
//
// Sets the 'HX-Refresh' header.
func (r response) Refresh(shouldRefresh bool) response {
	if shouldRefresh {
		r.headers[HeaderRefresh] = "true"
	} else {
		r.headers[HeaderRefresh] = "false"
	}
	return r
}

// ReplaceURL replaces the current URL in the browser location history.
//
// Sets the 'HX-Replace-Url' header.
func (r response) ReplaceURL(url string) response {
	r.headers[HeaderReplaceUrl] = url
	return r
}

// PreventReplaceURL prevents the browser’s current URL from being updated.
//
// Sets the 'HX-Replace-Url' header to 'false'.
func (r response) PreventReplaceURL(url string) response {
	r.headers[HeaderReplaceUrl] = falseString
	return r
}

// Reswap allows you to specify how the response will be swapped. Accepts Swap values from this library.
//
// Sets the 'HX-Reswap' header.
func (r response) Reswap(sv swapValue) response {
	r.headers[HeaderReswap] = sv.Swap()
	return r
}

// Retarget accepts a CSS selector that updates the target of the content update to a different element on the page.
//
// Sets the 'HX-Retarget' header.
func (r response) Retarget(selector string) response {
	r.headers[HeaderRetarget] = selector
	return r
}

// Reselect accepts a CSS selector that allows you to choose which part of the response is used to be swapped in.
// Overrides an existing hx-select on the triggering element.
//
// Sets the 'HX-Reselect' header.
func (r response) Reselect(selector string) response {
	r.headers[HeaderReselect] = selector
	return r
}

func testResponses() {
	NewResponse().Refresh(true).ReplaceURL("/eseee")
	NewResponse().Retarget("#errors")
	NewResponse().Reswap(SwapBeforeEnd)
}
