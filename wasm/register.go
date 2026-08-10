//go:build js && wasm

package wasm

import (
	"encoding/json"
	"syscall/js"

	"github.com/runik-3/builder/wikiBot"
)

func buildResponse(data any) (string, error) {
	res, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

// Wrapper around go func that return a promise-like shape so we can call
// async wasm functions from JS
func asyncFunc(fn func(_this js.Value, _args []js.Value) (string, error)) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		promiseHandler := js.FuncOf(func(promiseThis js.Value, promiseArgs []js.Value) any {
			resolve := promiseArgs[0]
			reject := promiseArgs[1]

			go func() {
				data, error := fn(this, args)
				if error != nil {
					reject.Invoke(error)
				}
				resolve.Invoke(data)
			}()
			return nil
		})
		return js.Global().Get("Promise").New(promiseHandler)
	})
}

func getWikiInfo(this js.Value,args []js.Value) (string, error) {
	url := args[0].String()
	data, err := wikiBot.GetWikiDetails(url)
	if err != nil {
		return "", err
	}
	return buildResponse(data)
}

// Registers WASM functions into the global JS context
func RegisterWasmFuncs() {
	js.Global().Set("getWikiInfo", asyncFunc(getWikiInfo))
}
