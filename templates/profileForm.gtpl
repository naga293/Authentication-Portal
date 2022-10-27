<html>
    <head>
        <script src="https://code.jquery.com/jquery-3.5.1.min.js"></script>   
        <script src="https://cdnjs.cloudflare.com/ajax/libs/underscore.js/1.8.3/underscore-min.js"></script>
    </head> 
    <body>
        <h3>Profile Details</h1>
        <form action="/profile_form" method="post">
            <div>
                Email/Username: {{.Email}}
            </div>
            &nbsp; 
            <div>
                Name:<input type="text" id='name' name="name" required>
            </div>
            &nbsp; 
            <div>
                City:<input type="text" id='city' name="city"  required>
            </div>
            &nbsp; 
            <div>
                <input type="submit" value="Submit">
            </div>
        </form>
    </body>
    <script>
        var name="{{.Name}}"
        var city="{{.City}}"
        if (name != "name") {
            document.getElementById('name').value=name
        }
        if(city!="city"){
            document.getElementById('city').value=city
        }
    </script>
</html>