// +build darwin linux

package main

import (
	"unsafe"
)

/*
#cgo LDFLAGS: -lpam -fPIC
#include <security/pam_appl.h>
#include <stdlib.h>

extern int pam_prompt(pam_handle_t *pamh, int style, char **response, const char *fmt);

*/
import "C"

func init() {
	if !disablePtrace() {
		pamLog("unable to disable ptrace")
	}
}

//export pam_sm_authenticate
func pam_sm_authenticate(pamh *C.pam_handle_t, flags, argc C.int, argv **C.char) C.int {
	var nameBuf *C.char
	defer C.free(unsafe.Pointer(nameBuf))
	C.pam_prompt(pamh, C.PAM_PROMPT_ECHO_ON, &nameBuf, C.CString("Your 6-digit token (or \"push\" or \"sms\"): "))
	// pam_info() is defined like this anyway
	C.pam_prompt(pamh, C.PAM_TEXT_INFO, nil, nameBuf)
	auth := Authy{
		APIKey:  "x",
		//BaseURL: "https://api.authy.com",
		BaseURL: "x",
	}
	res, err := auth.SendApprovalRequest("x", "say yes plsplspls")
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
