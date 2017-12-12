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


//get value from major check box
//store in variable
//print in onto the website
//then add api calls to get classes and use jquery to add checkboxes




function deptSubmit() {
  alert("in deptSubmit!");
  var ba_checked = document.getElementById("ba").checked;
  alert(ba_checked);
  if(ba_checked){
    document.getElementById("second-degree-title").innerHTML = "Selected major: COMP - B.A.";
    $('ba').prop('checked' , ba_checked);
    // document.getElementById("second-degree-title").attr("checked");
  }

}
