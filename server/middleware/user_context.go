package middleware

// type UserContext struct {
// 	currentUser    *utils.CurrentUserManager
// 	sessionManager *utils.SessionManager
// }

// func NewUserContext(currentUser *utils.CurrentUserManager, sessionManager *utils.SessionManager) *UserContext {
// 	return &UserContext{currentUser, sessionManager}
// }

// func (uc *UserContext) Middleware() negroni.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
// 		state, err := uc.sessionManager.GetState(r)
// 		if err != nil {
// 			log.Printf("UserContext Middleware: Failed to retrieve state from session: %s", err.Error())
// 		}
// 		if user := models.GetUser(state); state != "" && user != nil {
// 			// log.Printf("Current user with state %s", state)
// 			uc.currentUser.Set(r, user)
// 		} else {
// 			log.Printf("No current user found")
// 		}
// 		next(w, r)
// 	}
// }
