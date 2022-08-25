# Terraform Provider for Sendgrid

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
make build
```

In order to test the provider, you can simply run `make test`.

```sh
make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

```sh
make testacc
```

## [Documentation](docs/index.md)

The documentation is created thank's to a fork of [terraform-providers/terraform-provider-baiducloud](https://github.com/terraform-providers/terraform-provider-baiducloud/tree/master/gendocs).

## [Terraform Registry](https://registry.terraform.io/providers/anna-money/sendgrid)

## Known issues

The API KEY API is not completely documented: when you don't set scopes, you get all scopes. This is managed by the provider.

When you set one or multiple scopes, even if you don't set the scopes `sender_verification_eligible` and `2fa_required`, you will get them in the end. It's managed by the provider: if you don't add these scopes to the list of scopes, the provider does it for you.
