package main

// The default header to use if one is not available
var defaultHead = `<!doctype html>
<html>
<head>
    <title>gowik</title>
    <meta charset="UTF-8" />
    <meta name="author" content="Kyle Isom" />
    <link href="%s" rel="stylesheet" type="text/css" />
    <link rel="icon" type="image/png" href="/assets/img/favicon.png" />
    <style type="text/css">
   html,
    body {
        height: 100%;
        padding-top: 60px;

    }


    #wrap {
        min-height: 80%;
        height: auto !important;
        height: 90%;

        margin: 0 auto -60px;
    }


    #push,
    #footer {
        height: 60px;
    }
    #footer {
        background-color: #f5f5f5;
    }

    .mono {
        font-family: monospace;
    }

    h1,h2 {
        text-align:center;
    }


    @media (max-width: 767px) {
        #footer {
            margin-left: -20px;
            margin-right: -20px;
            padding-left: 20px;
            padding-right: 20px;
        }
    }
    </style>
</head>

<body>
`

var defaultNavBar = `
 <div class="navbar navbar-inverse navbar-fixed-top">
  <div class="navbar-inner">
    <div class="container">
      <a class="btn btn-navbar" data-toggle="collapse" data-target=".nav-collapse">
        <span class="icon-bar"></span>
        <span class="icon-bar"></span>
        <span class="icon-bar"></span>
    </a>
    <a class="brand" href="/">GoWik</a>
{{if .Authenticated}}
   <form class="navbar-form pull-right" action="/logout">
   <button type="submit" class="btn">Logout</button>
   </form>
{{else}}
    <form class="navbar-form pull-right" action="/login">
      <input class="span2" type="text" placeholder="Email" name="user">
      <input class="span2" type="password" placeholder="Password" name="pass">
    <button type="submit" class="btn">Log in</button>
    </form>
{{end}}
</div>
</div>
</div>
`

var defaultDisplayBody = `
<div class="container">
    <div class="row">
        <div class="span2"></div>
        <div class="span8">
            {{.Body}}
        </div>
        <div class="span2"></div>
    </div>
</div>
`

var defaultFooter = `
<div class="footer">
    <!-- the default version of this stylesheet uses the
            Bootstrap stylesheet (http://twitter.github.com/bootstrap/)
            which is licensed under the Apache License v2.0
            (http://www.apache.org/licenses/LICENSE-2.0). -->
    <div class="container">
        <p class="muted credit" style="text-align:center"><br>
                Built by <a href="http://gokyle.github.com/">Kyle Isom</a> circa 2012.<br>
            Powered by <a href="http://www.golang.org">Go</a>.</p>
            </div>
        </div>
    </body>
    </html>
`
