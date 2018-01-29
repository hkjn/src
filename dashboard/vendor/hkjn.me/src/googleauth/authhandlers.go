package googleauth

import (
	"html/template"
	"net/http"

	"github.com/golang/glog"
)

var (
	// Name of the top-level login template.
	TemplateName = "login"
	// Template to use for the login redirect.
	LoginTmpl = template.Must(template.New(TemplateName).Parse(tmpl))
	tmpl      = `{{define "scripts"}}
<script src="https://apis.google.com/js/client:platform.js" async defer></script>
<script src="https://ajax.googleapis.com/ajax/libs/jquery/2.1.1/jquery.min.js"></script>

<script>
  var helper = (function() {
    var BASE_API_PATH = 'plus/v1/';
    var authResult = undefined;

    return {
      /**
       * Hides the sign-in button and connects the server-side app after
       * the user successfully signs in.
       */
      signInCallback: function(authResult) {
        if (authResult['access_token']) {
          this.logIn(authResult.code);
        } else if (authResult["error"]) {
          $('#authResult').append('Logged out');
        }
      },
      logIn: function(code) {
        window.location.replace(
          window.location.origin + "/connect?state={{.StateToken}}&code=" + code);
      }
    };
  })();

function signInCallback(authResult) {
  helper.signInCallback(authResult);
}
</script>
{{end}}

{{define "body"}}
<div id="gConnect">
<span id="signinButton">
  <span
    class="g-signin"
    data-callback="signInCallback"
    data-clientid="{{.ClientId}}"
    data-accesstype="offline"
    data-cookiepolicy="single_host_origin"
    data-scope="profile">
  </span>
</span>
</div>
{{template "scripts" .}}
{{end}}

{{define "login"}}
<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
</head>
<body>
{{template "body" .}}
</body>
</html>
{{end}}`
)

// RequireLogin returns a wrapped HandlerFunc that enforces Google+ login.
//
// If the user is logged in, the specified HandlerFunc is called, otherwise the
// login page defined by LoginTmpl is served.
//
// RequireLogin returns HTTP 500 (Internal Server Error) if the
// template fails to render or the package has an internal error.
func RequireLogin(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		loggedIn, err := IsLoggedIn(r)
		if err != nil {
			glog.Errorf("failed to get login info: %v\n", err)
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}
		if loggedIn {
			glog.V(1).Infof("user is logged in, onward to original handler func\n")
			fn(w, r)
			return
		}
		glog.V(1).Infof("not logged in, fetching state token\n")
		li, err := LogIn(w, r)
		if err != nil {
			glog.Errorf("failed to get login info: %v\n", err)
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}
		err = LoginTmpl.ExecuteTemplate(w, TemplateName, li)
		if err != nil {
			glog.Errorf("failed to execute login template: %v\n", err)
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}
	}
}

// ConnectHandler finishes the connection process, exchanging the
// one-time authorization code for an access token and storing it in
// the session.
//
// ConnectHandler redirects to the request "referer" on successful login.
//
// ConnectHandler returns HTTP 401 (Unauthorized) if the user does not
// have access.
func ConnectHandler(w http.ResponseWriter, r *http.Request) {
	err := Connect(w, r)
	if err != nil {
		if isAccessDenied(err) {
			http.Error(w, "Access denied.", http.StatusUnauthorized)
		} else {
			glog.Errorf("error connecting to googleauth: %v\n", err)
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
		}
		return
	}
	glog.V(1).Infof("current user is connected, redirecting to %q\n", r.Referer())
	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}
