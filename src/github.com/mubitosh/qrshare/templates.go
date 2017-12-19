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
            -moz-appearance: none;
            outline: none;
        }
        
        .img-wrap:hover .checkbox-round {
            opacity: 1;
        }     

        .checkbox-round:checked {
            background-color: #F36B13;
            opacity: 1;
        }   

        #selected-files {
            visibility: hidden;
        }

        #btn-download-cancel {
            background-color: #e7e7e7;
            border: none;
            color: black;
            padding: 6px 8px;
            text-align: center;
            text-decoration: none;
            display: inline-block;
            font-size: 14px;
            margin: 4px 2px;
            cursor: pointer;
            border-radius: 2px;
        }
        
        #btn-download-do {
            background-color: #008CBA;
            border: none;
            color: white;
            padding: 6px 8px;
            text-align: center;
            text-decoration: none;
            display: inline-block;
            font-size: 14px;
            margin: 4px 2px;
            cursor: pointer;
            border-radius: 2px;
        }
        
        #label-download {
            padding: 6px;
            padding-right: 12px;
        }
        
        #snackbar {
            visibility: hidden;
            min-width: 100px;
            margin-left: -125px;
            background-color: #333;
            color: #fff;
            text-align: center;
            border-radius: 2px;
            padding: 6px;
            position: fixed;
            z-index: 1;
            left: 50%;
            bottom: 30px;
            font-size: 17px;
        }
        
        #snackbar.show {
            visibility: visible;
            -webkit-animation: fadein 0.5s;
            animation: fadein 0.5s;
        }
        
        @-webkit-keyframes fadein {
            from {
                bottom: 0;
                opacity: 0;
            }
            to {
                bottom: 30px;
                opacity: 1;
            }
        }
        
        @keyframes fadein {
            from {
                bottom: 0;
                opacity: 0;
            }
            to {
                bottom: 30px;
                opacity: 1;
            }
        }
        
        @-webkit-keyframes fadeout {
            from {
                top: 30px;
                opacity: 1;
            }
            to {
                top: 0;
                opacity: 0;
            }
        }
        
        @keyframes fadeout {
            from {
                top: 30px;
                opacity: 1;
            }
            to {
                top: 0;
                opacity: 0;
            }
        }
    </style>
</head>

<body onunload="resetSelect()">
    <form target="_blank" method="POST" action="/zip/" onsubmit="return downloadSelected()">

    <div>
        {{$name := .Name}} 
        {{range .ChildDirs}}
        <div class="file">
            <div class="img-wrap">
                <input type="checkbox" name="download" class="checkbox-round" value="{{$name}}/{{.Name}}" onclick="countSelected(this)">
                <a class="file-url" href="{{$name}}/{{.Name}}">
                    <div class="icon">
                        <img class="icon-image" src="data:image/svg+xml;base64,{{.Icon}}">
                    </div>
                    <div class="file-name">{{.Name}}</div>
                </a>
            </div>
        </div>
        {{end}} 
        {{range .ChildFiles}}
        <div class="file">
            <div class="img-wrap">
                <input type="checkbox" name="download" class="checkbox-round" value="{{$name}}/{{.Name}}" onclick="countSelected(this)">
                <a class="file-url" href="{{$name}}/{{.Name}}">
                    <div class="icon">
                        <img class="icon-image" src="data:image/svg+xml;base64,{{.Icon}}">
                    </div>
                    <div class="file-name">{{.Name}}</div>
                </a>
            </div>
        </div>
        {{end}}
    </div>

    <div id="snackbar">
        <label id="label-download">Download selected files</label>
        <input type="button" value="Cancel" name="btn-download-cancel" id="btn-download-cancel" onclick="cancelSelect()">
        <input type="submit" value="Download" name="btn-download-do" id="btn-download-do">
    </div>

    <input type="text" id="selected-files" name="selected-files">

    </form>

    <script type="text/javascript">
    var selCounter = 0;
    var selItems = new Set();

    function countSelected(cb){
        if (cb.checked) {
            selItems.add(cb.value);
            selCounter++;
        } else {
            selItems.delete(cb.value);
            selCounter--;
        }
        
        if (selCounter < 1) {
            selItems.clear();
            selCounter = 0;
            hideToast();
        } else {
            showToast();
        }
    }

    function cancelSelect() {
        hideToast();
        resetSelect();
    }

    function showToast() {
        var x = document.getElementById("snackbar");
        x.className = "show";
    }

    function hideToast() {
        var x = document.getElementById("snackbar");
        x.className = x.className.replace("show", "");
        setTimeout(function() {
            x.visibility = "hidden";
        }, 1000);
    }

    function resetSelect() {
        selItems.clear();
        selCounter = 0;
        var x = document.getElementsByTagName("input");
        for (var i=0; i<x.length; i++) {
            if (x[i].type == "checkbox") {
                x[i].checked = false;
            }
        }
    }

    function downloadSelected() {
        var x = document.getElementById("selected-files");
        x.value = JSON.stringify(Array.from(selItems));

        console.log(x.value);

        hideToast();
        resetSelect();

        return true;
    }
    </script>
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
