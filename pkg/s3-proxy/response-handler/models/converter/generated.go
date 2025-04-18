// Code generated by github.com/jmattheis/goverter, DO NOT EDIT.
//go:build !goverter

package converter

import (
	models "github.com/oxyno-zeta/s3-proxy/pkg/s3-proxy/response-handler/models"
	"net/http"
	"net/url"
)

func init() {
	ConvertAndSanitizeHTTPRequest = func(source *http.Request) *models.LightSanitizedRequest {
		var pModelsLightSanitizedRequest *models.LightSanitizedRequest
		if source != nil {
			var modelsLightSanitizedRequest models.LightSanitizedRequest
			modelsLightSanitizedRequest.URL = sanitizeURL((*source).URL)
			modelsLightSanitizedRequest.Header = httpHeaderToHttpHeader((*source).Header)
			modelsLightSanitizedRequest.Trailer = httpHeaderToHttpHeader((*source).Trailer)
			modelsLightSanitizedRequest.RemoteAddr = sanitizeString((*source).RemoteAddr)
			modelsLightSanitizedRequest.Method = sanitizeString((*source).Method)
			modelsLightSanitizedRequest.Proto = sanitizeString((*source).Proto)
			modelsLightSanitizedRequest.Pattern = sanitizeString((*source).Pattern)
			modelsLightSanitizedRequest.RequestURI = sanitizeString((*source).RequestURI)
			modelsLightSanitizedRequest.Host = sanitizeString((*source).Host)
			if (*source).TransferEncoding != nil {
				modelsLightSanitizedRequest.TransferEncoding = make([]string, len((*source).TransferEncoding))
				for i := 0; i < len((*source).TransferEncoding); i++ {
					modelsLightSanitizedRequest.TransferEncoding[i] = sanitizeString((*source).TransferEncoding[i])
				}
			}
			modelsLightSanitizedRequest.ProtoMajor = (*source).ProtoMajor
			modelsLightSanitizedRequest.ContentLength = (*source).ContentLength
			modelsLightSanitizedRequest.ProtoMinor = (*source).ProtoMinor
			pModelsLightSanitizedRequest = &modelsLightSanitizedRequest
		}
		return pModelsLightSanitizedRequest
	}
	sanitizeURL = func(source *url.URL) *url.URL {
		var pUrlURL *url.URL
		if source != nil {
			var urlURL url.URL
			urlURL.Scheme = sanitizeString((*source).Scheme)
			urlURL.Opaque = sanitizeString((*source).Opaque)
			urlURL.User = sanitizeURLUserInfo((*source).User)
			urlURL.Host = sanitizeString((*source).Host)
			urlURL.Path = sanitizeString((*source).Path)
			urlURL.RawPath = sanitizeString((*source).RawPath)
			urlURL.OmitHost = (*source).OmitHost
			urlURL.ForceQuery = (*source).ForceQuery
			urlURL.RawQuery = sanitizeString((*source).RawQuery)
			urlURL.Fragment = sanitizeString((*source).Fragment)
			urlURL.RawFragment = sanitizeString((*source).RawFragment)
			pUrlURL = &urlURL
		}
		return pUrlURL
	}
}
func httpHeaderToHttpHeader(source http.Header) http.Header {
	var httpHeader http.Header
	if source != nil {
		httpHeader = make(http.Header, len(source))
		for key, value := range source {
			var stringList []string
			if value != nil {
				stringList = make([]string, len(value))
				for i := 0; i < len(value); i++ {
					stringList[i] = sanitizeString(value[i])
				}
			}
			httpHeader[sanitizeString(key)] = stringList
		}
	}
	return httpHeader
}
func init() {
	ConvertSanitizedToHTTPRequest = func(source *models.LightSanitizedRequest) *http.Request {
		var pHttpRequest *http.Request
		if source != nil {
			var httpRequest http.Request
			httpRequest.Method = (*source).Method
			httpRequest.URL = identity((*source).URL)
			httpRequest.Proto = (*source).Proto
			httpRequest.ProtoMajor = (*source).ProtoMajor
			httpRequest.ProtoMinor = (*source).ProtoMinor
			httpRequest.Header = identity((*source).Header)
			httpRequest.ContentLength = (*source).ContentLength
			httpRequest.TransferEncoding = identity((*source).TransferEncoding)
			httpRequest.Host = (*source).Host
			httpRequest.Trailer = identity((*source).Trailer)
			httpRequest.RemoteAddr = (*source).RemoteAddr
			httpRequest.RequestURI = (*source).RequestURI
			httpRequest.Pattern = (*source).Pattern
			pHttpRequest = &httpRequest
		}
		return pHttpRequest
	}
}
