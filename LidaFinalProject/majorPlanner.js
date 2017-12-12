$(document).ready(function () {
  console.log("The js is hooked up");
});


//get values for username and password from website
//store them in variables
//print them to console
//then add api calls to get them and check them
function loginSubmit() {
  alert("in loginSubmit!");
  uname = document.getElementById("uname").value;
  alert(uname);
  var pass = document.getElementById("pass").value;
  alert(pass);
  API_URL = "http://ec2-18-217-72-185.us-east-2.compute.amazonaws.com:8080/Login/?username=" + uname + "&password=" + pass;
  sessionId = "";
  alert(API_URL);
  var xhr = createCORSRequest('GET', API_URL);
  xhr.responseType = 'text';
 if (!xhr) {
   alert('CORS not supported');
   return;
 }
 // Response handlers.
  xhr.onload = function() {
    var jsonResponse = JSON.parse(xhr.responseText);
    sessionId = jsonResponse.sessionId;
  };

  xhr.onerror = function() {
      alert('FAILURE');
  };

  xhr.send();

  //call API, check if they match
}

function registerSubmit() {
  alert("in registerSubmit!");
  uname = document.getElementById("uname").value;
  alert(uname);
  var pass = document.getElementById("pass").value;
  alert(pass);
  //call API, save these in the DB
}

function selectAdditionalDegree(){
  additionalDegree = document.getElementById("majors").value;
  alert(additionalDegree);

}


//get value from major check box
//store in variable
//print in onto the website
//then add api calls to get classes and use jquery to add checkboxes




function deptSubmit() {
  // alert("in deptSubmit!");
  var ba_checked = document.getElementById("ba").checked;
  var bs_checked = document.getElementById("bs").checked;
  var econ_checked = document.getElementById("econ").checked;
  var math_checked = document.getElementById("math").checked;
  if(ba_checked){
    document.getElementById("second-degree-title").innerHTML = "COMP - B.A. - Select the classes you have taken.";
  }else if(bs_checked){
    document.getElementById("second-degree-title").innerHTML = "COMP - B.S. - Select the classes you have taken.";
  }else if(econ_checked){
    document.getElementById("second-degree-title").innerHTML = "ECON - B.A. - Select the classes you have taken.";
  }else if(math_checked){
    document.getElementById("second-degree-title").innerHTML = "MATH - B.A. - Select the classes you have taken.";
  }

  var sems_input = '<input id="sems-left" type=number max=7 minimum=1 placeholder="Enter Number">';
  $("#sems-left-input-div").empty();
  $('#sems-left-input-div').append('Semesters Left: '+ sems_input);

  var submit_classes = '<input id="submt-classes" onclick="degreeFinder()" type="submit" value="Find Me Another Degree!">';
  $("#submt-classes-div").empty();
  $('#submt-classes-div').append(submit_classes);

}

function degreeFinder(){
  alert("in degreeFinder!");
  var sems_left = document.getElementById("sems-left").value;
  alert(sems_left);

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













/////bottom
