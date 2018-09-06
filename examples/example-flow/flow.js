From('timer:1?period=1s')
    .SetHeader('my-header', Math.random())
    .To('log:1?logHeaders=true')

From('timer:2?period=5s')
    .To('log:2?logHeaders=true')
