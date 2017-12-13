
$(document).ready(function () {
  console.log("The js is hooked up");
  sessionId=  window.localStorage.getItem('sessionId');
  alert("about to try populate");
  populateCoursesTaken();
  populatePossiblePrograms();
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
  var displayLength = 0;
  var jsonResponse = JSON.parse(xhr.responseText);
  for(var i = 0; i < jsonResponse.strictRemainingCourses.length; i++){
    if(displayLength == 5){
      $('#classesRemaining').append('<br>');
      displayLength = 0;
    }
  var dept = jsonResponse.strictRemainingCourses[i].program;
  var num = jsonResponse.strictRemainingCourses[i].number;
  $('#classesRemaining').append("<span class = 'remainingCourse'>" + dept + " " + num + "</span> ");
  displayLength++;
}
};


xhr.onerror = function() {
    alert('FAILURE');
};

xhr.send();
}




function populatePossiblePrograms(){
API_URL = "http://ec2-18-217-72-185.us-east-2.compute.amazonaws.com:8080/Result/?sessionId=" + sessionId;
alert(API_URL);
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
  alert("about to write");
  for(var i = 0; i < jsonResponse.Result.length; i++){
    // if(displayLength == 5){
    //   $('#classesRemaining').append('<br>');
    //   displayLength = 0;
    // }
  var dept = jsonResponse.possiblePrograms[i].dept;
  var type = jsonResponse.possiblePrograms[i].type;
  $('#majorsAvailable').append("<span class = 'remainingCourse'>" + dept + " " + type + "</span> ");
  // displayLength++;
}
};
xhr.onerror = function() {
    alert('FAILURE');
};
xhr.send();
}
