
$(document).ready(function () {
  console.log("The js is hooked up");
  sessionId=  window.localStorage.getItem('sessionId');
  alert("about to try populate");
  populateCoursesTaken();
});


function populateCoursesTaken(){
alert("first line of populate");
API_URL = "http://ec2-18-217-72-185.us-east-2.compute.amazonaws.com:8080/PossiblePrograms/?sessionId=" + sessionId;
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
  for(var i = 0; i < jsonResponse.strictRemainingCourses.length; i++){
  var dept = jsonResponse.strictRemainingCourses[i].program;
  var num = jsonResponse.strictRemainingCourses[i].number;
  $('#classesRemaining').append("<p class = 'remainingCourse'>" + dept + " " + num + "</p>");
}
};

xhr.onerror = function() {
    alert('FAILURE');
};

xhr.send();
}
