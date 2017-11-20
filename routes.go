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
        testLogin,
    },
    Route{
        "TestSignUp",
        "GET",
        "/SignUp/",
        testSignUp,
    },
    Route{
        "TestGetCoursesTaken",
        "GET",
        "/CoursesTaken/",
        testGetCoursesTaken,
    },
    Route{
        "TestPostUserInformation",
        "POST",
        "/UserInfo/",
        testPostUserInformation,
    },
    Route{
        "GetResult",
        "GET",
        "/PossiblePrograms/",
        testGetResult,
    },
}
