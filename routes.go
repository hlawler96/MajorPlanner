package main

import "net/http"

type Route struct {
    Name        string
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
    Route{
        "Course",
        "GET",
        "/Courses",
        getCourses,
    },
    Route{
        "Test",
        "GET",
        "/",
        test,
    },
    Route{
        "Login",
        "GET",
        "/Login/",
        Login,
    },
    Route{
        "SignUp",
        "GET",
        "/SignUp/",
        SignUp,
    },
    Route{
        "GetCoursesTaken",
        "GET",
        "/CoursesTaken/",
        GetCoursesTaken,
    },
    Route{
        "PostUserInformation",
        "POST",
        "/UserInfo/",
        PostUserInformation,
    },
    Route{
        "GetResult",
        "GET",
        "/PossiblePrograms/",
        GetResult,
    },
    Route{
        "GetUserInfo",
        "GET",
        "/UserInfo/",
        GetUserInfo,
    },
}
