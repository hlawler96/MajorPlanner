$(document).ready(function () {
  console.log("The js is hooked up");
  if(typeof(Storage)!=="undefined"){
    if(window.localStorage.getItem('sessionId')){
       // alert(window.localStorage.getItem('sessionId'))
    }else {
       // alert("NO Session Id");
    }
  } else{
  //alert("storage not supported by browser");
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
   //alert('CORS not supported');
   return;
 }
 // Response handlers
  xhr.onload = function() {
    var jsonResponse = JSON.parse(xhr.responseText);
    //alert(jsonResponse.sessionId);
    window.localStorage.setItem('sessionId',jsonResponse.sessionId);
    if( jsonResponse.sessionId == ""){
      //alert("Not a valid login");
    }else {
      //alert(window.localStorage.getItem('sessionId'));
    }
     window.location.href = "results.html";
  };
  xhr.onerror = function() {
      //alert('FAILURE');
  };
  xhr.send();
}

function selectAdditionalDegree(){
  additionalDegree = document.getElementById("majors").value;
  //alert(additionalDegree);

}

var dept = "";
var dept1 = "";
var type1 = "";

function deptSubmit() {
  // alert("in deptSubmit!");

  var ba_checked = document.getElementById("ba").checked;
  var bs_checked = document.getElementById("bs").checked;
  var econ_checked = document.getElementById("econ").checked;
  var math_checked = document.getElementById("math").checked;
  if(ba_checked){
    document.getElementById("second-degree-title").innerHTML = "COMP - B.A. - Select the classes you have taken.";
     dept = "COMP";
     dept1 = "COMP";
     type1 = "BA";
    var type = "BA";
  }else if(bs_checked){
    document.getElementById("second-degree-title").innerHTML = "COMP - B.S. - Select the classes you have taken.";
     dept = "COMP";
     dept1 = "COMP";
     type1 = "BS";
    var type = "BS";
  }else if(econ_checked){
    document.getElementById("second-degree-title").innerHTML = "ECON - B.A. - Select the classes you have taken.";
     dept = "ECON";
     dept1 = "ECON";
     type1 = "BA";
    var type = "BA";
  }else if(math_checked){
    document.getElementById("second-degree-title").innerHTML = "MATH - B.A. - Select the classes you have taken.";
     dept = "MATH";
     type1 = "BA";
    var type = "BA";
  }
  //inside deptSubmit

  $('#select-your-minor').empty();
  $('#minor-checkboxes').empty();
  $('#select-minor-classes-div').empty();


  $('#select-your-minor').append("Select Your Minor");
  if(ba_checked || bs_checked){
    var check1 = '<input id="math-min" onclick="minorSubmit(\'' + dept + '\')" class="menu" type="radio" name="dept" value="math-min"> MATH<br>';
    // alert("check1= " + check1);
    $('#minor-checkboxes').append(check1);
    var check2 = '<input id="stor-min" onclick="minorSubmit(\'' + dept + '\')" class="menu" type="radio" name="dept" value="stor-min"> STOR<br>';
    $('#minor-checkboxes').append(check2);
  }if(econ_checked){
    var check1 = '<input id="comp-min" onclick="minorSubmit(\'' + dept + '\')" class="menu" type="radio" name="dept" value="comp-min"> COMP<br>';
    $('#minor-checkboxes').append(check1);
    var check2 = '<input id="math-min" onclick="minorSubmit(\'' + dept + '\')" class="menu" type="radio" name="dept" value="math-min"> MATH<br>';
    $('#minor-checkboxes').append(check2);
    var check3 = '<input id="stor-min" onclick="minorSubmit(\'' + dept + '\')" class="menu" type="radio" name="dept" value="stor-min"> STOR<br>';
    $('#minor-checkboxes').append(check3);
  }if(math_checked){
    var check1 = '<input id="comp-min" onclick="minorSubmit(\'' + dept + '\')" class="menu" type="radio" name="dept" value="comp-min"> COMP<br>';
    $('#minor-checkboxes').append(check1);
    var check3 = '<input id="stor-min" onclick="minorSubmit(\'' + dept + '\')" class="menu" type="radio" name="dept" value="stor-min"> STOR<br>';
    $('#minor-checkboxes').append(check3);
  }

  API_URL = "http://ec2-18-217-72-185.us-east-2.compute.amazonaws.com:8080/Courses?dept=" + dept + "&type=" + type;
  sessionId = "";
  // alert(API_URL);
  var xhr = createCORSRequest('GET', API_URL);
  xhr.responseType = 'text';
 if (!xhr) {
   //alert('CORS not supported');
   return;
 }
 // Response handlers.
 checkBoxes = new Array();

  xhr.onload = function() {
    // alert("in function");
    $('#select-classes-div').empty();
    var jsonResponse = JSON.parse(xhr.responseText);
    var length = jsonResponse.length;
      // alert("length= " +length);
    var i = 0;
    var displayLength=0;

    for(i; i < length; i++){
      if(displayLength==5){
        $('#select-classes-div').append('<br>');
        displayLength=0;
      }
      var prog = jsonResponse[i].program;
      var num = jsonResponse[i].number;
      var id = prog + "-" + num;
      var checkbox = '<input onclick="" id="' + prog + '-' + num + '" class="menu" type="checkbox" name="dept" value="check"> '+ prog + ' ' + num + '';
      $('#select-classes-div').append(checkbox);
       checkBoxes[i] = id;
      //displayLength++;
     }

     // alert("checkBoxes = " + checkBoxes);
    // alert(jsonResponse);

    // courses = jsonResponse.Courses;
    // alert(courses);
  };

  xhr.onerror = function() {
      //alert('FAILURE');
  };

  xhr.send();



  var sems_input = '<input id="sems-left" type=number max=7 minimum=1 placeholder="Enter Number">';
  $("#sems-left-input-div").empty();
  $('#sems-left-input-div').append('Semesters Left: '+ sems_input);

  var gens_input = '<input id="gens-left" type=number max=20 minimum=1 placeholder="Enter Number">';
  $("#gen-eds-left").empty();
  $('#gen-eds-left').append('Gen Eds Left: '+ gens_input);

  var submit_classes = '<input id="submt-classes" onclick="degreeFinder()" type="submit" value="Find Me Another Degree!">';
  $("#submt-classes-div").empty();
  $('#submt-classes-div').append(submit_classes);

}
var type = "";
var dept2 = "";
var minorCheckBoxes = new Array();

function minorSubmit(dept){
  // alert("in minorSubmit");
  // alert("dept= " + dept);
   type = "Minor";
  document.getElementById("minor-degree-title").innerHTML = "Select the minor classes you have taken.";
  $('#select-minor-classes-div').empty();

  if(dept == "COMP"){

  var stor_checked = document.getElementById("stor-min").checked;
  var math_checked = document.getElementById("math-min").checked;
  // alert("math_checked= " + math_checked);

  // if(comp_checked){
  //    dept = "COMP";
  //    alert("dept= " + dept);
  // }
 if(stor_checked){
     dept = "STOR";
     dept2 = "STOR";
  } else if(math_checked){
     dept = "MATH";
     dept2 = "MATH";

     // alert("dept= " + dept);
  }
 }
 else if(dept == "ECON"){
  var comp_checked = document.getElementById("comp-min").checked;
  var stor_checked = document.getElementById("stor-min").checked;
  var math_checked = document.getElementById("math-min").checked;
  // alert("math_checked= " + math_checked);

  if(comp_checked){
     dept = "COMP";
     dept2 = "COMP";

     // alert("dept= " + dept);
  }
 if(stor_checked){
     dept = "STOR";
     dept2 = "STOR";

  } if(math_checked){
     dept = "MATH";
     dept2 = "MATH";

     // alert("dept= " + dept);

  }
}else if(dept == "MATH"){
  var comp_checked = document.getElementById("comp-min").checked;
  var stor_checked = document.getElementById("stor-min").checked;
  // var math_checked = document.getElementById("math-min").checked;
  // alert("math_checked= " + math_checked);

  if(comp_checked){
     dept = "COMP";
     dept2 = "COMP";

     // alert("dept= " + dept);
  }
 if(stor_checked){
     dept = "STOR";
     dept2 = "STOR";

  // } if(math_checked){
  //    dept = "MATH";
  //    alert("dept= " + dept);
  //
  // }
}
}


  API_URL = "http://ec2-18-217-72-185.us-east-2.compute.amazonaws.com:8080/Courses?dept=" + dept + "&type=" + type;
  sessionId = "";
  // alert(API_URL);
  var xhr = createCORSRequest('GET', API_URL);
  xhr.responseType = 'text';
 if (!xhr) {
   //alert('CORS not supported');
   return;
 }
 // Response handlers.

 // minorCheckBoxes= [];

  xhr.onload = function() {
    // alert("in function");
    var jsonResponse = JSON.parse(xhr.responseText);
    var length = jsonResponse.length;
      // alert("length= " +length);
    var i = 0;
    var displayLength=0;

    for(i; i < length; i++){
      if(displayLength==5){
        $('#select-minor-classes-div').append('<br>');
        displayLength=0;
      }
      var prog = jsonResponse[i].program;
      var num = jsonResponse[i].number;
      var id = prog + "-" + num;
      var checkbox = '<input onclick="" id="' + prog + '-' + num + '" class="menu" type="checkbox" name="dept" value="check"> '+ prog + ' ' + num + '';
      $('#select-minor-classes-div').append(checkbox);
       minorCheckBoxes[i] = id;
      displayLength++;
     }

    // alert(jsonResponse);

    // courses = jsonResponse.Courses;
    // alert(courses);
  };

  xhr.onerror = function() {
      //alert('FAILURE');
  };

  xhr.send();




}


function degreeFinder(){
  var sems_left = document.getElementById("sems-left").value;
  var gens_left = document.getElementById("gens-left").value;
  var checkedArray = new Array();
  var i = 0;
  var checkCount = 0;
  var type2 = "Minor";
  deptObj1 = {"name":dept1, "type":type1};
  if(dept2 != ""){
    deptObj2 = {"name":dept2, "type":type2};
    currDept = [deptObj1, deptObj2];
  }else{
    currDept = [deptObj1];
  }

  son = {};
  deptTaken = [];
  son ["sessionId"] = window.localStorage.getItem('sessionId');
  son ["deptTaken"] = deptTaken;
  son ["currDept"] = currDept;
  son ["semLeft"] = parseInt(sems_left);
  son ["genEdsLeft"] = parseInt(gens_left);
  for(i; i< checkBoxes.length; i++){
    var box_checked = document.getElementById(checkBoxes[i]).checked;
    if(box_checked){
      checkedArray[checkCount] = checkBoxes[i];
      checkCount++;
    }
  }
  var j = 0;
  var minorCheckCount = 0;
  var minorCheckedArray = new Array();
  for(j; j< minorCheckBoxes.length; j++){
    var box_checked = document.getElementById(minorCheckBoxes[j]).checked;
    if(box_checked){
      minorCheckedArray[minorCheckCount] = minorCheckBoxes[j];
      minorCheckCount++;
    }
  }
  var allCheckedClasses = checkedArray.concat(minorCheckedArray);
  // alert(allCheckedClasses.length);
  for(i=0; i<allCheckedClasses.length; i++){
    // alert(i);
    k = deptTaken.length;
    var found = false;
    d = allCheckedClasses[i].split('-')[0];
    n = allCheckedClasses[i].split('-')[1];
    // alert(d + " " + n);
    for(j=0; j<deptTaken.length; j++){
      if(deptTaken[j].name == d){
      found = true;
      k = j;
      }
    }
    // alert("Found:" + found);
    if(found){
      deptTaken[k].coursesTaken.push({"dept":d,"number":parseInt(n)}) ;
    }else{
      dT = {"name":d, "coursesTaken":[{"dept":d,"number":parseInt(n)}]};
      deptTaken.push(dT);
    }
  }
  //alert(JSON.stringify(son));

  API_URL = "http://localhost:8080/UserInfo/";

  var xhr = createCORSRequest('POST', API_URL);
  xhr.setRequestHeader("Content-type", "json");

 if (!xhr) {
   //alert('CORS not supported');
   return;
 }
 // Response handlers
  xhr.onload = function() {
    //alert("SUCCESS");
  };
  xhr.onerror = function() {
      //alert('FAILURE');
  };
  xhr.send(JSON.stringify(son));
  // window.location.replace("file:///Users/lahixson/Documents/GitHub/MajorPlanner/LidaFinalProject/results.html");

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
