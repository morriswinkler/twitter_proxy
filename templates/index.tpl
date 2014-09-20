<html>
  
  <link href="css/style.css" rel="stylesheet" type="text/css" media="all">
  <body>

    <div id="twitter-container">
      <ul class="statuses">
	{{range $element := .}}
	<li>
	  <a href="#"><img class="avatar" src={{ $element.User.profile_image_url}} width="48" height="48" alt="avatar" /></a>
	  <div class="date">	{{ formatDate $element.CreatedAt }} </div>
	  <div class="tweetTxt">


	    <strong><a href=http://twitter.com/{{ $element.User.name }}> @{{ $element.User.screen_name }} </a></strong>
	    {{ $element.Text }}     
	  </div>
	  <div class="clear"></div>
	</li>
	{{ end }}
      </ul>
    </div>
  </body>
</html>
