package api

import (
	"fmt"
	"strconv"

	"github.com/palantir/stacktrace"
)

// Structure of the JSON object that it's sent to request a new url to upload a new photo
type RequestUploadURL struct {
	ProtocolVersion      string               `json:"protocolVersion"`
	CreateSessionRequest CreateSessionRequest `json:"createSessionRequest"`
}

// Inner object of the request to get a new url to upload a photo.
type CreateSessionRequest struct {
	// The fields array is a slice that should contain only ExternalField or InternalField structs
	Fields []interface{} `json:"fields"`
}

// Possible field for the Fields slice in the CreateSessionRequest struct
type ExternalField struct {
	External ExternalFieldObject `json:"external"`
}

// Possible field for the Fields slice in the CreateSessionRequest struct
type InlinedField struct {
	Inlined InlinedFieldObject `json:"inlined"`
}

// Struct that describes the file that need to be uploaded. This objects should be contained in a ExternalField
type ExternalFieldObject struct {
	Name     string `json:"name"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

// Struct used to define parameters of the upload. This object should be contained in a InternalField
type InlinedFieldObject struct {
	Name        string `json:"name"`
	Content     string `json:"contentType"`
	ContentType string `json:"contentType"`
}

// Struct that represents the JSON response from the request to get an upload url
type UploadURLRequestResponse struct {
	SessionStatus SessionStatus `json:"sessionStatus"`
}

// Struct that represents the inner JSON object of the UploadURLRequestResponse
type SessionStatus struct {
	// Field used in the response for the request to get a new upload URL
	ExternalFieldTransfers []struct {
		Name    string `json:"name"`
		PutInfo struct {
			Url string `json:"url"`
		} `json:"putInfo"`
	} `json:"externalFieldTransfers"`

	// Field used in the UploadImageResponse
	AdditionalInfo struct {
		UploadService struct {
			CompletionInfo struct {
				CustomerSpecificInfo struct {
					UploadTokenBase64 string `json:"upload_token_base64"`
				} `json:"customerSpecificInfo"`
			} `json:"completionInfo"`
		} `json:"uploader_service.GoogleRupioAdditionalInfo"`
	} `json:"additionalInfo"`
}

// JSON representation of the response from the upload image request.
type UploadImageResponse struct {
	SessionStatus SessionStatus `json:"sessionStatus"`
}

type EnableImageRequest []interface{}

type FirstItemEnableImageRequest []InnerItemFirstItemEnableImageRequest

type InnerItemFirstItemEnableImageRequest interface{}

type SecondInnerArray []MapOfItemsToEnable

type MapOfItemsToEnable map[string]ItemToEnable

type ItemToEnable []ItemToEnableArray

type ItemToEnableArray []InnerItemToEnableArray

type InnerItemToEnableArray interface{}

type EnableImageResponse []interface{}

func (eir EnableImageResponse) getInner6Array() ([]interface{}, error) {
	var inner3Array []interface{}
	if len(eir) > 0 {
		if inner1Array, ok := eir[0].([]interface{}); ok && len(inner1Array) >= 2 {
			if inner2Map, ok := inner1Array[1].(map[string]interface{}); ok {
				inner3Array = inner2Map[strconv.Itoa(EnablePhotoKey)].([]interface{})
			}
		}
	}
	if len(inner3Array) > 0 {
		if inner4Array, ok := inner3Array[0].([]interface{}); ok && len(inner4Array) > 0 {
			if inner5Array, ok := inner4Array[0].([]interface{}); ok && len(inner5Array) >= 2 {
				return inner5Array[1].([]interface{}), nil
			}
		}
	}
	return nil, fmt.Errorf("no inner6Array")
}

func (eir EnableImageResponse) getEnabledImageId() (string, error) {
	inner6Array, err := eir.getInner6Array()
	if err != nil {
		return "", stacktrace.Propagate(err, "no enabledImageID")
	}
	if len(inner6Array) >= 1 {
		if enabledImageId, ok := inner6Array[0].(string); ok {
			return enabledImageId, nil
		}
	}

	return "", fmt.Errorf("no enabledImageID")
}

func (eir EnableImageResponse) getEnabledImageURL() (string, error) {
	inner6Array, err := eir.getInner6Array()
	if err != nil {
		return "", stacktrace.Propagate(err, "no enabledImageURL")
	}
	if len(inner6Array) >= 2 {
		inner7Array := inner6Array[1].([]interface{})
		if enabledImageURL, ok := inner7Array[0].(string); ok {
			return enabledImageURL, nil
		}
	}
	return "", fmt.Errorf("no enabledImageURL")
}
