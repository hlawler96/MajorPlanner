$(document).ready(function () {
  console.log("The js is hooked up");
  if(typeof(Storage)!=="undefined"){
    if(window.localStorage.getItem('sessionId')){
       // alert(window.localStorage.getItem('sessionId'))
    }else {
       // alert("NO Session Id");
    }
  } else{
  alert("storage not supported by browser");
  }
  var fun = function(e){
      if(e.keyCode==13) loginSubmit();
  }
  $('#uname').keypress(fun);
  $('#pass').keypress(fun);
});

function loginSubmit() {
  uname = document.getElementById("uname").value;
  var pass = document.getElementById("pass").value;
  API_URL = "http://ec2-18-217-72-185.us-east-2.compute.amazonaws.com:8080/Login/?username=" + uname + "&password=" + pass;
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
      alert("Not a valid login");
    }else {
      alert(window.localStorage.getItem('sessionId'));
    }
     window.location.href = "results.html";
  };
  xhr.onerror = function() {
      alert('FAILURE');
  };
  xhr.send();
}

function selectAdditionalDegree(){
  additionalDegree = document.getElementById("majors").value;
  alert(additionalDegree);

}

function deptSubmit() {
  // alert("in deptSubmit!");
  var ba_checked = document.getElementById("ba").checked;
  var bs_checked = document.getElementById("bs").checked;
  var econ_checked = document.getElementById("econ").checked;
  var math_checked = document.getElementById("math").checked;
  if(ba_checked){
    document.getElementById("second-degree-title").innerHTML = "COMP - B.A. - Select the classes you have taken.";
    var dept = "COMP";
    var type = "BA";
  }else if(bs_checked){
    document.getElementById("second-degree-title").innerHTML = "COMP - B.S. - Select the classes you have taken.";
    var dept = "COMP";
    var type = "BS";
  }else if(econ_checked){
    document.getElementById("second-degree-title").innerHTML = "ECON - B.A. - Select the classes you have taken.";
    var dept = "ECON";
    var type = "BA";
  }else if(math_checked){
    document.getElementById("second-degree-title").innerHTML = "MATH - B.A. - Select the classes you have taken.";
    var dept = "MATH";
    var type = "BA";
  }


  API_URL = "http://ec2-18-217-72-185.us-east-2.compute.amazonaws.com:8080/Courses?dept=" + dept + "&type=" + type;
  sessionId = "";
  alert(API_URL);
  var xhr = createCORSRequest('GET', API_URL);
  xhr.responseType = 'text';
 if (!xhr) {
   alert('CORS not supported');
   return;
 }
 // Response handlers.
 checkBoxes = new Array();

  xhr.onload = function() {
    alert("in function");
    $('#select-classes-div').empty();
    var jsonResponse = JSON.parse(xhr.responseText);
    var length = jsonResponse.length;
      alert("length= " +length);
    var i = 0;
    for(i; i < length; i++){
      var prog = jsonResponse[i].program;
      var num = jsonResponse[i].number;
      var id = prog + "-" + num;
      var checkbox = '<input onclick="" id="' + prog + '-' + num + '" class="menu" type="checkbox" name="dept" value="check"> '+ prog + ' ' + num + '';
      $('#select-classes-div').append(checkbox);
       checkBoxes[i] = id;
      //displayLength++;
     }

    alert(jsonResponse);

    // courses = jsonResponse.Courses;
    // alert(courses);
  };

  xhr.onerror = function() {
      alert('FAILURE');
  };

  xhr.send();




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
  //get the id's of all of the checkboxes that are checked

  var checkedArray = new Array();
  var i = 0;
  var checkCount = 0;
  for(i; i< checkBoxes.length; i++){
    var box_checked = document.getElementById(checkBoxes[i]).checked;
    if(box_checked){
      checkedArray[checkCount] = checkBoxes[i];
      checkCount++;
    }
  }
  alert("checkedArray= " + checkedArray);


  window.location.replace("file:///Users/lahixson/Documents/GitHub/MajorPlanner/LidaFinalProject/results.html");
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
