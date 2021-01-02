Shync - Shopify Email Template Sync Client
==========================================

Shync can checkout/download the Shopify email templates from your Shopify store admin to your local machine,
and push/sync your email templates on your local machine to Shopify.

## How to use

```
# Checkout email templates from Shopify to templates directory
% shync checkout --all templates

# Checkout one or more email templates from Shopify to templates directory
% shync checkout --id order_confirmation --id order_edited templates

# List the email template IDs
% shync ids

# Push all email templates in the templates directory to Shopify
% shync push --all templates

# Push one or more email templates in the templates directory to Shopify
% shync push --id order_confirmation --id order_edited templates
```

## Configuration

Shync needs to know some information to do it's job.
You can set them using following environment variables, or in a .env file in the current directory.

SHYNC_STORE: The Shopify store URL. eg: `https://{shopname}.myshopify.com`. Required.

SHYNC_USERNAME: The Shopify admin username. Required.

SHYNC_PASSWORD: The Shopify admin password. Required.

## Motivation

Because filling forms manually to apply a change is just not enough.

## How to contribute

