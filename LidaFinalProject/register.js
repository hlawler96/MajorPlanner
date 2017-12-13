$(document).ready(function () {
  console.log("The js is hooked up");
  if(typeof(Storage)!=="undefined"){
    if(window.localStorage.getItem('sessionId')){
       alert(window.localStorage.getItem('sessionId'))
    }else {
      alert("NO Session Id");
    }
  } else{
  alert("storage not supported by browser");
  }
  var fun = function(e){
      if(e.keyCode==13) registerSubmit();
  }
  $('#uname').keypress(fun);
  $('#pass').keypress(fun);
});

function registerSubmit() {
  var name = document.getElementById("name").value;
  var uname = document.getElementById("uname").value;
  var pass = document.getElementById("pass").value;
  API_URL = "http://ec2-18-217-72-185.us-east-2.compute.amazonaws.com:8080/SignUp/?name=" + name + "&username="+ uname + "&password=" + pass;
  sessionId = "";
  var xhr = createCORSRequest('GET', API_URL);
  xhr.responseType = 'text';
 if (!xhr) {
   alert('CORS not supported');
   return;
 }
 // Response handlers
  xhr.onload = function() {
    var jsonResponse = JSON.parse(xhr.responseText);
    alert(jsonResponse.sessionId);
    window.localStorage.setItem('sessionId',jsonResponse.sessionId);
    if( jsonResponse.sessionId == ""){
      alert("Username Already Taken");
    }else {
      alert(window.localStorage.getItem('sessionId'));
    }
     window.location.href = "planner.html";
  };
  xhr.onerror = function() {
      alert('FAILURE');
  };
  xhr.send();
  // window.location.href = "planner.html";

}
function createCORSRequest(method, url) {
  var xhr = new XMLHttpRequest();
  if ("withCredentials" in xhr) {

    // Check if the XMLHttpRequest object has a "withCredentials" property.
    // "withCredentials" only exists on XMLHTTPRequest2 objects.
    xhr.open(method, url, true);

  } else if (typeof XDomainRequest != "undefined") {

    // Otherwise, check if XDomainRequest.
    // XDomainRequest only exists in IE, and is IE's way of making CORS requests.
    xhr = new XDomainRequest();
    xhr.open(method, url);

  } else {

    // Otherwise, CORS is not supported by the browser.
    xhr = null;

  }
  return xhr;
}
