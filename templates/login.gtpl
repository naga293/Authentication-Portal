<html>
    <head>
        <style>
            .error {
                color: red;
                }
        </style>
    </head> 
    </head>
    <body>
        <form action="/login" method="POST">
        <div>
            Username:<input type="text" name="username" required>
        </div>
        &nbsp; 
        <div>
            Password:<input type="password" name="password" required>
        </div>
        &nbsp; 
        <div>
            <input type="submit" value="Login">
        </div>
        </form>
        &nbsp;
        <form action="/gmail_login" method="GET">
            <div>
                <input type="submit"  value="Gmail Login">
            </div>
        </form>
        &nbsp; 
        <form action="/signup" method="GET">
            {{ with . }}
            <p class="error" >{{ . }}</p>
            {{ end }}
            <div>
                <input type="submit"  value="Sign Up Page">
            </div>
        </form>
    </body>
</html>