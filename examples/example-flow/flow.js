
// ******************************************
//
// Functions
//
// ******************************************

function myProc(e) {
    e.SetBody(Math.random());
}

// ******************************************
//
// Route
//
// ******************************************

From('timer:1?period=1s')
    .SetHeader('my-header', Math.random())
    .Process().Fn(myProc)
    .To('log:1?logHeaders=true')

