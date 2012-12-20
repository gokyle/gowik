package main

// The default header to use if one is not available
defaultHead := `<!doctype html>
<html>
<head>
    <title>gowik</title>
    <meta charset="UTF-8" />
    <meta name="author" content="Kyle Isom" />
    <link href="/assets/css/bootstrap.css" rel="stylesheet" type="text/css" />
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

defaultNavBar := `
 <div class="navbar navbar-inverse navbar-fixed-top">
  <div class="navbar-inner">
    <div class="container">
      <a class="btn btn-navbar" data-toggle="collapse" data-target=".nav-collapse">
        <span class="icon-bar"></span>
        <span class="icon-bar"></span>
        <span class="icon-bar"></span>
    </a>
    <a class="brand" href="/">GoWik</a>
    <div class="nav-collapse collapse">
        <ul class="nav">
      </ul>
  </div><!--/.nav-collapse -->
</div>
</div>
</div>
`

defaultBody :=`
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

