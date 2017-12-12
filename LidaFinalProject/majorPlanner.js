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
  window.location.replace("file:///C:/Users/farmerma/Documents/GitHub/MajorPlanner/LidaFinalProject/results.html/");
  alert(sems_left);

}













/////bottom
