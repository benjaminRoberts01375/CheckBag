package main

// Reset the user's password. Authentication is done with the JWT.
// func userResetPassword(w http.ResponseWriter, r *http.Request) {
// 	_, password, err := checkUserRequest[string](r)
// 	if err != nil {
// 		requestRespondCode(w, http.StatusForbidden)
// 		return
// 	}
// 	newPasswordHash, err := createPasswordHash(*password)
// 	if err != nil {
// 		requestRespondCode(w, http.StatusInternalServerError)
// 		return
// 	}
// 	fileSystem.SetUserData(string(newPasswordHash))
// 	requestRespondCode(w, http.StatusOK)
// }
