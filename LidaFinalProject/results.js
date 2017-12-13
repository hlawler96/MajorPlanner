
$(document).ready(function () {
  console.log("The js is hooked up");
  sessionId=  window.localStorage.getItem('sessionId');
  // alert("about to try populate");
  response = ""
  populateCoursesTaken();
});


function populateCoursesTaken(){
// alert("first line of populate");
API_URL = "http://ec2-18-217-72-185.us-east-2.compute.amazonaws.com:8080/PossiblePrograms/?sessionId=" + sessionId;
// alert(API_URL);
var xhr = createCORSRequest('GET', API_URL);
xhr.responseType = 'text';
if (!xhr) {
 alert('CORS not supported');
 return;
}
// Response handlers.
xhr.onload = function() {
  var displayLength = 0;
  var jsonResponse = JSON.parse(xhr.responseText);
  response = jsonResponse;
  if(jsonResponse.strictRemainingCourses.length > 0){
    $('#coursesToTake').append("<tr> <th>Required</th> <th> <ul id = 'strict'> </ul> </th> </tr>")
  }

  //strict remaining courses
  for(var i = 0; i < jsonResponse.strictRemainingCourses.length; i++){
  var dept = jsonResponse.strictRemainingCourses[i].program;
  var num = jsonResponse.strictRemainingCourses[i].number;
  $('#strict').append("<li class = courses> " + dept + " " + num + " </li>");

}
//loose remaining courses
req = ""
j = 0;
for(var i = 0; i < jsonResponse.looseRemainingCourses.length; i++){
  var dept = jsonResponse.looseRemainingCourses[i].course.program;
  var num = jsonResponse.looseRemainingCourses[i].course.number;
  var temp = jsonResponse.looseRemainingCourses[i].requirement;
  if(temp != req){
      $('#coursesToTake').append("<tr> <th>" + temp + "</th> <th> <ul id = 'loose" + i + "'> </ul> </th> </tr>")
      req = temp;
      j = i;
  }
  //do CSS for this
  $('#loose' + j).append("<li class = 'courses'>" + dept + " " + num + "</li> ");
}

//possible programs
for(var i = 0; i < jsonResponse.possiblePrograms.length; i++){
  var dept = jsonResponse.possiblePrograms[i].dept;
  var type = jsonResponse.possiblePrograms[i].type;
  $('#majors').append("<option value='" + dept + " " + type + "'>" + dept + " " + type + "</option>");
}



};


xhr.onerror = function() {
    alert('FAILURE');
};

xhr.send();
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
