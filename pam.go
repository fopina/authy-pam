// +build darwin linux

package main

import (
	"unsafe"
)

/*
#cgo LDFLAGS: -lpam -fPIC
#include <security/pam_appl.h>
#include <stdlib.h>

int pam_prompt(pam_handle_t *pamh, int style, char **response, const char *fmt);
int pam_get_user(pam_handle_t *pamh, const char **user, const char *prompt);
*/
import "C"

func init() {
	if !disablePtrace() {
		pamLog("unable to disable ptrace")
	}
}

//export pam_sm_authenticate
func pam_sm_authenticate(pamh *C.pam_handle_t, flags, argc C.int, argv **C.char) C.int {
	var tempBuf *C.char
	defer C.free(unsafe.Pointer(tempBuf))

	config := Config{}
	config.LoadFromFile("/etc/pam_authy.conf")

	pam_err := C.pam_get_user(pamh, &tempBuf, nil)
	if (pam_err != C.PAM_SUCCESS) {
        return pam_err;
    }
	
	user := C.GoString(tempBuf)

	if config.Users[user] == "" {
		// Authy not setup for this user
		return C.PAM_IGNORE
	}

	C.pam_prompt(pamh, C.PAM_PROMPT_ECHO_ON, &tempBuf, C.CString("Your 6-digit token (or \"push\" or \"sms\"): "))
	// pam_info() is defined like this anyway
	C.pam_prompt(pamh, C.PAM_PROMPT_ECHO_ON, nil, tempBuf)

	authy := Authy{
		APIKey:  config.Authy.Token,
		BaseURL: config.Authy.URL,
	}
	res, err := authy.SendApprovalRequest(config.Users[user], "say yes plsplspls")
	if err != nil {
		panic(err)
	}
	if res == OneTouchStatusApproved {
		return C.PAM_SUCCESS
	} else {
    	return C.PAM_AUTH_ERR
	}
}

//export pam_sm_setcred
func pam_sm_setcred(pamh *C.pam_handle_t, flags, argc C.int, argv **C.char) C.int {
    return C.PAM_IGNORE
}
