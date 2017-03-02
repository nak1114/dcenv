/*

 Copyright 2017 Kii Corporation
 http://kii.com

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.

*/

// Called by the "Log In" button.
function performDeleteAccount() {
    // Get the username and password from the UI.
    var username = document.getElementById("username-field").value;
    var password = document.getElementById("password-field").value;

    // Authenticate the user asynchronously.
    KiiUser.authenticate(username, password, {
        // If the user was authenticated
        success: function(theUser) {
            console.log("User authenticated: " + JSON.stringify(theUser));
            //var user = KiiUser.getCurrentUser();
            theUser.delete({
              success: function(theUser) {
                console.log("User deleted: " + JSON.stringify(theUser));
                // Go to the main screen.
                openDeleteAccountPage();
              },
              failure: function(theUser, errorString) {
                alert("Unable to delete account: " + errorString);
                // Error handling
              }
            });

        },
        // If the user was not authenticated
        failure: function(theUser, errorString) {
            console.log("Unable to authenticate user: " + errorString);
            alert("Unable to authenticate user: " + errorString);
        }
    });
}

// Called by the "Sign Up" button.
function performSignUp() {
    // Get the username and password from the UI.
    var username = document.getElementById("username-field").value;
    var password = document.getElementById("password-field").value;
    var email = document.getElementById("email-field").value;

    // Create a KiiUser object.
    var user = KiiUser.userWithEmailAddressAndUsername(email, username, password);
    // Register the user asynchronously.
    user.register({
        // If the user was registered
        success: function(theUser) {
            console.log("User registered: " + JSON.stringify(theUser));

            // Go to the main screen.
            openListPage();
        },
        // If the user was not registered
        failure: function(theUser, errorString) {
            alert("Unable to register user: " + errorString);
            console.log("Unable to register user: " + errorString);
        }
    });
}
// Called by the "Reset Password" button.
function performResetPassword() {
    // Get the username and password from the UI.
    var email = document.getElementById("email-field").value;

    KiiUser.resetPasswordWithNotificationMethod(email, "EMAIL", {
      success: function() {
            openResetPage();
      },
      failure: function(errorString) {
        // Error handling
            alert("Unable to reset password: " + errorString);
            console.log("Unable to reset password: " + errorString);
      }
    });
}
