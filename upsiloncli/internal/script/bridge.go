package script

import (
	"fmt"
)

func (a *Agent) bindJSAPI() {
	upsilonObj := map[string]interface{}{
		"call":         a.jsCall,
		"waitForEvent": a.jsWaitForEvent,
		"getContext":   a.jsGetContext,
		"setContext":   a.jsSetContext,
		"log":          a.jsLog,
	}
	a.VM.Set("upsilon", upsilonObj)
}

func (a *Agent) jsLog(msg interface{}) {
	fmt.Fprintf(a.Logger, "[%s] %v\n", a.ID, msg)
}

func (a *Agent) jsCall(routeName string, params map[string]interface{}) (interface{}, error) {
	ep := a.Registry.Get(routeName)
	if ep == nil {
		return nil, fmt.Errorf("unknown route: %s", routeName)
	}

	// Convert JS params to string map expected by endpoint.Execute
	inputs := make(map[string]string)
	for k, v := range params {
		inputs[k] = fmt.Sprintf("%v", v)
	}

	resp, err := ep.ExecuteRaw(a.Client, a.Session, inputs)
	if err != nil {
		return nil, err
	}

	// Ensure WebSockets are synced if auth happened (token might have been set)
	a.Listener.Sync()

	return resp.Data, nil
}

func (a *Agent) jsGetContext(key string) string {
	val, _ := a.Session.Get(key)
	return val
}

func (a *Agent) jsSetContext(key, value string) {
	a.Session.Set(key, value)
}

func (a *Agent) jsWaitForEvent(eventName string, timeoutMs int) (interface{}, error) {
	return a.Listener.WaitForData(eventName, timeoutMs)
}
