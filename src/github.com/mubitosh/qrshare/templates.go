package main

// Listing of files and folders.
var listingHTML = `<!DOCTYPE html>
<html>

<head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta charset='utf-8'>
    <style type="text/css">
        a:hover,
        a:visited,
        a:link,
        a:active {
            text-decoration: none!important;
            -webkit-box-shadow: none!important;
            box-shadow: none!important;
        }
        
        .file {
            width: 80px;
            word-wrap: break-word;
            display: inline-block;
            margin: 10px;
            vertical-align: top;
        }
        
        .icon {
            display: flex;
            justify-content: center;
            margin-bottom: 5px;
        }
        
        .icon-image {
            max-width: 100%;
        }
        
        .file-name {
            text-align: center;
            color: black;
        	font-family: "Arial";
        	font-size: small;
        	color: rgb(80, 80, 80);
        }

        .img-wrap {
            position: relative;
            display: inline-block;
            font-size: 0;
        }
        
        .img-wrap .checkbox-round {
            position: absolute;
            top: -2px;
            left: -8px;
            z-index: 5;
            font-weight: bold;
            cursor: pointer;
            opacity: 0;
            text-align: center;
            width: 1.3em;
            height: 1.3em;
            background-color: gray;
            border-radius: 50%;
            vertical-align: middle;
            border: 1px solid #ddd;
            -webkit-appearance: none;
            outline: none;
        }
        
        .img-wrap:hover .checkbox-round {
            opacity: 1;
        }
        
        .checkbox-round:checked {
            background-color: #F36B13;
            opacity: 1;
        }
    </style>
</head>

<body>
    <div>

        {{$name := .Name}} 

        {{range .ChildDirs}}
        <div class="file">
            <a class="file-url" href="{{$name}}/{{.Name}}">
                <div class="icon">
                    <div class="img-wrap">
                        <input type="checkbox" name="" class="checkbox-round">
                        <img class="icon-image" src="data:image/svg+xml;base64,{{.Icon}}">
                    </div>
                </div>
                <div class="file-name">{{.Name}}</div>
            </a>
        </div>
        {{end}} 

        {{range .ChildFiles}}
        <div class="file">
            <a class="file-url" href="{{$name}}/{{.Name}}">
                <div class="icon">
                    <div class="img-wrap">
                        <input type="checkbox" name="" class="checkbox-round">
                        <img class="icon-image" src="data:image/svg+xml;base64,{{.Icon}}">
                    </div>
                </div>
                <div class="file-name">{{.Name}}</div>
            </a>
        </div>
        {{end}}

    </div>
</body>

</html>`

// Page for internal server error.
var internalErrorHTML = `<!DOCTYPE html>
<html>
<head>
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta charset='utf-8'>
	<title>Not found</title>
	<style type="text/css">
		#header-500 {
            text-align: center;
            color: black;
            font-size: 2em;
        	font-family: "Arial";
        	color: rgb(80, 80, 80)
        }
	</style>
</head>
<body>
<div>
	<div id="header-500">
		<i>"Something is not right</i>
		<br>
		<i>Stop the share and try sharing again"</i>
	</div>
</div>
</body>
</html>`

// Page for not found error.
var notFoundHTML = `<!DOCTYPE html>
<html>
<head>
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta charset='utf-8'>
	<title>Not found</title>
	<style type="text/css">
		#header-404 {
            text-align: center;
            color: black;
            font-size: 2em;
        	font-family: "Arial";
        	color: rgb(80, 80, 80)
        }
	</style>
</head>
<body>
<div>
	<div id="header-404">
		<i>"The requested file or folder is not found"</i>
	</div>
</div>
</body>
</html>`
