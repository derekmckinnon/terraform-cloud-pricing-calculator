Terraform Cloud Calculator
==========================

A crummy little calculator that queries your Terraform Cloud workspaces to generate cost estimates.

The calculations are based on Terraform Cloud's updated [pricing model](https://www.hashicorp.com/products/terraform/pricing).

# Usage

After you clone this project locally, run:

```shell
go run .
```

If you do not have the `TFE_TOKEN` environment variable set, you will be prompted for it.

If you have more than one organization, you will have the chance to select one interactively.
