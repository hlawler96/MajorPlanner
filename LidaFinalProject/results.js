
$(document).ready(function () {
  console.log("The js is hooked up");
  sessionId=  window.localStorage.getItem('sessionId');
  // alert("about to try populate");

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
    $('#coursesToTake').append("<tr class = 'table'> <th>Required</th> <th> <ul id = 'strict'> </ul> </th> </tr>")
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
      $('#coursesToTake').append("<tr class= 'table'> <th>" + temp + "</th> <th> <ul id = 'loose" + i + "'> </ul> </th> </tr>")
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

//prereqs
$('#currPrereqs').append("<table id = 'prereqs' class='tableClass' ><tr class='table'><th>Requirement</th><th>Courses</th></tr></table>");
for(var i = 0; i < jsonResponse.orderOfPrereqs.length; i++){
  var course = jsonResponse.orderOfPrereqs[i].Courses[0];
  var type = jsonResponse.orderOfPrereqs[i].Type;
    if(jsonResponse.orderOfPrereqs[i].Courses.length > 2){
      var program = course.program.split(' ')[0];
      var number = course.number;

      $('#prereqs').append("<tr class= 'table'> <th>" + course.program + " " + course.number + "</th> <th> <ul id = '" + program+ "_" + number + "'> </ul> </th> </tr>");
      for(var j = 1; j < jsonResponse.orderOfPrereqs[i].Courses.length ; j++){
        $('#'+program+"_"+number).append("<li>" + jsonResponse.orderOfPrereqs[i].Courses[j].program + " " + jsonResponse.orderOfPrereqs[i].Courses[j].number + "</li>");
      }
    }else{
        $('#prereqs').append("<tr class= 'table'> <th>" + course.program + " " + course.number + "</th> <th> " + jsonResponse.orderOfPrereqs[i].Courses[1].program + " " + jsonResponse.orderOfPrereqs[i].Courses[1].number + " </th> </tr>");
    }

}



};


xhr.onerror = function() {
    alert('FAILURE');
};

xhr.send();
}
function selectAdditionalDegree() {

  var major = $('#majors :selected').val().split(' ');
  var dept = major[0];
  var type = major[major.length-1];

  $('#selectedProgram').empty();

  $('#selectedProgram').append("<table id = 'possCoursesToTake' class='tableClass' ><tr class='table'><th>Requirement</th><th>Courses</th></tr></table>");

  for(var i = 0; i < response.possiblePrograms.length; i++){

      if(response.possiblePrograms[i].dept.split(' ')[0] == dept && response.possiblePrograms[i].type == type){

        //poss strict courses
        // alert("dept: " + dept + " type: " + type);
        $('#possCoursesToTake').append("<tr class='table' ><th>Required</th><th><ul id='possStrict'></ul></th></tr>")
        for(var j = 0 ; j < response.possiblePrograms[i].strictRemainingCourses.length ; j++){
          $('#possStrict').append("<li class = 'courses'>" + response.possiblePrograms[i].strictRemainingCourses[j].program + " " + response.possiblePrograms[i].strictRemainingCourses[j].number + " </li> ")
        }

        req = ""
        k = 0;
        //poss loose courses
        for(var j = 0; j < response.possiblePrograms[i].looseRemainingCourses.length; j++){
          var dept = response.possiblePrograms[i].looseRemainingCourses[j].course.program;
          var num = response.possiblePrograms[i].looseRemainingCourses[j].course.number;
          var temp = response.possiblePrograms[i].looseRemainingCourses[j].requirement;
          // alert("Dept: " + dept + " num: " + num);
          if(temp != req){
              $('#possCoursesToTake').append("<tr class='table' ><th>" + temp + "</th><th><ul id='possLoose" + j + "'></ul></th></tr>");
              req = temp;
              k = j;
          }
          $('#possLoose' + k).append("<li class = 'courses'>" + dept + " " + num + "</li> ");
        }

        //average hours per semester
        $('#hours').append("<span>" + response.possiblePrograms[i].avgHoursPerSem + "</span>")


      }
  }


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
