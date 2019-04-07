package mdtohtml

var HtmlPrefix = `
  <html>
    <head>
	  <meta charset="utf-8">
	<style>
	
body {
    background-color: #FFF;
    color: #172B4D;
    font-family: -apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,Oxygen,Ubuntu,"Fira Sans","Droid Sans","Helvetica Neue",sans-serif;
    font-size: 14px;
    font-style: normal;
    font-weight: 400;
    line-height: 1.42857142857143;
    -ms-overflow-style: -ms-autohiding-scrollbar;
    text-decoration-skip: ink;
}
	 p {
    display: block;
	magin-left: 1em;
    margin-block-start: 1em;
    margin-block-end: 1em;
    margin-inline-start: 0px;
    margin-inline-end: 0px;
}

h1 {
    font-size: 1.74285714em;
    font-style: inherit;
    font-weight: 600;
    line-height: 1.25;
    letter-spacing: -.006em;
    margin-top: 24px;
    color: #172B4D;
}

h2 {
    font-size: 1.52857143em;
    font-style: inherit;
    font-weight: 500;
    letter-spacing: -.008em;
    line-height: 1.2;
    margin-top: 28px;
    color: #172B4D;
}

 h3 {
    font-size: 1.34285714em;
    font-style: inherit;
    font-weight: 600;
    line-height: 1.25;
    letter-spacing: -.006em;
    margin-top: 24px;
    color: #172B4D;
}

h4 {
    font-size: 1.1;
    font-style: inherit;
    font-weight: 500;
    line-height: 1.25;
    letter-spacing: -.006em;
    margin-top: 24px;
    color: #172B4D;
}

h5 {
    font-size: 1.0;
    font-style: inherit;
    font-weight: 400;
    line-height: 1.25;
    letter-spacing: -.006em;
    margin-top: 24px;
    color: #172B4D;
}


pre {
    background: #F4F5F7;
    border: 1px solid #DFE1E6;
    border-radius: 3px;
    overflow-x: auto;
    padding: 5px 10px;
    word-wrap: normal;
	font-family: monospace; 
	font-size: 12px;
	font-weight: 300;
	
}


blockquote  {
    margin: 12px 12px 12px 12px;
	color: Green;
	font-weight: 700;
}

a:link {
  text-decoration: none;
  color: blue;
  font-weight: 500;
}

a:visited {
  text-decoration: none;
}

a:hover {
  text-decoration: underline;
  font-weight: 550;
}

a:active {
  text-decoration: underline;
}

	</style>
    </head>
	<body>
	  <article class="markdown-body entry-content" style="padding: 30px;">
  `

var HtmlJavscriptReload = `
  <script>
  function doReload(){   
    location.reload(); 
  }

  setInterval(doReload, loopDelay);
  </script>
  `

var HtmlSuffix = `
    </article>
	</body>
	</html>
  `
