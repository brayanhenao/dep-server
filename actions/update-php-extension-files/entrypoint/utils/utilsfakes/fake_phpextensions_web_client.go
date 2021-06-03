// Code generated by counterfeiter. DO NOT EDIT.
package utilsfakes

import (
	"sync"

	"github.com/paketo-buildpacks/dep-server/actions/update-php-extension-files/entrypoint/utils"
)

type FakePHPExtensionsWebClient struct {
	DownloadExtensionsSourceStub        func(string, string) error
	downloadExtensionsSourceMutex       sync.RWMutex
	downloadExtensionsSourceArgsForCall []struct {
		arg1 string
		arg2 string
	}
	downloadExtensionsSourceReturns struct {
		result1 error
	}
	downloadExtensionsSourceReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakePHPExtensionsWebClient) DownloadExtensionsSource(arg1 string, arg2 string) error {
	fake.downloadExtensionsSourceMutex.Lock()
	ret, specificReturn := fake.downloadExtensionsSourceReturnsOnCall[len(fake.downloadExtensionsSourceArgsForCall)]
	fake.downloadExtensionsSourceArgsForCall = append(fake.downloadExtensionsSourceArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	stub := fake.DownloadExtensionsSourceStub
	fakeReturns := fake.downloadExtensionsSourceReturns
	fake.recordInvocation("DownloadExtensionsSource", []interface{}{arg1, arg2})
	fake.downloadExtensionsSourceMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakePHPExtensionsWebClient) DownloadExtensionsSourceCallCount() int {
	fake.downloadExtensionsSourceMutex.RLock()
	defer fake.downloadExtensionsSourceMutex.RUnlock()
	return len(fake.downloadExtensionsSourceArgsForCall)
}

func (fake *FakePHPExtensionsWebClient) DownloadExtensionsSourceCalls(stub func(string, string) error) {
	fake.downloadExtensionsSourceMutex.Lock()
	defer fake.downloadExtensionsSourceMutex.Unlock()
	fake.DownloadExtensionsSourceStub = stub
}

func (fake *FakePHPExtensionsWebClient) DownloadExtensionsSourceArgsForCall(i int) (string, string) {
	fake.downloadExtensionsSourceMutex.RLock()
	defer fake.downloadExtensionsSourceMutex.RUnlock()
	argsForCall := fake.downloadExtensionsSourceArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakePHPExtensionsWebClient) DownloadExtensionsSourceReturns(result1 error) {
	fake.downloadExtensionsSourceMutex.Lock()
	defer fake.downloadExtensionsSourceMutex.Unlock()
	fake.DownloadExtensionsSourceStub = nil
	fake.downloadExtensionsSourceReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakePHPExtensionsWebClient) DownloadExtensionsSourceReturnsOnCall(i int, result1 error) {
	fake.downloadExtensionsSourceMutex.Lock()
	defer fake.downloadExtensionsSourceMutex.Unlock()
	fake.DownloadExtensionsSourceStub = nil
	if fake.downloadExtensionsSourceReturnsOnCall == nil {
		fake.downloadExtensionsSourceReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.downloadExtensionsSourceReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakePHPExtensionsWebClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.downloadExtensionsSourceMutex.RLock()
	defer fake.downloadExtensionsSourceMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakePHPExtensionsWebClient) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ utils.PHPExtensionsWebClient = new(FakePHPExtensionsWebClient)
