<html>
    <head>
        <script src="https://code.jquery.com/jquery-3.5.1.min.js"></script>   
        <script src="https://cdnjs.cloudflare.com/ajax/libs/underscore.js/1.8.3/underscore-min.js"></script>
        <style>
            .error {
                color: red;
                }
        </style>
    </head> 
    <body>
        <form action="/signup" method="post">
            <div>
                {{ with .Errors.Email }}
                    <p class="error">{{ . }}</p>
                {{ end }}
                Username:<input type="text" name="username" id="username" required >
            </div> 
            &nbsp; 
            <div>
                Password:<input type="password" name="password" id="password" required>
            </div>
            &nbsp; 
            <div>
                {{ with .Errors.confirm_password }}
                    <p class="error" >{{ . }}</p>
                {{ end }}
                Confirm Password:<input type="password" name="confirm_password" id="confirm_password" required>
            </div>
            &nbsp;
            <div>
                {{ with .Errors.signup }}
                    <p class="error" >{{ . }}</p>
                {{ end }}
                <input type="submit" value="Sign Up">
            </div>
        </form>
        <form action="/gmail_login" method="GET">
            <input type="submit" value="Gmail Sign Up">
        </form>
        <form action="/login" method="GET">
            <input type="submit" value="Login Page">
        </form>
        </body>
    </body>
    <script>
        $('#username,#password,#confirm_password').change(function () {
            $("p").remove();
        });
    </script>
</html>
