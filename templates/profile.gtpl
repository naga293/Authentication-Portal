<html>
    <body>
        <h1>Profile</h1>
        <p>
            Name : {{.Name}}
        <p>
        <p>
            Email/Username : {{.Email}}
        <p>
        <p>
            City : {{.City}}
        <p>
        <form action="/profile" method="post">
            <input type="submit" value="Logout">
        </form>
        <form action="/profile_form" method="GET">
            <input type="submit" value="Edit">
        </form>
    </body>
</html>