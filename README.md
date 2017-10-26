Terraform AutoScalr Provider
=========================

- Website: https://www.autoscalr.com

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.9 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/autoscalr/terraform-provider-autoscalr` by either:

```sh
$ go get github.com/autoscalr/terraform-provider-autoscalr
```

or

```sh
$ git clone git@github.com:autoscalr/terraform-provider-autoscalr $GOPATH/src/github.com/autoscalr/terraform-provider-autoscalr
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/autoscalr/terraform-provider-autoscalr
$ ./build.sh
```

Copy the resulting terraform-provider-autoscalr file in that directory to the terraform plugins directory
for your terraform workspace located in

```sh
$TERRAFORM_WKSP/terraform.d/plugins/{ARCH}/
```


Using the provider
----------------------
## terraform-provider-autoscalr

The AutoScalr provider requires that you specify the api_key associated with your AutoScalr account.
It is recommended to specify it via the environment variable AUTOSCALR_API_KEY.
Alternatively you can specify it as a parameter when initializing the provider in your .tf file as:

provider "autoscalr" {
  api_key = "your API key value"
}

If you do not know what your API key value is for your account contact support@autoscalr.com.

The example.tf file shows how to use the autoscalr provider to extend an AWS autoscaling group.
All the attributes available on the resource are documented in website/docs/autoscalr_autoscaling_group.txt

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.8+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make bin
...
$ $GOPATH/bin/terraform-provider-autoscalr
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```
