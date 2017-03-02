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

// Define global variables used for creating objects.
var objectCount = 0;
var BUCKET_NAME = "myBucket";
var OBJECT_KEY = "myObjectValue";

// Called when the user is logged in or registered.
function openListPage() {
    document.getElementById("login-page").style.display = "none";
    document.getElementById("list-page").style.display = "block";

}
// Called when the user is logged in or registered.
function openResetPage() {
    document.getElementById("login-page").style.display = "none";
    document.getElementById("reset-page").style.display = "block";

}

// Called when the user is logged in or registered.
function openDeleteAccountPage() {
    document.getElementById("login-page").style.display = "none";
    document.getElementById("delete-page").style.display = "block";

}

