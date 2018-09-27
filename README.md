# aws-github-email-webhook
Webhook replacement for GitHub's deprecated email service. Runs on AWS.

## Deprecation of GitHub Services

https://developer.github.com/changes/2018-04-25-github-services-deprecation/

For those of us who used GitHub's baked-in email service to monitor pushes to
our repository, we have to look for a webhook solution. I didn't want to own
and maintain another server, so I searched for a "serverless" solution, but
didn't turn up anything. So this is a quick effort to try to duplicate what the
existing GitHub service was doing, and in a way that reduces some of the usual
sysadmin duties/worries. I figure there might be others in a similar situation,
so I'm sharing the code. Improve it as you see fit.


## Download the Binary

Here's the link to the `handler.zip` that you can just install on AWS Lambda.

     todo link here

But of course you shouldn't necessarily trust me, in which case you can build
your own binary:

## Building the Binary

Set your `GOPATH` environment variable to the `aws-github-email-webhook` directory:

    export GOPATH=/path/to/aws-github-email-webhook

To build `handler` binary:

    cd $GOPATH
    go get -u github.com/aws/aws-lambda-go
    go get -u github.com/aws/aws-sdk-go
    GOOS=linux GOARCH=amd64 go build -o github_webhook_handler handler

If on windows, build a zip file with proper executable file permissions.

    build-lambda-zip -o github_webhook_handler.zip github_webhook_handler

To build `build-lambda-zip` executable, see:
https://github.com/aws/aws-lambda-go/blob/master/README.md#for-developers-on-windows


## Install the Binary on AWS Lambda

todo

## Setup Environment Variables

todo

## References

These two AWS docs show how using Go to send an email via SES, and how to write
a lambda function:

https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/ses-example-send-email.html
https://docs.aws.amazon.com/lambda/latest/dg/go-programming-model-handler-types.html

The following article was helpful in setting up the neccessary AWS services. I'm not
an AWS expert, so I'll leave it to you to set the proper permissions and policies
for your setup.

https://medium.com/@khlbrg/sending-emails-with-go-and-aws-lambda-35c4626446ed

Google's `go-github` project already did the work to define all the GitHub events
as data types / structs with the JSON marshaling annotations. So I lifted only the
parts needed to receive and unmarshal a GitHub Push Event. That amounted to what
you'll find in `few_event_types.go`, `stringify.go`, and `timestamp.go`.

https://github.com/google/go-github

