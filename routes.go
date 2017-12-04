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
        "TestLogin",
        "GET",
        "/Login/",
        Login,
    },
    Route{
        "TestSignUp",
        "GET",
        "/SignUp/",
        SignUp,
    },
    Route{
        "TestGetCoursesTaken",
        "GET",
        "/CoursesTaken/",
        GetCoursesTaken,
    },
    Route{
        "TestPostUserInformation",
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
