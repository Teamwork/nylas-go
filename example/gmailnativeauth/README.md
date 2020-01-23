# Example: Native Authentication (Gmail)

This is an example file that demonstrates how to connect to Nylas using the
[Native Authentication](https://docs.nylas.com/reference#native-authentication-1)
flow. Note that different email providers have different native authentication
processes; this example project *only* works with Gmail.

In order to successfully run this example, you need to do the following things:

## Get a client ID & client secret from Nylas

To do this, make a [Nylas Developer](https://developer.nylas.com/) account.
You should see your client ID and client secret on the dashboard,
once you've logged in on the
[Nylas Developer](https://developer.nylas.com/) website.

## Get a client ID & client secret from Google

To do this, go to the
[Google Developers Console](https://console.developers.google.com)
and create a project. Then go to the "Library" section and enable the
following APIs: "Gmail API", "Contacts API", "Google Calendar API".
Then go to the "Credentials" section and create a new OAuth client ID.
Select "Web application" for the application type, and click the "Create"
button.

Check out the
[Google OAuth Setup Guide](https://docs.nylas.com/v1.0/docs/native-auth-google-oauth-setup-guide)
on the Nylas support website, for more information.

## Update the constants in `main.go`

Open the `main.go` file int his directory and replace the constants
with the values you've obtained from the steps above.

## Set the Authorized Redirect URI for Google

Once you have a HTTPS URL that points to your computer, you'll need to tell
Google about it. On the
[Google Developer Console](https://console.developers.google.com),
click on the "Credentials" section, find the OAuth client that you
already created, and click on the "edit" button on the right side.
There is a section called "Authorized redirect URIs"; this is where
you need to tell Google about your URL.

Add `http://localhost:8080`

## Run the Example

Finally, run the example like this:

```
go run main.go
```

Once the server is running following the instructions printed.
